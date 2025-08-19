// Package blockchain 实现了简化比特币网络的区块链核心功能
// 本文件包含创世区块的创建、验证、管理和相关工具函数
// 创世区块是区块链的第一个区块，具有特殊的属性和固定的参数
package blockchain

import (
	"encoding/hex"
	"fmt"
	"time"

	"simplied-bitcoin-network-go/pkg/utils"
)

// GenesisBlock 全局创世区块实例
// 存储整个网络共享的创世区块，在程序启动时初始化
// 所有节点必须使用相同的创世区块以确保网络一致性
var GenesisBlock *Block

// genesisBlockHash 创世区块哈希缓存
// 预先计算并缓存创世区块的哈希值，避免重复计算
// 使用固定大小的32字节数组存储SHA-256哈希结果
var genesisBlockHash [32]byte

// init 包初始化函数，在程序启动时自动执行
// 功能：
// 1. 创建标准的比特币兼容创世区块
// 2. 计算并缓存创世区块的哈希值
// 3. 确保全局创世区块实例的可用性
func init() {
	GenesisBlock = CreateGenesisBlock()
	genesisBlockHash = GenesisBlock.Hash()
}

// CreateGenesisBlock 创建比特币兼容的创世区块
//
// 功能说明：
// 创建一个与比特币网络兼容的创世区块，包含固定的参数和结构
// 创世区块是区块链的起始点，具有以下特征：
// - 前块哈希为全零（因为没有前一个区块）
// - 包含一个特殊的Coinbase交易
// - 使用预定义的时间戳、难度和Nonce值
//
// 实现步骤：
// 1. 初始化前块哈希为全零数组
// 2. 创建创世区块的Coinbase交易
// 3. 计算交易的Merkle根哈希
// 4. 构建区块头结构
// 5. 组装完整的创世区块
//
// 返回值：
// *Block - 完整的创世区块实例，包含区块头和交易列表
func CreateGenesisBlock() *Block {
	// 创世区块的前块哈希为全零
	var prevBlockHash [32]byte

	// 创建创世区块的Coinbase交易
	coinbaseTx := createGenesisCoinbaseTransaction()

	// 计算Merkle根（创世区块只有一个交易）
	txHashes := [][]byte{coinbaseTx.Hash[:]}
	merkleRoot := utils.MerkleRoot(txHashes)
	var merkleRootArray [32]byte
	copy(merkleRootArray[:], merkleRoot)

	// 创建创世区块头
	header := NewBlockHeader(
		CurrentBlockVersion,   // Version
		prevBlockHash,         // PrevBlockHash (全零)
		merkleRootArray,       // MerkleRoot
		GenesisBlockTimestamp, // Timestamp
		GenesisBlockBits,      // Bits (难度)
		GenesisBlockNonce,     // Nonce
	)

	// 创建创世区块
	block := NewBlock(header, []*Transaction{coinbaseTx})

	return block
}

