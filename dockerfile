# 你可以选择合适的标签，例如 "1.16-alpine"
FROM golang:1.21-alpine

# 设置工作目录
WORKDIR /app

# 将 go.mod 和 go.sum 文件复制到容器中
# 如果你的项目中没有这些文件，请忽略这一步
COPY go.mod ./
COPY go.sum ./
COPY .env ./

# 下载依赖（如果有的话）
RUN go mod download

# 将代码复制到容器中
COPY *.go ./
COPY .env ./

# 构建 Go 程序为二进制可执行文件
RUN go build -o /translate-api

# 执行时运行该程序
CMD [ "/translate-api" ]
