# 构建阶段
FROM golang:1.21-alpine AS builder

WORKDIR /app
COPY . .

# 安装时区数据
RUN apk add --no-cache tzdata
# 设置环境变量（可选但推荐）
ENV TZ=Asia/Shanghai

# 下载依赖（如果使用go modules）
RUN go mod download

# 构建可执行文件
RUN CGO_ENABLED=0 GOOS=linux go build -o /ip-api

# 运行时阶段
FROM alpine:3.19

WORKDIR /app
COPY --from=builder /ip-api /app/ip-api

# 暴露端口
EXPOSE 18125

# 启动命令
CMD ["/app/ip-api"]
