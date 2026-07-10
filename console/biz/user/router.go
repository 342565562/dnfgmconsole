package user

import (
	role "dnf/biz/user/role/view"
	users "dnf/biz/user/users/view"
	"dnf/mods/ginx"
)

func InitUserRouter(g *ginx.RouterGroup) {
	{
		role.InitRoleRouter(g.Group("用户权限", "role"))
		users.InitUsersRouter(g.Group("用户信息", ""))
	}
}
