package utils

import (
	"encoding/binary"
	"fmt"
	"math/big"
	"strings"
)

// Base58字符集，去除了容易混淆的字符 0, O, I, l
const base58Alphabet = "123456789ABCDEFGHJKLMNPQRSTUVWXYZabcdefghijkmnopqrstuvwxyz"

// base58DecodeMap 用于快速查找字符对应的值
var base58DecodeMap = make(map[rune]int)

func init() {
	// 初始化Base58解码映射表
	for i, char := range base58Alphabet {
		base58DecodeMap[char] = i
	}
}

// Base58Encode 将字节数组编码为Base58字符串
func Base58Encode(data []byte) string {
	if len(data) == 0 {
		return ""
	}

	// 计算前导零字节数
	leadingZeros := 0
	for i := 0; i < len(data) && data[i] == 0; i++ {
		leadingZeros++
	}

	// 将字节数组转换为大整数
	num := new(big.Int)
	num.SetBytes(data)

	// Base58编码
	var encoded []byte
	base := big.NewInt(58)
	zero := big.NewInt(0)

	for num.Cmp(zero) > 0 {
		remainder := new(big.Int)
		num.DivMod(num, base, remainder)
		encoded = append(encoded, base58Alphabet[remainder.Int64()])
	}

	// 反转结果
	for i, j := 0, len(encoded)-1; i < j; i, j = i+1, j-1 {
		encoded[i], encoded[j] = encoded[j], encoded[i]
	}

	// 添加前导零对应的'1'字符
	result := strings.Repeat("1", leadingZeros) + string(encoded)

	return result
}

// Base58Decode 将Base58字符串解码为字节数组
func Base58Decode(encoded string) ([]byte, error) {
	if len(encoded) == 0 {
		return []byte{}, nil
	}

	// 验证字符有效性
	for _, char := range encoded {
		if _, exists := base58DecodeMap[char]; !exists {
			return nil, fmt.Errorf("Base58Decode: 无效的Base58字符: %c", char)
		}
	}

	// 计算前导'1'字符数
	leadingOnes := 0
	for i := 0; i < len(encoded) && encoded[i] == '1'; i++ {
		leadingOnes++
	}

	// 将Base58字符串转换为大整数
	num := big.NewInt(0)
	base := big.NewInt(58)

	for _, char := range encoded {
		value := base58DecodeMap[char]
		num.Mul(num, base)
		num.Add(num, big.NewInt(int64(value)))
	}

	// 转换为字节数组
	decoded := num.Bytes()

	// 添加前导零字节
	result := make([]byte, leadingOnes+len(decoded))
	copy(result[leadingOnes:], decoded)

	return result, nil
}

// Base58CheckEncode 使用Base58Check编码（带4字节校验和）
func Base58CheckEncode(data []byte, version byte) string {
	// 添加版本字节
	payload := make([]byte, 1+len(data))
	payload[0] = version
	copy(payload[1:], data)

	// 计算校验和
	checksum := Checksum(payload)

	// 组合payload和校验和
	fullPayload := make([]byte, len(payload)+4)
	copy(fullPayload, payload)
	copy(fullPayload[len(payload):], checksum)

	// Base58编码
	return Base58Encode(fullPayload)
}

// Base58CheckDecode 解码Base58Check编码并验证校验和
func Base58CheckDecode(encoded string) ([]byte, byte, error) {
	// Base58解码
	decoded, err := Base58Decode(encoded)
	if err != nil {
		return nil, 0, fmt.Errorf("Base58CheckDecode: Base58解码失败: %v", err)
	}

	// 检查最小长度（版本字节 + 校验和）
	if len(decoded) < 5 {
		return nil, 0, fmt.Errorf("Base58CheckDecode: 数据长度不足，至少需要5字节")
	}

	// 分离payload和校验和
	payloadLen := len(decoded) - 4
	payload := decoded[:payloadLen]
	checksum := decoded[payloadLen:]

	// 验证校验和
	if !VerifyChecksum(payload, checksum) {
		return nil, 0, fmt.Errorf("Base58CheckDecode: 校验和验证失败")
	}

	// 提取版本字节和数据
	if len(payload) < 1 {
		return nil, 0, fmt.Errorf("Base58CheckDecode: payload长度不足")
	}

	version := payload[0]
	data := payload[1:]

	return data, version, nil
}

