package blockchain

import (
	"bytes"
	"crypto/rand"
	"fmt"
	"testing"

	"simplied-bitcoin-network-go/pkg/blockchain"
	"simplied-bitcoin-network-go/pkg/utils"
)

// TestNewMerkleNode 测试Merkle节点创建
func TestNewMerkleNode(t *testing.T) {
	// 创建测试哈希
	hash := [32]byte{1, 2, 3, 4, 5}

	// 测试叶子节点创建
	leafNode := blockchain.NewMerkleNode(hash, nil, nil)
	if !leafNode.IsLeaf {
		t.Error("叶子节点IsLeaf标志应为true")
	}
	if leafNode.Hash != hash {
		t.Error("节点哈希值不匹配")
	}
	if !leafNode.IsRoot() {
		t.Error("单独节点应为根节点")
	}

	// 创建子节点测试父子关系
	leftHash := [32]byte{10, 11, 12}
	rightHash := [32]byte{20, 21, 22}
	leftChild := blockchain.NewMerkleNode(leftHash, nil, nil)
	rightChild := blockchain.NewMerkleNode(rightHash, nil, nil)

	parentNode := blockchain.NewMerkleNode(hash, leftChild, rightChild)
	if parentNode.IsLeaf {
		t.Error("父节点IsLeaf标志应为false")
	}
	if leftChild.Parent != parentNode {
		t.Error("左子节点的父节点指针设置错误")
	}
	if rightChild.Parent != parentNode {
		t.Error("右子节点的父节点指针设置错误")
	}
}

// TestMerkleNodeSibling 测试兄弟节点获取
func TestMerkleNodeSibling(t *testing.T) {
	leftHash := [32]byte{1}
	rightHash := [32]byte{2}
	parentHash := [32]byte{3}

	leftNode := blockchain.NewMerkleNode(leftHash, nil, nil)
	rightNode := blockchain.NewMerkleNode(rightHash, nil, nil)
	parentNode := blockchain.NewMerkleNode(parentHash, leftNode, rightNode)

	leftSibling := leftNode.GetSibling()
	if leftSibling != rightNode {
		t.Error("左节点的兄弟节点应为右节点")
	}

	rightSibling := rightNode.GetSibling()
	if rightSibling != leftNode {
		t.Error("右节点的兄弟节点应为左节点")
	}

	rootSibling := parentNode.GetSibling()
	if rootSibling != nil {
		t.Error("根节点不应有兄弟节点")
	}
}

// TestNewMerkleTreeEmpty 测试空Merkle树创建
func TestNewMerkleTreeEmpty(t *testing.T) {
	tree := blockchain.NewMerkleTree([][]byte{})

	if tree.Root != nil {
		t.Error("空树的根节点应为nil")
	}
	if tree.TransactionCount != 0 {
		t.Error("空树的交易数量应为0")
	}
	if tree.TreeDepth != 0 {
		t.Error("空树的深度应为0")
	}

	rootHash := tree.GetRoot()
	expectedHash := [32]byte{}
	if rootHash != expectedHash {
		t.Error("空树应返回全零哈希")
	}
}

// TestNewMerkleTreeSingle 测试单交易Merkle树
func TestNewMerkleTreeSingle(t *testing.T) {
	txHash := utils.DoubleSHA256([]byte("test transaction"))
	tree := blockchain.NewMerkleTree([][]byte{txHash})

	if tree.Root == nil {
		t.Error("单交易树应有根节点")
	}
	if tree.TransactionCount != 1 {
		t.Error("单交易树的交易数量应为1")
	}
	if tree.TreeDepth != 1 {
		t.Error("单交易树的深度应为1")
	}

	rootHash := tree.GetRoot()
	var expectedHash [32]byte
	copy(expectedHash[:], txHash)
	if rootHash != expectedHash {
		t.Error("单交易树的根哈希应等于交易哈希")
	}
}

// TestNewMerkleTreeMultiple 测试多交易Merkle树
func TestNewMerkleTreeMultiple(t *testing.T) {
	txHashes := [][]byte{
		utils.DoubleSHA256([]byte("transaction 1")),
		utils.DoubleSHA256([]byte("transaction 2")),
		utils.DoubleSHA256([]byte("transaction 3")),
		utils.DoubleSHA256([]byte("transaction 4")),
	}

	tree := blockchain.NewMerkleTree(txHashes)

	if tree.Root == nil {
		t.Error("多交易树应有根节点")
	}
	if tree.TransactionCount != 4 {
		t.Error("交易数量应为4")
	}
	if tree.TreeDepth != 3 {
		t.Error("4个交易的树深度应为3")
	}
	if len(tree.Leaves) != 4 {
		t.Error("应有4个叶子节点")
	}

	// 验证叶子节点哈希
	for i, expectedHash := range txHashes {
		leaf := tree.GetLeaf(i)
		if leaf == nil {
			t.Errorf("叶子节点%d不应为nil", i)
			continue
		}
		if !bytes.Equal(leaf.Hash[:], expectedHash) {
			t.Errorf("叶子节点%d哈希不匹配", i)
		}
	}
}

