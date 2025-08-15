package utils

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	"gopkg.in/yaml.v3"
)

// Config 主配置结构
type Config struct {
	App        AppConfig        `yaml:"app"`
	Network    NetworkConfig    `yaml:"network"`
	RPC        RPCConfig        `yaml:"rpc"`
	Mining     MiningConfig     `yaml:"mining"`
	Blockchain BlockchainConfig `yaml:"blockchain"`
	Database   DatabaseConfig   `yaml:"database"`
	Web        WebConfig        `yaml:"web"`
	Logging    LoggingConfig    `yaml:"logging"`
	Security   SecurityConfig   `yaml:"security"`
}

// AppConfig 应用配置
type AppConfig struct {
	Name     string `yaml:"name"`
	Version  string `yaml:"version"`
	Debug    bool   `yaml:"debug"`
	LogLevel string `yaml:"log_level"`
}

// NetworkConfig 网络配置
type NetworkConfig struct {
	Port              int      `yaml:"port"`
	MaxConnections    int      `yaml:"max_connections"`
	ConnectionTimeout int      `yaml:"connection_timeout"`
	HeartbeatInterval int      `yaml:"heartbeat_interval"`
	Seeds             []string `yaml:"seeds"`
}

// RPCConfig RPC配置
type RPCConfig struct {
	Port       int    `yaml:"port"`
	EnableCORS bool   `yaml:"enable_cors"`
	RateLimit  int    `yaml:"rate_limit"`
	EnableAuth bool   `yaml:"enable_auth"`
	AuthKey    string `yaml:"auth_key"`
}

// MiningConfig 挖矿配置
type MiningConfig struct {
	Enabled      bool   `yaml:"enabled"`
	MinerAddress string `yaml:"miner_address"`
	Threads      int    `yaml:"threads"`
	BlockTime    int    `yaml:"block_time"`
}

// BlockchainConfig 区块链配置
type BlockchainConfig struct {
	DataDir                      string `yaml:"data_dir"`
	MaxBlockSize                 int    `yaml:"max_block_size"`
	GenesisDifficulty            uint32 `yaml:"genesis_difficulty"`
	DifficultyAdjustmentInterval int    `yaml:"difficulty_adjustment_interval"`
	MaxSupply                    int64  `yaml:"max_supply"`
}

// DatabaseConfig 数据库配置
type DatabaseConfig struct {
	Type      string `yaml:"type"`
	Path      string `yaml:"path"`
	CacheSize int    `yaml:"cache_size"`
	BatchSize int    `yaml:"batch_size"`
}

// WebConfig Web配置
type WebConfig struct {
	StaticDir   string `yaml:"static_dir"`
	TemplateDir string `yaml:"template_dir"`
	EnableGzip  bool   `yaml:"enable_gzip"`
}

// LoggingConfig 日志配置
type LoggingConfig struct {
	Level      string `yaml:"level"`
	Format     string `yaml:"format"`
	Output     string `yaml:"output"`
	FilePath   string `yaml:"file_path"`
	MaxSize    int    `yaml:"max_size"`
	MaxBackups int    `yaml:"max_backups"`
	MaxAge     int    `yaml:"max_age"`
}

// SecurityConfig 安全配置
type SecurityConfig struct {
	TLS            TLSConfig `yaml:"tls"`
	MaxRequestSize int       `yaml:"max_request_size"`
	RequestTimeout int       `yaml:"request_timeout"`
}

// TLSConfig TLS配置
type TLSConfig struct {
	Enabled  bool   `yaml:"enabled"`
	CertFile string `yaml:"cert_file"`
	KeyFile  string `yaml:"key_file"`
}

