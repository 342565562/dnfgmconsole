@echo off
chcp 65001 >nul
setlocal enabledelayedexpansion

:: 项目配置
set IMAGE_NAME=gmwebconsole
set IMAGE_TAG=latest
set CONTAINER_NAME=gmwebconsole

echo ================================
echo GM Web Console Docker 构建脚本
echo ================================
echo.

:: 检查 Docker 是否安装
docker --version >nul 2>&1
if %errorlevel% neq 0 (
    echo [错误] Docker 未安装,请先安装 Docker
    pause
    exit /b 1
)

:: 检查前端文件是否存在
if not exist "console\dist\web\static" (
    echo [警告] 前端静态文件目录不存在: console\dist\web\static\
    echo [警告] 请确保已编译前端代码并放置在正确位置
    set /p continue="是否继续构建? (y/n): "
    if /i not "!continue!"=="y" exit /b 1
)

:: 检查配置文件是否存在
if not exist "console\dist\config\server.json" (
    echo [错误] 配置文件不存在: console\dist\config\server.json
    pause
    exit /b 1
)

echo [信息] 开始构建 Docker 镜像...
echo.

:: 构建镜像
docker build -t %IMAGE_NAME%:%IMAGE_TAG% .
if %errorlevel% neq 0 (
    echo [错误] Docker 镜像构建失败
    pause
    exit /b 1
)

echo.
echo [成功] Docker 镜像构建成功: %IMAGE_NAME%:%IMAGE_TAG%
echo.

:: 询问是否立即运行
set /p run="是否立即运行容器? (y/n): "
if /i not "%run%"=="y" goto :end

:: 检查是否已有同名容器在运行
docker ps -q -f name=%CONTAINER_NAME% >nul 2>&1
if %errorlevel% equ 0 (
    echo [信息] 容器 %CONTAINER_NAME% 正在运行,先停止它...
    docker stop %CONTAINER_NAME%
)

:: 删除已存在的容器
docker ps -aq -f name=%CONTAINER_NAME% >nul 2>&1
if %errorlevel% equ 0 (
    echo [信息] 删除已存在的容器 %CONTAINER_NAME%...
    docker rm %CONTAINER_NAME%
)

:: 运行容器
echo [信息] 启动容器...
docker run -d ^
    --name %CONTAINER_NAME% ^
    -p 8080:8080 ^
    -v "%cd%\console\dist\config\server.json:/app/config/server.json" ^
    -v "%cd%\console\dist\log:/app/log" ^
    -v "%cd%\console\dist\data:/app/data" ^
    -e TZ=Asia/Shanghai ^
    --restart unless-stopped ^
    %IMAGE_NAME%:%IMAGE_TAG%

if %errorlevel% neq 0 (
    echo [错误] 容器启动失败
    pause
    exit /b 1
)

echo.
echo [成功] 容器启动成功
echo [信息] 容器名称: %CONTAINER_NAME%
echo [信息] 访问地址: http://localhost:8080
echo.
echo [提示] 查看日志: docker logs -f %CONTAINER_NAME%
echo [提示] 停止容器: docker stop %CONTAINER_NAME%
echo [提示] 删除容器: docker rm %CONTAINER_NAME%

:end
echo.
echo ================================
echo 完成!
echo ================================
pause
