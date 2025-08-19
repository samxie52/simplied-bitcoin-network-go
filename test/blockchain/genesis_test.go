package blockchain_test

import (
	"bytes"
	"testing"
	"time"

	"simplied-bitcoin-network-go/pkg/blockchain"
)

// TestCreateGenesisBlock 测试创世区块创建
func TestCreateGenesisBlock(t *testing.T) {
	genesis := blockchain.CreateGenesisBlock()

	if genesis == nil {
		t.Fatal("创世区块创建失败")
	}

	if genesis.Header == nil {
		t.Fatal("创世区块头为空")
	}

	// 验证创世区块基本属性
	if genesis.Header.Version != blockchain.CurrentBlockVersion {
		t.Errorf("创世区块版本错误: 期望%d, 实际%d", blockchain.CurrentBlockVersion, genesis.Header.Version)
	}

	// 验证前块哈希为全零
	zeroHash := [32]byte{}
	if !bytes.Equal(genesis.Header.PrevBlockHash[:], zeroHash[:]) {
		t.Error("创世区块前块哈希应该为全零")
	}

	// 验证时间戳
	if genesis.Header.Timestamp != blockchain.GenesisBlockTimestamp {
		t.Errorf("创世区块时间戳错误: 期望%d, 实际%d", blockchain.GenesisBlockTimestamp, genesis.Header.Timestamp)
	}

	// 验证难度位
	if genesis.Header.Bits != blockchain.GenesisBlockBits {
		t.Errorf("创世区块难度位错误: 期望%d, 实际%d", blockchain.GenesisBlockBits, genesis.Header.Bits)
	}

	// 验证Nonce
	if genesis.Header.Nonce != blockchain.GenesisBlockNonce {
		t.Errorf("创世区块Nonce错误: 期望%d, 实际%d", blockchain.GenesisBlockNonce, genesis.Header.Nonce)
	}

	// 验证交易数量（应该只有一个Coinbase交易）
	if len(genesis.Transactions) != 1 {
		t.Errorf("创世区块应该只有一个交易, 实际有%d个", len(genesis.Transactions))
	}
}

// TestGetGenesisBlock 测试获取创世区块
func TestGetGenesisBlock(t *testing.T) {
	genesis1 := blockchain.GetGenesisBlock()
	genesis2 := blockchain.GetGenesisBlock()

	// 验证返回相同实例
	if genesis1 != genesis2 {
		t.Error("GetGenesisBlock应该返回相同的实例")
	}

	// 验证不为空
	if genesis1 == nil {
		t.Error("GetGenesisBlock不应该返回空值")
	}
}

// TestGetGenesisBlockHash 测试获取创世区块哈希
func TestGetGenesisBlockHash(t *testing.T) {
	hash1 := blockchain.GetGenesisBlockHash()
	hash2 := blockchain.GetGenesisBlockHash()

	// 验证哈希一致性
	if !bytes.Equal(hash1[:], hash2[:]) {
		t.Error("创世区块哈希应该保持一致")
	}

	// 验证哈希长度
	if len(hash1) != 32 {
		t.Errorf("哈希长度错误: 期望32, 实际%d", len(hash1))
	}

	// 验证哈希与实际计算的哈希一致
	genesis := blockchain.GetGenesisBlock()
	calculatedHash := genesis.Hash()
	if !bytes.Equal(hash1[:], calculatedHash[:]) {
		t.Error("缓存的哈希与计算的哈希不一致")
	}
}

