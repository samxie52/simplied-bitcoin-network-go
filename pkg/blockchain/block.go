package blockchain

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"time"

	"simplied-bitcoin-network-go/pkg/utils"
)

// block.go - 区块链核心数据结构定义文件
//
// 本文件定义了简化版比特币网络的核心数据结构，包括：
// - BlockHeader: 区块头结构，包含区块的元数据信息
// - Transaction: 简化的交易结构，用于Step 1.3阶段
// - Block: 完整的区块结构，包含区块头和交易列表
//
// 主要功能模块：
// 1. 数据结构定义：定义区块链的基本数据类型
// 2. 序列化/反序列化：支持区块数据的持久化存储
// 3. 哈希计算：实现区块和交易的哈希算法
// 4. 验证功能：提供区块完整性和有效性验证
// 5. 工具方法：提供各种便利的查询和操作方法

// BlockHeader 区块头结构，固定80字节
//
// 区块头是区块链中每个区块的核心元数据，包含了区块的所有关键信息。
// 设计为固定80字节的结构，确保网络传输和存储的效率。
//
// 字段说明：
// - Version: 区块版本号，用于协议升级和兼容性管理
// - PrevBlockHash: 前一个区块的哈希值，形成区块链的链式结构
// - MerkleRoot: 交易Merkle树的根哈希，确保交易数据的完整性
// - Timestamp: 区块创建的Unix时间戳，用于时间验证和难度调整
// - Bits: 难度目标的紧凑表示，定义了工作量证明的难度要求
// - Nonce: 工作量证明的随机数，矿工通过调整此值来满足难度要求
type BlockHeader struct {
	Version       uint32   // 4字节 - 区块版本号，当前使用版本1
	PrevBlockHash [32]byte // 32字节 - 前一个区块的SHA-256哈希值
	MerkleRoot    [32]byte // 32字节 - 交易Merkle树根的SHA-256哈希值
	Timestamp     uint32   // 4字节 - 区块创建时间的Unix时间戳
	Bits          uint32   // 4字节 - 难度目标的紧凑表示（nBits格式）
	Nonce         uint32   // 4字节 - 工作量证明随机数，范围0到2^32-1
}

// Transaction 简化的交易结构（Step 1.3阶段使用）
//
// 这是一个简化的交易实现，用于区块链开发的早期阶段。
// 在完整实现中，交易结构会包含输入、输出、脚本等复杂字段。
//
// 字段说明：
// - Hash: 交易的SHA-256哈希值，用作交易的唯一标识符
// - Data: 交易的原始数据，包含所有交易信息的序列化结果
type Transaction struct {
	Hash [32]byte // 交易的SHA-256哈希值，32字节固定长度
	Data []byte   // 交易的序列化数据，变长字段
}

// Block 区块结构
//
// 区块是区块链的基本单位，包含区块头和交易列表。
// 每个区块通过区块头中的PrevBlockHash字段与前一个区块链接，
// 形成不可篡改的链式数据结构。
//
// 字段说明：
// - Header: 指向区块头结构的指针，包含区块的元数据
// - Transactions: 交易列表，包含区块中的所有交易
type Block struct {
	Header       *BlockHeader   // 区块头指针，包含区块的核心元数据
	Transactions []*Transaction // 交易列表，包含区块中的所有交易
}

// NewBlockHeader 创建新的区块头
//
// 根据给定的参数创建新的区块头结构。
//
// 参数：
// - version: 区块版本号
// - prevBlockHash: 前一个区块的哈希值
// - merkleRoot: 交易Merkle树的根哈希
// - timestamp: 区块创建时间的Unix时间戳
// - bits: 难度目标的紧凑表示
// - nonce: 工作量证明的随机数
func NewBlockHeader(version uint32, prevBlockHash [32]byte, merkleRoot [32]byte, timestamp uint32, bits uint32, nonce uint32) *BlockHeader {
	return &BlockHeader{
		Version:       version,
		PrevBlockHash: prevBlockHash,
		MerkleRoot:    merkleRoot,
		Timestamp:     timestamp,
		Bits:          bits,
		Nonce:         nonce,
	}
}

