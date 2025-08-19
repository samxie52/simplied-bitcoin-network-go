package blockchain_test

import (
	"bytes"
	"testing"
	"time"

	"simplied-bitcoin-network-go/pkg/blockchain"
	"simplied-bitcoin-network-go/pkg/utils"
)

// TestBlockHeaderSerialization 测试区块头序列化和反序列化
func TestBlockHeaderSerialization(t *testing.T) {
	// 创建测试区块头
	var prevHash, merkleRoot [32]byte
	copy(prevHash[:], []byte("prev_block_hash_32_bytes_test___"))
	copy(merkleRoot[:], []byte("merkle_root_hash_32_bytes_test__"))

	original := blockchain.NewBlockHeader(
		1,          // version
		prevHash,   // prevBlockHash
		merkleRoot, // merkleRoot
		1231006505, // timestamp
		0x1d00ffff, // bits
		2083236893, // nonce
	)

	// 序列化
	data := original.Serialize()
	if len(data) != blockchain.BlockHeaderSize {
		t.Errorf("序列化后大小错误: 期望%d, 实际%d", blockchain.BlockHeaderSize, len(data))
	}

	// 反序列化
	deserialized := &blockchain.BlockHeader{}
	err := deserialized.Deserialize(data)
	if err != nil {
		t.Fatalf("反序列化失败: %v", err)
	}

	// 验证字段
	if deserialized.Version != original.Version {
		t.Errorf("版本不匹配: 期望%d, 实际%d", original.Version, deserialized.Version)
	}
	if !bytes.Equal(deserialized.PrevBlockHash[:], original.PrevBlockHash[:]) {
		t.Errorf("前块哈希不匹配")
	}
	if !bytes.Equal(deserialized.MerkleRoot[:], original.MerkleRoot[:]) {
		t.Errorf("Merkle根不匹配")
	}
	if deserialized.Timestamp != original.Timestamp {
		t.Errorf("时间戳不匹配: 期望%d, 实际%d", original.Timestamp, deserialized.Timestamp)
	}
	if deserialized.Bits != original.Bits {
		t.Errorf("难度位不匹配: 期望%x, 实际%x", original.Bits, deserialized.Bits)
	}
	if deserialized.Nonce != original.Nonce {
		t.Errorf("Nonce不匹配: 期望%d, 实际%d", original.Nonce, deserialized.Nonce)
	}
}

// TestBlockHeaderHash 测试区块头哈希计算
func TestBlockHeaderHash(t *testing.T) {
	// 创建测试区块头
	var prevHash, merkleRoot [32]byte
	header := blockchain.NewBlockHeader(1, prevHash, merkleRoot, 1231006505, 0x1d00ffff, 2083236893)

	// 计算哈希
	hash1 := header.Hash()
	hash2 := header.Hash()

	// 验证哈希一致性
	if !bytes.Equal(hash1[:], hash2[:]) {
		t.Error("相同区块头的哈希应该相同")
	}

	// 验证哈希长度
	if len(hash1) != 32 {
		t.Errorf("哈希长度错误: 期望32, 实际%d", len(hash1))
	}
}

// TestBlockHeaderValidation 测试区块头验证
func TestBlockHeaderValidation(t *testing.T) {
	tests := []struct {
		name      string
		header    *blockchain.BlockHeader
		shouldErr bool
	}{
		{
			name: "有效区块头",
			header: blockchain.NewBlockHeader(
				1,
				[32]byte{},
				[32]byte{},
				uint32(time.Now().Unix()),
				0x1d00ffff,
				12345,
			),
			shouldErr: false,
		},
		{
			name: "版本过旧",
			header: blockchain.NewBlockHeader(
				0, // 无效版本
				[32]byte{},
				[32]byte{},
				uint32(time.Now().Unix()),
				0x1d00ffff,
				12345,
			),
			shouldErr: true,
		},
		{
			name: "未来时间戳",
			header: blockchain.NewBlockHeader(
				1,
				[32]byte{},
				[32]byte{},
				uint32(time.Now().Unix()+7200+1), // 超过2小时
				0x1d00ffff,
				12345,
			),
			shouldErr: true,
		},
		{
			name: "无效难度位",
			header: blockchain.NewBlockHeader(
				1,
				[32]byte{},
				[32]byte{},
				uint32(time.Now().Unix()),
				0, // 无效难度位
				12345,
			),
			shouldErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			isValid := tt.header.IsValid()
			if tt.shouldErr && isValid {
				t.Error("应该验证失败但通过了")
			}
			if !tt.shouldErr && !isValid {
				t.Error("应该验证通过但失败了")
			}
		})
	}
}

