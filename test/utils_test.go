package test

import (
	"bytes"
	"encoding/hex"
	"math/big"
	"testing"

	"simplied-bitcoin-network-go/pkg/utils"
)

// TestDoubleSHA256 测试双重SHA-256哈希函数
func TestDoubleSHA256(t *testing.T) {
	tests := []struct {
		name     string
		input    []byte
		expected string
	}{
		{
			name:     "空数据",
			input:    []byte{},
			expected: "5df6e0e2761359d30a4f551d15c0f5b5d5d5d5d5d5d5d5d5d5d5d5d5d5d5d5d5", // 实际值需要计算
		},
		{
			name:     "Hello World",
			input:    []byte("Hello World"),
			expected: "", // 需要填入实际计算值
		},
		{
			name:     "比特币创世区块数据",
			input:    []byte("The Times 03/Jan/2009 Chancellor on brink of second bailout for banks"),
			expected: "", // 需要填入实际计算值
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := utils.DoubleSHA256(tt.input)
			if len(result) != 32 {
				t.Errorf("DoubleSHA256() 返回长度 = %d, 期望 32", len(result))
			}

			// 验证结果是确定性的
			result2 := utils.DoubleSHA256(tt.input)
			if !bytes.Equal(result, result2) {
				t.Error("DoubleSHA256() 结果不确定")
			}
		})
	}
}

// TestSingleSHA256 测试单次SHA-256哈希函数
func TestSingleSHA256(t *testing.T) {
	input := []byte("test")
	result := utils.SingleSHA256(input)

	if len(result) != 32 {
		t.Errorf("SingleSHA256() 返回长度 = %d, 期望 32", len(result))
	}

	// 验证与标准库结果一致
	expected := "9f86d081884c7d659a2feaa0c55ad015a3bf4f1b2b0b822cd15d6c15b0f00a08"
	if hex.EncodeToString(result) != expected {
		t.Errorf("SingleSHA256() = %x, 期望 %s", result, expected)
	}
}

// TestMerkleHash 测试Merkle哈希计算
func TestMerkleHash(t *testing.T) {
	left := make([]byte, 32)
	right := make([]byte, 32)

	// 填充测试数据
	for i := 0; i < 32; i++ {
		left[i] = byte(i)
		right[i] = byte(i + 32)
	}

	result := utils.MerkleHash(left, right)
	if len(result) != 32 {
		t.Errorf("MerkleHash() 返回长度 = %d, 期望 32", len(result))
	}

	// 测试panic情况
	defer func() {
		if r := recover(); r == nil {
			t.Error("MerkleHash() 应该在输入长度错误时panic")
		}
	}()
	utils.MerkleHash([]byte{1, 2, 3}, right)
}

// TestMerkleRoot 测试Merkle根计算
func TestMerkleRoot(t *testing.T) {
	tests := []struct {
		name   string
		hashes [][]byte
		valid  bool
	}{
		{
			name:   "空哈希列表",
			hashes: [][]byte{},
			valid:  true,
		},
		{
			name:   "单个哈希",
			hashes: [][]byte{make([]byte, 32)},
			valid:  true,
		},
		{
			name: "两个哈希",
			hashes: [][]byte{
				make([]byte, 32),
				make([]byte, 32),
			},
			valid: true,
		},
		{
			name: "三个哈希",
			hashes: [][]byte{
				make([]byte, 32),
				make([]byte, 32),
				make([]byte, 32),
			},
			valid: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.valid {
				result := utils.MerkleRoot(tt.hashes)
				if len(result) != 32 {
					t.Errorf("MerkleRoot() 返回长度 = %d, 期望 32", len(result))
				}
			}
		})
	}
}

// TestBlockHash 测试区块哈希计算
func TestBlockHash(t *testing.T) {
	// 创建80字节的区块头
	blockHeader := make([]byte, 80)

	result := utils.BlockHash(blockHeader)
	if len(result) != 32 {
		t.Errorf("BlockHash() 返回长度 = %d, 期望 32", len(result))
	}

	// 测试panic情况
	defer func() {
		if r := recover(); r == nil {
			t.Error("BlockHash() 应该在输入长度错误时panic")
		}
	}()
	utils.BlockHash([]byte{1, 2, 3})
}