// NewBlock 创建新的区块
//
// 根据给定的区块头和交易列表创建新的区块结构。
//
// 参数：
// - header: 指向区块头结构的指针
// - transactions: 交易列表
func NewBlock(header *BlockHeader, transactions []*Transaction) *Block {
	return &Block{
		Header:       header,
		Transactions: transactions,
	}
}

// Hash 计算区块头哈希
//
// 根据区块头结构计算SHA-256哈希值。
func (bh *BlockHeader) Hash() [32]byte {
	data := bh.Serialize()
	hash := utils.DoubleSHA256(data)
	var result [32]byte
	copy(result[:], hash)
	return result
}

// Serialize 序列化区块头为80字节数组
//
// 将区块头结构序列化为固定长度的字节数组。
func (bh *BlockHeader) Serialize() []byte {
	buf := make([]byte, BlockHeaderSize)
	offset := 0

	// Version (4字节，小端序)
	binary.LittleEndian.PutUint32(buf[offset:], bh.Version)
	offset += 4

	// PrevBlockHash (32字节)
	copy(buf[offset:], bh.PrevBlockHash[:])
	offset += 32

	// MerkleRoot (32字节)
	copy(buf[offset:], bh.MerkleRoot[:])
	offset += 32

	// Timestamp (4字节，小端序)
	binary.LittleEndian.PutUint32(buf[offset:], bh.Timestamp)
	offset += 4

	// Bits (4字节，小端序)
	binary.LittleEndian.PutUint32(buf[offset:], bh.Bits)
	offset += 4

	// Nonce (4字节，小端序)
	binary.LittleEndian.PutUint32(buf[offset:], bh.Nonce)

	return buf
}

// Deserialize 从字节数组反序列化区块头
//
// 根据给定的字节数组反序列化区块头结构。
//
// 参数：
// - data: 字节数组
func (bh *BlockHeader) Deserialize(data []byte) error {
	if len(data) != BlockHeaderSize {
		return fmt.Errorf("区块头数据长度错误: 期望%d字节, 实际%d字节", BlockHeaderSize, len(data))
	}

	offset := 0

	// Version (4字节，小端序)
	bh.Version = binary.LittleEndian.Uint32(data[offset:])
	offset += 4

	// PrevBlockHash (32字节)
	copy(bh.PrevBlockHash[:], data[offset:offset+32])
	offset += 32

	// MerkleRoot (32字节)
	copy(bh.MerkleRoot[:], data[offset:offset+32])
	offset += 32

	// Timestamp (4字节，小端序)
	bh.Timestamp = binary.LittleEndian.Uint32(data[offset:])
	offset += 4

	// Bits (4字节，小端序)
	bh.Bits = binary.LittleEndian.Uint32(data[offset:])
	offset += 4

	// Nonce (4字节，小端序)
	bh.Nonce = binary.LittleEndian.Uint32(data[offset:])

	return nil
}

// IsValid 验证区块头基础有效性
//
// 验证区块头的版本号、时间戳和难度目标是否有效。
func (bh *BlockHeader) IsValid() bool {
	// 检查版本号
	if bh.Version < BlockVersion1 {
		return false
	}

	// 检查时间戳（不能太远的未来）
	now := uint32(time.Now().Unix())
	maxFutureTime := now + uint32(MaxTimeOffset.Seconds())
	if bh.Timestamp > maxFutureTime {
		return false
	}

	// 检查难度位格式
	if bh.Bits == 0 {
		return false
	}

	return true
}

// GetDifficulty 获取区块难度值
//
// 根据区块头的难度目标计算难度值。
func (bh *BlockHeader) GetDifficulty() float64 {
	target := utils.BitsToTarget(bh.Bits)
	return utils.CalculateDifficulty(target)
}

