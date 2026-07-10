package model

import (
	"dnf/biz/gm/model"
	"dnf/mods/game_db"
	"time"
)

// ActivationCode 激活码模型
type ActivationCode struct {
	Id        int64      `json:"id" gorm:"primaryKey"`
	Code      string     `json:"code" gorm:"type:string;size:64;unique;not null"` // 激活码
	IsUsed    bool       `json:"is_used" gorm:"default:false"`                   // 是否已使用
	UserId    *int64     `json:"user_id" gorm:"type:bigint"`                     // 绑定的用户ID
	UsedAt    *time.Time `json:"used_at"`                                        // 使用时间
	CreatedAt time.Time  `json:"created_at" gorm:"autoCreateTime"`               // 创建时间
}

func init() {
	game_db.RegTables(model.WebServer, &ActivationCode{})
}
