// Package blockchain 实现了简化比特币网络的区块链核心功能
// 本文件包含Merkle树的完整实现，包括树构建、证明生成、验证和优化功能
// Merkle树是区块链中用于高效验证交易存在性的关键数据结构
package blockchain

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"sync"

	"simplied-bitcoin-network-go/pkg/utils"
)

// MerkleNode Merkle树节点结构
// 表示Merkle树中的单个节点，包含哈希值和树结构信息
type MerkleNode struct {
	Hash    [32]byte    // 节点的哈希值，32字节SHA-256哈希
	Left    *MerkleNode // 左子节点指针
	Right   *MerkleNode // 右子节点指针
	Parent  *MerkleNode // 父节点指针，用于向上遍历
	IsLeaf  bool        // 是否为叶子节点标志
	TxIndex int         // 交易索引，仅对叶子节点有效
}

// NewMerkleNode 创建新的Merkle树节点
//
// 功能说明：
// 创建一个新的Merkle节点，设置哈希值和子节点
// 自动设置子节点的父节点指针以维护树结构完整性
//
// 参数：
// hash [32]byte - 节点的哈希值
// left *MerkleNode - 左子节点，可为nil
// right *MerkleNode - 右子节点，可为nil
//
// 返回值：
// *MerkleNode - 新创建的节点实例
func NewMerkleNode(hash [32]byte, left, right *MerkleNode) *MerkleNode {
	node := &MerkleNode{
		Hash:   hash,
		Left:   left,
		Right:  right,
		IsLeaf: left == nil && right == nil,
	}

	// 设置子节点的父节点指针
	if left != nil {
		left.Parent = node
	}
	if right != nil {
		right.Parent = node
	}

	return node
}

// IsRoot 判断节点是否为根节点
//
// 功能说明：
// 通过检查父节点指针是否为nil来判断是否为根节点
//
// 返回值：
// bool - true表示是根节点，false表示不是
func (n *MerkleNode) IsRoot() bool {
	return n.Parent == nil
}

// GetSibling 获取节点的兄弟节点
//
// 功能说明：
// 返回同一父节点下的另一个子节点
// 如果当前节点是根节点或父节点只有一个子节点，返回nil
//
// 返回值：
// *MerkleNode - 兄弟节点，如果不存在则返回nil
func (n *MerkleNode) GetSibling() *MerkleNode {
	if n.Parent == nil {
		return nil
	}

	if n.Parent.Left == n {
		return n.Parent.Right
	}
	return n.Parent.Left
}

// GetPath 获取从当前节点到根节点的路径
//
// 功能说明：
// 返回从当前节点向上到根节点的所有节点列表
// 路径不包含当前节点本身，从父节点开始
//
// 返回值：
// []*MerkleNode - 路径上的节点列表，按从下到上的顺序
func (n *MerkleNode) GetPath() []*MerkleNode {
	var path []*MerkleNode
	current := n.Parent

	for current != nil {
		path = append(path, current)
		current = current.Parent
	}

	return path
}

// MerkleTree Merkle树主结构
// 包含完整的树结构和相关元数据
type MerkleTree struct {
	Root             *MerkleNode     // 根节点
	Leaves           []*MerkleNode   // 叶子节点列表，按交易索引顺序排列
	Levels           [][]*MerkleNode // 各层节点，levels[0]为叶子层
	TransactionCount int             // 交易总数
	TreeDepth        int             // 树的深度（从叶子到根的层数）
	mutex            sync.RWMutex    // 读写锁，保证并发安全
}

