#!/bin/bash

# 开发环境启动脚本
echo "启动开发环境..."

# 设置GO环境变量
export GO111MODULE=on
export GOPROXY=https://goproxy.cn,direct

# 检查并安装air
if ! command -v air &> /dev/null; then
    echo "正在安装air..."
    go install github.com/cosmtrek/air@latest
fi

# 确保在项目根目录运行
cd "$(dirname "$0")/.."

# 创建临时目录
mkdir -p tmp

# 使用air启动项目
air -c scripts/.air.toml