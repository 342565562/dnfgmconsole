@echo off
chcp 65001 >nul
setlocal enabledelayedexpansion

echo ================================
echo 清理 Docker 缓存并重新构建
echo ================================
echo.

:: 1. 清理 Docker 构建缓存
echo [步骤 1] 清理 Docker 构建缓存...
docker builder prune -f
echo [成功] 缓存清理完成
echo.

:: 2. 删除旧镜像（可选）
echo [步骤 2] 删除旧镜像...
docker rmi gmwebconsole:latest 2>nul || echo 旧镜像不存在，跳过
echo.

:: 3. 重新构建镜像（不使用缓存）
echo [步骤 3] 重新构建镜像（不使用缓存）...
docker build --no-cache -t gmwebconsole:latest .

if %errorlevel% neq 0 (
    echo.
    echo [错误] 构建失败
    pause
    exit /b 1
)

echo.
echo [成功] 镜像构建成功！
echo.

:: 询问是否立即运行
set /p run="是否立即运行容器? (y/n): "
if /i not "%run%"=="y" goto :end

:: 停止并删除旧容器
docker ps -q -f name=gmwebconsole >nul 2>&1
if %errorlevel% equ 0 (
    echo [信息] 停止旧容器...
    docker stop gmwebconsole
)

docker ps -aq -f name=gmwebconsole >nul 2>&1
if %errorlevel% equ 0 (
    echo [信息] 删除旧容器...
    docker rm gmwebconsole
)

:: 运行新容器
echo [信息] 启动新容器...
docker run -d ^
    --name gmwebconsole ^
    -p 8080:8080 ^
    -v "%cd%\console\dist\config\server.json:/app/config/server.json" ^
    -v "%cd%\console\dist\log:/app/log" ^
    -e TZ=Asia/Shanghai ^
    --restart unless-stopped ^
    gmwebconsole:latest

if %errorlevel% neq 0 (
    echo [错误] 容器启动失败
    pause
    exit /b 1
)

echo.
echo [成功] 容器启动成功
echo [信息] 访问地址: http://localhost:8080
echo.
echo [提示] 查看日志: docker logs -f gmwebconsole

:end
echo.
echo ================================
echo 完成！
echo ================================
pause
