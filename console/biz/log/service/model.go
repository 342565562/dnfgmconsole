package service

import (
	"dnf/biz/log/model"
	"github.com/localhostjason/webserver/db"
	"gorm.io/gorm"
)

type RechargeQ struct {
	Uid string `form:"uid"`
}

func (q RechargeQ) FilterQuery(dbx *gorm.DB) (tx *gorm.DB) {
	tx = dbx.Model(&model.RechargeLog{})

	if q.Uid != "" {
		tx = tx.Where("uid like ?", db.Like(q.Uid))
	}
	// 默认最新在前(主键倒序，走主键索引，数据量大也快)
	tx = tx.Order("id desc")
	return
}

type OperateQ struct {
	Action string `form:"action"`
}

func (q OperateQ) FilterQuery(dbx *gorm.DB) (tx *gorm.DB) {
	tx = dbx.Model(&model.OperateLog{})

	if q.Action != "" {
		tx = tx.Where("action like ?", db.Like(q.Action))
	}
	// 默认最新在前
	tx = tx.Order("id desc")
	return
}
