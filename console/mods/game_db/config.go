package game_db

import "github.com/localhostjason/webserver/server/config"

const _key = "game_db"

type MysqlDBConfig struct {
	Key             string `json:"key"`
	User            string `json:"user"`
	Password        string `json:"password"`
	Host            string `json:"host"`
	Port            int    `json:"port"`
	DB              string `json:"db"`
	Charset         string `json:"charset"`
	Timeout         int    `json:"timeout"`
	MultiStatements bool   `json:"multi_statements"`
	Debug           bool   `json:"debug"`
}

type DbConfig struct {
	Enable bool            `json:"enable"`
	Mysql  []MysqlDBConfig `json:"mysql"`
}

func init() {
	mc := MysqlDBConfig{
		Key:             "TestMysql1",
		User:            "gmweb",
		Password:        "GKhGNc6vtKERHEfiX5GW",
		Host:            "1.94.14.10",
		Port:            3306,
		DB:              "test",
		Charset:         "utf8mb4",
		MultiStatements: false,
		Timeout:         5,
		Debug:           false,
	}

	c := DbConfig{
		Enable: true,
		Mysql:  []MysqlDBConfig{mc},
	}
	_ = config.RegConfig(_key, c)
}

func getDbConfig() DbConfig {
	var c DbConfig
	_ = config.GetConfig(_key, &c)
	return c
}