// NewMerkleTree 从交易哈希列表创建新的Merkle树
//
// 功能说明：
// 根据给定的交易哈希列表构建完整的Merkle树
// 实现比特币标准的Merkle树算法，包括奇数节点处理
//
// 算法步骤：
// 1. 验证输入哈希的有效性
// 2. 创建叶子节点层
// 3. 逐层向上构建，直到根节点
// 4. 处理奇数节点情况（最后节点自我复制）
//
// 参数：
// txHashes [][]byte - 交易哈希列表，每个哈希必须是32字节
//
// 返回值：
// *MerkleTree - 构建完成的Merkle树实例
func NewMerkleTree(txHashes [][]byte) *MerkleTree {
	if len(txHashes) == 0 {
		// 空交易列表，创建空树
		return &MerkleTree{
			Root:             nil,
			Leaves:           []*MerkleNode{},
			Levels:           [][]*MerkleNode{},
			TransactionCount: 0,
			TreeDepth:        0,
		}
	}

	// 验证所有哈希长度
	for i, hash := range txHashes {
		if len(hash) != 32 {
			panic(fmt.Sprintf("NewMerkleTree: 交易哈希%d长度无效，期望32字节，实际%d字节", i, len(hash)))
		}
	}

	tree := &MerkleTree{
		TransactionCount: len(txHashes),
		Levels:           make([][]*MerkleNode, 0),
	}

	// 创建叶子节点层
	leaves := make([]*MerkleNode, len(txHashes))
	for i, hash := range txHashes {
		var hashArray [32]byte
		copy(hashArray[:], hash)

		node := NewMerkleNode(hashArray, nil, nil)
		node.TxIndex = i
		leaves[i] = node
	}

	tree.Leaves = leaves
	tree.Levels = append(tree.Levels, leaves)

	// 逐层构建树结构
	currentLevel := leaves
	for len(currentLevel) > 1 {
		nextLevel := make([]*MerkleNode, 0, (len(currentLevel)+1)/2)

		for i := 0; i < len(currentLevel); i += 2 {
			left := currentLevel[i]
			var right *MerkleNode

			if i+1 < len(currentLevel) {
				right = currentLevel[i+1]
			} else {
				// 奇数个节点时，最后一个节点与自己配对
				right = currentLevel[i]
			}

			// 计算父节点哈希
			parentHash := tree.calculateParentHash(left.Hash, right.Hash)
			parentNode := NewMerkleNode(parentHash, left, right)

			nextLevel = append(nextLevel, parentNode)
		}

		tree.Levels = append(tree.Levels, nextLevel)
		currentLevel = nextLevel
	}

	// 设置根节点和树深度
	tree.Root = currentLevel[0]
	tree.TreeDepth = len(tree.Levels)

	return tree
}

// calculateParentHash 计算父节点哈希
//
// 功能说明：
// 根据左右子节点的哈希值计算父节点哈希
// 使用比特币标准的双重SHA-256算法
//
// 参数：
// leftHash [32]byte - 左子节点哈希
// rightHash [32]byte - 右子节点哈希
//
// 返回值：
// [32]byte - 计算得出的父节点哈希
func (mt *MerkleTree) calculateParentHash(leftHash, rightHash [32]byte) [32]byte {
	// 连接左右哈希
	combined := make([]byte, 64)
	copy(combined[:32], leftHash[:])
	copy(combined[32:], rightHash[:])

	// 计算双重SHA-256哈希
	hash := utils.DoubleSHA256(combined)
	var result [32]byte
	copy(result[:], hash)
	return result
}

// GetRoot 获取Merkle树根哈希
//
// 功能说明：
// 返回树根节点的哈希值，这是整个交易集合的唯一标识
// 如果树为空，返回全零哈希
//
// 返回值：
// [32]byte - 根节点哈希值
func (mt *MerkleTree) GetRoot() [32]byte {
	mt.mutex.RLock()
	defer mt.mutex.RUnlock()

	if mt.Root == nil {
		return [32]byte{} // 空树返回全零哈希
	}
	return mt.Root.Hash
}

