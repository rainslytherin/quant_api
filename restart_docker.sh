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
    image_name="${image_name}_pre"
elif [ "$1" == "prod" ]; then
    image_name="${image_name}_prod"
elif [ "$1" == "test" ]; then
    image_name="${image_name}_test"
else
    echo "参数错误: 请选择 test, pre 或 prod"
    exit 1
fi

docker stop ${image_name} 
docker rm ${image_name} 
docker run --name ${image_name} -d ${image_name}