// DefaultConfig 返回默认配置
func DefaultConfig() *Config {
	return &Config{
		App: AppConfig{
			Name:     "Simplied Bitcoin Network",
			Version:  "1.0.0",
			Debug:    false,
			LogLevel: "info",
		},
		Network: NetworkConfig{
			Port:              8080,
			MaxConnections:    50,
			ConnectionTimeout: 30,
			HeartbeatInterval: 60,
			Seeds:             []string{},
		},
		RPC: RPCConfig{
			Port:       8545,
			EnableCORS: true,
			RateLimit:  1000,
			EnableAuth: false,
			AuthKey:    "",
		},
		Mining: MiningConfig{
			Enabled:      false,
			MinerAddress: "",
			Threads:      0,
			BlockTime:    10,
		},
		Blockchain: BlockchainConfig{
			DataDir:                      "./data",
			MaxBlockSize:                 1048576,
			GenesisDifficulty:            0x1d00ffff,
			DifficultyAdjustmentInterval: 2016,
			MaxSupply:                    21000000,
		},
		Database: DatabaseConfig{
			Type:      "bolt",
			Path:      "./data/blockchain.db",
			CacheSize: 100,
			BatchSize: 1000,
		},
		Web: WebConfig{
			StaticDir:   "./web",
			TemplateDir: "./web/templates",
			EnableGzip:  true,
		},
		Logging: LoggingConfig{
			Level:      "info",
			Format:     "text",
			Output:     "stdout",
			FilePath:   "./logs/bitcoin-network.log",
			MaxSize:    100,
			MaxBackups: 7,
			MaxAge:     30,
		},
		Security: SecurityConfig{
			TLS: TLSConfig{
				Enabled:  false,
				CertFile: "",
				KeyFile:  "",
			},
			MaxRequestSize: 10485760,
			RequestTimeout: 30,
		},
	}
}

// LoadConfig 从文件加载配置
func LoadConfig(configPath string) (*Config, error) {
	config := DefaultConfig()

	// 检查配置文件是否存在
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		return config, fmt.Errorf("配置文件不存在: %s", configPath)
	}

	// 读取配置文件
	data, err := os.ReadFile(configPath)
	if err != nil {
		return config, fmt.Errorf("读取配置文件失败: %v", err)
	}

	// 解析YAML配置
	if err := yaml.Unmarshal(data, config); err != nil {
		return config, fmt.Errorf("解析配置文件失败: %v", err)
	}

	// 验证配置
	if err := config.Validate(); err != nil {
		return config, fmt.Errorf("配置验证失败: %v", err)
	}

	// 创建必要的目录
	if err := config.CreateDirectories(); err != nil {
		return config, fmt.Errorf("创建目录失败: %v", err)
	}

	return config, nil
}

// Validate 验证配置有效性
func (c *Config) Validate() error {
	if c.Network.Port <= 0 || c.Network.Port > 65535 {
		return fmt.Errorf("无效的网络端口: %d", c.Network.Port)
	}

	if c.RPC.Port <= 0 || c.RPC.Port > 65535 {
		return fmt.Errorf("无效的RPC端口: %d", c.RPC.Port)
	}

	if c.Blockchain.MaxBlockSize <= 0 {
		return fmt.Errorf("无效的最大区块大小: %d", c.Blockchain.MaxBlockSize)
	}

	if c.Mining.BlockTime <= 0 {
		return fmt.Errorf("无效的目标出块时间: %d", c.Mining.BlockTime)
	}

	return nil
}

// CreateDirectories 创建必要的目录
func (c *Config) CreateDirectories() error {
	dirs := []string{
		c.Blockchain.DataDir,
		filepath.Dir(c.Database.Path),
	}

	if c.Logging.Output == "file" {
		dirs = append(dirs, filepath.Dir(c.Logging.FilePath))
	}

	for _, dir := range dirs {
		if err := os.MkdirAll(dir, 0755); err != nil {
			return fmt.Errorf("创建目录 %s 失败: %v", dir, err)
		}
	}

	return nil
}

// GetConnectionTimeout 获取连接超时时间
func (c *Config) GetConnectionTimeout() time.Duration {
	return time.Duration(c.Network.ConnectionTimeout) * time.Second
}

// GetHeartbeatInterval 获取心跳间隔
func (c *Config) GetHeartbeatInterval() time.Duration {
	return time.Duration(c.Network.HeartbeatInterval) * time.Second
}

// GetRequestTimeout 获取请求超时时间
func (c *Config) GetRequestTimeout() time.Duration {
	return time.Duration(c.Security.RequestTimeout) * time.Second
}

// GetBlockTime 获取目标出块时间
func (c *Config) GetBlockTime() time.Duration {
	return time.Duration(c.Mining.BlockTime) * time.Second
}
