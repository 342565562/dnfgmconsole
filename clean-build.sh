#!/bin/bash

# 颜色定义
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

echo -e "${GREEN}================================${NC}"
echo -e "${GREEN}清理 Docker 缓存并重新构建${NC}"
echo -e "${GREEN}================================${NC}"
echo

# 1. 清理 Docker 构建缓存
echo -e "${YELLOW}步骤 1: 清理 Docker 构建缓存...${NC}"
docker builder prune -f
echo -e "${GREEN}✓ 缓存清理完成${NC}"
echo

# 2. 删除旧镜像（可选）
echo -e "${YELLOW}步骤 2: 删除旧镜像...${NC}"
docker rmi webconsole:lk70s2a1 2>/dev/null || echo "旧镜像不存在，跳过"
echo

# 3. 重新构建镜像（不使用缓存）
echo -e "${YELLOW}步骤 3: 重新构建镜像（不使用缓存）...${NC}"
docker build --no-cache -t webconsole:lk70s2a1 .

if [ $? -eq 0 ]; then
    echo
    echo -e "${GREEN}✓ 镜像构建成功！${NC}"
    echo

    # 询问是否立即运行
    read -p "是否立即运行容器? (y/n): " -n 1 -r
    echo
    if [[ $REPLY =~ ^[Yy]$ ]]; then
        # 停止并删除旧容器
        if [ "$(docker ps -q -f name=gmwebconsole)" ]; then
            echo -e "${YELLOW}停止旧容器...${NC}"
            docker stop gmwebconsole
        fi

        if [ "$(docker ps -aq -f name=gmwebconsole)" ]; then
            echo -e "${YELLOW}删除旧容器...${NC}"
            docker rm gmwebconsole
        fi

        # 运行新容器
        echo -e "${GREEN}启动新容器...${NC}"
        docker run -d \
            --name gmwebconsole \
            -p 8080:8080 \
            -v $(pwd)/console/dist/config/server.json:/app/config/server.json \
            -v $(pwd)/console/dist/log:/app/log \
            -e TZ=Asia/Shanghai \
            --restart unless-stopped \
            gmwebconsole:latest

        echo
        echo -e "${GREEN}✓ 容器启动成功${NC}"
        echo -e "${GREEN}访问地址: http://localhost:8080${NC}"
        echo
        echo -e "${YELLOW}查看日志: docker logs -f gmwebconsole${NC}"
    fi
else
    echo
    echo -e "${RED}✗ 构建失败${NC}"
    exit 1
fi

echo
echo -e "${GREEN}================================${NC}"
echo -e "${GREEN}完成！${NC}"
echo -e "${GREEN}================================${NC}"
