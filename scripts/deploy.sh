#!/bin/bash

GREEN='\033[0;32m'
RED='\033[0;31m'
NC='\033[0m' # No Color

error_exit() {
    echo -e "${RED}错误: $1${NC}" >&2
    exit 1
}

# 切换到项目根目录
cd "$(dirname "$0")/.." || error_exit "无法切换到项目根目录"

if [[ ! -w ./ ]]; then
    error_exit "当前目录没有写入权限"
fi

export CGO_ENABLED=0
export GOOS=linux
export GOARCH=amd64

PROJECT_NAME=$(basename $(pwd))
REMOTE_HOST="root@vtyug.cn"
REMOTE_PATH="/home/rpc-master-serve"
BINARY_NAME="rpc-master-serve"

echo -e "${GREEN}开始构建...${NC}"
go build -o ${BINARY_NAME} . || error_exit "构建失败"
echo -e "${GREEN}构建完成${NC}"

echo -e "${GREEN}开始上传...${NC}"
ssh ${REMOTE_HOST} "mkdir -p ${REMOTE_PATH}" || error_exit "创建远程目录失败"
rsync ${BINARY_NAME} ${REMOTE_HOST}:"${REMOTE_PATH}" || error_exit "上传失败"
echo -e "${GREEN}上传完成${NC}"

echo -e "${GREEN}重启远程服务...${NC}"
ssh ${REMOTE_HOST} "cd ${REMOTE_PATH} && ./restart.sh" || error_exit "重启服务失败"
echo -e "${GREEN}服务重启完成${NC}"

echo -e "${GREEN}清理本地构建文件...${NC}"
rm -f ${BINARY_NAME}
echo -e "${GREEN}部署完成${NC}"