# Step 1.3: 区块数据结构定义

## 📋 概述

**目标**: 实现比特币区块的完整数据结构，包括区块头、区块体、序列化机制和基础验证功能。

**前置条件**: Step 1.2 (密码学工具) 完成  
**预计时间**: 2-3天  
**难度级别**: ⭐⭐⭐

## 🎯 功能要求

### 核心功能
- **区块头结构**: 实现标准比特币区块头格式
- **区块结构**: 完整的区块数据结构定义
- **序列化机制**: 二进制序列化和反序列化
- **哈希计算**: 区块哈希和验证机制
- **创世区块**: 网络初始化区块生成
- **大小限制**: 区块大小验证和限制

### 技术规格
- **区块头大小**: 固定80字节
- **最大区块大小**: 1MB (1,000,000字节)
- **哈希算法**: 双重SHA-256
- **字节序**: 小端序 (Little Endian)
- **时间戳精度**: Unix时间戳 (秒)

## 🏗️ 架构设计

### 数据结构层次
```
Block
├── Header (BlockHeader)
│   ├── Version (4字节)
│   ├── PrevBlockHash (32字节)
│   ├── MerkleRoot (32字节)
│   ├── Timestamp (4字节)
│   ├── Bits (4字节)
│   └── Nonce (4字节)
└── Transactions ([]Transaction)
    └── TransactionCount (VarInt)
```

### 区块头字段详解

| 字段 | 大小 | 类型 | 描述 |
|------|------|------|------|
| Version | 4字节 | uint32 | 区块版本号，用于软分叉升级 |
| PrevBlockHash | 32字节 | [32]byte | 前一个区块的哈希值 |
| MerkleRoot | 32字节 | [32]byte | 交易Merkle树根哈希 |
| Timestamp | 4字节 | uint32 | 区块创建时间戳 |
| Bits | 4字节 | uint32 | 难度目标的紧凑表示 |
| Nonce | 4字节 | uint32 | 工作量证明随机数 |

## 📁 文件结构

```
pkg/blockchain/
├── block.go          # 区块结构定义和方法
├── genesis.go        # 创世区块创建
└── constants.go      # 区块链常量定义

test/blockchain/
├── block_test.go     # 区块结构测试
└── genesis_test.go   # 创世区块测试
```

## 🔧 实现内容

### 1. 区块头结构 (BlockHeader)

**字段定义**:
```go
type BlockHeader struct {
    Version       uint32    // 区块版本
    PrevBlockHash [32]byte  // 前块哈希
    MerkleRoot    [32]byte  // Merkle根
    Timestamp     uint32    // 时间戳
    Bits          uint32    // 难度目标
    Nonce         uint32    // 随机数
}
```

**核心方法**:
- `Hash() [32]byte` - 计算区块头哈希
- `Serialize() []byte` - 序列化为字节数组
- `Deserialize([]byte) error` - 从字节数组反序列化
- `IsValid() bool` - 基础有效性验证

### 2. 区块结构 (Block)

**字段定义**:
```go
type Block struct {
    Header       *BlockHeader     // 区块头
    Transactions []*Transaction   // 交易列表
}
```

**核心方法**:
- `Hash() [32]byte` - 获取区块哈希
- `Size() int` - 计算区块大小
- `Serialize() []byte` - 完整序列化
- `Deserialize([]byte) error` - 完整反序列化
- `Validate() error` - 区块验证
- `GetMerkleRoot() [32]byte` - 计算Merkle根

### 3. 创世区块 (Genesis Block)

**特殊属性**:
- PrevBlockHash: 全零哈希
- 固定时间戳: 比特币创世区块时间
- 特殊交易: Coinbase交易
- 固定难度: 最大难度目标

**创建函数**:
```go
func CreateGenesisBlock() *Block
func GetGenesisBlockHash() [32]byte
func IsGenesisBlock(block *Block) bool
```

## 🧪 测试策略

### 单元测试覆盖

**区块头测试**:
- ✅ 序列化和反序列化一致性
- ✅ 哈希计算正确性
- ✅ 字段验证和边界条件
- ✅ 无效数据处理

**区块测试**:
- ✅ 完整区块序列化
- ✅ 大小计算准确性
- ✅ Merkle根计算验证
- ✅ 区块大小限制验证

**创世区块测试**:
- ✅ 创世区块生成
- ✅ 哈希值固定性
- ✅ 特殊属性验证

### 性能测试

**基准测试**:
- 区块序列化性能
- 哈希计算性能
- 大区块处理性能
- 内存使用优化

## 📊 验证标准

### 功能验证
1. **序列化一致性**: 序列化后反序列化数据完全一致
2. **哈希正确性**: 区块哈希计算符合比特币标准
3. **大小限制**: 正确执行1MB区块大小限制
4. **创世区块**: 生成固定且可重现的创世区块

### 性能指标
- **序列化速度**: >1000 blocks/sec
- **哈希计算**: <1ms per block
- **内存使用**: <10MB for 1000 blocks
- **并发安全**: 支持多goroutine访问

## 🔗 依赖关系

### 输入依赖
- `pkg/utils/hash.go` - DoubleSHA256, MerkleRoot
- `pkg/utils/encoding.go` - 字节序转换, VarInt
- `pkg/utils/crypto.go` - 随机数生成
- `time` - 时间戳处理

### 输出接口
- 为Step 1.4 (Merkle树) 提供区块结构
- 为Step 1.5 (区块链) 提供区块定义
- 为挖矿模块提供区块模板

## 📝 使用示例

### 创建新区块
```go
// 创建区块头
header := &BlockHeader{
    Version:       1,
    PrevBlockHash: prevHash,
    MerkleRoot:    merkleRoot,
    Timestamp:     uint32(time.Now().Unix()),
    Bits:          difficulty,
    Nonce:         0,
}

// 创建区块
block := &Block{
    Header:       header,
    Transactions: transactions,
}

// 计算区块哈希
blockHash := block.Hash()
```

### 序列化和反序列化
```go
// 序列化区块
data := block.Serialize()

// 反序列化区块
newBlock := &Block{}
err := newBlock.Deserialize(data)
if err != nil {
    log.Fatal("反序列化失败:", err)
}
```

### 创世区块使用
```go
// 获取创世区块
genesis := CreateGenesisBlock()

// 验证是否为创世区块
if IsGenesisBlock(genesis) {
    fmt.Println("创世区块验证成功")
}
```

## ⚠️ 注意事项

### 安全考虑
- **哈希验证**: 所有区块哈希必须验证
- **大小限制**: 严格执行区块大小限制
- **时间戳验证**: 防止时间戳攻击
- **输入验证**: 所有外部输入必须验证

### 性能优化
- **内存池化**: 重用序列化缓冲区
- **并发安全**: 使用只读操作避免锁
- **缓存机制**: 缓存计算结果
- **批量处理**: 支持批量序列化

### 兼容性
- **字节序**: 严格遵循小端序
- **版本控制**: 支持未来版本升级
- **向后兼容**: 保持API稳定性

## 🚀 下一步

完成Step 1.3后，将具备：
- ✅ 完整的区块数据结构
- ✅ 可靠的序列化机制
- ✅ 标准的哈希计算
- ✅ 创世区块支持

**准备进入Step 1.4**: Merkle树实现，利用已定义的区块结构构建交易验证机制。
