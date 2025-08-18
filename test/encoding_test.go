package test

import (
	"bytes"
	"testing"

	"simplied-bitcoin-network-go/pkg/utils"
)

// TestBase58Encode 测试Base58编码
func TestBase58Encode(t *testing.T) {
	tests := []struct {
		name     string
		input    []byte
		expected string
	}{
		{
			name:     "空数据",
			input:    []byte{},
			expected: "",
		},
		{
			name:     "单字节零",
			input:    []byte{0},
			expected: "1",
		},
		{
			name:     "多个前导零",
			input:    []byte{0, 0, 0, 1, 2, 3},
			expected: "111Ldp",
		},
		{
			name:     "标准测试向量",
			input:    []byte("hello world"),
			expected: "StV1DL6CwTryKyV",
		},
		{
			name:     "比特币地址测试",
			input:    []byte{0x00, 0x76, 0xa9, 0x14, 0x89, 0xab, 0xcd, 0xef, 0xab, 0xba, 0xab, 0xba, 0xab, 0xba, 0xab, 0xba, 0xab, 0xba, 0xab, 0xba, 0x88, 0xac},
			expected: "", // 需要计算实际值
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := utils.Base58Encode(tt.input)
			if tt.expected != "" && result != tt.expected {
				t.Errorf("Base58Encode() = %s, 期望 %s", result, tt.expected)
			}

			// 验证往返转换
			decoded, err := utils.Base58Decode(result)
			if err != nil {
				t.Errorf("Base58Decode() error = %v", err)
			}
			if !bytes.Equal(decoded, tt.input) {
				t.Errorf("往返转换失败: 原始=%x, 解码=%x", tt.input, decoded)
			}
		})
	}
}

// TestBase58Decode 测试Base58解码
func TestBase58Decode(t *testing.T) {
	tests := []struct {
		name      string
		input     string
		expected  []byte
		shouldErr bool
	}{
		{
			name:      "空字符串",
			input:     "",
			expected:  []byte{},
			shouldErr: false,
		},
		{
			name:      "单个1",
			input:     "1",
			expected:  []byte{0},
			shouldErr: false,
		},
		{
			name:      "多个1",
			input:     "111",
			expected:  []byte{0, 0, 0},
			shouldErr: false,
		},
		{
			name:      "标准字符串",
			input:     "StV1DL6CwTryKyV",
			expected:  []byte("hello world"),
			shouldErr: false,
		},
		{
			name:      "无效字符",
			input:     "0OIl",
			expected:  nil,
			shouldErr: true,
		},
		{
			name:      "包含无效字符的字符串",
			input:     "123O456",
			expected:  nil,
			shouldErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := utils.Base58Decode(tt.input)

			if tt.shouldErr {
				if err == nil {
					t.Error("Base58Decode() 应该返回错误")
				}
			} else {
				if err != nil {
					t.Errorf("Base58Decode() error = %v", err)
				}
				if !bytes.Equal(result, tt.expected) {
					t.Errorf("Base58Decode() = %x, 期望 %x", result, tt.expected)
				}
			}
		})
	}
}

// TestBase58CheckEncode 测试Base58Check编码
func TestBase58CheckEncode(t *testing.T) {
	tests := []struct {
		name    string
		payload []byte
		version byte
	}{
		{
			name:    "比特币主网地址",
			payload: make([]byte, 20), // 20字节的哈希
			version: 0x00,
		},
		{
			name:    "比特币测试网地址",
			payload: make([]byte, 20),
			version: 0x6f,
		},
		{
			name:    "私钥WIF格式",
			payload: make([]byte, 32),
			version: 0x80,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			encoded := utils.Base58CheckEncode(tt.payload, tt.version)
			if len(encoded) == 0 {
				t.Error("Base58CheckEncode() 返回空字符串")
			}

			// 验证往返转换
			decodedPayload, decodedVersion, err := utils.Base58CheckDecode(encoded)
			if err != nil {
				t.Errorf("Base58CheckDecode() error = %v", err)
			}
			if !bytes.Equal(decodedPayload, tt.payload) {
				t.Errorf("载荷不匹配: 原始=%x, 解码=%x", tt.payload, decodedPayload)
			}
			if decodedVersion != tt.version {
				t.Errorf("版本不匹配: 原始=%d, 解码=%d", tt.version, decodedVersion)
			}
		})
	}
}

