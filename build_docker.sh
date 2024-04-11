#!/bin/bash

# 检查参数是否正确
if [ $# -ne 1 ]; then
    echo "使用方式: $0 <test|pre|prod>"
    exit 1
fi

# 设置 Dockerfile 文件路径和生成的镜像名称
dockerfile=""
image_name="quant_core"

# 根据参数选择 Dockerfile 和生成的镜像名称
if [ "$1" == "pre" ]; then
    dockerfile="Dockerfile.pre"
    image_name="${image_name}_pre"
elif [ "$1" == "prod" ]; then
    dockerfile="Dockerfile.prod"
    image_name="${image_name}_prod"
elif [ "$1" == "test" ]; then
    dockerfile="Dockerfile.test"
    image_name="${image_name}_test"
else
    echo "参数错误: 请选择 test, pre 或 prod"
    exit 1
fi

# 构建镜像
docker build -f "$dockerfile" -t "$image_name" .

# 标记镜像
docker tag "$image_name" swr.cn-north-4.myhuaweicloud.com/quant-group/"$image_name":latest

# 推送镜像
docker push swr.cn-north-4.myhuaweicloud.com/quant-group/"$image_name":latest

