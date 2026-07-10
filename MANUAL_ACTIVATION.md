# 手工激活已注册用户指南

## 📋 概述

对于已经注册但未激活的用户，可以通过以下方式手工激活。

---

## 🎯 方法一：SQL 直接激活（推荐）

### 步骤详解

#### 1. 查看未激活的用户

```sql
USE webserver;

SELECT id, username, is_activated, time AS registered_at
FROM user
WHERE is_activated = 0
ORDER BY time DESC;
```

#### 2. 查看可用的激活码

```sql
SELECT id, code, is_used, created_at
FROM activation_code
WHERE is_used = 0
ORDER BY created_at DESC
LIMIT 20;
```

#### 3. 激活指定用户

假设要激活用户 `admin123`（用户ID为 `5`），使用激活码 `Xyz123AbcDef456G`：

```sql
-- 激活用户
UPDATE user
SET is_activated = 1
WHERE username = 'admin123';

-- 绑定激活码
UPDATE activation_code
SET is_used = 1,
    user_id = 5,
    used_at = NOW()
WHERE code = 'Xyz123AbcDef456G';
```

#### 4. 验证激活结果

```sql
SELECT
    u.id,
    u.username,
    u.is_activated,
    a.code AS activation_code,
    a.used_at AS activated_at
FROM user u
LEFT JOIN activation_code a ON a.user_id = u.id
WHERE u.username = 'admin123';
```

**预期结果**：
```
+----+-----------+--------------+----------------------+---------------------+
| id | username  | is_activated | activation_code      | activated_at        |
+----+-----------+--------------+----------------------+---------------------+
| 5  | admin123  | 1            | Xyz123AbcDef456G     | 2026-02-06 12:30:00 |
+----+-----------+--------------+----------------------+---------------------+
```

---

## 🛠️ 方法二：使用自动化脚本

### Windows 批处理脚本

```bash
cd console\tools
手工激活用户.bat
```

**脚本会自动**：
1. 连接数据库
2. 显示未激活用户列表
3. 显示可用激活码列表
4. 提示输入用户名和激活码
5. 执行激活并验证

### SQL 脚本文件

```bash
# 1. 编辑 activate_user.sql，设置用户名和激活码
# 2. 执行脚本
mysql -h 115.191.0.22 -u game -p webserver < console/tools/activate_user.sql
```

---

## 📊 方法三：批量激活多个用户

### 单条 SQL 批量激活

```sql
USE webserver;

-- 激活多个指定用户
UPDATE user
SET is_activated = 1
WHERE username IN ('user1', 'user2', 'user3');

-- 为特定用户绑定激活码（需要逐个执行）
UPDATE activation_code
SET is_used = 1,
    user_id = (SELECT id FROM user WHERE username = 'user1'),
    used_at = NOW()
WHERE code = 'Code1ForUser1';

UPDATE activation_code
SET is_used = 1,
    user_id = (SELECT id FROM user WHERE username = 'user2'),
    used_at = NOW()
WHERE code = 'Code2ForUser2';

UPDATE activation_code
SET is_used = 1,
    user_id = (SELECT id FROM user WHERE username = 'user3'),
    used_at = NOW()
WHERE code = 'Code3ForUser3';
```

### 批量激活所有用户（谨慎使用）

```sql
USE webserver;

-- ⚠️ 警告：这会激活所有未激活的用户，不绑定激活码
UPDATE user SET is_activated = 1 WHERE is_activated = 0;

-- 验证
SELECT COUNT(*) AS total_users,
       SUM(is_activated) AS activated_users,
       COUNT(*) - SUM(is_activated) AS unactivated_users
FROM user;
```

---

## 🔧 方法四：使用存储过程（高级）

### 创建存储过程

```sql
USE webserver;

DELIMITER //

CREATE PROCEDURE activate_user(
    IN p_username VARCHAR(64),
    IN p_activation_code VARCHAR(64)
)
BEGIN
    DECLARE v_user_id BIGINT;
    DECLARE v_code_exists INT;
    DECLARE v_code_used BOOLEAN;

    -- 开始事务
    START TRANSACTION;

    -- 获取用户ID
    SELECT id INTO v_user_id
    FROM user
    WHERE username = p_username;

    IF v_user_id IS NULL THEN
        ROLLBACK;
        SIGNAL SQLSTATE '45000' SET MESSAGE_TEXT = '用户不存在';
    END IF;

    -- 检查激活码是否存在和是否已使用
    SELECT COUNT(*), IFNULL(MAX(is_used), 1)
    INTO v_code_exists, v_code_used
    FROM activation_code
    WHERE code = p_activation_code;

    IF v_code_exists = 0 THEN
        ROLLBACK;
        SIGNAL SQLSTATE '45000' SET MESSAGE_TEXT = '激活码不存在';
    END IF;

    IF v_code_used = 1 THEN
        ROLLBACK;
        SIGNAL SQLSTATE '45000' SET MESSAGE_TEXT = '激活码已被使用';
    END IF;

    -- 激活用户
    UPDATE user
    SET is_activated = 1
    WHERE id = v_user_id;

    -- 绑定激活码
    UPDATE activation_code
    SET is_used = 1,
        user_id = v_user_id,
        used_at = NOW()
    WHERE code = p_activation_code;

    -- 提交事务
    COMMIT;

    -- 返回结果
    SELECT
        v_user_id AS user_id,
        p_username AS username,
        p_activation_code AS activation_code,
        '激活成功' AS status;
END //

DELIMITER ;
```

### 使用存储过程