// MeetsTarget 检查区块哈希是否满足难度目标
//
// 根据区块头的难度目标检查区块哈希是否满足难度要求。
func (bh *BlockHeader) MeetsTarget() bool {
	blockHash := bh.Hash()
	target := utils.BitsToTarget(bh.Bits)
	return utils.IsValidTarget(blockHash[:], target)
}

// Hash 获取区块哈希（等同于区块头哈希）
//
// 获取区块的SHA-256哈希值。
func (b *Block) Hash() [32]byte {
	return b.Header.Hash()
}

// Size 计算区块总大小
//
// 计算区块的总大小，包括区块头和交易列表。
func (b *Block) Size() int {
	size := BlockHeaderSize // 区块头大小

	// 添加交易数量的VarInt大小
	txCount := uint64(len(b.Transactions))
	size += len(utils.EncodeVarInt(txCount))

	// 添加所有交易的大小
	for _, tx := range b.Transactions {
		size += len(tx.Data)
	}

	return size
}

// Serialize 序列化完整区块
//
// 将区块结构序列化为字节数组。
func (b *Block) Serialize() []byte {
	var buf bytes.Buffer

	// 序列化区块头
	headerData := b.Header.Serialize()
	buf.Write(headerData)

	// 序列化交易数量
	txCount := uint64(len(b.Transactions))
	txCountBytes := utils.EncodeVarInt(txCount)
	buf.Write(txCountBytes)

	// 序列化所有交易
	for _, tx := range b.Transactions {
		buf.Write(tx.Data)
	}

	return buf.Bytes()
}

// Deserialize 反序列化完整区块
//
// 根据给定的字节数组反序列化区块结构。
//
// 参数：
// - data: 字节数组
func (b *Block) Deserialize(data []byte) error {
	if len(data) < BlockHeaderSize {
		return fmt.Errorf("数据长度不足，无法包含区块头")
	}

	// 反序列化区块头
	b.Header = &BlockHeader{}
	if err := b.Header.Deserialize(data[:BlockHeaderSize]); err != nil {
		return fmt.Errorf("反序列化区块头失败: %v", err)
	}

	offset := BlockHeaderSize

	// 反序列化交易数量
	if offset >= len(data) {
		return fmt.Errorf("数据长度不足，无法读取交易数量")
	}

	txCount, bytesRead, err := utils.DecodeVarInt(data[offset:])
	if err != nil {
		return fmt.Errorf("解码交易数量失败: %v", err)
	}
	offset += bytesRead

	// 验证交易数量限制
	if txCount > MaxTransactionsPerBlock {
		return fmt.Errorf("交易数量超过限制: %d > %d", txCount, MaxTransactionsPerBlock)
	}

	// 初始化交易列表
	b.Transactions = make([]*Transaction, txCount)

	// 反序列化所有交易（简化版本）
	for i := uint64(0); i < txCount; i++ {
		if offset >= len(data) {
			return fmt.Errorf("数据长度不足，无法读取交易%d", i)
		}

		// 简化处理：假设每个交易数据长度为剩余数据平均分配
		// 在完整实现中，每个交易都有自己的长度字段
		remainingTxs := txCount - i
		remainingData := len(data) - offset
		txDataLen := remainingData / int(remainingTxs)

		if txDataLen <= 0 || offset+txDataLen > len(data) {
			return fmt.Errorf("交易%d数据长度计算错误", i)
		}

		txData := data[offset : offset+txDataLen]
		tx := &Transaction{
			Data: make([]byte, len(txData)),
		}
		copy(tx.Data, txData)

		// 计算交易哈希
		hash := utils.DoubleSHA256(txData)
		copy(tx.Hash[:], hash)

		b.Transactions[i] = tx
		offset += txDataLen
	}

	return nil
}

