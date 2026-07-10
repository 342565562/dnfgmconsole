# 激活码相关改动与 Docker 构建说明

## 一、结论：Docker 构建失败与「激活码业务代码」无关

报错 `package console/biz/... is not in GOROOT` 是 **Go 模块解析** 问题，不是激活码逻辑导致的。  
激活码只改了业务代码和数据库访问方式，没有改 `go.mod`、模块名或目录结构。

---

## 二、激活码功能相关改动（不影响 Docker 能否 build）

### 1. 新增文件
- `biz/user/users/model/activation_code.go`（激活码模型）

### 2. 模型与表注册
- `biz/user/users/model/table.go`：User 增加 `IsActivated`，表改为 `game_db.RegTables(WebServer, &User{})`
- `biz/log/model/table.go`：改为 `game_db.RegTables(WebServer, ...)`
- `mods/ginx/table.go`、`mods/casbinx/table.go`：改为 `game_db.RegTables(WebServer, ...)`
- `activation_code.go`：`game_db.RegTables(WebServer, &ActivationCode{})`

### 3. 数据库从 db 改为 game_db
- 所有原来用 `db.DB` 的地方改为 `game_db.DBPools.Get(WebServer)` 或 `Get(gmModel.WebServer)`
- 涉及：user_init、user、auth、role、log、dash、casbinx、ginx 等

### 4. 为消除「model 重复声明」做的 import 修改（与激活码无直接关系）
- `biz/user/role/service/role.go`：`"console/biz/gm/model"` → `gmModel "console/biz/gm/model"`
- `biz/log/service/recharge_log.go`：拆成 `gmModel`、`logModel`
- `biz/user/users/service/user.go`、`user_init.go`：`gmModel` + 保留 `model`（users）
- `cmds/config.go`：删除未使用的 `errors`、`db` 导入

以上都是 **包内引用方式** 的调整，import 路径仍是 `console/biz/...`、`console/mods/...`，没有改模块名或目录结构，**不会** 导致「package is not in GOROOT」。

---

## 三、真正影响 Docker 构建的：Dockerfile 与 go.mod 的 module 名

### 1. Dockerfile 被改成了「按服务器 dnf 目录」构建

- **之前（可能）**：在 **项目根目录（console）** 里 build，例如：
  - 构建上下文 = `console/`
  - `COPY go.mod go.sum ./` 再 `COPY . .`
  - 此时 `/app` 下就是 `go.mod`（module console）+ 所有源码，模块解析正常。

- **现在**：按你提供的服务器目录来写：
  - 构建上下文 = **含 `dnf/` 的那一层**（与 dnf、Dockerfile 同级）
  - `COPY dnf/go.mod dnf/go.sum ./`，再 `COPY dnf/ .`
  - 若服务器上的 **`dnf/go.mod` 里是 `module dnf`**（或其它名字），而代码里全是 `import "console/..."`，Go 就会去找「module console」，找不到就报 **not in GOROOT**。

也就是说：**Docker 失败是因为「用 dnf 目录构建 + dnf 里 go.mod 的 module 名不是 console」**，而不是因为加了激活码。

### 2. 已做的应对：在 Dockerfile 里强制 module 名

在 Dockerfile 里加了：

```dockerfile
RUN sed -i 's/^module .*$/module console/' go.mod
```

这样无论服务器上 `dnf/go.mod` 写的是 `module dnf` 还是别的，镜像里都会变成 `module console`，build 就能通过。

---

## 四、修改项汇总

| 类型 | 修改内容 | 是否影响 Docker 构建 |
|------|----------|----------------------|
| 激活码功能 | 新模型、登录校验、IsActivated、game_db 等 | 否 |
| 编译 bug 修复 | gmModel / logModel 别名、去掉未使用 import | 否 |
| Dockerfile | 改为从 `dnf/` 拷贝、增加 `sed` 修正 module 名 | 是（解决 not in GOROOT） |
| db_mysql.go | COUNT(*) 查库存在、Scan 到 int | 否（只影响运行时） |

---

## 五、若希望「和加激活码之前一样」能 build

如果之前是在 **console 目录内** 执行 `docker build`（上下文是 console 自己），可以保留一个「不依赖 dnf」的 Dockerfile，例如：

- 在 **console 目录** 执行：`docker build -t webconsole:lk70s2a1 .`
- Dockerfile 使用：`COPY go.mod go.sum ./` 和 `COPY . .`（不用 `dnf/`）

这样就和「加激活码前」的用法一致，不依赖服务器上的 `dnf` 目录和 module 名。  
当前仓库里的 Dockerfile 是按「服务器上有 dnf 目录」的用法写的；若要同时支持「在 console 下直接 build」，可以再加一个 `Dockerfile.local` 或用构建参数区分两种方式。

---

**总结**：  
- **激活码相关改动**（含 model 别名、config 去未使用 import）**没有**改变模块名或目录结构，**不会**导致 Docker build 失败。  
- 不能正常 docker build 是因为 **Dockerfile 改为按 dnf 目录构建**，且 **dnf 内 go.mod 的 module 名不是 console**；  
- 通过 Dockerfile 里的 **`sed` 把 module 统一成 console** 已针对这一问题做了修复。