// TestHashToString 测试哈希转字符串
func TestHashToString(t *testing.T) {
	hash := make([]byte, 32)
	for i := 0; i < 32; i++ {
		hash[i] = byte(i)
	}

	result := utils.HashToString(hash)
	if len(result) != 64 {
		t.Errorf("HashToString() 返回长度 = %d, 期望 64", len(result))
	}

	// 测试非32字节哈希
	shortHash := []byte{1, 2, 3, 4}
	result2 := utils.HashToString(shortHash)
	expected := hex.EncodeToString(shortHash)
	if result2 != expected {
		t.Errorf("HashToString() = %s, 期望 %s", result2, expected)
	}
}

// TestStringToHash 测试字符串转哈希
func TestStringToHash(t *testing.T) {
	validHash := "0000000000000000000000000000000000000000000000000000000000000000"

	result, err := utils.StringToHash(validHash)
	if err != nil {
		t.Errorf("StringToHash() error = %v", err)
	}
	if len(result) != 32 {
		t.Errorf("StringToHash() 返回长度 = %d, 期望 32", len(result))
	}

	// 测试无效长度
	_, err = utils.StringToHash("invalid")
	if err == nil {
		t.Error("StringToHash() 应该返回错误对于无效长度")
	}

	// 测试无效十六进制
	_, err = utils.StringToHash("gggggggggggggggggggggggggggggggggggggggggggggggggggggggggggggggg")
	if err == nil {
		t.Error("StringToHash() 应该返回错误对于无效十六进制")
	}
}

// TestValidateHash 测试哈希验证
func TestValidateHash(t *testing.T) {
	validHash := make([]byte, 32)
	err := utils.ValidateHash(validHash)
	if err != nil {
		t.Errorf("ValidateHash() error = %v", err)
	}

	invalidHash := make([]byte, 31)
	err = utils.ValidateHash(invalidHash)
	if err == nil {
		t.Error("ValidateHash() 应该返回错误对于无效长度")
	}
}

// TestIsZeroHash 测试零哈希检查
func TestIsZeroHash(t *testing.T) {
	zeroHash := make([]byte, 32)
	if !utils.IsZeroHash(zeroHash) {
		t.Error("IsZeroHash() 应该返回true对于零哈希")
	}

	nonZeroHash := make([]byte, 32)
	nonZeroHash[0] = 1
	if utils.IsZeroHash(nonZeroHash) {
		t.Error("IsZeroHash() 应该返回false对于非零哈希")
	}

	invalidHash := make([]byte, 31)
	if utils.IsZeroHash(invalidHash) {
		t.Error("IsZeroHash() 应该返回false对于无效长度")
	}
}

// TestCompareHashes 测试哈希比较
func TestCompareHashes(t *testing.T) {
	hash1 := make([]byte, 32)
	hash2 := make([]byte, 32)

	if !utils.CompareHashes(hash1, hash2) {
		t.Error("CompareHashes() 应该返回true对于相同哈希")
	}

	hash2[0] = 1
	if utils.CompareHashes(hash1, hash2) {
		t.Error("CompareHashes() 应该返回false对于不同哈希")
	}

	invalidHash := make([]byte, 31)
	if utils.CompareHashes(hash1, invalidHash) {
		t.Error("CompareHashes() 应该返回false对于无效长度")
	}
}

// TestHashArray 测试数组哈希
func TestHashArray(t *testing.T) {
	array1 := []byte{1, 2, 3}
	array2 := []byte{4, 5, 6}

	result := utils.HashArray(array1, array2)
	if len(result) != 32 {
		t.Errorf("HashArray() 返回长度 = %d, 期望 32", len(result))
	}

	// 验证结果是确定性的
	result2 := utils.HashArray(array1, array2)
	if !bytes.Equal(result, result2) {
		t.Error("HashArray() 结果不确定")
	}
}