```sql
-- 激活用户
CALL activate_user('admin123', 'Xyz123AbcDef456G');

-- 查看结果
-- 输出: user_id | username | activation_code      | status
--       5       | admin123 | Xyz123AbcDef456G     | 激活成功
```

### 删除存储过程

```sql
DROP PROCEDURE IF EXISTS activate_user;
```

---

## 📋 快速参考命令

### 常用查询

```sql
USE webserver;

-- 1. 查看所有未激活用户
SELECT id, username, time AS registered_at
FROM user
WHERE is_activated = 0;

-- 2. 查看所有可用激活码
SELECT code
FROM activation_code
WHERE is_used = 0
LIMIT 20;

-- 3. 查看用户激活状态
SELECT
    u.id,
    u.username,
    u.is_activated,
    a.code AS activation_code,
    a.used_at AS activated_at
FROM user u
LEFT JOIN activation_code a ON a.user_id = u.id
WHERE u.username = 'your_username';

-- 4. 查看激活码使用情况
SELECT
    COUNT(*) AS total,
    SUM(CASE WHEN is_used = 1 THEN 1 ELSE 0 END) AS used,
    SUM(CASE WHEN is_used = 0 THEN 1 ELSE 0 END) AS available
FROM activation_code;

-- 5. 查看最近激活的用户
SELECT
    u.username,
    a.code AS activation_code,
    a.used_at AS activated_at
FROM activation_code a
JOIN user u ON a.user_id = u.id
WHERE a.is_used = 1
ORDER BY a.used_at DESC
LIMIT 10;
```

### 快速激活模板

```sql
USE webserver;

-- 设置变量（替换为实际值）
SET @username = 'your_username';
SET @code = 'YourActivationCode';

-- 执行激活
UPDATE user SET is_activated = 1 WHERE username = @username;
UPDATE activation_code
SET is_used = 1,
    user_id = (SELECT id FROM user WHERE username = @username),
    used_at = NOW()
WHERE code = @code;

-- 验证结果
SELECT u.username, u.is_activated, a.code, a.used_at
FROM user u
LEFT JOIN activation_code a ON a.user_id = u.id
WHERE u.username = @username;
```

---

## ⚠️ 注意事项

### 1. 激活前检查

```sql
-- 检查用户是否存在
SELECT * FROM user WHERE username = 'your_username';

-- 检查用户是否已激活
SELECT is_activated FROM user WHERE username = 'your_username';
-- 如果返回 1，说明已激活，无需重复激活

-- 检查激活码是否有效
SELECT * FROM activation_code WHERE code = 'YourCode';
-- 确保 is_used = 0（未使用）
```

### 2. 数据一致性

激活用户需要同时更新两个表：
- ✅ `user` 表：`is_activated = 1`
- ✅ `activation_code` 表：`is_used = 1`, `user_id = 用户ID`, `used_at = NOW()`

**错误示例**（不完整的激活）：
```sql
-- ❌ 只激活用户，不绑定激活码
UPDATE user SET is_activated = 1 WHERE username = 'admin123';
-- 这会导致激活记录不完整
```

**正确示例**（完整的激活）：
```sql
-- ✅ 激活用户并绑定激活码
UPDATE user SET is_activated = 1 WHERE username = 'admin123';
UPDATE activation_code
SET is_used = 1, user_id = (SELECT id FROM user WHERE username = 'admin123'), used_at = NOW()
WHERE code = 'Xyz123AbcDef456G';
```

### 3. 事务安全

建议在事务中执行激活操作：

```sql
START TRANSACTION;

UPDATE user SET is_activated = 1 WHERE username = 'admin123';
UPDATE activation_code
SET is_used = 1, user_id = (SELECT id FROM user WHERE username = 'admin123'), used_at = NOW()
WHERE code = 'Xyz123AbcDef456G';

-- 检查是否成功，如果有问题则回滚
-- ROLLBACK;

-- 确认无误后提交
COMMIT;
```

---

## 🔍 故障排查

### Q1: 激活后用户仍然提示需要激活码？

**检查**：
```sql
SELECT is_activated FROM user WHERE username = 'your_username';
```

**解决**：
```sql
UPDATE user SET is_activated = 1 WHERE username = 'your_username';
```

### Q2: 提示"激活码已被使用"？

**检查**：
```sql
SELECT is_used, user_id FROM activation_code WHERE code = 'YourCode';
```

**如果激活码确实未使用，但状态错误**：
```sql
-- 重置激活码状态（谨慎使用）
UPDATE activation_code
SET is_used = 0, user_id = NULL, used_at = NULL
WHERE code = 'YourCode';
```

### Q3: 如何解绑已使用的激活码？

```sql
-- 1. 取消用户激活
UPDATE user SET is_activated = 0 WHERE id = 5;

-- 2. 重置激活码
UPDATE activation_code
SET is_used = 0, user_id = NULL, used_at = NULL
WHERE user_id = 5;
```

---

## 📝 工具文件

项目提供了以下工具文件协助激活：

| 文件 | 说明 | 使用方法 |
|------|------|---------|
| [console/tools/手工激活用户.bat](console/tools/手工激活用户.bat) | Windows 交互式激活脚本 | 双击运行 |
| [console/tools/activate_user.sql](console/tools/activate_user.sql) | SQL 激活脚本模板 | 编辑后用 mysql 执行 |

---

## 📞 技术支持

如遇问题，请提供：
1. 用户名和用户ID
2. 激活码
3. 数据库查询结果：
   ```sql
   SELECT * FROM user WHERE username = 'xxx';
   SELECT * FROM activation_code WHERE code = 'xxx';
   ```

---

**文档版本**：v1.0
**更新时间**：2026-02-06
