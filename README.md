# DNF GM Web 控制台 (dnfgmconsole)

DNF 游戏 GM 管理后台，包含 Go 后端(`console/`)与 Vue3 前端(`web/`)。
功能：账号/角色管理、点券/D点充值、PK 段位/QP 修改、邮件发送(普通/时装/宠物)、
一键恢复、充值/登录日志、概览统计、后台激活码、物品代码库(gold)导入等。

---

## 一、项目结构

```
dnfgmconsole/
├── Dockerfile                   # 一体化镜像(容器内同时编译前端+后端)
├── .dockerignore
├── DNF.ico                      # 导入 GUI 图标
├── console/                     # 【后端】Go 源码(module: dnf)
│   ├── main.go                  # 入口
│   ├── config/
│   │   ├── server.json          # 配置文件
│   │   └── server.json.example  # 逐行注释的配置说明
│   ├── biz/                     # 业务模块(gm/user/dash/log...)
│   ├── mods/game_db/            # 多游戏库连接与迁移
│   ├── dist/                    # (手动部署时的产物目录)
│   └── tools/                   # 独立工具
│       ├── gold_import/         # 物品导入 CLI
│       ├── gold_import_gui/     # 物品导入 图形工具(仅 Windows)
│       ├── goldimport/          # 导入公共逻辑(CLI/GUI 共用)
│       └── generate_activation_codes.go  # 激活码生成工具
└── web/                         # 【前端】Vue3 源码
```

后端源码放在 `console/`，前端源码放在 `web/`，`Dockerfile` 在项目根目录。

---

## 二、Docker 一键构建(推荐)

根目录的 `Dockerfile` 采用**三阶段构建**，在容器内**同时编译前端和后端**：
Node 阶段编前端 → Go 阶段编后端 → Alpine 运行镜像(约 30–40MB)。

> 好处：**本机无需安装 Node / Go，也不再依赖 Windows**。在 **CentOS 7** 或任意装了 Docker 的机器上，直接在项目根目录一条命令即可构建。
> 编译好的前端会被自动放进镜像的 `/app/web/static`，无需手工拷贝。

```bash
# 在项目根目录(与 Dockerfile 同级)执行
# 1) 构建
docker build -t dnfgmconsole:latest .

# 2) 运行(挂载 server.json 便于改配置不重建；映射 8088)
docker run -d --name gmconsole \
  -p 8088:8088 \
  -v $(pwd)/console/config/server.json:/app/config/server.json \
  --restart unless-stopped \
  dnfgmconsole:latest

# 日志 / 停止 / 删除
docker logs -f gmconsole
docker stop gmconsole && docker rm gmconsole
```

- 构建上下文 = 项目根目录；`Dockerfile` 自动到 `console/`(后端)、`web/`(前端)取源码。
- 镜像内：`/app/gmserver`、`/app/config/`、`/app/web/static/`(工作目录 `/app`，端口 8088)。
- 容器需能连到 `server.json` 里配置的 MySQL(远程库注意网络连通)。
- 访问：`http://<服务器IP>:8088/`，默认账号 `admin / 123`。

---

## 三、手动编译(不使用 Docker)

### 3.1 后端(Go 1.19+，CentOS 7 可用)
```bash
cd console
export GO111MODULE=on
export GOPROXY=https://goproxy.cn,direct
go build -o dist/gmserver .
```
运行(需在 `dist` 目录内，静态页面依赖工作目录)：
```bash
cd console/dist
./gmserver -p config/server.json -x
```

### 3.2 前端(Node 16/18+，**CentOS 7 同样可以编译**)
前端是标准 Vue3/npm 工程，**不限 Windows**，任何装了 Node 的系统(含 CentOS 7)都能编译：
```bash
cd web
npm install          # 或 npm ci
npm run build        # 产物在 web/dist
```
把 `web/dist/*` 部署到后端可托管的静态目录 `console/dist/web/static/`。
> 后端把 `<可执行文件所在目录>/web/static` 挂到 `/static`，根路径 `/` 重定向到登录页。
> 若嫌本机装 Node 麻烦，直接用第二节的 Docker 构建即可(容器内自带 Node)。

### 常用运行参数
| 参数 | 说明 |
|------|------|
| `-p <path>` | 指定配置文件 |
| `-x` | 单机前台模式(日志输出到控制台) |
| `-i` | 初始化数据库(建表 + 默认管理员 `admin/123`) |
| `-d` | 打印默认配置 |
| `-k start\|stop\|status` | 服务方式管理 |

---

## 四、配置文件 server.json

