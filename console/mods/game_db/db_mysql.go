// Package game_db Package db 处理数据库连接信息，实现Model对象的存储读取
package game_db

import (
	"errors"
	"fmt"
	log "github.com/sirupsen/logrus"
	"sync"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
)

// ConnectWithMysqlConfig 连接，检验配置是否正确
func ConnectWithMysqlConfig(cfgs []MysqlDBConfig) error {
	if len(cfgs) == 0 {
		return nil
	}

	var wg = &sync.WaitGroup{}
	for _, c := range cfgs {
		wg.Add(1)

		go func(c MysqlDBConfig) {
			defer wg.Done()

			err := connectMysqlOne(c)
			if err != nil {
				log.Fatal(err)
			}
		}(c)

	}
	wg.Wait()
	return nil
}

func connectMysqlOne(c MysqlDBConfig) error {
	// 先连接到MySQL服务器（不指定数据库）以检查数据库是否存在
	serverDSN := fmt.Sprintf("%s:%s@tcp(%s:%d)/?charset=%s&parseTime=True&loc=Local&timeout=%ds",
		c.User, c.Password, c.Host, c.Port, c.Charset, c.Timeout)
	serverDB, err := gorm.Open(mysql.Open(serverDSN), &gorm.Config{})
	if err != nil {
		return errors.New(fmt.Sprintf("failed to connect to MySQL server %s:%d: %v", c.Host, c.Port, err))
	}

	// 检查数据库是否存在（COUNT 返回数字，可正确 Scan 到 int）
	var count int
	err = serverDB.Raw("SELECT COUNT(*) FROM INFORMATION_SCHEMA.SCHEMATA WHERE SCHEMA_NAME = ?", c.DB).Scan(&count).Error
	if err != nil {
		serverDB = nil
		return errors.New(fmt.Sprintf("failed to check database existence: %v", err))
	}

	// 如果数据库不存在，创建它（使用配置中的字符集）
	if count == 0 {
		createSQL := fmt.Sprintf("CREATE DATABASE IF NOT EXISTS `%s` DEFAULT CHARACTER SET %s", c.DB, c.Charset)
		err = serverDB.Exec(createSQL).Error
		if err != nil {
			serverDB = nil
			return errors.New(fmt.Sprintf("failed to create database %s: %v", c.DB, err))
		}
		log.Infof("Database %s created successfully", c.DB)
	}

	// 关闭服务器连接
	sqlDB, _ := serverDB.DB()
	if sqlDB != nil {
		sqlDB.Close()
	}

	// 连接到指定的数据库
	multiStatements := "false"
	if c.MultiStatements {
		multiStatements = "true"
	}
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?multiStatements=%s&charset=%s&parseTime=True&loc=Local&timeout=%ds&allowAllFiles=true",
		c.User, c.Password, c.Host, c.Port, c.DB, multiStatements, c.Charset, c.Timeout)
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		FullSaveAssociations:   true,
		SkipDefaultTransaction: true,
		NamingStrategy:         schema.NamingStrategy{SingularTable: true},
	})

	if c.Debug {
		db = db.Debug()
	}

	if err != nil {
		return errors.New(fmt.Sprintf("failed to connect databse %s:%v", dsn, err))
	}

	DBPools.Add(c.Key, db)
	return nil
}

// Close todo 关闭db连接
func Close() error {
	return nil
}