// GetLeaf 获取指定索引的叶子节点
//
// 功能说明：
// 根据交易索引返回对应的叶子节点
// 提供边界检查以防止索引越界
//
// 参数：
// index int - 交易索引，从0开始
//
// 返回值：
// *MerkleNode - 对应的叶子节点，如果索引无效返回nil
func (mt *MerkleTree) GetLeaf(index int) *MerkleNode {
	mt.mutex.RLock()
	defer mt.mutex.RUnlock()

	if index < 0 || index >= len(mt.Leaves) {
		return nil
	}
	return mt.Leaves[index]
}

// MerkleProof Merkle证明结构
// 包含验证交易存在性所需的所有信息
type MerkleProof struct {
	TransactionHash  [32]byte // 目标交易的哈希值
	TransactionIndex int      // 交易在区块中的索引位置
	MerkleRoot       [32]byte // 树根哈希值
	ProofHashes      [][]byte // 证明路径上的哈希值列表
	ProofFlags       []bool   // 路径方向标志，true表示右侧，false表示左侧
	TreeDepth        int      // 树的深度
}

// NewMerkleProof 创建新的Merkle证明
//
// 功能说明：
// 创建一个新的Merkle证明实例，初始化基本信息
//
// 参数：
// txHash [32]byte - 目标交易哈希
// txIndex int - 交易索引
// root [32]byte - Merkle根哈希
//
// 返回值：
// *MerkleProof - 新的证明实例
func NewMerkleProof(txHash [32]byte, txIndex int, root [32]byte) *MerkleProof {
	return &MerkleProof{
		TransactionHash:  txHash,
		TransactionIndex: txIndex,
		MerkleRoot:       root,
		ProofHashes:      make([][]byte, 0),
		ProofFlags:       make([]bool, 0),
	}
}

// GenerateProof 为指定交易生成Merkle存在性证明
//
// 功能说明：
// 生成证明指定交易存在于Merkle树中的密码学证明
// 证明包含从叶子节点到根节点路径上所有兄弟节点的哈希
//
// 算法步骤：
// 1. 验证交易索引的有效性
// 2. 从目标叶子节点开始向上遍历
// 3. 收集每一层的兄弟节点哈希和方向信息
// 4. 构建完整的证明结构
//
// 参数：
// txIndex int - 目标交易的索引
//
// 返回值：
// *MerkleProof - 生成的Merkle证明，如果索引无效返回nil
func (mt *MerkleTree) GenerateProof(txIndex int) *MerkleProof {
	mt.mutex.RLock()
	defer mt.mutex.RUnlock()

	// 验证索引有效性
	if txIndex < 0 || txIndex >= len(mt.Leaves) {
		return nil
	}

	if mt.Root == nil {
		return nil
	}

	leaf := mt.Leaves[txIndex]
	proof := NewMerkleProof(leaf.Hash, txIndex, mt.Root.Hash)
	proof.TreeDepth = mt.TreeDepth

	// 从叶子节点向上遍历到根节点
	current := leaf
	for current.Parent != nil {
		sibling := current.GetSibling()
		if sibling != nil {
			// 添加兄弟节点哈希到证明路径
			siblingHash := make([]byte, 32)
			copy(siblingHash, sibling.Hash[:])
			proof.ProofHashes = append(proof.ProofHashes, siblingHash)

			// 确定兄弟节点的方向（true=右侧，false=左侧）
			isRightSibling := current.Parent.Right == sibling
			proof.ProofFlags = append(proof.ProofFlags, isRightSibling)
		}
		current = current.Parent
	}

	return proof
}

// VerifyProof 验证Merkle证明的有效性
//
// 功能说明：
// 验证给定的Merkle证明是否有效
// 通过重新计算从叶子到根的哈希路径来验证
//
// 参数：
// proof *MerkleProof - 待验证的Merkle证明
//
// 返回值：
// bool - true表示证明有效，false表示无效
func (mt *MerkleTree) VerifyProof(proof *MerkleProof) bool {
	if proof == nil {
		return false
	}

	return VerifyMerkleProof(
		proof.TransactionHash,
		proof.MerkleRoot,
		proof.ProofHashes,
		proof.ProofFlags,
		proof.TransactionIndex,
	)
}