```jsonc
{
  "site": {
    "title": "金华DNF",                  // 浏览器标题
    "login_name": "金华DNF"              // 登录页显示名称
  },
  "auth": {
    "activation_code_enable": 0          // 后台激活码：1=开启，0=关闭(默认)
  },
  "game_db": {
    "enable": true,
    "mysql": [
      { "key": "taiwan_cain_2nd", "db": "taiwan_cain_2nd",
        "user": "game", "password": "***", "host": "127.0.0.1",
        "port": 3306, "charset": "utf8", "timeout": 5 }
      // ... 其余库
    ]
  },
  "item_classify": {                     // 物品分类规则(供导入工具使用，可增删改)
    "pet_types": ["creature", "creature expitem", "feed"],
    "fashion_types": ["disguise", "disguise random", "dye"],
    "fashion_rarity_keyword": "时装"
  }
}
```
> 标准 JSON 不支持注释；完整**逐行注释说明**见 [`console/config/server.json.example`](console/config/server.json.example)。

- **站点标题/登录名** `site.*`：前端启动读取(免鉴权接口 `GET /api/site-config`)。
- **激活码开关** `auth.activation_code_enable`(详见第五节)。
- **物品分类规则** `item_classify`(详见第七节)。

---

## 五、后台激活码功能

用于**限制后台账号必须凭激活码才能首次登录**，防止账号被随意使用。

- **开关**：`server.json` 的 `auth.activation_code_enable`
  - `0`(默认)：关闭。所有账号正常登录，无需激活码。
  - `1`：开启。**未激活**的账号登录时，除用户名/密码外还需填写一个有效激活码；
    校验通过后该账号被标记为已激活，并把该激活码标记为已使用(一码一号)。
- **登录逻辑**：账号已激活 → 直接登录；未激活 →
  - 未填激活码：提示“该账号未激活，请输入激活码”(前端自动弹出激活码输入框)；
  - 激活码无效/已用：提示“激活码无效或已被使用”；
  - 激活码有效：激活成功并登录。
- **初始化建的默认管理员 `admin` 默认视为已激活**，不受影响。

### 生成激活码(工具)
激活码存放在数据库 `webserver.activation_codes` 表，用工具批量生成：
```bash
cd console
go build -o gen_codes ./tools/generate_activation_codes.go
# 读取 server.json 的数据库配置，运行后按提示输入生成数量
./gen_codes -config config/server.json
```
> 也可用 `tools/activate_user.sql` 直接在数据库层手动激活/处理。

---

## 六、邮件发送与自动分类

发送邮件页可选邮件类型：**普通 / 时装 / 宠物**。

- 前端选“普通”时，后端按物品代码在 `gold` 表查分类：
  - 命中**时装** → 自动按**时装邮件**发送；
  - 命中**宠物** → 自动按**宠物邮件**发送；
  - 查不到 / 普通 → 按普通邮件发送。
- 前端**明确**选了时装/宠物 → 后端按所选类型发送，不再自动判断。

> 分类来自 `gold` 表的 `category` 字段，由导入工具在导入时写入。

---

## 七、物品代码库(gold)导入工具

`gold` 表：`code`(物品代码/it_id)、`name`(名称)、`rarity`(稀有度)、`category`(0普通/1时装/2宠物)。
名称/稀有度列为 `utf8mb4`，即使库为 latin1 也不乱码。程序启动时**表不存在才创建，已存在则跳过**。

### 分类规则(默认，可在 server.json 的 item_classify 修改)
- **宠物**：道具表原始种类 ∈ `creature / creature expitem / feed`
- **时装**：道具表原始种类 ∈ `disguise / disguise random / dye`，**或**装备表稀有度含关键字“时装”
- **普通**：其余全部(徽章、各类礼包、装备、未知/新类型自动归普通)

### 数据过滤
自动跳过**名称为空、纯问号(？/?)、乱码**(含替换符或西里尔字母)的行。

### CSV 格式(GBK 编码，工具自动转码)
- 装备 CSV 表头：`名称,ID,大类,小类,子类,等级,稀有度`
- 道具 CSV 表头：`名称,ID,原始种类,种类,使用等级,稀有度`

### 7.1 命令行工具(CLI，跨平台)
```bash
cd console
go build -o gold_import ./tools/gold_import/
# 复用 server.json 的连接与分类规则
./gold_import -config config/server.json \
  -equip "装备.csv" -item "道具.csv" -mode truncate -yes
```
`-mode truncate`(清空重导) / `append`(追加)；`-yes` 跳过确认；也可用 `-host/-user/-pass/-db` 手动指定连接。

### 7.2 图形工具(GUI，仅 Windows)
图标为项目根目录 `DNF.ico`。**注意：GUI 依赖 `lxn/walk`，只能在 Windows 编译**(前端和后端不受此限制)。
```bash
cd console
go build -ldflags="-H windowsgui" -o gold_import_gui.exe ./tools/gold_import_gui/
```
运行后：可“读取配置”从 server.json 一键填充 → 填连接 → 选装备/道具 CSV → 选模式 → 开始导入(下方显示日志)。

---

## 八、Git

```bash
git add .
git commit -m "your message"
git push
```
