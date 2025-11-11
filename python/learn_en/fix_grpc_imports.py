#!/usr/bin/env python3
# coding=utf-8
"""
自动修复 protobuf 生成的 Python gRPC 文件中的导入语句
将绝对导入改为相对导入，以避免模块导入错误
"""

import os
import re
import sys


def fix_grpc_imports(file_path):
    """
    修复 gRPC 生成文件中的导入语句
    将 'import xxx_pb2' 改为 'from . import xxx_pb2'
    """
    if not os.path.exists(file_path):
        print(f"文件不存在: {file_path}")
        return False
    
    with open(file_path, 'r', encoding='utf-8') as f:
        content = f.read()
    
    # 匹配模式: import xxx_pb2 as xxx__pb2
    # 替换为: from . import xxx_pb2 as xxx__pb2
    pattern = r'^import (\w+_pb2) as (\w+__pb2)$'
    replacement = r'from . import \1 as \2'
    
    new_content, count = re.subn(pattern, replacement, content, flags=re.MULTILINE)
    
    if count > 0:
        with open(file_path, 'w', encoding='utf-8') as f:
            f.write(new_content)
        print(f"✓ 已修复 {file_path} ({count} 处导入)")
        return True
    else:
        print(f"- {file_path} 无需修复")
        return False


def fix_all_grpc_files(proto_dir):
    """
    修复目录下所有的 *_grpc.py 文件
    """
    if not os.path.exists(proto_dir):
        print(f"目录不存在: {proto_dir}")
        return
    
    print(f"正在扫描目录: {proto_dir}")
    
    fixed_count = 0
    for filename in os.listdir(proto_dir):
        if filename.endswith('_grpc.py'):
            file_path = os.path.join(proto_dir, filename)
            if fix_grpc_imports(file_path):
                fixed_count += 1
    
    print(f"\n共修复 {fixed_count} 个文件")


if __name__ == '__main__':
    # 默认处理当前目录下的 proto 目录
    script_dir = os.path.dirname(os.path.abspath(__file__))
    proto_dir = os.path.join(script_dir, 'proto')
    
    # 也可以从命令行参数指定目录
    if len(sys.argv) > 1:
        proto_dir = sys.argv[1]
    
    print("=" * 60)
    print("gRPC Python 导入修复工具")
    print("=" * 60)
    
    fix_all_grpc_files(proto_dir)
    
    print("\n完成!")
