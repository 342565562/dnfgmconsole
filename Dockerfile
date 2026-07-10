# 使用较新的 Go 版本基础镜像
FROM golang:1.19

# 设置工作目录
WORKDIR /app

# 设置 Go 环境变量
ENV GO111MODULE=on
ENV GOPROXY=https://goproxy.cn,direct

# 复制 go.mod 和 go.sum（如果有）以实现依赖缓存
COPY dnf/go.mod dnf/go.sum ./

# 下载依赖
RUN go mod download

# 复制源代码目录到工作目录
COPY dnf/ .

# 检查文件内容
RUN ls -la

# 编译 Go 程序
RUN go build -o /app/dist/gmserver  .

# 设置工作目录
WORKDIR /app/dist/

# 运行程序
CMD ["/app/dist/gmserver", "-x"]