// TestIsGenesisBlock 测试创世区块识别
func TestIsGenesisBlock(t *testing.T) {
	genesis := blockchain.GetGenesisBlock()

	// 测试真正的创世区块
	if !blockchain.IsGenesisBlock(genesis) {
		t.Error("应该识别为创世区块")
	}

	// 测试空区块
	if blockchain.IsGenesisBlock(nil) {
		t.Error("空区块不应该被识别为创世区块")
	}

	// 测试普通区块
	tx := blockchain.NewTransaction([]byte("test"))
	header := blockchain.NewBlockHeader(1, [32]byte{}, [32]byte{}, uint32(time.Now().Unix()), 0x1d00ffff, 12345)
	normalBlock := blockchain.NewBlock(header, []*blockchain.Transaction{tx})

	if blockchain.IsGenesisBlock(normalBlock) {
		t.Error("普通区块不应该被识别为创世区块")
	}
}

// TestIsGenesisBlockHash 测试创世区块哈希识别
func TestIsGenesisBlockHash(t *testing.T) {
	genesisHash := blockchain.GetGenesisBlockHash()

	// 测试真正的创世区块哈希
	if !blockchain.IsGenesisBlockHash(genesisHash) {
		t.Error("应该识别为创世区块哈希")
	}

	// 测试随机哈希
	randomHash := [32]byte{1, 2, 3, 4, 5}
	if blockchain.IsGenesisBlockHash(randomHash) {
		t.Error("随机哈希不应该被识别为创世区块哈希")
	}

	// 测试零哈希
	zeroHash := [32]byte{}
	if blockchain.IsGenesisBlockHash(zeroHash) {
		t.Error("零哈希不应该被识别为创世区块哈希")
	}
}

// TestValidateGenesisBlock 测试创世区块验证
func TestValidateGenesisBlock(t *testing.T) {
	genesis := blockchain.GetGenesisBlock()

	// 测试有效的创世区块
	err := blockchain.ValidateGenesisBlock(genesis)
	if err != nil {
		t.Errorf("有效创世区块验证失败: %v", err)
	}

	// 测试空区块
	err = blockchain.ValidateGenesisBlock(nil)
	if err == nil {
		t.Error("空区块应该验证失败")
	}

	// 测试无效版本的创世区块
	invalidGenesis := blockchain.CreateGenesisBlock()
	invalidGenesis.Header.Version = 999
	err = blockchain.ValidateGenesisBlock(invalidGenesis)
	if err == nil {
		t.Error("无效版本的创世区块应该验证失败")
	}

	// 测试无效前块哈希的创世区块
	invalidGenesis2 := blockchain.CreateGenesisBlock()
	invalidGenesis2.Header.PrevBlockHash = [32]byte{1, 2, 3}
	err = blockchain.ValidateGenesisBlock(invalidGenesis2)
	if err == nil {
		t.Error("无效前块哈希的创世区块应该验证失败")
	}

	// 测试无效时间戳的创世区块
	invalidGenesis3 := blockchain.CreateGenesisBlock()
	invalidGenesis3.Header.Timestamp = 999999999
	err = blockchain.ValidateGenesisBlock(invalidGenesis3)
	if err == nil {
		t.Error("无效时间戳的创世区块应该验证失败")
	}

	// 测试无效难度位的创世区块
	invalidGenesis4 := blockchain.CreateGenesisBlock()
	invalidGenesis4.Header.Bits = 0x12345678
	err = blockchain.ValidateGenesisBlock(invalidGenesis4)
	if err == nil {
		t.Error("无效难度位的创世区块应该验证失败")
	}

	// 测试无效Nonce的创世区块
	invalidGenesis5 := blockchain.CreateGenesisBlock()
	invalidGenesis5.Header.Nonce = 999999999
	err = blockchain.ValidateGenesisBlock(invalidGenesis5)
	if err == nil {
		t.Error("无效Nonce的创世区块应该验证失败")
	}
}