// VerifyMerkleProof 独立验证Merkle证明（不需要完整树结构）
//
// 功能说明：
// 独立验证Merkle证明的有效性，不依赖于完整的树结构
// 这个函数可以被轻节点使用来验证交易存在性
//
// 验证算法：
// 1. 从交易哈希开始
// 2. 按照证明路径逐层向上计算哈希
// 3. 根据方向标志确定左右子树位置
// 4. 最终计算结果与根哈希比较
//
// 参数：
// txHash [32]byte - 交易哈希
// merkleRoot [32]byte - 预期的根哈希
// proofHashes [][]byte - 证明路径哈希列表
// proofFlags []bool - 路径方向标志
// txIndex int - 交易索引
//
// 返回值：
// bool - true表示证明有效，false表示无效
func VerifyMerkleProof(txHash [32]byte, merkleRoot [32]byte, proofHashes [][]byte, proofFlags []bool, txIndex int) bool {
	// 验证输入参数
	if len(proofHashes) != len(proofFlags) {
		return false
	}

	// 从交易哈希开始计算
	currentHash := txHash
	currentIndex := txIndex

	// 逐层向上计算哈希
	for i, siblingHashBytes := range proofHashes {
		if len(siblingHashBytes) != 32 {
			return false
		}

		var siblingHash [32]byte
		copy(siblingHash[:], siblingHashBytes)

		// 根据索引的奇偶性和方向标志确定左右位置
		var leftHash, rightHash [32]byte
		if currentIndex%2 == 0 {
			// 当前节点是左子节点
			leftHash = currentHash
			rightHash = siblingHash
		} else {
			// 当前节点是右子节点
			leftHash = siblingHash
			rightHash = currentHash
		}

		// 验证方向标志的一致性
		expectedFlag := currentIndex%2 == 0 // 如果当前是左节点，兄弟应该在右侧
		if proofFlags[i] != expectedFlag {
			return false
		}

		// 计算父节点哈希
		combined := make([]byte, 64)
		copy(combined[:32], leftHash[:])
		copy(combined[32:], rightHash[:])

		parentHashBytes := utils.DoubleSHA256(combined)
		copy(currentHash[:], parentHashBytes)

		// 更新索引到父层
		currentIndex = currentIndex / 2
	}

	// 最终哈希应该等于根哈希
	return bytes.Equal(currentHash[:], merkleRoot[:])
}

// Serialize 序列化Merkle证明
//
// 功能说明：
// 将Merkle证明序列化为字节数组，用于网络传输或存储
//
// 序列化格式：
// - TransactionHash (32字节)
// - TransactionIndex (4字节)
// - MerkleRoot (32字节)
// - TreeDepth (4字节)
// - ProofHashesCount (4字节)
// - ProofHashes (每个32字节)
// - ProofFlags (每个1字节)
//
// 返回值：
// []byte - 序列化后的字节数组
func (mp *MerkleProof) Serialize() []byte {
	buf := make([]byte, 0, 1024) // 预分配缓冲区

	// TransactionHash (32字节)
	buf = append(buf, mp.TransactionHash[:]...)

	// TransactionIndex (4字节)
	indexBytes := make([]byte, 4)
	binary.LittleEndian.PutUint32(indexBytes, uint32(mp.TransactionIndex))
	buf = append(buf, indexBytes...)

	// MerkleRoot (32字节)
	buf = append(buf, mp.MerkleRoot[:]...)

	// TreeDepth (4字节)
	depthBytes := make([]byte, 4)
	binary.LittleEndian.PutUint32(depthBytes, uint32(mp.TreeDepth))
	buf = append(buf, depthBytes...)

	// ProofHashesCount (4字节)
	countBytes := make([]byte, 4)
	binary.LittleEndian.PutUint32(countBytes, uint32(len(mp.ProofHashes)))
	buf = append(buf, countBytes...)

	// ProofHashes
	for _, hash := range mp.ProofHashes {
		buf = append(buf, hash...)
	}

	// ProofFlags
	for _, flag := range mp.ProofFlags {
		if flag {
			buf = append(buf, 1)
		} else {
			buf = append(buf, 0)
		}
	}

	return buf
}

