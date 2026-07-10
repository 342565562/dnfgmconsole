package service

import (
	gmModel "dnf/biz/gm/model"
	"dnf/biz/user/users/model"
	"dnf/mods/game_db"
	"errors"
	uuid "github.com/satori/go.uuid"
	"time"
)

const (
	Normal = 1
	Locked = -1
)

const (
	ADMIN    = "admin"
	SECURITY = "security"
	AUDITOR  = "auditor"
)

const _defaultPassword = "123"

func InitUser() error {
	now := time.Now()
	users := []model.User{
		{
			Username:     "admin",
			Role:         ADMIN,
			Email:        "",
			Desc:         "超级管理员",
			Time:         now,
			IsSuperAdmin: true,
		},
	}

	dbx := game_db.DBPools.Get(gmModel.WebServer)
	if dbx == nil {
		return errors.New("webserver database not connected")
	}

	for i := range users {
		u := &users[i]

		var userList []model.User
		if dbx.Limit(1).Where("username = ? ", u.Username).Find(&userList); len(userList) == 0 {
			u.SetPassword(_defaultPassword)
			u.JwtKey = uuid.NewV4()
			u.IsActivated = true // 初始化时创建的用户默认为已激活（旧用户）
			dbx.Create(u)
		} else {
			// 如果用户已存在，确保IsActivated为true（旧用户免激活）
			if len(userList) > 0 {
				existingUser := &userList[0]
				if !existingUser.IsActivated {
					existingUser.IsActivated = true
					dbx.Save(existingUser)
				}
			}
		}
	}
	return nil
}

func init() {
	game_db.AddInitHook(InitUser)
}
