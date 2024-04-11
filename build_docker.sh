#!/bin/bash

# 设置 Dockerfile 文件路径和生成的镜像名称
dockerfile="Dockerfile"
image_name="quant_api"

# 构建镜像
docker build -f "$dockerfile" -t "$image_name" .

# 标记镜像
docker tag "$image_name" swr.cn-north-4.myhuaweicloud.com/quant-group/"$image_name":latest

# 推送镜像
docker push swr.cn-north-4.myhuaweicloud.com/quant-group/"$image_name":latest