// LittleEndianToUint32 将小端字节序转换为uint32
func LittleEndianToUint32(data []byte) uint32 {
	if len(data) < 4 {
		// 不足4字节时，用零填充
		padded := make([]byte, 4)
		copy(padded, data)
		return binary.LittleEndian.Uint32(padded)
	}
	return binary.LittleEndian.Uint32(data[:4])
}

// Uint32ToLittleEndian 将uint32转换为小端字节序
func Uint32ToLittleEndian(value uint32) []byte {
	bytes := make([]byte, 4)
	binary.LittleEndian.PutUint32(bytes, value)
	return bytes
}

// LittleEndianToUint64 将小端字节序转换为uint64
func LittleEndianToUint64(data []byte) uint64 {
	if len(data) < 8 {
		// 不足8字节时，用零填充
		padded := make([]byte, 8)
		copy(padded, data)
		return binary.LittleEndian.Uint64(padded)
	}
	return binary.LittleEndian.Uint64(data[:8])
}

// Uint64ToLittleEndian 将uint64转换为小端字节序
func Uint64ToLittleEndian(value uint64) []byte {
	bytes := make([]byte, 8)
	binary.LittleEndian.PutUint64(bytes, value)
	return bytes
}

// BigEndianToUint32 将大端字节序转换为uint32
func BigEndianToUint32(data []byte) uint32 {
	if len(data) < 4 {
		// 不足4字节时，用零填充
		padded := make([]byte, 4)
		copy(padded[4-len(data):], data)
		return binary.BigEndian.Uint32(padded)
	}
	return binary.BigEndian.Uint32(data[:4])
}

// Uint32ToBigEndian 将uint32转换为大端字节序
func Uint32ToBigEndian(value uint32) []byte {
	bytes := make([]byte, 4)
	binary.BigEndian.PutUint32(bytes, value)
	return bytes
}

// BigEndianToUint64 将大端字节序转换为uint64
func BigEndianToUint64(data []byte) uint64 {
	if len(data) < 8 {
		// 不足8字节时，用零填充
		padded := make([]byte, 8)
		copy(padded[8-len(data):], data)
		return binary.BigEndian.Uint64(padded)
	}
	return binary.BigEndian.Uint64(data[:8])
}

// Uint64ToBigEndian 将uint64转换为大端字节序
func Uint64ToBigEndian(value uint64) []byte {
	bytes := make([]byte, 8)
	binary.BigEndian.PutUint64(bytes, value)
	return bytes
}

// ReverseBytes 反转字节数组
func ReverseBytes(data []byte) []byte {
	reversed := make([]byte, len(data))
	for i, j := 0, len(data)-1; i <= j; i, j = i+1, j-1 {
		reversed[i] = data[j]
		reversed[j] = data[i]
	}
	return reversed
}

// PadBytes 填充字节数组到指定长度
func PadBytes(data []byte, length int, padLeft bool) []byte {
	if len(data) >= length {
		return data[:length]
	}

	padded := make([]byte, length)
	if padLeft {
		copy(padded[length-len(data):], data)
	} else {
		copy(padded, data)
	}

	return padded
}

// VarIntEncode 编码变长整数（比特币协议使用）
func VarIntEncode(value uint64) []byte {
	if value < 0xfd {
		return []byte{byte(value)}
	} else if value <= 0xffff {
		result := make([]byte, 3)
		result[0] = 0xfd
		binary.LittleEndian.PutUint16(result[1:], uint16(value))
		return result
	} else if value <= 0xffffffff {
		result := make([]byte, 5)
		result[0] = 0xfe
		binary.LittleEndian.PutUint32(result[1:], uint32(value))
		return result
	} else {
		result := make([]byte, 9)
		result[0] = 0xff
		binary.LittleEndian.PutUint64(result[1:], value)
		return result
	}
}