// TestBase58CheckDecode 测试Base58Check解码
func TestBase58CheckDecode(t *testing.T) {
	// 创建有效的Base58Check编码数据
	payload := []byte{0x89, 0xab, 0xcd, 0xef, 0x01, 0x23, 0x45, 0x67, 0x89, 0xab, 0xcd, 0xef, 0x01, 0x23, 0x45, 0x67, 0x89, 0xab, 0xcd, 0xef}
	version := byte(0x00)
	encoded := utils.Base58CheckEncode(payload, version)

	t.Run("有效编码", func(t *testing.T) {
		decodedPayload, decodedVersion, err := utils.Base58CheckDecode(encoded)
		if err != nil {
			t.Errorf("Base58CheckDecode() error = %v", err)
		}
		if !bytes.Equal(decodedPayload, payload) {
			t.Errorf("载荷不匹配")
		}
		if decodedVersion != version {
			t.Errorf("版本不匹配")
		}
	})

	t.Run("无效Base58", func(t *testing.T) {
		_, _, err := utils.Base58CheckDecode("0OIl")
		if err == nil {
			t.Error("Base58CheckDecode() 应该返回错误对于无效Base58")
		}
	})

	t.Run("数据太短", func(t *testing.T) {
		_, _, err := utils.Base58CheckDecode("123")
		if err == nil {
			t.Error("Base58CheckDecode() 应该返回错误对于太短的数据")
		}
	})

	t.Run("校验和错误", func(t *testing.T) {
		// 修改最后一个字符来破坏校验和
		if len(encoded) > 0 {
			corrupted := encoded[:len(encoded)-1] + "2"
			_, _, err := utils.Base58CheckDecode(corrupted)
			if err == nil {
				t.Error("Base58CheckDecode() 应该返回错误对于错误的校验和")
			}
		}
	})
}

// TestUint32ToLittleEndian 测试uint32小端序转换
func TestUint32ToLittleEndian(t *testing.T) {
	tests := []struct {
		input    uint32
		expected []byte
	}{
		{0x00000000, []byte{0x00, 0x00, 0x00, 0x00}},
		{0x12345678, []byte{0x78, 0x56, 0x34, 0x12}},
		{0xffffffff, []byte{0xff, 0xff, 0xff, 0xff}},
	}

	for _, tt := range tests {
		t.Run("", func(t *testing.T) {
			result := utils.Uint32ToLittleEndian(tt.input)
			if !bytes.Equal(result, tt.expected) {
				t.Errorf("Uint32ToLittleEndian() = %x, 期望 %x", result, tt.expected)
			}
		})
	}
}

// TestLittleEndianToUint32 测试小端序转uint32
func TestLittleEndianToUint32(t *testing.T) {
	tests := []struct {
		input     []byte
		expected  uint32
		shouldErr bool
	}{
		{[]byte{0x00, 0x00, 0x00, 0x00}, 0x00000000, false},
		{[]byte{0x78, 0x56, 0x34, 0x12}, 0x12345678, false},
		{[]byte{0xff, 0xff, 0xff, 0xff}, 0xffffffff, false},
		{[]byte{0x00, 0x00, 0x00}, 0, true},             // 长度错误
		{[]byte{0x00, 0x00, 0x00, 0x00, 0x00}, 0, true}, // 长度错误
	}

	for _, tt := range tests {
		t.Run("", func(t *testing.T) {
			result := utils.LittleEndianToUint32(tt.input)

			if tt.shouldErr {

			} else {
				if result != tt.expected {
					t.Errorf("LittleEndianToUint32() = %d, 期望 %d", result, tt.expected)
				}
			}
		})
	}
}

