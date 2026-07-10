# 激活码生成工具使用说明

## 功能说明

此工具用于批量生成16位激活码，激活码只包含数字（0-9）和英文字母（大小写a-z, A-Z）。

## 使用方法

### 1. 编译脚本

在 `console` 目录下执行：

```bash
go build -o generate_activation_codes.exe tools/generate_activation_codes.go
```

或者在 `console/tools` 目录下执行：

```bash
cd tools
go build -o generate_activation_codes.exe generate_activation_codes.go
```

### 2. 运行脚本

```bash
# Windows
generate_activation_codes.exe

# Linux/Mac
./generate_activation_codes
```

### 3. 输入生成数量

运行后，程序会：
1. 默认读取“当前目录”的 `server.json`（也可用 `-config` 指定路径）
2. 连接 webserver 数据库
3. 提示输入要生成的激活码数量
4. 自动生成并保存到数据库

## 特性

- ✅ 自动读取 `server.json` 配置文件中的数据库信息
- ✅ 自动连接 webserver 数据库
- ✅ 生成16位随机激活码（数字+大小写字母）
- ✅ 自动检查激活码唯一性，避免重复
- ✅ 批量插入数据库，提高效率
- ✅ 显示生成的激活码列表

## 激活码格式

- **长度**: 16位
- **字符集**: 数字 0-9，小写字母 a-z，大写字母 A-Z
- **示例**: `a3B7c9D2e4F6g8H1`

## 注意事项

1. 确保当前目录存在 `server.json` 且配置正确（或使用 `-config` 指定路径）
2. 确保 webserver 数据库已创建且可连接
3. 工具会在连接成功后自动创建 `activation_codes` 表（需要账号具备建表权限）
4. 如果生成数量较大（>10000），程序会提示确认

## 示例输出

```
==========================================
    激活码生成工具
==========================================
配置文件路径: D:\代码项目\2.dev\gmwebconsole\console\config\server.json
数据库配置: game@123.56.165.124:3306/webserver

正在连接数据库...
数据库连接成功！

请输入要生成的激活码数量: 10

开始生成 10 个激活码...
已生成 10/10 个激活码...

正在将激活码保存到数据库...
成功保存 10 个激活码到数据库

生成的激活码列表：
==================================================
1. a3B7c9D2e4F6g8H1
2. x9Y2z5W8v1U4t7S0
3. m6N3p0Q5r2T8s1V4
4. k7J4h9G2f5E8d1C6
5. b8A1c4D7e0F3g6H9
6. i2L5m8N1o4P7q0R3
7. s6T9u2V5w8X1y4Z7
8. a0B3c6D9e2F5g8H1
9. j4K7l0M3n6O9p2Q5
10. r8S1t4U7v0W3x6Y9
==================================================

==========================================
激活码生成完成！
==========================================
```
