package ginx

import (
	"dnf/biz/gm/model"
	"dnf/mods/game_db"
)

type Authz struct {
	ID        int    `json:"id"`
	GroupName string `json:"group_name"`
	ApiName   string `json:"api_name"`
	Url       string `json:"url" description:"对象 obj"`
	Method    string `json:"method" description:"act"`
}

func init() {
	game_db.RegTables(model.WebServer, &Authz{})
}