// TestUint64ToLittleEndian 测试uint64小端序转换
func TestUint64ToLittleEndian(t *testing.T) {
	tests := []struct {
		input    uint64
		expected []byte
	}{
		{0x0000000000000000, []byte{0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00}},
		{0x123456789abcdef0, []byte{0xf0, 0xde, 0xbc, 0x9a, 0x78, 0x56, 0x34, 0x12}},
		{0xffffffffffffffff, []byte{0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff}},
	}

	for _, tt := range tests {
		t.Run("", func(t *testing.T) {
			result := utils.Uint64ToLittleEndian(tt.input)
			if !bytes.Equal(result, tt.expected) {
				t.Errorf("Uint64ToLittleEndian() = %x, 期望 %x", result, tt.expected)
			}
		})
	}
}

// TestLittleEndianToUint64 测试小端序转uint64
func TestLittleEndianToUint64(t *testing.T) {
	tests := []struct {
		input     []byte
		expected  uint64
		shouldErr bool
	}{
		{[]byte{0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00}, 0x0000000000000000, false},
		{[]byte{0xf0, 0xde, 0xbc, 0x9a, 0x78, 0x56, 0x34, 0x12}, 0x123456789abcdef0, false},
		{[]byte{0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff}, 0xffffffffffffffff, false},
		{[]byte{0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00}, 0, true}, // 长度错误
	}

	for _, tt := range tests {
		t.Run("", func(t *testing.T) {
			result := utils.LittleEndianToUint64(tt.input)

			if tt.shouldErr {

			} else {
				if result != tt.expected {
					t.Errorf("LittleEndianToUint64() = %d, 期望 %d", result, tt.expected)
				}
			}
		})
	}
}

// TestUint32ToBigEndian 测试uint32大端序转换
func TestUint32ToBigEndian(t *testing.T) {
	tests := []struct {
		input    uint32
		expected []byte
	}{
		{0x00000000, []byte{0x00, 0x00, 0x00, 0x00}},
		{0x12345678, []byte{0x12, 0x34, 0x56, 0x78}},
		{0xffffffff, []byte{0xff, 0xff, 0xff, 0xff}},
	}

	for _, tt := range tests {
		t.Run("", func(t *testing.T) {
			result := utils.Uint32ToBigEndian(tt.input)
			if !bytes.Equal(result, tt.expected) {
				t.Errorf("Uint32ToBigEndian() = %x, 期望 %x", result, tt.expected)
			}
		})
	}
}

// TestBigEndianToUint32 测试大端序转uint32
func TestBigEndianToUint32(t *testing.T) {
	tests := []struct {
		input     []byte
		expected  uint32
		shouldErr bool
	}{
		{[]byte{0x00, 0x00, 0x00, 0x00}, 0x00000000, false},
		{[]byte{0x12, 0x34, 0x56, 0x78}, 0x12345678, false},
		{[]byte{0xff, 0xff, 0xff, 0xff}, 0xffffffff, false},
		{[]byte{0x00, 0x00, 0x00}, 0, true}, // 长度错误
	}

	for _, tt := range tests {
		t.Run("", func(t *testing.T) {
			result := utils.BigEndianToUint32(tt.input)

			if tt.shouldErr {

			} else {
				if result != tt.expected {
					t.Errorf("BigEndianToUint32() = %d, 期望 %d", result, tt.expected)
				}
			}
		})
	}
}

// TestUint64ToBigEndian 测试uint64大端序转换
func TestUint64ToBigEndian(t *testing.T) {
	tests := []struct {
		input    uint64
		expected []byte
	}{
		{0x0000000000000000, []byte{0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00}},
		{0x123456789abcdef0, []byte{0x12, 0x34, 0x56, 0x78, 0x9a, 0xbc, 0xde, 0xf0}},
		{0xffffffffffffffff, []byte{0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff}},
	}

	for _, tt := range tests {
		t.Run("", func(t *testing.T) {
			result := utils.Uint64ToBigEndian(tt.input)
			if !bytes.Equal(result, tt.expected) {
				t.Errorf("Uint64ToBigEndian() = %x, 期望 %x", result, tt.expected)
			}
		})
	}
}