// TestBlockSerialization 测试完整区块序列化
func TestBlockSerialization(t *testing.T) {
	// 创建测试交易
	tx1 := blockchain.NewTransaction([]byte("transaction 1 data"))
	tx2 := blockchain.NewTransaction([]byte("transaction 2 data"))
	transactions := []*blockchain.Transaction{tx1, tx2}

	// 计算Merkle根
	txHashes := make([][]byte, len(transactions))
	for i, tx := range transactions {
		txHashes[i] = tx.Hash[:]
	}
	merkleRoot := utils.MerkleRoot(txHashes)
	var merkleRootArray [32]byte
	copy(merkleRootArray[:], merkleRoot)

	// 创建区块头
	header := blockchain.NewBlockHeader(
		1,
		[32]byte{},
		merkleRootArray,
		uint32(time.Now().Unix()),
		0x1d00ffff,
		12345,
	)

	// 创建区块
	original := blockchain.NewBlock(header, transactions)

	// 序列化
	data := original.Serialize()

	// 反序列化
	deserialized := &blockchain.Block{}
	err := deserialized.Deserialize(data)
	if err != nil {
		t.Fatalf("反序列化失败: %v", err)
	}

	// 验证区块头
	if deserialized.Header.Version != original.Header.Version {
		t.Error("区块头版本不匹配")
	}

	// 验证交易数量
	if len(deserialized.Transactions) != len(original.Transactions) {
		t.Errorf("交易数量不匹配: 期望%d, 实际%d", len(original.Transactions), len(deserialized.Transactions))
	}

	// 验证区块哈希
	originalHash := original.Hash()
	deserializedHash := deserialized.Hash()
	if !bytes.Equal(originalHash[:], deserializedHash[:]) {
		t.Error("区块哈希不匹配")
	}
}

// TestBlockSizeCalculation 测试区块大小计算
func TestBlockSizeCalculation(t *testing.T) {
	// 创建不同大小的交易
	smallTx := blockchain.NewTransaction([]byte("small"))
	largeTx := blockchain.NewTransaction(make([]byte, 1000))

	transactions := []*blockchain.Transaction{smallTx, largeTx}

	// 计算Merkle根
	txHashes := make([][]byte, len(transactions))
	for i, tx := range transactions {
		txHashes[i] = tx.Hash[:]
	}
	merkleRoot := utils.MerkleRoot(txHashes)
	var merkleRootArray [32]byte
	copy(merkleRootArray[:], merkleRoot)

	header := blockchain.NewBlockHeader(1, [32]byte{}, merkleRootArray, uint32(time.Now().Unix()), 0x1d00ffff, 12345)
	block := blockchain.NewBlock(header, transactions)

	// 计算大小
	size := block.Size()

	// 验证大小包含区块头
	if size < blockchain.BlockHeaderSize {
		t.Errorf("区块大小应该至少包含区块头大小: %d", blockchain.BlockHeaderSize)
	}

	// 验证大小合理性
	expectedMinSize := blockchain.BlockHeaderSize + 1 + 5 + 1000 // 头部 + VarInt + 小交易 + 大交易
	if size < expectedMinSize {
		t.Errorf("区块大小过小: 期望至少%d, 实际%d", expectedMinSize, size)
	}
}

