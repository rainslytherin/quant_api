FROM golang:1.22.0-alpine3.19 AS build

ARG app
ENV GOPROXY="https://goproxy.cn,direct"
WORKDIR /
COPY go.mod go.sum ./
RUN go mod download
COPY . ./

RUN set -eux && sed -i 's/dl-cdn.alpinelinux.org/mirrors.ustc.edu.cn/g' /etc/apk/repositories
RUN apk update && apk --no-cache add tzdata
RUN go build -o quant_api main.go

FROM alpine:3.19

ARG env
ARG app
WORKDIR /app
COPY --from=build /quant_api ./quant_api
COPY ./etc ./etc

# 安装 tzdata 包
RUN set -eux && sed -i 's/dl-cdn.alpinelinux.org/mirrors.ustc.edu.cn/g' /etc/apk/repositories
RUN apk update && apk --no-cache add tzdata
RUN ln -sf /usr/share/zoneinfo/Asia/Shanghai /etc/localtime
RUN chmod 755 /app/quant_api

ENTRYPOINT /app/quant_api -c /app/etc/config.json 
