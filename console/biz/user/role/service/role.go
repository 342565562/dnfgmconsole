package service

import (
	gmModel "dnf/biz/gm/model"
	jwtService "dnf/biz/user/auth/service"
	"dnf/biz/user/users/model"
	"dnf/mods/casbinx"
	"dnf/mods/game_db"
	"errors"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func CreateRole(role, path, method string) error {
	if success, _ := casbinx.E.AddPolicy(role, path, method); !success {
		return errors.New("存在相同的策略，添加失败")
	}
	return nil
}

func UpdateRoleDesc(id int, desc string) error {
	policy, err := checkPolicyInDb(id)
	if err != nil {
		return err
	}

	dbx := game_db.DBPools.Get(gmModel.WebServer)
	if dbx == nil {
		return errors.New("webserver database not connected")
	}

	policy.Desc = desc
	return dbx.Save(policy).Error
}

func UpdateRole(id int, role, path, method string) error {
	policy, err := checkPolicyInDb(id)
	if err != nil {
		return err
	}

	updated, err := casbinx.E.UpdatePolicy([]string{policy.Role, policy.Path, policy.Method}, []string{role, path, method})
	if err != nil {
		return err
	}
	if !updated {
		return errors.New("更新策略失败")
	}
	return nil
}

func DeleteRole(id int) error {
	policy, err := checkPolicyInDb(id)
	if err != nil {
		return err
	}

	ok, err := casbinx.E.RemoveFilteredNamedPolicy(policy.PType, 0, policy.Role, policy.Path, policy.Method)
	if err != nil {
		return err
	}
	if !ok {
		return errors.New("已删除策略")
	}
	return nil
}

func checkPolicyInDb(id int) (*casbinx.CasbinRule, error) {
	dbx := game_db.DBPools.Get(gmModel.WebServer)
	if dbx == nil {
		return nil, errors.New("webserver database not connected")
	}

	var rule casbinx.CasbinRule
	err := dbx.First(&rule, &casbinx.CasbinRule{ID: uint(id)}).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return &rule, errors.New("未找到策略")
	}
	return &rule, nil
}

func GetRolePolicy() [][]string {
	return casbinx.E.GetPolicy()
}

func GetCurrentUser(ctx *gin.Context) *model.User {
	currentUser := jwtService.CurrentUser(ctx)
	return currentUser
}
