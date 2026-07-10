# CentOS 7 编译说明

## 报错原因

```
package console/biz/client/service is not in GOROOT (/usr/lib/golang/src/console/biz/client/service)
```

表示 **Go 没有按模块模式编译**，把 `console/...` 当成标准库路径去 GOROOT 里找，所以报错。

常见原因：
1. **未开启 Go modules**：`GO111MODULE` 为 off 或 auto 且当前目录未被识别为模块根
2. **在错误目录执行**：在「没有 go.mod 的目录」执行了 `go build`

---

## 正确做法

### 1. 进入「含 go.mod 的目录」

项目根是 **console**（和 go.mod、main.go、biz、cmds、mods 同级）：

```bash
cd /path/to/console
# 或若你拷贝成了 dnf：cd /path/to/dnf
```

确认当前目录有 go.mod：

```bash
pwd
ls go.mod main.go biz cmds mods
```

### 2. 开启 Go modules 再编译

```bash
export GO111MODULE=on
export GOPROXY=https://goproxy.cn,direct
go build -o dist/gmserver .
```

或一行：

```bash
GO111MODULE=on GOPROXY=https://goproxy.cn,direct go build -o dist/gmserver .
```

### 3. 使用脚本（推荐）

在 **console 目录** 下执行：

```bash
chmod +x build.sh
./build.sh
```

脚本会自动切到脚本所在目录并设置 `GO111MODULE=on` 再执行 `go build`。

---

## 检查清单

| 检查项 | 说明 |
|--------|------|
| 当前目录 | 必须在 **含 go.mod 的目录**（console 或 dnf） |
| go.mod 第一行 | 必须是 `module console`，不能是 `module dnf` 等 |
| GO111MODULE | 编译前执行 `export GO111MODULE=on` |
| Go 版本 | 建议 1.19+，`go version` 查看 |

---

## 错误示例

```bash
# 错误：在项目父目录执行
cd /path/to/gmwebconsole
go build -o dist/gmserver .          # 当前目录没有 go.mod，失败

# 错误：在 console 目录但未开 modules（CentOS 7 默认可能 off）
cd /path/to/console
go build -o dist/gmserver .          # 可能报 not in GOROOT
```

## 正确示例

```bash
cd /path/to/console
export GO111MODULE=on
go build -o dist/gmserver .
```