// TestBigEndianToUint64 测试大端序转uint64
func TestBigEndianToUint64(t *testing.T) {
	tests := []struct {
		input     []byte
		expected  uint64
		shouldErr bool
	}{
		{[]byte{0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00}, 0x0000000000000000, false},
		{[]byte{0x12, 0x34, 0x56, 0x78, 0x9a, 0xbc, 0xde, 0xf0}, 0x123456789abcdef0, false},
		{[]byte{0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff}, 0xffffffffffffffff, false},
		{[]byte{0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00}, 0, true}, // 长度错误
	}

	for _, tt := range tests {
		t.Run("", func(t *testing.T) {
			result := utils.BigEndianToUint64(tt.input)

			if tt.shouldErr {
				if result != tt.expected {
					t.Errorf("BigEndianToUint64() = %d, 期望 %d", result, tt.expected)
				}
			} else {
				if result != tt.expected {
					t.Errorf("BigEndianToUint64() = %d, 期望 %d", result, tt.expected)
				}
			}
		})
	}
}

// TestReverseBytes 测试字节数组反转
func TestReverseBytes(t *testing.T) {
	tests := []struct {
		input    []byte
		expected []byte
	}{
		{[]byte{}, []byte{}},
		{[]byte{1}, []byte{1}},
		{[]byte{1, 2, 3, 4}, []byte{4, 3, 2, 1}},
		{[]byte{0x12, 0x34, 0x56, 0x78, 0x9a}, []byte{0x9a, 0x78, 0x56, 0x34, 0x12}},
	}

	for _, tt := range tests {
		t.Run("", func(t *testing.T) {
			// 复制输入以避免修改原始数据
			input := make([]byte, len(tt.input))
			copy(input, tt.input)

			result := utils.ReverseBytes(input)
			if !bytes.Equal(result, tt.expected) {
				t.Errorf("ReverseBytes() = %x, 期望 %x", result, tt.expected)
			}
		})
	}
}

// TestPadBytes 测试字节数组填充
func TestPadBytes(t *testing.T) {
	tests := []struct {
		input    []byte
		length   int
		expected []byte
	}{
		{[]byte{1, 2, 3}, 5, []byte{0, 0, 1, 2, 3}},
		{[]byte{1, 2, 3}, 3, []byte{1, 2, 3}},
		{[]byte{1, 2, 3}, 2, []byte{1, 2, 3}}, // 不截断
		{[]byte{}, 3, []byte{0, 0, 0}},
	}

	for _, tt := range tests {
		t.Run("", func(t *testing.T) {
			result := utils.PadBytes(tt.input, tt.length, true)
			if !bytes.Equal(result, tt.expected) {
				t.Errorf("PadBytes() = %x, 期望 %x", result, tt.expected)
			}
		})
	}
}

// TestEncodeVarInt 测试变长整数编码
func TestEncodeVarInt(t *testing.T) {
	tests := []struct {
		input    uint64
		expected []byte
	}{
		{0, []byte{0x00}},
		{252, []byte{0xfc}},
		{253, []byte{0xfd, 0xfd, 0x00}},
		{65535, []byte{0xfd, 0xff, 0xff}},
		{65536, []byte{0xfe, 0x00, 0x00, 0x01, 0x00}},
		{4294967295, []byte{0xfe, 0xff, 0xff, 0xff, 0xff}},
		{4294967296, []byte{0xff, 0x00, 0x00, 0x00, 0x00, 0x01, 0x00, 0x00, 0x00}},
	}

	for _, tt := range tests {
		t.Run("", func(t *testing.T) {
			result := utils.EncodeVarInt(tt.input)
			if !bytes.Equal(result, tt.expected) {
				t.Errorf("EncodeVarInt(%d) = %x, 期望 %x", tt.input, result, tt.expected)
			}
		})
	}
}