// VarIntDecode 解码变长整数
func VarIntDecode(data []byte) (uint64, int, error) {
	if len(data) == 0 {
		return 0, 0, fmt.Errorf("VarIntDecode: 数据为空")
	}

	first := data[0]

	if first < 0xfd {
		return uint64(first), 1, nil
	} else if first == 0xfd {
		if len(data) < 3 {
			return 0, 0, fmt.Errorf("VarIntDecode: 数据长度不足，需要3字节")
		}
		value := binary.LittleEndian.Uint16(data[1:3])
		return uint64(value), 3, nil
	} else if first == 0xfe {
		if len(data) < 5 {
			return 0, 0, fmt.Errorf("VarIntDecode: 数据长度不足，需要5字节")
		}
		value := binary.LittleEndian.Uint32(data[1:5])
		return uint64(value), 5, nil
	} else { // first == 0xff
		if len(data) < 9 {
			return 0, 0, fmt.Errorf("VarIntDecode: 数据长度不足，需要9字节")
		}
		value := binary.LittleEndian.Uint64(data[1:9])
		return value, 9, nil
	}
}

// VarIntSize 计算变长整数编码后的字节数
func VarIntSize(value uint64) int {
	if value < 0xfd {
		return 1
	} else if value <= 0xffff {
		return 3
	} else if value <= 0xffffffff {
		return 5
	} else {
		return 9
	}
}

// CompactSizeEncode 编码紧凑大小（与VarInt相同，但用于不同场景）
func CompactSizeEncode(value uint64) []byte {
	return VarIntEncode(value)
}

// CompactSizeDecode 解码紧凑大小
func CompactSizeDecode(data []byte) (uint64, int, error) {
	return VarIntDecode(data)
}

// ValidateBase58 验证字符串是否为有效的Base58编码
func ValidateBase58(encoded string) error {
	if len(encoded) == 0 {
		return fmt.Errorf("ValidateBase58: 空字符串")
	}

	for i, char := range encoded {
		if _, exists := base58DecodeMap[char]; !exists {
			return fmt.Errorf("ValidateBase58: 位置%d处包含无效字符: %c", i, char)
		}
	}

	return nil
}

// IsValidBase58 检查字符串是否为有效的Base58编码
func IsValidBase58(encoded string) bool {
	return ValidateBase58(encoded) == nil
}

// HexToBytes 将十六进制字符串转换为字节数组
func HexToBytes(hexStr string) ([]byte, error) {
	// 移除可能的0x前缀
	if strings.HasPrefix(hexStr, "0x") || strings.HasPrefix(hexStr, "0X") {
		hexStr = hexStr[2:]
	}

	// 确保长度为偶数
	if len(hexStr)%2 != 0 {
		hexStr = "0" + hexStr
	}

	bytes := make([]byte, len(hexStr)/2)
	for i := 0; i < len(hexStr); i += 2 {
		var b byte
		_, err := fmt.Sscanf(hexStr[i:i+2], "%02x", &b)
		if err != nil {
			return nil, fmt.Errorf("HexToBytes: 无效的十六进制字符: %s", hexStr[i:i+2])
		}
		bytes[i/2] = b
	}

	return bytes, nil
}

// BytesToHex 将字节数组转换为十六进制字符串
func BytesToHex(data []byte) string {
	return fmt.Sprintf("%x", data)
}

// BytesToHexWithPrefix 将字节数组转换为带0x前缀的十六进制字符串
func BytesToHexWithPrefix(data []byte) string {
	return fmt.Sprintf("0x%x", data)
}

// EncodeVarInt 编码变长整数
func EncodeVarInt(value uint64) []byte {
	return VarIntEncode(value)
}

// DecodeVarInt 解码变长整数
func DecodeVarInt(data []byte) (uint64, int, error) {
	return VarIntDecode(data)
}