// Deserialize 反序列化Merkle证明
//
// 功能说明：
// 从字节数组反序列化Merkle证明
//
// 参数：
// data []byte - 序列化的字节数组
//
// 返回值：
// error - 反序列化错误，成功时返回nil
func (mp *MerkleProof) Deserialize(data []byte) error {
	if len(data) < 76 { // 最小长度：32+4+32+4+4 = 76字节
		return fmt.Errorf("Merkle证明数据长度不足")
	}

	offset := 0

	// TransactionHash (32字节)
	copy(mp.TransactionHash[:], data[offset:offset+32])
	offset += 32

	// TransactionIndex (4字节)
	mp.TransactionIndex = int(binary.LittleEndian.Uint32(data[offset : offset+4]))
	offset += 4

	// MerkleRoot (32字节)
	copy(mp.MerkleRoot[:], data[offset:offset+32])
	offset += 32

	// TreeDepth (4字节)
	mp.TreeDepth = int(binary.LittleEndian.Uint32(data[offset : offset+4]))
	offset += 4

	// ProofHashesCount (4字节)
	hashCount := int(binary.LittleEndian.Uint32(data[offset : offset+4]))
	offset += 4

	// 验证剩余数据长度
	expectedLength := hashCount*32 + hashCount // 哈希数据 + 标志数据
	if len(data)-offset < expectedLength {
		return fmt.Errorf("Merkle证明数据长度不匹配")
	}

	// ProofHashes
	mp.ProofHashes = make([][]byte, hashCount)
	for i := 0; i < hashCount; i++ {
		mp.ProofHashes[i] = make([]byte, 32)
		copy(mp.ProofHashes[i], data[offset:offset+32])
		offset += 32
	}

	// ProofFlags
	mp.ProofFlags = make([]bool, hashCount)
	for i := 0; i < hashCount; i++ {
		mp.ProofFlags[i] = data[offset] == 1
		offset++
	}

	return nil
}

// GetProofSize 获取证明的字节大小
//
// 功能说明：
// 计算Merkle证明序列化后的字节大小
//
// 返回值：
// int - 证明的字节大小
func (mp *MerkleProof) GetProofSize() int {
	// 固定部分：32+4+32+4+4 = 76字节
	// 可变部分：每个证明哈希32字节 + 每个标志1字节
	return 76 + len(mp.ProofHashes)*33
}

// TreeInfo 树统计信息结构
type TreeInfo struct {
	TransactionCount int // 交易总数
	TreeDepth        int // 树深度
	LeafCount        int // 叶子节点数
	TotalNodes       int // 总节点数
	ProofSize        int // 平均证明大小
	MemoryUsage      int // 内存使用估算（字节）
}

// GetTreeInfo 获取Merkle树的统计信息
//
// 功能说明：
// 返回树的详细统计信息，用于监控和调试
//
// 返回值：
// *TreeInfo - 树的统计信息
func (mt *MerkleTree) GetTreeInfo() *TreeInfo {
	mt.mutex.RLock()
	defer mt.mutex.RUnlock()

	totalNodes := 0
	for _, level := range mt.Levels {
		totalNodes += len(level)
	}

	// 估算内存使用（每个节点约100字节）
	memoryUsage := totalNodes * 100

	return &TreeInfo{
		TransactionCount: mt.TransactionCount,
		TreeDepth:        mt.TreeDepth,
		LeafCount:        len(mt.Leaves),
		TotalNodes:       totalNodes,
		ProofSize:        mt.TreeDepth * 33, // 每层33字节（32字节哈希+1字节标志）
		MemoryUsage:      memoryUsage,
	}
}
