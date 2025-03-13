#!/bin/bash
# 1. 删除旧的编译文件
rm -rf build
mkdir -p build

# 2. 编译 Golang 程序
go build -o build/myapp .


echo "打包完成"
