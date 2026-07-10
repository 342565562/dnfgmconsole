package service

import (
	gmModel "dnf/biz/gm/model"
	roleService "dnf/biz/user/role/service"
	"dnf/biz/user/users/model"
	"dnf/mods/game_db"
	"errors"
	"github.com/gin-gonic/gin"
)

func GetUserInfo(c *gin.Context) *model.User {
	return roleService.GetCurrentUser(c)
}

func UpdateUserInfo(c *gin.Context, desc string) error {
	dbx := game_db.DBPools.Get(gmModel.WebServer)
	if dbx == nil {
		return errors.New("webserver database not connected")
	}
	currentUser := roleService.GetCurrentUser(c)
	currentUser.Desc = desc

	return dbx.Save(currentUser).Error
}

func GetUsers() []model.User {
	dbx := game_db.DBPools.Get(gmModel.WebServer)
	if dbx == nil {
		return make([]model.User, 0)
	}
	var users = make([]model.User, 0)
	dbx.Find(&users)
	return users
}
