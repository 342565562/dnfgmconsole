#!/bin/bash

# 颜色定义
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# 项目配置
IMAGE_NAME="gmwebconsole"
IMAGE_TAG="latest"
CONTAINER_NAME="gmwebconsole"

echo -e "${GREEN}================================${NC}"
echo -e "${GREEN}GM Web Console Docker 构建脚本${NC}"
echo -e "${GREEN}================================${NC}"

# 检查 Docker 是否安装
if ! command -v docker &> /dev/null; then
    echo -e "${RED}错误: Docker 未安装,请先安装 Docker${NC}"
    exit 1
fi

# 检查前端文件是否存在
if [ ! -d "console/dist/web/static" ]; then
    echo -e "${YELLOW}警告: 前端静态文件目录不存在: console/dist/web/static/${NC}"
    echo -e "${YELLOW}请确保已编译前端代码并放置在正确位置${NC}"
    read -p "是否继续构建? (y/n): " -n 1 -r
    echo
    if [[ ! $REPLY =~ ^[Yy]$ ]]; then
        exit 1
    fi
fi

# 检查配置文件是否存在
if [ ! -f "console/dist/config/server.json" ]; then
    echo -e "${RED}错误: 配置文件不存在: console/dist/config/server.json${NC}"
    exit 1
fi

echo -e "${GREEN}开始构建 Docker 镜像...${NC}"

# 构建镜像
docker build -t ${IMAGE_NAME}:${IMAGE_TAG} . || {
    echo -e "${RED}错误: Docker 镜像构建失败${NC}"
    exit 1
}

echo -e "${GREEN}✓ Docker 镜像构建成功: ${IMAGE_NAME}:${IMAGE_TAG}${NC}"

# 询问是否立即运行
read -p "是否立即运行容器? (y/n): " -n 1 -r
echo
if [[ $REPLY =~ ^[Yy]$ ]]; then
    # 检查是否已有同名容器在运行
    if [ "$(docker ps -q -f name=${CONTAINER_NAME})" ]; then
        echo -e "${YELLOW}容器 ${CONTAINER_NAME} 正在运行,先停止它...${NC}"
        docker stop ${CONTAINER_NAME}
    fi

    # 删除已存在的容器
    if [ "$(docker ps -aq -f name=${CONTAINER_NAME})" ]; then
        echo -e "${YELLOW}删除已存在的容器 ${CONTAINER_NAME}...${NC}"
        docker rm ${CONTAINER_NAME}
    fi

    # 运行容器
    echo -e "${GREEN}启动容器...${NC}"
    docker run -d \
        --name ${CONTAINER_NAME} \
        -p 8080:8080 \
        -v $(pwd)/console/dist/config/server.json:/app/config/server.json \
        -v $(pwd)/console/dist/log:/app/log \
        -v $(pwd)/console/dist/data:/app/data \
        -e TZ=Asia/Shanghai \
        --restart unless-stopped \
        ${IMAGE_NAME}:${IMAGE_TAG} || {
        echo -e "${RED}错误: 容器启动失败${NC}"
        exit 1
    }

    echo -e "${GREEN}✓ 容器启动成功${NC}"
    echo -e "${GREEN}容器名称: ${CONTAINER_NAME}${NC}"
    echo -e "${GREEN}访问地址: http://localhost:8080${NC}"
    echo ""
    echo -e "${YELLOW}查看日志: docker logs -f ${CONTAINER_NAME}${NC}"
    echo -e "${YELLOW}停止容器: docker stop ${CONTAINER_NAME}${NC}"
    echo -e "${YELLOW}删除容器: docker rm ${CONTAINER_NAME}${NC}"
fi

echo -e "${GREEN}================================${NC}"
echo -e "${GREEN}完成!${NC}"
echo -e "${GREEN}================================${NC}"
