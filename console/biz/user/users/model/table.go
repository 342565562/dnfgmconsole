package model

import (
	"dnf/biz/gm/model"
	"dnf/mods/game_db"
	uuid "github.com/satori/go.uuid"
	"golang.org/x/crypto/bcrypt"
	"time"
)

type User struct {
	// 基本信息
	Id            int64      `json:"id" gorm:"primaryKey"`
	Username      string     `json:"username" gorm:"type:string;size:64;unique;not null"`
	Password      string     `json:"-" gorm:"column:_password;type:string;size:128"`
	LastLoginTime *time.Time `json:"last_login_time"`
	Time          time.Time  `json:"time"`                           // 创建时间
	JwtKey        uuid.UUID  `json:"-" gorm:"type:string;size:128;"` // 为每个用户存一个唯一的jwt key (通用唯一识别码)

	Role  string `json:"role"`
	Email string `json:"email" gorm:"type:string;size:64"`
	Desc  string `json:"desc" gorm:"type:string;size:256"`

	IsSuperAdmin bool `json:"is_super_admin"`
	IsActivated  bool `json:"is_activated" gorm:"default:false"` // 是否已激活
}

func (u *User) SetPassword(password string) {
	b, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	if err != nil {
		return
	}
	u.Password = string(b)
}

func (u *User) CheckPassword(password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(password))
	return err == nil
}

func (u *User) GetInfo() map[string]interface{} {
	info := map[string]interface{}{
		"id":       u.Id,
		"username": u.Username,
		"role":     u.Role,
		"email":    u.Email,
		"desc":     u.Desc,
		"time":     u.Time.Format("2006-01-02T15:04:05Z07:00"),
	}
	
	// 处理 last_login_time，可能为 nil
	if u.LastLoginTime != nil {
		info["last_login_time"] = u.LastLoginTime.Format("2006-01-02T15:04:05Z07:00")
	} else {
		info["last_login_time"] = nil
	}
	
	return info
}

func init() {
	game_db.RegTables(model.WebServer, &User{})
}
