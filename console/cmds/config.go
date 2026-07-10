package cmds

import (
	"dnf/biz/client/service"
	"fmt"

	"github.com/localhostjason/webserver/server/config"
	"github.com/sirupsen/logrus"
)

func DumpDefaultConfig() {
	content, err := config.GeneDefaultConfig()
	if err != nil {
		fmt.Println("failed to generate default config")
	} else {
		fmt.Println(string(content))
	}
}

func SyncDB() (err error) {
	// webserver/db已迁移到game_db，此函数不再使用
	// 数据库迁移现在通过game_db自动完成
	logrus.Info("webserver/db已迁移到game_db，数据库迁移通过game_db自动完成")
	return nil
}

func AutoMigrate() (err error) {
	// webserver/db已迁移到game_db，此函数不再使用
	// 数据库迁移现在通过game_db自动完成
	logrus.Info("webserver/db已迁移到game_db，数据库迁移通过game_db自动完成")
	return nil
}

func CreatePem() error {
	rsa := service.NewRsa()
	return rsa.GenRsaKey(2048)
}