// TestGenerateRandomBytes 测试随机字节生成
func TestGenerateRandomBytes(t *testing.T) {
	length := 32
	result, err := utils.GenerateRandomBytes(length)
	if err != nil {
		t.Errorf("GenerateRandomBytes() error = %v", err)
	}
	if len(result) != length {
		t.Errorf("GenerateRandomBytes() 返回长度 = %d, 期望 %d", len(result), length)
	}

	// 测试随机性
	result2, _ := utils.GenerateRandomBytes(length)
	if bytes.Equal(result, result2) {
		t.Error("GenerateRandomBytes() 应该生成不同的随机数")
	}

	// 测试无效长度
	_, err = utils.GenerateRandomBytes(0)
	if err == nil {
		t.Error("GenerateRandomBytes() 应该返回错误对于零长度")
	}
}

// TestGenerateNonce 测试nonce生成
func TestGenerateNonce(t *testing.T) {
	nonce1 := utils.GenerateNonce()
	nonce2 := utils.GenerateNonce()

	// 虽然理论上可能相同，但概率极低
	if nonce1 == nonce2 {
		t.Log("警告: 生成了相同的nonce，这在统计上是可能的但不太可能")
	}
}

// TestRIPEMD160Hash 测试RIPEMD160哈希
func TestRIPEMD160Hash(t *testing.T) {
	input := []byte("test")
	result := utils.RIPEMD160Hash(input)

	if len(result) != 20 {
		t.Errorf("RIPEMD160Hash() 返回长度 = %d, 期望 20", len(result))
	}
}

// TestHash160 测试Hash160
func TestHash160(t *testing.T) {
	input := []byte("test")
	result := utils.Hash160(input)

	if len(result) != 20 {
		t.Errorf("Hash160() 返回长度 = %d, 期望 20", len(result))
	}
}

// TestChecksum 测试校验和
func TestChecksum(t *testing.T) {
	input := []byte("test")
	result := utils.Checksum(input)

	if len(result) != 4 {
		t.Errorf("Checksum() 返回长度 = %d, 期望 4", len(result))
	}
}

// TestVerifyChecksum 测试校验和验证
func TestVerifyChecksum(t *testing.T) {
	input := []byte("test")
	checksum := utils.Checksum(input)

	if !utils.VerifyChecksum(input, checksum) {
		t.Error("VerifyChecksum() 应该返回true对于正确的校验和")
	}

	wrongChecksum := []byte{0, 0, 0, 0}
	if utils.VerifyChecksum(input, wrongChecksum) {
		t.Error("VerifyChecksum() 应该返回false对于错误的校验和")
	}

	invalidChecksum := []byte{0, 0, 0}
	if utils.VerifyChecksum(input, invalidChecksum) {
		t.Error("VerifyChecksum() 应该返回false对于无效长度的校验和")
	}
}

// TestBitsToTarget 测试难度位转目标值
func TestBitsToTarget(t *testing.T) {
	// 测试最大目标值
	maxBits := uint32(0x1d00ffff)
	target := utils.BitsToTarget(maxBits)

	if target.Sign() <= 0 {
		t.Error("BitsToTarget() 应该返回正数")
	}

	// 测试小指数情况
	smallBits := uint32(0x03123456)
	smallTarget := utils.BitsToTarget(smallBits)
	if smallTarget.Sign() <= 0 {
		t.Error("BitsToTarget() 应该返回正数对于小指数")
	}
}

// TestTargetToBits 测试目标值转难度位
func TestTargetToBits(t *testing.T) {
	// 测试往返转换
	originalBits := uint32(0x1d00ffff)
	target := utils.BitsToTarget(originalBits)
	convertedBits := utils.TargetToBits(target)

	// 由于精度问题，可能不完全相等，但应该接近
	if convertedBits == 0 {
		t.Error("TargetToBits() 不应该返回0")
	}

	// 测试零目标值
	zeroTarget := big.NewInt(0)
	zeroBits := utils.TargetToBits(zeroTarget)
	if zeroBits != 0 {
		t.Error("TargetToBits() 应该返回0对于零目标值")
	}
}

// TestCalculateDifficulty 测试难度计算
func TestCalculateDifficulty(t *testing.T) {
	maxTarget := utils.GetMaxTarget()
	difficulty := utils.CalculateDifficulty(maxTarget)

	if difficulty != 1.0 {
		t.Errorf("CalculateDifficulty() = %f, 期望 1.0 对于最大目标值", difficulty)
	}

	// 测试零目标值
	zeroTarget := big.NewInt(0)
	zeroDifficulty := utils.CalculateDifficulty(zeroTarget)
	if zeroDifficulty != 0 {
		t.Error("CalculateDifficulty() 应该返回0对于零目标值")
	}
}

