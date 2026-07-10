package main

import (
	"dnf/cmds"

	"github.com/gin-gonic/gin"
)

func main() {
	// 生产模式：关闭 gin 调试日志与路由打印，降低每请求开销
	gin.SetMode(gin.ReleaseMode)
	cmds.Run()
}
