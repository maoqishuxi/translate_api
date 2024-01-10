# 你可以选择合适的标签，例如 "1.21-alpine"
FROM golang:1.21-alpine AS builder

# 设置工作目录
WORKDIR /app

# 将 go.mod 和 go.sum 文件复制到容器中
COPY go.mod go.sum ./

# 下载依赖（如果有的话）
RUN go mod download

# 将代码复制到容器中
COPY *.go ./

# 构建 Go 程序为二进制可执行文件
RUN go build -o /translate-api

# 使用多阶段构建
FROM alpine:latest  
WORKDIR /root/
COPY --from=builder /app/translate-api .
COPY .env .

# 执行时运行该程序
CMD [ "./translate-api" ]
