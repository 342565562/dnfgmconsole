package model

import (
	"dnf/biz/gm/model"
	"dnf/biz/user/auth/service"
	"dnf/mods/game_db"
	"github.com/gin-gonic/gin"
	"time"
)

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
		Ip:     c.RemoteIP(),
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
		Ip:       c.RemoteIP(),
	}
	dbx.Create(&log)
	return
}
