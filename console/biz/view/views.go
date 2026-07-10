package view

import (
	client "dnf/biz/client/view"
	dash "dnf/biz/dash/view"
	gm "dnf/biz/gm/view"
	log "dnf/biz/log/view"
	"dnf/biz/middleware"
	"dnf/biz/static"
	"dnf/biz/user"
	auth "dnf/biz/user/auth/view"
	"dnf/mods/ginx"
	"github.com/gin-gonic/gin"

	_ "dnf/biz/user/users/service"
)

func SetView(r *gin.Engine) error {
	c, err := GetConfig()
	if err != nil {
		return err
	}
	if c.CORS {
		setCORS(r)
	}

	static.AddStaticToRouter(r)
	r.Use(middleware.OperateHandler)

	// 公开接口(免鉴权)：站点标题/登录名，供登录页读取
	r.GET("api/site-config", getSiteConfig)

	apiAuth := r.Group("api/auth")
	api := r.Group("api")

	err = auth.AddJwtAuth(apiAuth, api)
	if err != nil {
		return err
	}

	// load casbin
	api.Use(middleware.CasbinHandler, middleware.ErrorHandler)
	routeGroup := ginx.NewRouterGroup(api)
	{
		gm.InitGmRouter(routeGroup.Group("GM管理", "gm"))
		client.InitClientRouter(routeGroup.Group("客户端", "client"))

		user.InitUserRouter(routeGroup.Group("用户管理", "user"))
		log.InitLogRouter(routeGroup.Group("日志", "log"))
		dash.InitDashRouter(routeGroup.Group("统计", "dash"))
	}

	return nil
}
