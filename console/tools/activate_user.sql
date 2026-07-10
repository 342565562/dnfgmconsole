-- =============================================
-- 手工激活用户脚本
-- 使用方法: mysql -h HOST -u USER -p DATABASE < activate_user.sql
-- =============================================

USE webserver;

-- 设置用户名和激活码（请根据实际情况修改）
SET @username = 'your_username';      -- 修改为实际用户名
SET @activation_code = 'YourCodeHere'; -- 修改为实际激活码

-- 显示当前要激活的用户信息
SELECT '===== 当前用户信息 =====' AS '';
SELECT id, username, is_activated, time AS registered_at
FROM user
WHERE username = @username;

-- 显示要使用的激活码信息
SELECT '===== 激活码信息 =====' AS '';
SELECT id, code, is_used, created_at
FROM activation_code
WHERE code = @activation_code;

-- 执行激活（需要手动取消注释才会执行）
-- UPDATE user SET is_activated = 1 WHERE username = @username;
-- UPDATE activation_code
-- SET is_used = 1,
--     user_id = (SELECT id FROM user WHERE username = @username),
--     used_at = NOW()
-- WHERE code = @activation_code;

-- 显示激活后的结果（取消上面注释后执行）
-- SELECT '===== 激活结果 =====' AS '';
-- SELECT
--     u.id,
--     u.username,
--     u.is_activated,
--     a.code AS activation_code,
--     a.used_at AS activated_at
-- FROM user u
-- LEFT JOIN activation_code a ON a.user_id = u.id
-- WHERE u.username = @username;

-- =============================================
-- 快速查询命令（可直接使用）
-- =============================================

-- 查看所有未激活用户
-- SELECT id, username, is_activated FROM user WHERE is_activated = 0;

-- 查看所有可用激活码
-- SELECT code FROM activation_code WHERE is_used = 0 LIMIT 10;

-- 查看所有用户的激活状态
-- SELECT
--     u.id, u.username, u.is_activated,
--     a.code AS activation_code
-- FROM user u
-- LEFT JOIN activation_code a ON a.user_id = u.id;
