package view

import (
	gmModel "console/biz/gm/model"
	"console/biz/user/users/model"
	"console/biz/user/users/service"
	"console/mods/game_db"
	"github.com/gin-gonic/gin"
	"github.com/localhostjason/webserver/server/util/uv"
	"strings"
)

func getUserInfo(c *gin.Context) {
	//c.Set(operateKey, uv.OP(I_OP, "1", "hello"))

	user := service.GetUserInfo(c)
	info := user.GetInfo()
	
	// 如果是游戏账号，查询并返回 UID
	if strings.HasPrefix(user.Username, "game_") {
		// 从用户名中提取账号名（去掉 game_ 前缀）
		accountName := strings.TrimPrefix(user.Username, "game_")
		
		// 查询 accounts 表获取 UID
		dbx := game_db.DBPools.Get(gmModel.DTaiwan)
		var account gmModel.Accounts
		if err := dbx.Where("accountname = ?", accountName).First(&account).Error; err == nil {
			info["game_uid"] = account.Uid
			info["is_game_account"] = true
		}
	} else {
		info["is_game_account"] = false
	}
	
	c.JSON(200, info)
}

func updateUserInfo(c *gin.Context) {
	userQ := &model.UserInfoPut{}
	uv.PB(c, userQ)

	err := service.UpdateUserInfo(c, userQ.Desc)
	uv.PEIf(E_USER_INFO_UPDATE, err)
	c.Status(201)
}

func getUsers(c *gin.Context) {
	data := service.GetUsers()
	c.JSON(200, data)
}
