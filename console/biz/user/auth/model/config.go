package model

import "github.com/localhostjason/webserver/server/config"

const _key = "auth"

func regConfig() {
	c := ConfigAuth{
		Realm:                "test zone",
		Secret:               "f450a7bdbde3416d22474b9fdc2a3636",
		IDKey:                "username",
		Timeout:              12 * 3600,
		MaxRefresh:           3600,
		ActivationCodeEnable: 0,
	}
	_ = config.RegConfig(_key, c)
}

type ConfigAuth struct {
	Realm      string `json:"realm"`
	Secret     string `json:"secret"`
	IDKey      string `json:"id_key"`
	Timeout    int    `json:"timeout"`
	MaxRefresh int    `json:"max_refresh"`
	// ActivationCodeEnable 激活码功能开关：1=开启（登录需激活码），0=关闭（无需激活码）
	ActivationCodeEnable int `json:"activation_code_enable"`
}

func GetConfig() (ConfigAuth, error) {
	var c ConfigAuth
	err := config.GetConfig(_key, &c)
	return c, err
}

func init() {
	regConfig()
}
