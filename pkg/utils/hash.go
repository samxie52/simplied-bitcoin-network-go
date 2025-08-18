package utils

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
)

// DoubleSHA256 执行双重SHA-256哈希，这是比特币网络的标准哈希方法
// 用于区块哈希、交易哈希等关键场景
func DoubleSHA256(data []byte) []byte {
	first := sha256.Sum256(data)
	second := sha256.Sum256(first[:])
	return second[:]
}

// SingleSHA256 执行单次SHA-256哈希
func SingleSHA256(data []byte) []byte {
	hash := sha256.Sum256(data)
	return hash[:]
}

// MerkleHash 计算两个哈希的Merkle父节点哈希
// 这是构建Merkle树的基础函数
func MerkleHash(left, right []byte) []byte {
	if len(left) != 32 || len(right) != 32 {
		panic("MerkleHash: 输入哈希长度必须为32字节")
	}

	// 连接左右哈希
	combined := make([]byte, 64)
	copy(combined[:32], left)
	copy(combined[32:], right)

	// 返回双重SHA-256哈希
	return DoubleSHA256(combined)
}

// MerkleRoot 计算交易哈希列表的Merkle根
// 实现比特币标准的Merkle树算法
func MerkleRoot(hashes [][]byte) []byte {
	if len(hashes) == 0 {
		// 空交易列表返回全零哈希
		return make([]byte, 32)
	}

	if len(hashes) == 1 {
		return hashes[0]
	}

	// 验证所有哈希长度
	for i, hash := range hashes {
		if len(hash) != 32 {
			panic(fmt.Sprintf("MerkleRoot: 哈希%d长度无效，期望32字节，实际%d字节", i, len(hash)))
		}
	}

	// 复制哈希列表避免修改原始数据
	currentLevel := make([][]byte, len(hashes))
	for i, hash := range hashes {
		currentLevel[i] = make([]byte, 32)
		copy(currentLevel[i], hash)
	}

	// 递归构建Merkle树
	for len(currentLevel) > 1 {
		nextLevel := make([][]byte, 0, (len(currentLevel)+1)/2)

		for i := 0; i < len(currentLevel); i += 2 {
			var left, right []byte
			left = currentLevel[i]

			if i+1 < len(currentLevel) {
				right = currentLevel[i+1]
			} else {
				// 奇数个节点时，最后一个节点与自己配对
				right = currentLevel[i]
			}

			parentHash := MerkleHash(left, right)
			nextLevel = append(nextLevel, parentHash)
		}

		currentLevel = nextLevel
	}

	return currentLevel[0]
}

// BlockHash 计算区块头哈希
// 区块头包含：版本、前块哈希、Merkle根、时间戳、难度目标、Nonce
func BlockHash(blockHeader []byte) []byte {
	if len(blockHeader) != 80 {
		panic(fmt.Sprintf("BlockHash: 区块头长度无效，期望80字节，实际%d字节", len(blockHeader)))
	}

	return DoubleSHA256(blockHeader)
}

// HashToString 将哈希转换为十六进制字符串（小端格式显示）
func HashToString(hash []byte) string {
	if len(hash) != 32 {
		return hex.EncodeToString(hash)
	}

	// 比特币显示哈希时使用小端字节序（反转显示）
	reversed := make([]byte, 32)
	for i := 0; i < 32; i++ {
		reversed[i] = hash[31-i]
	}

	return hex.EncodeToString(reversed)
}

// StringToHash 将十六进制字符串转换为哈希（处理小端格式）
func StringToHash(hashStr string) ([]byte, error) {
	if len(hashStr) != 64 {
		return nil, fmt.Errorf("StringToHash: 哈希字符串长度无效，期望64字符，实际%d字符", len(hashStr))
	}

	// 解码十六进制字符串
	decoded, err := hex.DecodeString(hashStr)
	if err != nil {
		return nil, fmt.Errorf("StringToHash: 解码十六进制失败: %v", err)
	}

	// 反转字节序（从显示格式转为内部格式）
	hash := make([]byte, 32)
	for i := 0; i < 32; i++ {
		hash[i] = decoded[31-i]
	}

	return hash, nil
}