// TestGetGenesisBlockInfo 测试获取创世区块信息
func TestGetGenesisBlockInfo(t *testing.T) {
	info := blockchain.GetGenesisBlockInfo()

	// 验证信息完整性
	requiredFields := []string{
		"hash", "version", "prevHash", "merkleRoot",
		"timestamp", "timestampStr", "bits", "nonce",
		"difficulty", "size", "txCount",
	}

	for _, field := range requiredFields {
		if _, exists := info[field]; !exists {
			t.Errorf("创世区块信息缺少字段: %s", field)
		}
	}

	// 验证具体值
	if version, ok := info["version"].(uint32); !ok || version != blockchain.CurrentBlockVersion {
		t.Errorf("版本信息错误: 期望%d, 实际%v", blockchain.CurrentBlockVersion, info["version"])
	}

	if timestamp, ok := info["timestamp"].(uint32); !ok || timestamp != blockchain.GenesisBlockTimestamp {
		t.Errorf("时间戳信息错误: 期望%d, 实际%v", blockchain.GenesisBlockTimestamp, info["timestamp"])
	}

	if bits, ok := info["bits"].(uint32); !ok || bits != blockchain.GenesisBlockBits {
		t.Errorf("难度位信息错误: 期望%d, 实际%v", blockchain.GenesisBlockBits, info["bits"])
	}

	if nonce, ok := info["nonce"].(uint32); !ok || nonce != blockchain.GenesisBlockNonce {
		t.Errorf("Nonce信息错误: 期望%d, 实际%v", blockchain.GenesisBlockNonce, info["nonce"])
	}

	if txCount, ok := info["txCount"].(int); !ok || txCount != 1 {
		t.Errorf("交易数量错误: 期望1, 实际%v", info["txCount"])
	}

	// 验证哈希字符串格式
	hashStr, ok := info["hash"].(string)
	if !ok {
		t.Error("哈希应该是字符串类型")
	}
	if len(hashStr) != 64 { // 32字节 = 64个十六进制字符
		t.Errorf("哈希字符串长度错误: 期望64, 实际%d", len(hashStr))
	}

	// 验证时间戳字符串格式
	timestampStr, ok := info["timestampStr"].(string)
	if !ok {
		t.Error("时间戳字符串应该是字符串类型")
	}
	if len(timestampStr) == 0 {
		t.Error("时间戳字符串不应该为空")
	}

	// 验证难度值
	difficulty, ok := info["difficulty"].(float64)
	if !ok {
		t.Error("难度应该是float64类型")
	}
	if difficulty <= 0 {
		t.Error("难度值应该大于0")
	}

	// 验证大小
	size, ok := info["size"].(int)
	if !ok {
		t.Error("大小应该是int类型")
	}
	if size <= blockchain.BlockHeaderSize {
		t.Errorf("区块大小应该大于区块头大小: %d", blockchain.BlockHeaderSize)
	}
}

// TestCreateTestGenesisBlock 测试创建测试创世区块
func TestCreateTestGenesisBlock(t *testing.T) {
	testGenesis := blockchain.CreateTestGenesisBlock()

	if testGenesis == nil {
		t.Fatal("测试创世区块创建失败")
	}

	// 验证基本属性
	if testGenesis.Header == nil {
		t.Fatal("测试创世区块头为空")
	}

	// 验证版本
	if testGenesis.Header.Version != blockchain.CurrentBlockVersion {
		t.Errorf("测试创世区块版本错误: 期望%d, 实际%d", blockchain.CurrentBlockVersion, testGenesis.Header.Version)
	}

	// 验证前块哈希为全零
	zeroHash := [32]byte{}
	if !bytes.Equal(testGenesis.Header.PrevBlockHash[:], zeroHash[:]) {
		t.Error("测试创世区块前块哈希应该为全零")
	}

	// 验证难度位（应该是最大难度目标）
	if testGenesis.Header.Bits != blockchain.MaxTargetBits {
		t.Errorf("测试创世区块难度位错误: 期望%d, 实际%d", blockchain.MaxTargetBits, testGenesis.Header.Bits)
	}

	// 验证交易数量
	if len(testGenesis.Transactions) != 1 {
		t.Errorf("测试创世区块应该只有一个交易, 实际有%d个", len(testGenesis.Transactions))
	}

	// 验证时间戳合理性（应该是最近的时间）
	now := uint32(time.Now().Unix())
	if testGenesis.Header.Timestamp > now || testGenesis.Header.Timestamp < now-60 {
		t.Error("测试创世区块时间戳不合理")
	}

	// 验证与标准创世区块不同
	standardGenesis := blockchain.GetGenesisBlock()
	testHash := testGenesis.Hash()
	standardHash := standardGenesis.Hash()
	if bytes.Equal(testHash[:], standardHash[:]) {
		t.Error("测试创世区块不应该与标准创世区块相同")
	}
}

