#!/bin/bash
# 生成 Proto 文件并自动修复 Python gRPC 导入

echo "=========================================="
echo "生成 Protocol Buffer 文件"
echo "=========================================="

# 进入 proto 目录
cd "$(dirname "$0")"

# 运行 buf generate
buf generate

if [ $? -eq 0 ]; then
    echo ""
    echo "✓ Proto 文件生成成功"
    echo ""
    echo "=========================================="
    echo "修复 Python gRPC 导入语句"
    echo "=========================================="
    
    # 运行 Python 修复脚本
    python ../../python/learn_en/fix_grpc_imports.py ../../python/learn_en/proto
    
    echo ""
    echo "=========================================="
    echo "全部完成!"
    echo "=========================================="
else
    echo ""
    echo "✗ Proto 文件生成失败"
    exit 1
fi
