package view

import (
	"github.com/gin-gonic/gin"
	"github.com/localhostjason/webserver/server/config"
)

const _siteKey = "site"

// SiteConfig 站点信息，供前端展示(免鉴权获取)
type SiteConfig struct {
	Title     string `json:"title"`      // 浏览器标题
	LoginName string `json:"login_name"` // 登录页显示的名称
}

func init() {
	_ = config.RegConfig(_siteKey, SiteConfig{Title: "金华DNF", LoginName: "金华DNF"})
}

// GetSiteConfig 读取站点配置
func GetSiteConfig() SiteConfig {
	var c SiteConfig
	_ = config.GetConfig(_siteKey, &c)
	if c.Title == "" {
		c.Title = "金华DNF"
	}
	if c.LoginName == "" {
		c.LoginName = c.Title
	}
	return c
}

// getSiteConfig 公开接口 GET /api/site-config
func getSiteConfig(c *gin.Context) {
	c.JSON(200, GetSiteConfig())
}
