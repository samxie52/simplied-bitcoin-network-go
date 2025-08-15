package utils

import "time"

// 版本信息
const (
	AppName    = "Simplified Bitcoin Network"
	AppVersion = "1.0.0"

	// 协议版本
	ProtocolVersion = 1

	// 网络魔数
	MainNetMagic = 0xD9B4BEF9
	TestNetMagic = 0xDAB5BFFA
	RegTestMagic = 0xFABFB5DA
)

// 区块链常量
const (
	// 创世区块时间戳
	GenesisTimestamp = 1635724800 // 2021-11-01 00:00:00 UTC

	// 最大区块大小
	MaxBlockSize = 1 * 1024 * 1024 // 1MB

	// 目标出块时间
	TargetBlockTime = 10 * time.Minute

	// 难度调整间隔
	DifficultyAdjustmentInterval = 2016

	// 最大供应量
	MaxSupply = 21_000_000

	// 初始区块奖励
	InitialBlockReward = 50

	// 奖励减半间隔
	HalvingInterval = 210_000
)

// 网络常量
const (
	// 默认网络端口
	DefaultNetworkPort = 8080

	// 默认RPC端口
	DefaultRPCPort = 8545

	// 最大连接数
	MaxConnections = 125

	// 连接超时时间
	ConnectionTimeout = 30 * time.Second

	// 心跳间隔
	HeartbeatInterval = 60 * time.Second

	// 消息最大大小
	MaxMessageSize = 32 * 1024 * 1024 // 32MB
)

// 挖矿常量
const (
	// 最大目标值（最低难度）
	MaxTarget = 0x1d00ffff

	// 最大nonce值
	MaxNonce = 0xffffffff

	// Coinbase成熟确认数
	CoinbaseMaturity = 100
)

// 交易常量
const (
	// 最大交易大小
	MaxTransactionSize = 100 * 1024 // 100KB

	// 最小交易费
	MinTransactionFee = 1000 // satoshis

	// 灰尘阈值
	DustThreshold = 546 // satoshis

	// 最大输入数
	MaxTransactionInputs = 10000

	// 最大输出数
	MaxTransactionOutputs = 10000
)

// 地址常量
const (
	// 地址版本
	MainNetAddressVersion = 0x00
	TestNetAddressVersion = 0x6F

	// 地址长度
	AddressLength = 25

	// Base58字符集
	Base58Alphabet = "123456789ABCDEFGHJKLMNPQRSTUVWXYZabcdefghijkmnopqrstuvwxyz"
)

// 数据库常量
const (
	// 数据库桶名
	BlocksBucket     = "blocks"
	ChainStateBucket = "chainstate"
	UTXOBucket       = "utxo"
	WalletsBucket    = "wallets"
	PeersBucket      = "peers"

	// 缓存大小
	DefaultCacheSize = 100 * 1024 * 1024 // 100MB
)

// HTTP常量
const (
	// API版本
	APIVersion = "v1"

	// API基础路径
	APIBasePath = "/api/" + APIVersion

	// 请求限制
	MaxRequestSize = 10 * 1024 * 1024 // 10MB
	RequestTimeout = 30 * time.Second

	// CORS
	CORSMaxAge = 86400 // 24小时
)

// 错误码
const (
	ErrCodeSuccess            = 0
	ErrCodeInvalidParameter   = 1001
	ErrCodeNotFound           = 1002
	ErrCodeInternalError      = 1003
	ErrCodeUnauthorized       = 1004
	ErrCodeRateLimited        = 1005
	ErrCodeInvalidTransaction = 2001
	ErrCodeInvalidBlock       = 2002
	ErrCodeChainError         = 2003
)