// TestNewMerkleTreeOdd 测试奇数交易Merkle树
func TestNewMerkleTreeOdd(t *testing.T) {
	txHashes := [][]byte{
		utils.DoubleSHA256([]byte("transaction 1")),
		utils.DoubleSHA256([]byte("transaction 2")),
		utils.DoubleSHA256([]byte("transaction 3")),
	}

	tree := blockchain.NewMerkleTree(txHashes)

	if tree.Root == nil {
		t.Error("奇数交易树应有根节点")
	}
	if tree.TransactionCount != 3 {
		t.Error("交易数量应为3")
	}
	if len(tree.Leaves) != 3 {
		t.Error("应有3个叶子节点")
	}

	// 验证根哈希计算正确性
	expectedRoot := calculateExpectedMerkleRoot(txHashes)
	actualRoot := tree.GetRoot()
	if !bytes.Equal(actualRoot[:], expectedRoot) {
		t.Error("奇数交易树的根哈希计算错误")
	}
}

// TestMerkleProofGeneration 测试Merkle证明生成
func TestMerkleProofGeneration(t *testing.T) {
	txHashes := [][]byte{
		utils.DoubleSHA256([]byte("transaction 1")),
		utils.DoubleSHA256([]byte("transaction 2")),
		utils.DoubleSHA256([]byte("transaction 3")),
		utils.DoubleSHA256([]byte("transaction 4")),
	}

	tree := blockchain.NewMerkleTree(txHashes)

	// 为每个交易生成证明
	for i := 0; i < len(txHashes); i++ {
		proof := tree.GenerateProof(i)
		if proof == nil {
			t.Errorf("交易%d的证明生成失败", i)
			continue
		}

		if proof.TransactionIndex != i {
			t.Errorf("交易%d的证明索引错误", i)
		}

		var expectedTxHash [32]byte
		copy(expectedTxHash[:], txHashes[i])
		if proof.TransactionHash != expectedTxHash {
			t.Errorf("交易%d的证明哈希错误", i)
		}

		if proof.MerkleRoot != tree.GetRoot() {
			t.Errorf("交易%d的证明根哈希错误", i)
		}
	}

	// 测试无效索引
	if tree.GenerateProof(-1) != nil {
		t.Error("负索引应返回nil证明")
	}
	if tree.GenerateProof(len(txHashes)) != nil {
		t.Error("超出范围的索引应返回nil证明")
	}
}

// TestMerkleProofVerification 测试Merkle证明验证
func TestMerkleProofVerification(t *testing.T) {
	txHashes := [][]byte{
		utils.DoubleSHA256([]byte("transaction 1")),
		utils.DoubleSHA256([]byte("transaction 2")),
		utils.DoubleSHA256([]byte("transaction 3")),
		utils.DoubleSHA256([]byte("transaction 4")),
	}

	tree := blockchain.NewMerkleTree(txHashes)

	// 测试所有交易的证明验证
	for i := 0; i < len(txHashes); i++ {
		proof := tree.GenerateProof(i)
		if proof == nil {
			t.Errorf("交易%d证明生成失败", i)
			continue
		}

		// 使用树方法验证
		if !tree.VerifyProof(proof) {
			t.Errorf("交易%d的证明验证失败（树方法）", i)
		}

		// 使用独立方法验证
		isValid := blockchain.VerifyMerkleProof(
			proof.TransactionHash,
			proof.MerkleRoot,
			proof.ProofHashes,
			proof.ProofFlags,
			proof.TransactionIndex,
		)
		if !isValid {
			t.Errorf("交易%d的证明验证失败（独立方法）", i)
		}
	}
}

// TestMerkleProofInvalid 测试无效Merkle证明
func TestMerkleProofInvalid(t *testing.T) {
	txHashes := [][]byte{
		utils.DoubleSHA256([]byte("transaction 1")),
		utils.DoubleSHA256([]byte("transaction 2")),
	}

	tree := blockchain.NewMerkleTree(txHashes)
	proof := tree.GenerateProof(0)

	// 测试nil证明
	if tree.VerifyProof(nil) {
		t.Error("nil证明应验证失败")
	}

	// 测试篡改交易哈希
	originalHash := proof.TransactionHash
	proof.TransactionHash = [32]byte{99, 99, 99}
	if tree.VerifyProof(proof) {
		t.Error("篡改交易哈希的证明应验证失败")
	}
	proof.TransactionHash = originalHash

	// 测试篡改根哈希
	originalRoot := proof.MerkleRoot
	proof.MerkleRoot = [32]byte{88, 88, 88}
	if tree.VerifyProof(proof) {
		t.Error("篡改根哈希的证明应验证失败")
	}
	proof.MerkleRoot = originalRoot
}

