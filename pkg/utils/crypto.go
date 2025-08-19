package utils

import (
	"crypto/rand"
	"crypto/sha256"
	"fmt"
	"math/big"

	"golang.org/x/crypto/ripemd160"
)

// GenerateRandomBytes 生成指定长度的密码学安全随机字节
func GenerateRandomBytes(length int) ([]byte, error) {
	if length <= 0 {
		return nil, fmt.Errorf("GenerateRandomBytes: 长度必须大于0")
	}

	bytes := make([]byte, length)
	_, err := rand.Read(bytes)
	if err != nil {
		return nil, fmt.Errorf("GenerateRandomBytes: 生成随机字节失败: %v", err)
	}

	return bytes, nil
}

// GenerateNonce 生成32位随机nonce值
func GenerateNonce() uint32 {
	bytes, err := GenerateRandomBytes(4)
	if err != nil {
		// 如果随机数生成失败，使用时间戳作为备选方案
		// 注意：这不是密码学安全的，仅用于紧急情况
		panic(fmt.Sprintf("GenerateNonce: 无法生成安全随机数: %v", err))
	}

	return LittleEndianToUint32(bytes)
}

// RIPEMD160Hash 计算RIPEMD160哈希
func RIPEMD160Hash(data []byte) []byte {
	hasher := ripemd160.New()
	hasher.Write(data)
	return hasher.Sum(nil)
}

// Hash160 执行SHA-256后再RIPEMD160，用于比特币地址生成
// 这是比特币地址生成的标准哈希方法
func Hash160(data []byte) []byte {
	// 先执行SHA-256
	sha256Hash := sha256.Sum256(data)

	// 再执行RIPEMD160
	return RIPEMD160Hash(sha256Hash[:])
}

// Checksum 计算4字节校验和（双重SHA-256的前4字节）
func Checksum(data []byte) []byte {
	hash := DoubleSHA256(data)
	checksum := make([]byte, 4)
	copy(checksum, hash[:4])
	return checksum
}

// VerifyChecksum 验证校验和是否正确
func VerifyChecksum(data []byte, checksum []byte) bool {
	if len(checksum) != 4 {
		return false
	}

	expectedChecksum := Checksum(data)
	for i := 0; i < 4; i++ {
		if checksum[i] != expectedChecksum[i] {
			return false
		}
	}
	return true
}

// BitsToTarget 将压缩难度位转换为256位目标值
// 比特币使用紧凑格式存储难度目标：4字节表示256位数
func BitsToTarget(bits uint32) *big.Int {
	// 提取指数和尾数
	exponent := bits >> 24
	mantissa := bits & 0x00ffffff

	// 检查有效性
	if exponent <= 3 {
		return big.NewInt(int64(mantissa >> (8 * (3 - exponent))))
	}

	// 构造目标值
	target := big.NewInt(int64(mantissa))
	target.Lsh(target, uint(8*(exponent-3)))

	// 确保不超过最大目标值
	maxTarget := GetMaxTarget()
	if target.Cmp(maxTarget) > 0 {
		return maxTarget
	}

	return target
}

// TargetToBits 将256位目标值转换为压缩难度位
func TargetToBits(target *big.Int) uint32 {
	if target.Sign() <= 0 {
		return 0
	}

	// 获取目标值的字节表示
	targetBytes := target.Bytes()

	// 计算字节长度
	length := len(targetBytes)
	if length == 0 {
		return 0
	}

	var mantissa uint32
	var exponent uint32

	if length <= 3 {
		// 目标值较小，可以直接编码
		mantissa = 0
		for i := 0; i < length; i++ {
			mantissa |= uint32(targetBytes[length-1-i]) << (8 * uint(i))
		}
		exponent = uint32(length)
	} else {
		// 目标值较大，需要截取前3字节
		mantissa = uint32(targetBytes[0])<<16 | uint32(targetBytes[1])<<8 | uint32(targetBytes[2])
		exponent = uint32(length)

		// 如果最高位为1，需要调整以避免符号位问题
		if targetBytes[0] >= 0x80 {
			mantissa >>= 8
			exponent++
		}
	}

	// 组合指数和尾数
	return (exponent << 24) | mantissa
}

// GetMaxTarget 获取最大目标值（最低难度）
func GetMaxTarget() *big.Int {
	// 直接构造最大目标值，避免调用BitsToTarget造成循环依赖
	// 0x1d00ffff 对应的目标值是 0x00ffff * 2^(8*(0x1d-3))
	// = 0x00ffff * 2^(8*26) = 0x00ffff * 2^208
	target := big.NewInt(0x00ffff)
	target.Lsh(target, 208) // 左移208位
	return target
}

