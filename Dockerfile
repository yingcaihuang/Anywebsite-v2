# 使用官方 Go 镜像作为构建环境 (使用固定版本)
FROM golang:1.21-alpine3.18 AS builder

# 设置 Alpine 镜像源为阿里云镜像
RUN sed -i 's/dl-cdn.alpinelinux.org/mirrors.aliyun.com/g' /etc/apk/repositories

# 设置工作目录
WORKDIR /app

# 安装必要的系统依赖
RUN apk add --no-cache git ca-certificates tzdata

# 设置 Go 代理（使用国内代理加速）
ENV GOPROXY=https://goproxy.cn,direct
ENV GO111MODULE=on

# 复制 go 模块文件
COPY go.mod go.sum ./

# 下载依赖
RUN go mod download

# 复制源代码
COPY . .

# 构建应用
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main ./cmd/server

# 使用 alpine 作为最终镜像 (使用国内镜像加速)
FROM alpine:3.18

# 设置 Alpine 镜像源为阿里云镜像
RUN sed -i 's/dl-cdn.alpinelinux.org/mirrors.aliyun.com/g' /etc/apk/repositories

# 安装 ca-certificates 用于 HTTPS
RUN apk --no-cache add ca-certificates tzdata

WORKDIR /root/

# 从构建阶段复制二进制文件
COPY --from=builder /app/main .

# 复制配置文件和模板
COPY --from=builder /app/configs ./configs
COPY --from=builder /app/templates ./templates

# 创建必要的目录
RUN mkdir -p static uploads certs

# 暴露端口
EXPOSE 8080 8443

# 运行应用
CMD ["./main"]
