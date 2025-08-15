#!/bin/bash

echo "=== 简化版比特币网络 - 环境验证 ==="

# 检查Go版本
echo "检查Go版本..."
if ! command -v go &> /dev/null; then
    echo "❌ Go未安装"
    exit 1
fi

GO_VERSION=$(go version | grep -o 'go[0-9]\+\.[0-9]\+' | sed 's/go//')
MIN_VERSION="1.24"

if [ "$(printf '%s\n' "$MIN_VERSION" "$GO_VERSION" | sort -V | head -n1)" != "$MIN_VERSION" ]; then
    echo "❌ Go版本过低，需要1.24+，当前版本：$GO_VERSION"
    exit 1
fi
echo "✅ Go版本：$GO_VERSION"

# 检查Make
echo "检查Make..."
if ! command -v make &> /dev/null; then
    echo "❌ Make未安装"
    exit 1
fi
echo "✅ Make已安装"

# 检查Git
echo "检查Git..."
if ! command -v git &> /dev/null; then
    echo "❌ Git未安装"
    exit 1
fi
echo "✅ Git已安装"

# 验证Go模块
echo "验证Go模块..."
if [ ! -f "go.mod" ]; then
    echo "❌ go.mod文件不存在"
    exit 1
fi
echo "✅ go.mod文件存在"

# 下载依赖
echo "下载依赖..."
go mod download
if [ $? -ne 0 ]; then
    echo "❌ 依赖下载失败"
    exit 1
fi
echo "✅ 依赖下载成功"

# 验证项目结构
echo "验证项目结构..."
REQUIRED_DIRS=("cmd" "pkg" "web" "test" "config" "docs")
for dir in "${REQUIRED_DIRS[@]}"; do
    if [ ! -d "$dir" ]; then
        echo "❌ 目录不存在: $dir"
        exit 1
    fi
done
echo "✅ 项目结构正确"

# 验证配置文件
echo "验证配置文件..."
if [ ! -f "config/config.yaml" ]; then
    echo "❌ 配置文件不存在: config/config.yaml"
    exit 1
fi
echo "✅ 配置文件存在"

# 测试构建
echo "测试构建..."
make clean > /dev/null 2>&1
if ! make build > /dev/null 2>&1; then
    echo "❌ 构建失败"
    exit 1
fi
echo "✅ 构建成功"

echo ""
echo "🎉 环境配置验证完成！所有检查项目都已通过。"
echo "📝 下一步可以开始开发区块链核心功能。"