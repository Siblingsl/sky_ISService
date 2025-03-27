#!/bin/bash
set -e  # 发生错误时终止脚本

# 1. 删除旧的编译文件
rm -rf build
mkdir -p build

# 2. 编译总服务
echo "编译总服务..."
go build -trimpath -o build/sky .

# 3. 查找并编译所有子服务
echo "查找并编译子服务..."
find services -type f -path "*/cmd/main.go" | while read -r main_file; do
    # 获取子服务名称 (services/auth/cmd/main.go -> auth)
    service_name=$(basename $(dirname $(dirname "$main_file")))

    echo "编译子服务: $service_name"
    go build -trimpath -o "build/$service_name" "$main_file"
done

echo "所有服务打包完成 ✅"

# 4. 打印当前目录
pwd

# 5. 运行总服务
echo "启动总服务..."
./build/sky
