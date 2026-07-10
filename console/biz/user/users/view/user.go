package view

import (
	gmModel "dnf/biz/gm/model"
	"dnf/biz/user/users/model"
	"dnf/biz/user/users/service"
	"dnf/mods/game_db"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/localhostjason/webserver/server/util/uv"
)

func getUserInfo(c *gin.Context) {
	//c.Set(operateKey, uv.OP(I_OP, "1", "hello"))

	user := service.GetUserInfo(c)
	if user == nil {
		c.JSON(401, gin.H{"msg": "用户信息获取失败"})
		return
	}

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
		} else {
			// 如果查询失败，仍然标记为游戏账号，但不设置 UID
			info["is_game_account"] = true
			info["game_uid"] = nil
		}
	} else {
		info["is_game_account"] = false
		info["game_uid"] = nil
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