// CalculateDifficulty 根据目标值计算难度
func CalculateDifficulty(target *big.Int) float64 {
	if target.Sign() <= 0 {
		return 0
	}

	maxTarget := GetMaxTarget()

	// 难度 = 最大目标值 / 当前目标值
	difficulty := new(big.Rat)
	difficulty.SetFrac(maxTarget, target)

	result, _ := difficulty.Float64()
	return result
}

// IsValidTarget 验证哈希是否满足目标难度要求
func IsValidTarget(hash []byte, target *big.Int) bool {
	if len(hash) != 32 {
		return false
	}

	// 将哈希转换为大整数（小端字节序）
	hashInt := new(big.Int)

	// 反转字节序（比特币内部使用小端，但big.Int期望大端）
	reversed := make([]byte, 32)
	for i := 0; i < 32; i++ {
		reversed[i] = hash[31-i]
	}

	hashInt.SetBytes(reversed)

	// 检查哈希是否小于等于目标值
	return hashInt.Cmp(target) <= 0
}

// GetDifficultyFromBits 从压缩难度位直接计算难度值
func GetDifficultyFromBits(bits uint32) float64 {
	target := BitsToTarget(bits)
	return CalculateDifficulty(target)
}

// ValidateDifficultyBits 验证难度位格式是否有效
func ValidateDifficultyBits(bits uint32) error {
	if bits == 0 {
		return fmt.Errorf("ValidateDifficultyBits: 难度位不能为0")
	}

	exponent := bits >> 24
	mantissa := bits & 0x00ffffff

	// 检查指数范围
	if exponent == 0 || exponent > 32 {
		return fmt.Errorf("ValidateDifficultyBits: 无效的指数值: %d", exponent)
	}

	// 检查尾数
	if mantissa == 0 {
		return fmt.Errorf("ValidateDifficultyBits: 尾数不能为0")
	}

	// 检查是否超过最大目标
	target := BitsToTarget(bits)
	maxTarget := GetMaxTarget()
	if target.Cmp(maxTarget) > 0 {
		return fmt.Errorf("ValidateDifficultyBits: 目标值超过最大值")
	}

	return nil
}

// CompareTargets 比较两个目标值的大小
// 返回值：-1 表示 target1 < target2，0 表示相等，1 表示 target1 > target2
func CompareTargets(target1, target2 *big.Int) int {
	return target1.Cmp(target2)
}

// AdjustDifficulty 根据实际出块时间调整难度
// actualTime: 实际出块时间（秒）
// targetTime: 目标出块时间（秒）
// currentBits: 当前难度位
func AdjustDifficulty(actualTime, targetTime int64, currentBits uint32) uint32 {
	if actualTime <= 0 || targetTime <= 0 {
		return currentBits
	}

	// 计算时间比率
	timeRatio := new(big.Rat)
	timeRatio.SetFrac(big.NewInt(actualTime), big.NewInt(targetTime))

	// 获取当前目标值
	currentTarget := BitsToTarget(currentBits)

	// 调整目标值：新目标 = 当前目标 * (实际时间 / 目标时间)
	newTarget := new(big.Int)
	newTarget.Mul(currentTarget, big.NewInt(actualTime))
	newTarget.Div(newTarget, big.NewInt(targetTime))

	// 限制调整幅度（比特币限制为4倍）
	maxTarget := GetMaxTarget()
	minAdjustment := new(big.Int)
	maxAdjustment := new(big.Int)

	minAdjustment.Div(currentTarget, big.NewInt(4)) // 最多增加4倍难度
	maxAdjustment.Mul(currentTarget, big.NewInt(4)) // 最多降低4倍难度

	if newTarget.Cmp(minAdjustment) < 0 {
		newTarget = minAdjustment
	}
	if newTarget.Cmp(maxAdjustment) > 0 {
		newTarget = maxAdjustment
	}
	if newTarget.Cmp(maxTarget) > 0 {
		newTarget = maxTarget
	}

	return TargetToBits(newTarget)
}

// GetHashRate 根据难度和出块时间估算网络哈希率（H/s）
func GetHashRate(difficulty float64, blockTime int64) float64 {
	if blockTime <= 0 {
		return 0
	}

	// 哈希率 = 难度 * 2^32 / 出块时间
	// 2^32 是因为难度基于32位nonce空间
	hashesPerSecond := difficulty * 4294967296.0 / float64(blockTime)
	return hashesPerSecond
}

// SecureCompare 安全比较两个字节数组，防止时序攻击
func SecureCompare(a, b []byte) bool {
	if len(a) != len(b) {
		return false
	}

	var result byte
	for i := 0; i < len(a); i++ {
		result |= a[i] ^ b[i]
	}

	return result == 0
}

// ZeroBytes 安全清零字节数组
func ZeroBytes(data []byte) {
	for i := range data {
		data[i] = 0
	}
}
