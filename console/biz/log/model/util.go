package model

import (
	"dnf/biz/gm/model"
	"dnf/biz/user/auth/service"
	"dnf/mods/game_db"
	"github.com/gin-gonic/gin"
	"strings"
	"time"
)

// realIP 获取用户真实 IP：优先反向代理(nginx)注入的头，其次 gin 解析
func realIP(c *gin.Context) string {
	if ip := strings.TrimSpace(c.GetHeader("X-Real-IP")); ip != "" {
		return ip
	}
	if xff := c.GetHeader("X-Forwarded-For"); xff != "" {
		if ip := strings.TrimSpace(strings.Split(xff, ",")[0]); ip != "" {
			return ip
		}
	}
	return c.ClientIP()
}

func CreateRechargeLog(uid int, action int, number int, c *gin.Context) {
	dbx := game_db.DBPools.Get(model.WebServer)
	if dbx == nil {
		return
	}

	log := &RechargeLog{
		Uid:    uid,
		Time:   time.Now(),
		Action: action,
		Number: number,
		Ip:     realIP(c),
	}

	dbx.Create(log)
	return
}

func CreateOperateLog(action string, msg string, c *gin.Context) {
	dbx := game_db.DBPools.Get(model.WebServer)
	if dbx == nil {
		return
	}

	currentUser := service.CurrentUser(c)
	log := &OperateLog{
		Username: currentUser.Username,
		Time:     time.Now(),
		Action:   action,
		Result:   msg,
		Ip:       realIP(c),
	}
	dbx.Create(&log)
	return
}
