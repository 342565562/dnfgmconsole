# 本地运行指南

本文档说明如何在本地运行后端和前端服务。

## 环境要求

- **Go**: 1.19 或更高版本
- **Node.js**: 建议使用 Node.js 14 或更高版本
- **npm**: 建议使用 npm 6 或更高版本

## 后端运行步骤

### 1. 进入后端目录

```bash
cd console
```

### 2. 安装依赖

```bash
go mod download
```

### 3. 初始化数据库（首次运行需要）

```bash
go run main.go -i
```

### 4. 启动后端服务

在前台模式运行（推荐开发时使用）：

```bash
go run main.go -x
```

或者使用指定的配置文件：

```bash
go run main.go -x -p config/server.json
```

后端服务启动后，会在控制台输出监听的端口地址，请记录此端口用于前端配置。

**提示**: 如果后端服务没有指定端口，通常会使用默认端口。请查看后端启动日志中的端口信息。

## 前端运行步骤

### 1. 进入前端目录

```bash
cd web
```

### 2. 安装依赖

```bash
npm install
```

### 3. 配置后端 API 地址

在 `web` 目录下创建 `.env.development` 文件，配置后端 API 地址：

```bash
cp .env.development.example .env.development
```

然后编辑 `.env.development` 文件，根据后端实际启动的端口修改 API 地址：

```env
VUE_APP_BASE_API=http://localhost:8081/api
```

**注意**: 
- 请根据后端实际启动的端口调整上述地址
- 如果后端使用了其他端口（例如 3000、8080 等），请相应修改
- 确保后端已经启动并可以访问后再启动前端

### 4. 启动前端开发服务器

```bash
npm run serve
```

前端开发服务器会在 `8080` 端口启动，浏览器访问 `http://localhost:8080`。

## 同时运行后端和前端

推荐使用两个终端窗口分别运行：

**终端 1 - 运行后端：**
```bash
cd console
go run main.go -x
```

**终端 2 - 运行前端：**
```bash
cd web
npm run serve
```

## 常见问题

### 1. 后端启动失败

- 检查数据库配置是否正确（`console/config/server.json`）
- 检查端口是否被占用
- 确认已运行数据库初始化命令 `go run main.go -i`

### 2. 前端无法连接到后端

- 检查 `.env.development` 文件中的 `VUE_APP_BASE_API` 配置是否正确
- 确认后端服务已经启动
- 检查后端端口是否与前端配置的端口一致
- 检查浏览器控制台的网络请求错误

### 3. 前端依赖安装失败

```bash
# 清除缓存后重新安装
rm -rf node_modules package-lock.json
npm install
```

### 4. Go 依赖下载失败

如果在中国大陆，可能需要配置 Go 代理：

```bash
go env -w GOPROXY=https://goproxy.cn,direct
```

## 其他有用命令

### 后端

- 查看默认配置：`go run main.go -d`
- 生成 RSA 密钥对：`go run main.go -pem`
- 查看帮助：`go run main.go -h`

### 前端

- 构建生产版本：`npm run build`
- 代码检查：`npm run lint`