// ValidateHash 验证哈希格式是否正确
func ValidateHash(hash []byte) error {
	if len(hash) != 32 {
		return fmt.Errorf("ValidateHash: 哈希长度无效，期望32字节，实际%d字节", len(hash))
	}
	return nil
}

// IsZeroHash 检查是否为全零哈希
func IsZeroHash(hash []byte) bool {
	if len(hash) != 32 {
		return false
	}

	for _, b := range hash {
		if b != 0 {
			return false
		}
	}
	return true
}

// CompareHashes 比较两个哈希是否相等
func CompareHashes(hash1, hash2 []byte) bool {
	if len(hash1) != 32 || len(hash2) != 32 {
		return false
	}

	for i := 0; i < 32; i++ {
		if hash1[i] != hash2[i] {
			return false
		}
	}
	return true
}

// HashArray 计算字节数组的哈希
func HashArray(arrays ...[]byte) []byte {
	totalLen := 0
	for _, arr := range arrays {
		totalLen += len(arr)
	}

	combined := make([]byte, 0, totalLen)
	for _, arr := range arrays {
		combined = append(combined, arr...)
	}

	return DoubleSHA256(combined)
}

// MerkleProof 表示Merkle路径证明
type MerkleProof struct {
	Index  int      // 叶子节点索引
	Hashes [][]byte // 证明路径上的哈希
	Flags  []bool   // 标记哈希在左侧还是右侧
}

// GenerateMerkleProof 为指定索引的交易生成Merkle路径证明
func GenerateMerkleProof(hashes [][]byte, index int) (*MerkleProof, error) {
	if index < 0 || index >= len(hashes) {
		return nil, fmt.Errorf("GenerateMerkleProof: 索引超出范围")
	}

	if len(hashes) == 0 {
		return nil, fmt.Errorf("GenerateMerkleProof: 空哈希列表")
	}

	proof := &MerkleProof{
		Index:  index,
		Hashes: make([][]byte, 0),
		Flags:  make([]bool, 0),
	}

	currentLevel := make([][]byte, len(hashes))
	copy(currentLevel, hashes)
	currentIndex := index

	// 从叶子节点向根节点构建证明路径
	for len(currentLevel) > 1 {
		nextLevel := make([][]byte, 0, (len(currentLevel)+1)/2)

		// 找到当前节点的兄弟节点
		var siblingIndex int
		var isLeft bool

		if currentIndex%2 == 0 {
			// 当前节点是左子节点
			siblingIndex = currentIndex + 1
			isLeft = false
		} else {
			// 当前节点是右子节点
			siblingIndex = currentIndex - 1
			isLeft = true
		}

		// 添加兄弟节点哈希到证明路径
		if siblingIndex < len(currentLevel) {
			proof.Hashes = append(proof.Hashes, currentLevel[siblingIndex])
		} else {
			// 奇数个节点时，最后一个节点与自己配对
			proof.Hashes = append(proof.Hashes, currentLevel[currentIndex])
		}
		proof.Flags = append(proof.Flags, isLeft)

		// 构建下一层
		for i := 0; i < len(currentLevel); i += 2 {
			var left, right []byte
			left = currentLevel[i]

			if i+1 < len(currentLevel) {
				right = currentLevel[i+1]
			} else {
				right = currentLevel[i]
			}

			parentHash := MerkleHash(left, right)
			nextLevel = append(nextLevel, parentHash)
		}

		currentLevel = nextLevel
		currentIndex = currentIndex / 2
	}

	return proof, nil
}

// VerifyMerkleProof 验证Merkle路径证明
func VerifyMerkleProof(leafHash []byte, proof *MerkleProof, merkleRoot []byte) bool {
	if proof == nil {
		return false
	}

	if len(proof.Hashes) != len(proof.Flags) {
		return false
	}

	currentHash := make([]byte, 32)
	copy(currentHash, leafHash)

	// 沿着证明路径向上计算
	for i, siblingHash := range proof.Hashes {
		if proof.Flags[i] {
			// 兄弟节点在左侧
			currentHash = MerkleHash(siblingHash, currentHash)
		} else {
			// 兄弟节点在右侧
			currentHash = MerkleHash(currentHash, siblingHash)
		}
	}

	return CompareHashes(currentHash, merkleRoot)
}
