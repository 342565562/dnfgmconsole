package view

import (
	"dnf/biz/gm/model"
	"dnf/biz/gm/service"
	roleService "dnf/biz/user/role/service"
	"dnf/mods/game_db"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/localhostjason/webserver/server/util/uv"
)

func getAccounts(c *gin.Context) {
	var pi, order, q = &uv.PagingIn{}, &uv.Order{}, &service.AccountFilter{}
	uv.PQ(c, pi, order, q)
	
	// 如果是游戏账号登录，只查询当前用户的 UID
	currentUser := roleService.GetCurrentUser(c)
	if currentUser != nil && strings.HasPrefix(currentUser.Username, "game_") {
		// 从用户名中提取账号名（去掉 game_ 前缀）
		accountName := strings.TrimPrefix(currentUser.Username, "game_")
		
		// 查询 accounts 表获取 UID
		dbx := game_db.DBPools.Get(model.DTaiwan)
		var account model.Accounts
		if err := dbx.Where("accountname = ?", accountName).First(&account).Error; err == nil {
			// 设置 UID 精确匹配（不使用 LIKE）
			q.Uid = strconv.Itoa(account.Uid)
		}
	}
	
	lst, po, err := service.GetAccounts(q, pi, order)
	uv.PEIf(E_ACCOUNT_GET, err)

	c.JSON(200, uv.PagedOut(lst, po))
}

func rechargeAccount(c *gin.Context) {
	uid := uv.PPID(c, "id")

	data := &service.RechargeReq{}
	uv.PB(c, data)

	err := service.RechargeAccount(uid, data, c)
	uv.PEIf(E_RECHARGE_POST, err)
	c.Status(201)
}

func resetCreateCharac(c *gin.Context) {
	uid := uv.PPID(c, "id")
	err := service.ResetCreateCharac(uid)
	uv.PEIf(E_RESET_CREATE_CHARAC, err)
	c.Status(201)
}

func deleteAccount(c *gin.Context) {
	uid := uv.PPID(c, "id")
	err := service.DeleteAccount(uid)
	uv.PEIf(E_ACCOUNT_DELETE, err)
	c.Status(201)
}

// 清空宠物栏（并清空邮件）
func clearCreatures(c *gin.Context) {
	characNo := uv.PPID(c, "id")
	err := service.ClearCreaturesNotEquipped(characNo)
	uv.PEIf(E_ACCOUNT_DELETE, err)
	c.Status(204)
}

// 清空时装栏（并清空邮件）
func clearAvatars(c *gin.Context) {
	characNo := uv.PPID(c, "id")
	err := service.ClearAvatarsInBag(characNo)
	uv.PEIf(E_ACCOUNT_DELETE, err)
	c.Status(204)
}

// 一键恢复功能：同时执行删除邮件、删除宠物、删除时装
func restoreAccount(c *gin.Context) {
	characNo := uv.PPID(c, "id")
	err := service.RestoreAccount(characNo)
	uv.PEIf(E_ACCOUNT_DELETE, err)
	c.Status(204)
}

func changePassword(c *gin.Context) {
	args := &service.PasswordReq{}
	uv.PB(c, args)

	uid := uv.PPID(c, "id")
	err := service.ChangeAccountPassword(uid, args)
	uv.PEIf(E_ACCOUNT_CHANGE_PASSWORD, err)
	c.Status(201)
}

func createAccount(c *gin.Context) {
	args := &service.CreateAccountReq{}
	uv.PB(c, args)

	err := service.CreateAccount(args)
	uv.PEIf(E_ACCOUNT_CREARE, err)
	c.Status(201)
}

func updateAccount(c *gin.Context) {
	args := &service.UpdateAccountReq{}
	uv.PB(c, args)

	uid := uv.PPID(c, "id")
	err := service.UpdateAccountInfo(uid, args)
	uv.PEIf(E_ACCOUNT_UPDATE, err)
	c.Status(201)
}