// TestMerkleProofSerialization 测试Merkle证明序列化
func TestMerkleProofSerialization(t *testing.T) {
	txHashes := [][]byte{
		utils.DoubleSHA256([]byte("transaction 1")),
		utils.DoubleSHA256([]byte("transaction 2")),
		utils.DoubleSHA256([]byte("transaction 3")),
	}

	tree := blockchain.NewMerkleTree(txHashes)
	originalProof := tree.GenerateProof(1)

	// 序列化证明
	data := originalProof.Serialize()
	if len(data) == 0 {
		t.Error("序列化数据不应为空")
	}

	// 反序列化证明
	deserializedProof := &blockchain.MerkleProof{}
	err := deserializedProof.Deserialize(data)
	if err != nil {
		t.Errorf("反序列化失败: %v", err)
	}

	// 验证反序列化结果
	if deserializedProof.TransactionHash != originalProof.TransactionHash {
		t.Error("反序列化后交易哈希不匹配")
	}
	if deserializedProof.TransactionIndex != originalProof.TransactionIndex {
		t.Error("反序列化后交易索引不匹配")
	}
	if deserializedProof.MerkleRoot != originalProof.MerkleRoot {
		t.Error("反序列化后根哈希不匹配")
	}

	// 验证反序列化的证明仍然有效
	if !tree.VerifyProof(deserializedProof) {
		t.Error("反序列化后的证明验证失败")
	}
}

// calculateExpectedMerkleRoot 计算期望的Merkle根（用于验证）
func calculateExpectedMerkleRoot(txHashes [][]byte) []byte {
	if len(txHashes) == 0 {
		return make([]byte, 32)
	}
	if len(txHashes) == 1 {
		return txHashes[0]
	}

	currentLevel := make([][]byte, len(txHashes))
	for i, hash := range txHashes {
		currentLevel[i] = make([]byte, 32)
		copy(currentLevel[i], hash)
	}

	for len(currentLevel) > 1 {
		nextLevel := make([][]byte, 0, (len(currentLevel)+1)/2)

		for i := 0; i < len(currentLevel); i += 2 {
			left := currentLevel[i]
			var right []byte

			if i+1 < len(currentLevel) {
				right = currentLevel[i+1]
			} else {
				right = currentLevel[i]
			}

			combined := make([]byte, 64)
			copy(combined[:32], left)
			copy(combined[32:], right)
			parentHash := utils.DoubleSHA256(combined)

			nextLevel = append(nextLevel, parentHash)
		}

		currentLevel = nextLevel
	}

	return currentLevel[0]
}

// BenchmarkMerkleTreeConstruction 基准测试：Merkle树构建性能
func BenchmarkMerkleTreeConstruction(b *testing.B) {
	testSizes := []int{10, 100, 1000}

	for _, size := range testSizes {
		b.Run(fmt.Sprintf("Size_%d", size), func(b *testing.B) {
			txHashes := make([][]byte, size)
			for i := 0; i < size; i++ {
				data := make([]byte, 32)
				rand.Read(data)
				txHashes[i] = utils.DoubleSHA256(data)
			}

			b.ResetTimer()
			b.ReportAllocs()

			for i := 0; i < b.N; i++ {
				_ = blockchain.NewMerkleTree(txHashes)
			}
		})
	}
}

// BenchmarkMerkleProofGeneration 基准测试：Merkle证明生成性能
func BenchmarkMerkleProofGeneration(b *testing.B) {
	txHashes := make([][]byte, 1000)
	for i := 0; i < 1000; i++ {
		data := make([]byte, 32)
		rand.Read(data)
		txHashes[i] = utils.DoubleSHA256(data)
	}

	tree := blockchain.NewMerkleTree(txHashes)

	b.ResetTimer()
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		idx := i % len(txHashes)
		_ = tree.GenerateProof(idx)
	}
}

// BenchmarkMerkleProofVerification 基准测试：Merkle证明验证性能
func BenchmarkMerkleProofVerification(b *testing.B) {
	txHashes := make([][]byte, 1000)
	for i := 0; i < 1000; i++ {
		data := make([]byte, 32)
		rand.Read(data)
		txHashes[i] = utils.DoubleSHA256(data)
	}

	tree := blockchain.NewMerkleTree(txHashes)
	proofs := make([]*blockchain.MerkleProof, 100)
	for i := 0; i < 100; i++ {
		proofs[i] = tree.GenerateProof(i)
	}

	b.ResetTimer()
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		proof := proofs[i%len(proofs)]
		_ = tree.VerifyProof(proof)
	}
}
