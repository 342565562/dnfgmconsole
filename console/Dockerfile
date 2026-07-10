# 按服务器目录构建：构建上下文为「含 dnf/ 的目录」，即与 dnf、Dockerfile 同级
# 目录示例：
#   .
#   ├── dnf
#   │   ├── biz
#   │   ├── cmds
#   │   ├── config
#   │   ├── go.mod
#   │   ├── go.sum
#   │   ├── main.go
#   │   ├── mods
#   │   └── ...
#   ├── Dockerfile
#   ├── run.sh
#   └── server.json
# 构建命令：docker build -t webconsole:lk70s2a1 .

FROM golang:1.19

WORKDIR /app

ENV GO111MODULE=on
ENV GOPROXY=https://goproxy.cn,direct

# 复制 go.mod、go.sum（来自 dnf/）以利用依赖缓存
COPY dnf/go.mod dnf/go.sum ./

RUN go mod download

# 复制 dnf 下全部源码到 /app（go.mod 已为 module dnf，import 均为 dnf/...）
COPY dnf/ .

RUN go build -o /app/dist/gmserver .

WORKDIR /app/dist

CMD ["/app/dist/gmserver", "-x"]
