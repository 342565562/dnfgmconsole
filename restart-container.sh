#!/bin/bash

# 颜色定义
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m'

echo -e "${GREEN}================================${NC}"
echo -e "${GREEN}重启 Docker 容器${NC}"
echo -e "${GREEN}================================${NC}"
echo

# 1. 停止容器
echo -e "${YELLOW}步骤 1: 停止容器...${NC}"
docker stop gmwebconsole
echo -e "${GREEN}✓ 容器已停止${NC}"
echo

# 2. 删除容器
echo -e "${YELLOW}步骤 2: 删除容器...${NC}"
docker rm gmwebconsole
echo -e "${GREEN}✓ 容器已删除${NC}"
echo

# 3. 重新运行容器
echo -e "${YELLOW}步骤 3: 启动新容器...${NC}"
docker run -d \
  --name gmwebconsole \
  -p 8080:8080 \
  -v $(pwd)/console/dist/config/server.json:/app/config/server.json \
  -v $(pwd)/console/dist/log:/app/log \
  -e TZ=Asia/Shanghai \
  --restart unless-stopped \
  gmwebconsole:latest

if [ $? -eq 0 ]; then
    echo
    echo -e "${GREEN}✓ 容器启动成功${NC}"
    echo -e "${GREEN}访问地址: http://localhost:8080${NC}"
    echo
    echo -e "${YELLOW}查看日志: docker logs -f gmwebconsole${NC}"

    # 等待 2 秒后显示日志
    sleep 2
    echo
    echo -e "${YELLOW}最近日志:${NC}"
    docker logs --tail 20 gmwebconsole
else
    echo
    echo -e "${RED}✗ 容器启动失败${NC}"
    exit 1
fi

echo
echo -e "${GREEN}================================${NC}"
echo -e "${GREEN}完成！${NC}"
echo -e "${GREEN}================================${NC}"
