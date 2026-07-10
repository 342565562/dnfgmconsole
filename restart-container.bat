@echo off
chcp 65001 >nul
setlocal enabledelayedexpansion

echo ================================
echo 重启 Docker 容器
echo ================================
echo.

:: 1. 停止容器
echo [步骤 1] 停止容器...
docker stop gmwebconsole
echo [成功] 容器已停止
echo.

:: 2. 删除容器
echo [步骤 2] 删除容器...
docker rm gmwebconsole
echo [成功] 容器已删除
echo.

:: 3. 重新运行容器
echo [步骤 3] 启动新容器...
docker run -d ^
  --name gmwebconsole ^
  -p 8080:8080 ^
  -v "%cd%\console\dist\config\server.json:/app/config/server.json" ^
  -v "%cd%\console\dist\log:/app/log" ^
  -e TZ=Asia/Shanghai ^
  --restart unless-stopped ^
  gmwebconsole:latest

if %errorlevel% neq 0 (
    echo.
    echo [错误] 容器启动失败
    pause
    exit /b 1
)

echo.
echo [成功] 容器启动成功
echo [信息] 访问地址: http://localhost:8080
echo.
echo [提示] 查看日志: docker logs -f gmwebconsole

:: 等待 2 秒后显示日志
timeout /t 2 /nobreak >nul
echo.
echo [最近日志]
docker logs --tail 20 gmwebconsole

echo.
echo ================================
echo 完成！
echo ================================
pause