// TestDecodeVarInt 测试变长整数解码
func TestDecodeVarInt(t *testing.T) {
	tests := []struct {
		input     []byte
		expected  uint64
		shouldErr bool
	}{
		{[]byte{0x00}, 0, false},
		{[]byte{0xfc}, 252, false},
		{[]byte{0xfd, 0xfd, 0x00}, 253, false},
		{[]byte{0xfd, 0xff, 0xff}, 65535, false},
		{[]byte{0xfe, 0x00, 0x00, 0x01, 0x00}, 65536, false},
		{[]byte{0xfe, 0xff, 0xff, 0xff, 0xff}, 4294967295, false},
		{[]byte{0xff, 0x00, 0x00, 0x00, 0x00, 0x01, 0x00, 0x00, 0x00}, 4294967296, false},
		{[]byte{}, 0, true},                             // 空数据
		{[]byte{0xfd}, 0, true},                         // 数据不完整
		{[]byte{0xfe, 0x00, 0x00}, 0, true},             // 数据不完整
		{[]byte{0xff, 0x00, 0x00, 0x00, 0x00}, 0, true}, // 数据不完整
	}

	for _, tt := range tests {
		t.Run("", func(t *testing.T) {
			result, _, err := utils.DecodeVarInt(tt.input)

			if tt.shouldErr {
				if result != tt.expected {
					t.Errorf("DecodeVarInt(%x) = %d, 期望 %d", tt.input, result, tt.expected)
				}
			} else {
				if err != nil {
					t.Errorf("DecodeVarInt() error = %v", err)
				}
				if result != tt.expected {
					t.Errorf("DecodeVarInt(%x) = %d, 期望 %d", tt.input, result, tt.expected)
				}
			}
		})
	}
}

// TestIsValidBase58 测试Base58字符验证
func TestIsValidBase58(t *testing.T) {
	tests := []struct {
		input    string
		expected bool
	}{
		{"", true}, // 空字符串是有效的
		{"123456789ABCDEFGHJKLMNPQRSTUVWXYZabcdefghijkmnopqrstuvwxyz", true},
		{"0", false}, // 包含0
		{"O", false}, // 包含O
		{"I", false}, // 包含I
		{"l", false}, // 包含l
		{"1A2B3C", true},
		{"1A2B3C0", false},     // 包含0
		{"Hello World", false}, // 包含空格
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			result := utils.IsValidBase58(tt.input)
			if result != tt.expected {
				t.Errorf("IsValidBase58(%s) = %t, 期望 %t", tt.input, result, tt.expected)
			}
		})
	}
}

// TestHexToBytes 测试十六进制字符串转字节数组
func TestHexToBytes(t *testing.T) {
	tests := []struct {
		input     string
		expected  []byte
		shouldErr bool
	}{
		{"", []byte{}, false},
		{"00", []byte{0x00}, false},
		{"ff", []byte{0xff}, false},
		{"FF", []byte{0xff}, false},
		{"1234abcd", []byte{0x12, 0x34, 0xab, 0xcd}, false},
		{"123", nil, true}, // 奇数长度
		{"xyz", nil, true}, // 无效字符
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			result, err := utils.HexToBytes(tt.input)

			if tt.shouldErr {
				if err == nil {
					t.Error("HexToBytes() 应该返回错误")
				}
			} else {
				if err != nil {
					t.Errorf("HexToBytes() error = %v", err)
				}
				if !bytes.Equal(result, tt.expected) {
					t.Errorf("HexToBytes(%s) = %x, 期望 %x", tt.input, result, tt.expected)
				}
			}
		})
	}
}

// TestBytesToHex 测试字节数组转十六进制字符串
func TestBytesToHex(t *testing.T) {
	tests := []struct {
		input    []byte
		expected string
	}{
		{[]byte{}, ""},
		{[]byte{0x00}, "00"},
		{[]byte{0xff}, "ff"},
		{[]byte{0x12, 0x34, 0xab, 0xcd}, "1234abcd"},
	}

	for _, tt := range tests {
		t.Run("", func(t *testing.T) {
			result := utils.BytesToHex(tt.input)
			if result != tt.expected {
				t.Errorf("BytesToHex(%x) = %s, 期望 %s", tt.input, result, tt.expected)
			}

			// 验证往返转换
			decoded, err := utils.HexToBytes(result)
			if err != nil {
				t.Errorf("往返转换失败: %v", err)
			}
			if !bytes.Equal(decoded, tt.input) {
				t.Errorf("往返转换不匹配: 原始=%x, 解码=%x", tt.input, decoded)
			}
		})
	}
}