// TestBlockValidation 测试区块验证
func TestBlockValidation(t *testing.T) {
	// 创建有效区块
	tx := blockchain.NewTransaction([]byte("test transaction"))
	txHashes := [][]byte{tx.Hash[:]}
	merkleRoot := utils.MerkleRoot(txHashes)
	var merkleRootArray [32]byte
	copy(merkleRootArray[:], merkleRoot)

	header := blockchain.NewBlockHeader(1, [32]byte{}, merkleRootArray, uint32(time.Now().Unix()), 0x1d00ffff, 12345)
	validBlock := blockchain.NewBlock(header, []*blockchain.Transaction{tx})

	// 测试有效区块
	err := validBlock.Validate()
	if err != nil {
		t.Errorf("有效区块验证失败: %v", err)
	}

	// 测试空区块
	emptyBlock := blockchain.NewBlock(header, []*blockchain.Transaction{})
	err = emptyBlock.Validate()
	if err == nil {
		t.Error("空区块应该验证失败")
	}

	// 测试Merkle根不匹配
	wrongMerkleRoot := [32]byte{1, 2, 3}
	wrongHeader := blockchain.NewBlockHeader(1, [32]byte{}, wrongMerkleRoot, uint32(time.Now().Unix()), 0x1d00ffff, 12345)
	wrongBlock := blockchain.NewBlock(wrongHeader, []*blockchain.Transaction{tx})
	err = wrongBlock.Validate()
	if err == nil {
		t.Error("Merkle根不匹配的区块应该验证失败")
	}
}

// TestBlockMerkleRoot 测试Merkle根计算
func TestBlockMerkleRoot(t *testing.T) {
	// 创建测试交易
	tx1 := blockchain.NewTransaction([]byte("tx1"))
	tx2 := blockchain.NewTransaction([]byte("tx2"))
	tx3 := blockchain.NewTransaction([]byte("tx3"))

	transactions := []*blockchain.Transaction{tx1, tx2, tx3}

	// 使用utils包计算期望的Merkle根
	txHashes := make([][]byte, len(transactions))
	for i, tx := range transactions {
		txHashes[i] = tx.Hash[:]
	}
	expectedMerkleRoot := utils.MerkleRoot(txHashes)

	// 创建区块
	var merkleRootArray [32]byte
	copy(merkleRootArray[:], expectedMerkleRoot)
	header := blockchain.NewBlockHeader(1, [32]byte{}, merkleRootArray, uint32(time.Now().Unix()), 0x1d00ffff, 12345)
	block := blockchain.NewBlock(header, transactions)

	// 计算区块的Merkle根
	calculatedMerkleRoot := block.GetMerkleRoot()

	// 验证Merkle根
	if !bytes.Equal(calculatedMerkleRoot[:], expectedMerkleRoot) {
		t.Errorf("Merkle根计算错误")
	}
}

// TestBlockTransactionOperations 测试区块交易操作
func TestBlockTransactionOperations(t *testing.T) {
	// 创建测试交易
	tx1 := blockchain.NewTransaction([]byte("tx1"))
	tx2 := blockchain.NewTransaction([]byte("tx2"))
	tx3 := blockchain.NewTransaction([]byte("tx3"))

	transactions := []*blockchain.Transaction{tx1, tx2, tx3}

	// 创建区块
	txHashes := make([][]byte, len(transactions))
	for i, tx := range transactions {
		txHashes[i] = tx.Hash[:]
	}
	merkleRoot := utils.MerkleRoot(txHashes)
	var merkleRootArray [32]byte
	copy(merkleRootArray[:], merkleRoot)

	header := blockchain.NewBlockHeader(1, [32]byte{}, merkleRootArray, uint32(time.Now().Unix()), 0x1d00ffff, 12345)
	block := blockchain.NewBlock(header, transactions)

	// 测试HasTransaction
	if !block.HasTransaction(tx1.Hash) {
		t.Error("区块应该包含tx1")
	}

	// 测试不存在的交易
	nonExistentTx := blockchain.NewTransaction([]byte("non-existent"))
	if block.HasTransaction(nonExistentTx.Hash) {
		t.Error("区块不应该包含不存在的交易")
	}

	// 测试GetTransaction
	retrievedTx := block.GetTransaction(tx2.Hash)
	if retrievedTx == nil {
		t.Error("应该能够获取tx2")
	}
	if !bytes.Equal(retrievedTx.Hash[:], tx2.Hash[:]) {
		t.Error("获取的交易哈希不匹配")
	}

	// 测试GetTransactionHashes
	hashes := block.GetTransactionHashes()
	if len(hashes) != len(transactions) {
		t.Errorf("交易哈希数量不匹配: 期望%d, 实际%d", len(transactions), len(hashes))
	}
}

