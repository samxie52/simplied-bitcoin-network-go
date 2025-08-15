# 变量定义
BINARY_NAME=bitcoin-node
CLI_NAME=bitcoin-cli
MINER_NAME=bitcoin-miner
VERSION=1.0.0
BUILD_DIR=bin
GO_FILES=$(shell find . -name "*.go" -type f)

# 默认目标
.DEFAULT_GOAL := build

# 帮助信息
.PHONY: help
help:
	@echo "可用的构建目标:"
	@echo "  build       - 编译所有程序"
	@echo "  build-node  - 编译节点程序"
	@echo "  build-cli   - 编译CLI工具"
	@echo "  build-miner - 编译挖矿程序"
	@echo "  test        - 运行所有测试"
	@echo "  bench       - 运行基准测试"
	@echo "  clean       - 清理构建文件"
	@echo "  fmt         - 格式化代码"
	@echo "  lint        - 代码检查"
	@echo "  coverage    - 生成测试覆盖率"
	@echo "  docker      - 构建Docker镜像"

# 编译所有程序
.PHONY: build
build: build-node build-cli build-miner

# 编译节点程序
.PHONY: build-node
build-node:
	@echo "编译节点程序..."
	@mkdir -p $(BUILD_DIR)
	@go build -ldflags "-X main.Version=$(VERSION)" -o $(BUILD_DIR)/$(BINARY_NAME) cmd/node/main.go

# 编译CLI工具
.PHONY: build-cli  
build-cli:
	@echo "编译CLI工具..."
	@mkdir -p $(BUILD_DIR)
	@go build -ldflags "-X main.Version=$(VERSION)" -o $(BUILD_DIR)/$(CLI_NAME) cmd/cli/main.go

# 编译挖矿程序
.PHONY: build-miner
build-miner:
	@echo "编译挖矿程序..."
	@mkdir -p $(BUILD_DIR)
	@go build -ldflags "-X main.Version=$(VERSION)" -o $(BUILD_DIR)/$(MINER_NAME) cmd/miner/main.go

# 运行测试
.PHONY: test
test:
	@echo "运行所有测试..."
	@go test -v ./...

# 运行基准测试
.PHONY: bench
bench:
	@echo "运行基准测试..."
	@go test -bench=. -benchmem ./...

# 清理构建文件
.PHONY: clean
clean:
	@echo "清理构建文件..."
	@rm -rf $(BUILD_DIR)
	@go clean

# 格式化代码
.PHONY: fmt
fmt:
	@echo "格式化代码..."
	@go fmt ./...

# 代码检查
.PHONY: lint
lint:
	@echo "运行代码检查..."
	@golangci-lint run

# 生成测试覆盖率
.PHONY: coverage
coverage:
	@echo "生成测试覆盖率报告..."
	@go test -coverprofile=coverage.out ./...
	@go tool cover -html=coverage.out -o coverage.html
	@echo "覆盖率报告已生成: coverage.html"

# 构建Docker镜像
.PHONY: docker
docker:
	@echo "构建Docker镜像..."
	@docker build -t bitcoin-network:$(VERSION) .
	@docker tag bitcoin-network:$(VERSION) bitcoin-network:latest

# 安装依赖
.PHONY: deps
deps:
	@echo "下载并整理依赖..."
	@go mod download
	@go mod tidy

# 开发模式运行
.PHONY: dev
dev: build-node
	@echo "启动开发模式..."
	@./$(BUILD_DIR)/$(BINARY_NAME) --config config/config.yaml --debug

# 检查Go版本
.PHONY: check-go
check-go:
	@echo "检查Go版本..."
	@go version | grep -q "go1.24" || (echo "需要Go 1.24+版本" && exit 1)