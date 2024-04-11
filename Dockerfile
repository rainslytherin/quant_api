FROM golang:1.22.0-alpine3.19 AS build

ARG app
ENV GOPROXY="https://goproxy.cn,direct"
WORKDIR /
COPY go.mod go.sum ./
RUN go mod download
COPY . ./

RUN set -eux && sed -i 's/dl-cdn.alpinelinux.org/mirrors.ustc.edu.cn/g' /etc/apk/repositories
RUN apk update && apk --no-cache add tzdata
RUN go build -o quant_core main.go

FROM alpine:3.19

ARG env
ARG app
WORKDIR /app
COPY --from=build /quant_core ./quant_core
COPY ./etc ./etc

# 安装 tzdata 包
RUN set -eux && sed -i 's/dl-cdn.alpinelinux.org/mirrors.ustc.edu.cn/g' /etc/apk/repositories
RUN apk update && apk --no-cache add tzdata
RUN ln -sf /usr/share/zoneinfo/Asia/Shanghai /etc/localtime
RUN chmod 755 /app/quant_core

ENTRYPOINT /app/quant_core -c /app/etc/config_prod.json 
