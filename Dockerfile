# 多阶段构建 Dockerfile
FROM golang:1.21-alpine AS builder

# 设置工作目录
WORKDIR /build

# 安装依赖
RUN apk add --no-cache git ca-certificates

# 复制 go.mod 和 go.sum
COPY go.mod go.sum ./
RUN go mod download

# 复制源代码
COPY . .

# 编译
ARG VERSION=dev
ARG BUILD_TIME=unknown
RUN CGO_ENABLED=0 GOOS=linux go build \
    -ldflags "-s -w -X main.version=${VERSION} -X main.buildTime=${BUILD_TIME}" \
    -o uinetd ./cmd/uinetd

# 运行时镜像
FROM alpine:latest

# 安装 ca-certificates
RUN apk --no-cache add ca-certificates

# 创建非 root 用户
RUN addgroup -g 1000 uinetd && \
    adduser -D -u 1000 -G uinetd uinetd

# 复制二进制文件
COPY --from=builder /build/uinetd /usr/local/bin/uinetd

# 复制配置文件示例
COPY configs/uinetd.conf.example /etc/uinetd.conf

# 确保二进制可执行
RUN chmod +x /usr/local/bin/uinetd

# 如果需要监听 1024 以下端口，取消注释以下行（以 root 运行）
# USER root

# 否则使用非 root 用户（推荐）
USER uinetd

# 暴露常用端口（根据实际需求调整）
EXPOSE 8080 8443

# 健康检查（可选）
HEALTHCHECK --interval=30s --timeout=3s --start-period=5s --retries=3 \
  CMD pidof uinetd || exit 1

# 启动命令
ENTRYPOINT ["/usr/local/bin/uinetd"]
CMD ["-c", "/etc/uinetd.conf"]