// createGenesisCoinbaseTransaction 创建创世区块的Coinbase交易
//
// 功能说明：
// Coinbase交易是每个区块的第一个交易，用于奖励矿工
// 创世区块的Coinbase交易具有特殊性质：
// - 输入引用全零的前置交易哈希（表示这是新创建的币）
// - 输入索引为0xFFFFFFFF（Coinbase交易的标识）
// - 包含特殊的脚本消息
// - 输出50 BTC的初始奖励
//
// 实现细节：
// 1. 构建交易版本号（4字节小端序）
// 2. 设置输入数量为1
// 3. 添加Coinbase输入（全零哈希 + 0xFFFFFFFF索引）
// 4. 插入包含创世消息的脚本
// 5. 设置输出数量为1
// 6. 添加50 BTC的输出金额
// 7. 设置锁定时间为0
//
// 返回值：
// *Transaction - 创世区块的Coinbase交易实例
func createGenesisCoinbaseTransaction() *Transaction {
	// 创世区块Coinbase交易的输入脚本包含特殊消息
	coinbaseScript := []byte(GenesisCoinbaseMessage)

	// 创建简化的Coinbase交易数据
	// 在完整实现中，这里会包含完整的交易结构
	coinbaseData := make([]byte, 0, len(coinbaseScript)+100)

	// 添加版本号 (4字节)
	coinbaseData = append(coinbaseData, utils.Uint32ToLittleEndian(1)...)

	// 添加输入数量 (1字节，VarInt)
	coinbaseData = append(coinbaseData, utils.EncodeVarInt(1)...)

	// 添加输入：前置交易哈希 (32字节全零)
	coinbaseData = append(coinbaseData, make([]byte, 32)...)

	// 添加输入：输出索引 (4字节，0xFFFFFFFF表示Coinbase)
	coinbaseData = append(coinbaseData, 0xFF, 0xFF, 0xFF, 0xFF)

	// 添加脚本长度
	coinbaseData = append(coinbaseData, utils.EncodeVarInt(uint64(len(coinbaseScript)))...)

	// 添加脚本内容
	coinbaseData = append(coinbaseData, coinbaseScript...)

	// 添加序列号 (4字节)
	coinbaseData = append(coinbaseData, utils.Uint32ToLittleEndian(0xFFFFFFFF)...)

	// 添加输出数量 (1字节)
	coinbaseData = append(coinbaseData, utils.EncodeVarInt(1)...)

	// 添加输出：金额 (8字节，50 BTC = 5000000000 satoshis)
	coinbaseData = append(coinbaseData, utils.Uint64ToLittleEndian(5000000000)...)

	// 添加输出脚本长度 (简化为空脚本)
	coinbaseData = append(coinbaseData, utils.EncodeVarInt(0)...)

	// 添加锁定时间 (4字节)
	coinbaseData = append(coinbaseData, utils.Uint32ToLittleEndian(0)...)

	// 创建交易
	return NewTransaction(coinbaseData)
}

// GetGenesisBlock 获取全局创世区块实例
//
// 功能说明：
// 提供对全局创世区块实例的只读访问
// 确保所有代码使用相同的创世区块引用
//
// 使用场景：
// - 区块链初始化时获取起始区块
// - 验证区块链完整性时的参考点
// - 网络同步时的基准区块
//
// 返回值：
// *Block - 全局创世区块实例的指针
func GetGenesisBlock() *Block {
	return GenesisBlock
}

// GetGenesisBlockHash 获取创世区块的哈希值
//
// 功能说明：
// 返回预先计算并缓存的创世区块哈希
// 避免每次调用时重新计算哈希，提高性能
//
// 实现优势：
// - 性能优化：哈希值在初始化时计算一次
// - 内存效率：使用固定大小的数组存储
// - 一致性保证：所有调用返回相同的哈希值
//
// 返回值：
// [32]byte - 创世区块的SHA-256双重哈希值
func GetGenesisBlockHash() [32]byte {
	return genesisBlockHash
}

// IsGenesisBlock 检查给定区块是否为创世区块
//
// 功能说明：
// 通过比较区块哈希来判断给定区块是否为创世区块
// 提供安全的创世区块识别机制
//
// 验证逻辑：
// 1. 检查输入区块和区块头的有效性
// 2. 计算给定区块的哈希值
// 3. 与缓存的创世区块哈希进行比较
//
// 参数：
// block *Block - 待检查的区块实例
//
// 返回值：
// bool - true表示是创世区块，false表示不是
func IsGenesisBlock(block *Block) bool {
	if block == nil || block.Header == nil {
		return false
	}

	// 检查区块哈希
	blockHash := block.Hash()
	return blockHash == genesisBlockHash
}

// IsGenesisBlockHash 检查给定哈希是否为创世区块哈希
//
// 功能说明：
// 直接通过哈希值判断是否为创世区块
// 适用于只有哈希值而没有完整区块数据的场景
//
// 使用场景：
// - 网络消息中的区块哈希验证
// - 区块链索引中的快速查找
// - P2P网络中的区块标识验证
//
// 参数：
// hash [32]byte - 待检查的32字节哈希值
//
// 返回值：
// bool - true表示是创世区块哈希，false表示不是
func IsGenesisBlockHash(hash [32]byte) bool {
	return hash == genesisBlockHash
}

