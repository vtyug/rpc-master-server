#!/bin/bash

# 数据库迁移脚本
echo "开始数据库迁移..."

# 检查数据库连接
mysql -h localhost -u root -p fastgo -e "SELECT 1;" > /dev/null 2>&1
if [ $? -ne 0 ]; then
    echo "数据库连接失败"
    exit 1
fi

# 执行迁移SQL
mysql -h localhost -u root -p fastgo < ./scripts/sql/init.sql

echo "数据库迁移完成"