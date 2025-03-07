# 构建阶段
FROM golang:1.21-alpine AS builder

WORKDIR /app
COPY . .

# 下载依赖（如果使用go modules）
RUN go mod download

# 构建可执行文件
RUN CGO_ENABLED=0 GOOS=linux go build -o /ipapi

# 运行时阶段
FROM alpine:3.19

WORKDIR /app
COPY --from=builder /ipapi /app/ipapi

# 暴露端口
EXPOSE 8080

# 启动命令
CMD ["/app/ipapi"]