// ValidateGenesisBlock 验证创世区块的正确性和完整性
//
// 功能说明：
// 对创世区块进行全面的验证，确保其符合比特币协议标准
// 验证包括结构完整性、参数正确性和数据一致性
//
// 验证项目：
// 1. 基本结构验证（区块和区块头非空）
// 2. 版本号验证（必须匹配当前区块版本）
// 3. 前块哈希验证（必须为全零）
// 4. 时间戳验证（必须匹配预定义值）
// 5. 难度位验证（必须匹配创世区块难度）
// 6. Nonce值验证（必须匹配预定义值）
// 7. 交易数量验证（必须只有一个Coinbase交易）
// 8. Merkle根验证（计算值必须匹配区块头中的值）
// 9. 区块哈希验证（最终哈希必须匹配创世区块哈希）
//
// 参数：
// block *Block - 待验证的区块实例
//
// 返回值：
// error - 验证失败时返回具体错误信息，成功时返回nil
func ValidateGenesisBlock(block *Block) error {
	if block == nil {
		return fmt.Errorf("创世区块不能为空")
	}

	if block.Header == nil {
		return fmt.Errorf("创世区块头不能为空")
	}

	// 验证版本
	if block.Header.Version != CurrentBlockVersion {
		return fmt.Errorf("创世区块版本错误: 期望%d, 实际%d", CurrentBlockVersion, block.Header.Version)
	}

	// 验证前块哈希（应该为全零）
	zeroHash := [32]byte{}
	if block.Header.PrevBlockHash != zeroHash {
		return fmt.Errorf("创世区块前块哈希必须为全零")
	}

	// 验证时间戳
	if block.Header.Timestamp != GenesisBlockTimestamp {
		return fmt.Errorf("创世区块时间戳错误: 期望%d, 实际%d", GenesisBlockTimestamp, block.Header.Timestamp)
	}

	// 验证难度位
	if block.Header.Bits != GenesisBlockBits {
		return fmt.Errorf("创世区块难度位错误: 期望%x, 实际%x", GenesisBlockBits, block.Header.Bits)
	}

	// 验证Nonce
	if block.Header.Nonce != GenesisBlockNonce {
		return fmt.Errorf("创世区块Nonce错误: 期望%d, 实际%d", GenesisBlockNonce, block.Header.Nonce)
	}

	// 验证交易数量（创世区块应该只有一个Coinbase交易）
	if len(block.Transactions) != 1 {
		return fmt.Errorf("创世区块应该只有一个交易, 实际有%d个", len(block.Transactions))
	}

	// 验证Merkle根
	calculatedMerkleRoot := block.GetMerkleRoot()
	if calculatedMerkleRoot != block.Header.MerkleRoot {
		return fmt.Errorf("创世区块Merkle根验证失败")
	}

	// 验证区块哈希
	if !IsGenesisBlock(block) {
		return fmt.Errorf("创世区块哈希验证失败")
	}

	return nil
}

