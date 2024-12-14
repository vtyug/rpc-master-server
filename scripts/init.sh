#!/bin/bash

# 项目初始化脚本
echo "开始初始化项目..."

# 创建必要的目录
mkdir -p logs
mkdir -p build

# 安装依赖
go mod tidy

# 创建配置文件
if [ ! -f "./config/config.yaml" ]; then
    cp ./config/config.example.yaml ./config/config.yaml
fi

echo "项目初始化完成"