// TestBlockHeaderDifficulty 测试区块头难度相关功能
func TestBlockHeaderDifficulty(t *testing.T) {
	header := blockchain.NewBlockHeader(1, [32]byte{}, [32]byte{}, uint32(time.Now().Unix()), 0x1d00ffff, 12345)

	// 测试GetDifficulty
	difficulty := header.GetDifficulty()
	if difficulty <= 0 {
		t.Error("难度值应该大于0")
	}

	// 测试MeetsTarget (这个测试可能失败，因为随机nonce不太可能满足难度要求)
	// 这里主要测试函数不会panic
	_ = header.MeetsTarget()
}

// TestBlockString 测试字符串表示
func TestBlockString(t *testing.T) {
	tx := blockchain.NewTransaction([]byte("test"))
	txHashes := [][]byte{tx.Hash[:]}
	merkleRoot := utils.MerkleRoot(txHashes)
	var merkleRootArray [32]byte
	copy(merkleRootArray[:], merkleRoot)

	header := blockchain.NewBlockHeader(1, [32]byte{}, merkleRootArray, uint32(time.Now().Unix()), 0x1d00ffff, 12345)
	block := blockchain.NewBlock(header, []*blockchain.Transaction{tx})

	// 测试区块字符串表示
	blockStr := block.String()
	if len(blockStr) == 0 {
		t.Error("区块字符串表示不应该为空")
	}

	// 测试区块头字符串表示
	headerStr := header.String()
	if len(headerStr) == 0 {
		t.Error("区块头字符串表示不应该为空")
	}
}

// TestNewTransaction 测试交易创建
func TestNewTransaction(t *testing.T) {
	data := []byte("test transaction data")
	tx := blockchain.NewTransaction(data)

	if tx == nil {
		t.Fatal("交易创建失败")
	}

	if !bytes.Equal(tx.Data, data) {
		t.Error("交易数据不匹配")
	}

	// 验证哈希计算
	expectedHash := utils.DoubleSHA256(data)
	if !bytes.Equal(tx.Hash[:], expectedHash) {
		t.Error("交易哈希计算错误")
	}
}

// TestBlockDeserializationErrors 测试反序列化错误处理
func TestBlockDeserializationErrors(t *testing.T) {
	tests := []struct {
		name string
		data []byte
	}{
		{
			name: "数据太短",
			data: []byte{1, 2, 3},
		},
		{
			name: "区块头大小错误",
			data: make([]byte, blockchain.BlockHeaderSize-1),
		},
		{
			name: "无效VarInt",
			data: append(make([]byte, blockchain.BlockHeaderSize), 0xff, 0xff),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			block := &blockchain.Block{}
			err := block.Deserialize(tt.data)
			if err == nil {
				t.Error("应该返回反序列化错误")
			}
		})
	}

	// 测试区块头反序列化错误
	header := &blockchain.BlockHeader{}
	err := header.Deserialize([]byte{1, 2, 3})
	if err == nil {
		t.Error("区块头反序列化应该返回错误")
	}
}
