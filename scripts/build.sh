#!/bin/bash

# 构建脚本
echo "开始构建项目..."

# 设置GO环境变量
export GO111MODULE=on
export GOPROXY=https://goproxy.cn,direct

# 清理之前的构建
rm -rf ./build

# 创建构建目录
mkdir -p ./build

# 编译
go build -o ./build/app ./cmd/main.go

# 复制配置文件
cp -r ./config ./build/

echo "构建完成，输出目录: ./build"