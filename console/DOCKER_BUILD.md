# Docker 构建说明

## 错误原因

报错 `package console/biz/... is not in GOROOT` 表示 Go 没有把 `console` 当作当前模块，常见原因：

1. **COPY 后目录结构不对**：`/app` 下没有与 `go.mod` 同级的 `biz`、`cmds`、`mods`、`main.go`。
2. **dnf 是“整仓”拷贝**：若 `dnf` 是整仓（例如含 `dnf/console/`），则实际 Go 工程在 `dnf/console/`，需要按下面“方式二”拷贝。

## 方式一：在仓库内从 console 目录构建（推荐）

在 **本机** 进入 `console` 目录再构建，上下文就是当前目录，无需 `dnf`：

```bash
cd /path/to/gmwebconsole/console
docker build -t webconsole:lk70s2a1 .
```

此时 Dockerfile 里用 `COPY go.mod go.sum ./` 和 `COPY . .` 即可，无需改路径。

## 方式二：在服务器用 dnf 目录构建

若你在服务器上把「整个 console 项目」拷贝成了 `dnf` 目录，且 **dnf 内部是“扁平”的**（和 console 一样），应满足：

- `dnf/go.mod`、`dnf/main.go`、`dnf/biz/`、`dnf/cmds/`、`dnf/mods/`、`dnf/config/` 等都在 **dnf 根下**，没有再多一层 `dnf/console/`。
- `dnf/go.mod` 第一行必须是：`module console`（不能是 `module dnf` 等）。

Dockerfile 保持：

```dockerfile
COPY dnf/go.mod dnf/go.sum ./
COPY dnf/ .
RUN go build -o /app/dist/gmserver  .
```

若 **dnf 里还有一层 console**（即实际是 `dnf/console/go.mod`、`dnf/console/biz/`…），则需改成从 `dnf/console/` 拷贝，让 `/app` 下直接是 go.mod + main.go + biz + cmds + mods：

```dockerfile
# 从 dnf/console 拷贝，保证 /app 下是扁平结构
COPY dnf/console/go.mod dnf/console/go.sum ./
COPY dnf/console/ .
RUN go build -o /app/dist/gmserver  .
```

## 检查 dnf 目录结构

在服务器上执行：

```bash
ls -la dnf/
# 应看到：go.mod  main.go  biz  cmds  mods  config 等（扁平）

# 若看到的是：
# dnf/console/...
# 则说明是嵌套结构，Dockerfile 要用上面的 dnf/console/ 拷贝方式

head -1 dnf/go.mod
# 必须输出：module console
```

## 本仓库自带的 Dockerfile

仓库根下的 `console/Dockerfile` 是按「在 console 目录内构建」写的，直接：

```bash
cd console
docker build -t webconsole:lk70s2a1 .
```

即可。