// TestResetGenesisBlock 测试重置创世区块
func TestResetGenesisBlock(t *testing.T) {
	// 获取原始创世区块哈希
	originalHash := blockchain.GetGenesisBlockHash()

	// 重置创世区块
	blockchain.ResetGenesisBlock()

	// 验证哈希保持一致
	newHash := blockchain.GetGenesisBlockHash()
	if !bytes.Equal(originalHash[:], newHash[:]) {
		t.Error("重置后创世区块哈希应该保持一致")
	}

	// 验证创世区块仍然有效
	genesis := blockchain.GetGenesisBlock()
	err := blockchain.ValidateGenesisBlock(genesis)
	if err != nil {
		t.Errorf("重置后创世区块验证失败: %v", err)
	}
}

// TestGenesisBlockConsistency 测试创世区块一致性
func TestGenesisBlockConsistency(t *testing.T) {
	// 多次创建创世区块，验证结果一致
	genesis1 := blockchain.CreateGenesisBlock()
	genesis2 := blockchain.CreateGenesisBlock()

	hash1 := genesis1.Hash()
	hash2 := genesis2.Hash()

	if !bytes.Equal(hash1[:], hash2[:]) {
		t.Error("多次创建的创世区块哈希应该一致")
	}

	// 验证序列化结果一致
	data1 := genesis1.Serialize()
	data2 := genesis2.Serialize()

	if !bytes.Equal(data1, data2) {
		t.Error("多次创建的创世区块序列化结果应该一致")
	}
}

// TestGenesisBlockMerkleRoot 测试创世区块Merkle根
func TestGenesisBlockMerkleRoot(t *testing.T) {
	genesis := blockchain.GetGenesisBlock()

	// 计算Merkle根
	calculatedMerkleRoot := genesis.GetMerkleRoot()

	// 验证与区块头中的Merkle根一致
	if !bytes.Equal(calculatedMerkleRoot[:], genesis.Header.MerkleRoot[:]) {
		t.Error("创世区块计算的Merkle根与区块头中的不一致")
	}

	// 验证Merkle根不为零
	zeroHash := [32]byte{}
	if bytes.Equal(calculatedMerkleRoot[:], zeroHash[:]) {
		t.Error("创世区块Merkle根不应该为零")
	}
}

// TestGenesisBlockSerialization 测试创世区块序列化
func TestGenesisBlockSerialization(t *testing.T) {
	genesis := blockchain.GetGenesisBlock()

	// 序列化
	data := genesis.Serialize()
	if len(data) == 0 {
		t.Error("创世区块序列化结果不应该为空")
	}

	// 反序列化
	deserialized := &blockchain.Block{}
	err := deserialized.Deserialize(data)
	if err != nil {
		t.Fatalf("创世区块反序列化失败: %v", err)
	}

	// 验证哈希一致
	originalHash := genesis.Hash()
	deserializedHash := deserialized.Hash()
	if !bytes.Equal(originalHash[:], deserializedHash[:]) {
		t.Error("序列化后的创世区块哈希不一致")
	}

	// 验证是创世区块
	if !blockchain.IsGenesisBlock(deserialized) {
		t.Error("反序列化后应该仍然被识别为创世区块")
	}
}