// Validate 验证区块完整性
//
// 验证区块的有效性，包括区块头、交易数量和Merkle根。
func (b *Block) Validate() error {
	// 验证区块头
	if b.Header == nil {
		return fmt.Errorf(ErrInvalidBlockHeader)
	}

	if !b.Header.IsValid() {
		return fmt.Errorf(ErrInvalidBlockHeader)
	}

	// 验证区块大小
	if b.Size() > MaxBlockSize {
		return fmt.Errorf(ErrBlockTooLarge)
	}

	// 验证交易数量
	if len(b.Transactions) == 0 {
		return fmt.Errorf(ErrEmptyBlock)
	}

	if len(b.Transactions) > MaxTransactionsPerBlock {
		return fmt.Errorf(ErrTooManyTransactions)
	}

	// 验证Merkle根
	calculatedMerkleRoot := b.GetMerkleRoot()
	if !bytes.Equal(calculatedMerkleRoot[:], b.Header.MerkleRoot[:]) {
		return fmt.Errorf(ErrInvalidMerkleRoot)
	}

	return nil
}

// GetMerkleRoot 计算区块的Merkle根
//
// 根据交易列表计算Merkle根哈希值。
func (b *Block) GetMerkleRoot() [32]byte {
	if len(b.Transactions) == 0 {
		return [32]byte{} // 空区块返回零哈希
	}

	// 收集所有交易哈希
	txHashes := make([][]byte, len(b.Transactions))
	for i, tx := range b.Transactions {
		txHashes[i] = tx.Hash[:]
	}

	// 计算Merkle根
	merkleRoot := utils.MerkleRoot(txHashes)
	var result [32]byte
	copy(result[:], merkleRoot)
	return result
}

// GetTransactionHashes 获取所有交易哈希
//
// 获取区块中的所有交易哈希值。
func (b *Block) GetTransactionHashes() [][32]byte {
	hashes := make([][32]byte, len(b.Transactions))
	for i, tx := range b.Transactions {
		hashes[i] = tx.Hash
	}
	return hashes
}

// HasTransaction 检查区块是否包含指定交易
//
// 根据交易哈希值检查区块是否包含指定交易。
func (b *Block) HasTransaction(txHash [32]byte) bool {
	for _, tx := range b.Transactions {
		if bytes.Equal(tx.Hash[:], txHash[:]) {
			return true
		}
	}
	return false
}

// GetTransaction 根据哈希获取交易
//
// 根据交易哈希值获取交易结构。
func (b *Block) GetTransaction(txHash [32]byte) *Transaction {
	for _, tx := range b.Transactions {
		if bytes.Equal(tx.Hash[:], txHash[:]) {
			return tx
		}
	}
	return nil
}

// String 返回区块的字符串表示
//
// 返回区块的字符串表示，包括哈希值、交易数量和大小。
func (b *Block) String() string {
	return fmt.Sprintf("Block{Hash: %x, Transactions: %d, Size: %d bytes}",
		b.Hash(), len(b.Transactions), b.Size())
}

// String 返回区块头的字符串表示
//
// 返回区块头的字符串表示，包括版本号、前一个区块哈希值、Merkle根、时间戳、难度目标和随机数。
func (bh *BlockHeader) String() string {
	return fmt.Sprintf("BlockHeader{Version: %d, PrevHash: %x, MerkleRoot: %x, Timestamp: %d, Bits: %x, Nonce: %d}",
		bh.Version, bh.PrevBlockHash, bh.MerkleRoot, bh.Timestamp, bh.Bits, bh.Nonce)
}

// NewTransaction 创建新交易（简化版本）
//
// 根据给定的交易数据创建新交易结构。
//
// 参数：
// - data: 交易数据
func NewTransaction(data []byte) *Transaction {
	hash := utils.DoubleSHA256(data)
	var txHash [32]byte
	copy(txHash[:], hash)

	return &Transaction{
		Hash: txHash,
		Data: data,
	}
}