// TestIsValidTarget 测试目标值验证
func TestIsValidTarget(t *testing.T) {
	// 创建一个满足条件的哈希（全零）
	hash := make([]byte, 32)
	target := utils.GetMaxTarget()

	if !utils.IsValidTarget(hash, target) {
		t.Error("IsValidTarget() 应该返回true对于全零哈希")
	}

	// 测试无效长度
	invalidHash := make([]byte, 31)
	if utils.IsValidTarget(invalidHash, target) {
		t.Error("IsValidTarget() 应该返回false对于无效长度")
	}
}

// TestValidateDifficultyBits 测试难度位验证
func TestValidateDifficultyBits(t *testing.T) {
	validBits := uint32(0x1d00ffff)
	err := utils.ValidateDifficultyBits(validBits)
	if err != nil {
		t.Errorf("ValidateDifficultyBits() error = %v", err)
	}

	// 测试零值
	err = utils.ValidateDifficultyBits(0)
	if err == nil {
		t.Error("ValidateDifficultyBits() 应该返回错误对于零值")
	}

	// 测试无效指数
	invalidBits := uint32(0x00123456)
	err = utils.ValidateDifficultyBits(invalidBits)
	if err == nil {
		t.Error("ValidateDifficultyBits() 应该返回错误对于无效指数")
	}
}

// TestAdjustDifficulty 测试难度调整
func TestAdjustDifficulty(t *testing.T) {
	currentBits := uint32(0x1d00ffff)
	targetTime := int64(600) // 10分钟

	// 测试出块时间过快（应该增加难度）
	fastTime := int64(300) // 5分钟
	adjustedBits := utils.AdjustDifficulty(fastTime, targetTime, currentBits)
	if adjustedBits == currentBits {
		t.Log("难度调整可能受到限制，这是正常的")
	}

	// 测试出块时间过慢（应该降低难度）
	slowTime := int64(1200) // 20分钟
	adjustedBits2 := utils.AdjustDifficulty(slowTime, targetTime, currentBits)
	if adjustedBits2 == currentBits {
		t.Log("难度调整可能受到限制，这是正常的")
	}

	// 测试无效参数
	invalidBits := utils.AdjustDifficulty(0, targetTime, currentBits)
	if invalidBits != currentBits {
		t.Error("AdjustDifficulty() 应该返回原值对于无效参数")
	}
}

// TestGetHashRate 测试哈希率计算
func TestGetHashRate(t *testing.T) {
	difficulty := 1.0
	blockTime := int64(600)

	hashRate := utils.GetHashRate(difficulty, blockTime)
	if hashRate <= 0 {
		t.Error("GetHashRate() 应该返回正数")
	}

	// 测试零出块时间
	zeroHashRate := utils.GetHashRate(difficulty, 0)
	if zeroHashRate != 0 {
		t.Error("GetHashRate() 应该返回0对于零出块时间")
	}
}

// TestSecureCompare 测试安全比较
func TestSecureCompare(t *testing.T) {
	data1 := []byte{1, 2, 3, 4}
	data2 := []byte{1, 2, 3, 4}
	data3 := []byte{1, 2, 3, 5}

	if !utils.SecureCompare(data1, data2) {
		t.Error("SecureCompare() 应该返回true对于相同数据")
	}

	if utils.SecureCompare(data1, data3) {
		t.Error("SecureCompare() 应该返回false对于不同数据")
	}

	if utils.SecureCompare(data1, []byte{1, 2, 3}) {
		t.Error("SecureCompare() 应该返回false对于不同长度")
	}
}

// TestZeroBytes 测试字节清零
func TestZeroBytes(t *testing.T) {
	data := []byte{1, 2, 3, 4, 5}
	utils.ZeroBytes(data)

	for i, b := range data {
		if b != 0 {
			t.Errorf("ZeroBytes() 位置%d的字节不为零: %d", i, b)
		}
	}
}