// GetGenesisBlockInfo 获取创世区块的详细信息
//
// 功能说明：
// 返回创世区块的完整信息映射，用于API响应、调试和监控
// 提供人类可读的格式和程序可处理的数据结构
//
// 信息包含：
// - hash: 区块哈希的十六进制字符串表示
// - version: 区块版本号
// - prevHash: 前块哈希的十六进制字符串（创世区块为全零）
// - merkleRoot: Merkle根哈希的十六进制字符串
// - timestamp: Unix时间戳（秒）
// - timestampStr: 人类可读的时间字符串（RFC3339格式）
// - bits: 难度位的十六进制表示
// - nonce: 工作量证明的随机数
// - difficulty: 计算得出的难度值
// - size: 区块的字节大小
// - txCount: 区块中的交易数量
//
// 返回值：
// map[string]interface{} - 包含所有创世区块信息的映射
func GetGenesisBlockInfo() map[string]interface{} {
	genesis := GetGenesisBlock()
	hash := GetGenesisBlockHash()

	return map[string]interface{}{
		"hash":         hex.EncodeToString(hash[:]),
		"version":      genesis.Header.Version,
		"prevHash":     hex.EncodeToString(genesis.Header.PrevBlockHash[:]),
		"merkleRoot":   hex.EncodeToString(genesis.Header.MerkleRoot[:]),
		"timestamp":    genesis.Header.Timestamp,
		"timestampStr": time.Unix(int64(genesis.Header.Timestamp), 0).UTC().Format(time.RFC3339),
		"bits":         genesis.Header.Bits,
		"nonce":        genesis.Header.Nonce,
		"difficulty":   genesis.Header.GetDifficulty(),
		"size":         genesis.Size(),
		"txCount":      len(genesis.Transactions),
	}
}

// CreateTestGenesisBlock 创建用于测试的创世区块
//
// 功能说明：
// 创建一个用于测试环境的创世区块，与主网创世区块不同
// 使用动态参数（当前时间戳和随机Nonce）以避免与主网冲突
//
// 测试区块特征：
// - 使用当前时间作为时间戳（而非固定的历史时间）
// - 生成随机Nonce值（而非预定义值）
// - 使用最低难度设置（便于测试挖矿）
// - 包含测试专用的Coinbase消息
// - 保持与主网相同的基本结构
//
// 实现步骤：
// 1. 获取当前Unix时间戳
// 2. 生成加密安全的随机Nonce
// 3. 创建测试专用的Coinbase交易
// 4. 计算交易的Merkle根
// 5. 构建测试区块头（使用最低难度）
// 6. 组装完整的测试创世区块
//
// 使用场景：
// - 单元测试中的区块链初始化
// - 集成测试中的独立网络环境
// - 开发环境中的快速原型验证
//
// 返回值：
// *Block - 测试用创世区块实例
func CreateTestGenesisBlock() *Block {
	// 使用当前时间作为时间戳
	timestamp := uint32(time.Now().Unix())

	// 生成随机Nonce
	nonce := utils.GenerateNonce()

	// 创建测试Coinbase交易
	testMessage := "Test Genesis Block for Simplified Bitcoin Network"
	coinbaseData := []byte(testMessage)
	coinbaseTx := NewTransaction(coinbaseData)

	// 计算Merkle根
	txHashes := [][]byte{coinbaseTx.Hash[:]}
	merkleRoot := utils.MerkleRoot(txHashes)
	var merkleRootArray [32]byte
	copy(merkleRootArray[:], merkleRoot)

	// 创建测试创世区块头
	var prevBlockHash [32]byte
	header := NewBlockHeader(
		CurrentBlockVersion,
		prevBlockHash,
		merkleRootArray,
		timestamp,
		MaxTargetBits, // 使用最低难度
		nonce,
	)

	// 创建测试创世区块
	return NewBlock(header, []*Transaction{coinbaseTx})
}

// ResetGenesisBlock 重置全局创世区块实例
//
// 功能说明：
// 重新创建并设置全局创世区块实例，主要用于测试场景
// 确保测试之间的状态隔离和一致性
//
// 操作步骤：
// 1. 调用CreateGenesisBlock()创建新的创世区块
// 2. 更新全局GenesisBlock变量
// 3. 重新计算并缓存创世区块哈希
//
// 使用场景：
// - 测试用例之间的状态重置
// - 测试环境的初始化
// - 调试过程中的状态恢复
//
// 注意事项：
// - 仅应在测试环境中使用
// - 生产环境中创世区块应保持不变
// - 重置后所有依赖创世区块的状态都会改变
func ResetGenesisBlock() {
	GenesisBlock = CreateGenesisBlock()
	genesisBlockHash = GenesisBlock.Hash()
}
