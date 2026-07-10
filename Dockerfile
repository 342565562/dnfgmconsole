# ============================================================================
# 一体化多阶段构建：在容器内同时编译【前端】和【后端】，无需本机装 Node/Go，
# 也不再依赖 Windows —— 在 CentOS 7 / 任意有 Docker 的机器上都能构建。
#
# 目录约定(本文件放在项目根目录，构建上下文 = 项目根目录)：
#   ./console/   后端 Go 源码(module dnf)
#   ./web/       前端 Vue3 源码
#   ./Dockerfile 本文件
#
# 构建：  docker build -t dnfgmconsole:latest .
# 运行：  见 README.md
# ============================================================================

# ---- 阶段 1：编译前端(Node) ----
FROM node:18 AS webbuilder
WORKDIR /web
# 先装依赖(利用缓存)
COPY web/package.json web/package-lock.json ./
RUN npm ci --registry=https://registry.npmmirror.com || npm install --registry=https://registry.npmmirror.com
# 编译前端 -> /web/dist
COPY web/ .
RUN npm run build

# ---- 阶段 2：编译后端(Go) ----
FROM golang:1.19-alpine AS gobuilder
WORKDIR /src
ENV GO111MODULE=on \
    GOPROXY=https://goproxy.cn,direct \
    CGO_ENABLED=0
COPY console/go.mod console/go.sum ./
RUN go mod download
COPY console/ .
RUN go build -ldflags="-s -w" -o /out/gmserver .

# ---- 阶段 3：运行镜像(Alpine，约 30–40MB) ----
FROM alpine:3.19
RUN apk add --no-cache tzdata ca-certificates \
    && ln -sf /usr/share/zoneinfo/Asia/Shanghai /etc/localtime \
    && echo "Asia/Shanghai" > /etc/timezone
WORKDIR /app
# 后端二进制
COPY --from=gobuilder /out/gmserver ./gmserver
# 配置(可运行时用 -v 覆盖 config/server.json)
COPY console/config ./config
# 前端编译产物 -> /app/web/static (后端把工作目录下 web/static 挂到 /static)
COPY --from=webbuilder /web/dist ./web/static
EXPOSE 8088
CMD ["./gmserver", "-p", "config/server.json", "-x"]
