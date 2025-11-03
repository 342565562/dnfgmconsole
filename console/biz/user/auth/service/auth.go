package service

import (
	"console/biz/user/auth/model"
	m "console/biz/user/users/model"
	gmModel "console/biz/gm/model"
	"console/mods/game_db"
	"crypto/md5"
	"encoding/hex"
	"errors"
	"fmt"
	jwt "github.com/appleboy/gin-jwt/v2"
	"github.com/gin-gonic/gin"
	"github.com/localhostjason/webserver/db"
	"github.com/localhostjason/webserver/server/util/uv"
	uuid "github.com/satori/go.uuid"
	"gorm.io/gorm"
	"net/http"
	"time"
)

const currentUserKey = "current_user"
const currentPassword = "current_password"
const loginFailedKey = "___login_failed"

var _c model.ConfigAuth

func PayloadFunc(data interface{}) jwt.MapClaims {
	if v, ok := data.(*m.User); ok {
		return jwt.MapClaims{
			_c.IDKey: v.JwtKey,
		}
	}
	return jwt.MapClaims{}
}

func IdHandler(c *gin.Context) interface{} {
	claims := jwt.ExtractClaims(c)
	jwtKey, ok := claims[_c.IDKey].(string)
	if !ok {
		return nil
	}

	var user = &m.User{}
	err := db.DB.Where("jwt_key = ?", jwtKey).First(user).Error
	if err != nil {
		return nil
	}
	return user
}

type loginArgs struct {
	Username string `json:"username" binding:"required,lte=64"`
	Password string `json:"password" binding:"required,printascii,lte=128"`
}

// ToMd5 MD5加密函数
func ToMd5(str string) string {
	h := md5.New()
	h.Write([]byte(str))
	return hex.EncodeToString(h.Sum(nil))
}

func Authenticator(c *gin.Context) (interface{}, error) {
	var loginValues loginArgs
	if err := c.ShouldBind(&loginValues); err != nil {
		return nil, errors.New("请正确输入用户名密码")
	}

	userName := loginValues.Username
	password := loginValues.Password

	c.Set(operateKey, uv.OP(I_OP_LOGIN, userName))

	// 直接记录下来， 不管成功与否， 后面看情况使用
	c.Set(loginFailedKey, &m.User{Username: userName})

	// 1. 首先尝试默认 users 表登录（支持 admin/123）
	var user m.User
	err := db.DB.Where("username = ?", userName).First(&user).Error
	if err == nil && user.CheckPassword(password) {
		// 默认用户登录成功
		c.Set(currentUserKey, &user)
		c.Set(currentPassword, password)
		return &user, nil
	}

	// 2. 如果默认 users 表找不到或密码错误，尝试 d_taiwan.accounts 表
	if errors.Is(err, gorm.ErrRecordNotFound) || !user.CheckPassword(password) {
		dbx := game_db.DBPools.Get(gmModel.DTaiwan)
		var account gmModel.Accounts
		err := dbx.Where("accountname = ? AND password = ?", userName, ToMd5(password)).First(&account).Error
		
		if err == nil {
			// accounts 表登录成功，创建或获取对应的 User 记录
			// 查找是否已存在对应的用户（使用 game_ 前缀避免与现有用户冲突）
			gameUsername := fmt.Sprintf("game_%s", userName)
			var gameUser m.User
			err = db.DB.Where("username = ?", gameUsername).First(&gameUser).Error
			
			if errors.Is(err, gorm.ErrRecordNotFound) {
				// 如果不存在，创建新用户记录
				gameUser = m.User{
					Username:     gameUsername,
					Role:         "user",
					Email:        "",
					Desc:         fmt.Sprintf("游戏账号: %s (UID: %d)", userName, account.Uid),
					Time:         time.Now(),
					IsSuperAdmin: false,
					JwtKey:       uuid.NewV4(),
				}
				// 设置一个随机密码（实际不会用到，因为验证通过 accounts 表）
				gameUser.SetPassword(fmt.Sprintf("%s_%d", userName, account.Uid))
				db.DB.Create(&gameUser)
			} else {
				// 如果已存在，更新最后登录时间
				now := time.Now()
				gameUser.LastLoginTime = &now
				db.DB.Save(&gameUser)
			}
			
			c.Set(currentUserKey, &gameUser)
			c.Set(currentPassword, password)
			return &gameUser, nil
		}
	}

	// 所有验证都失败
	return nil, errors.New("用户名或者密码填写不对")
}

func Authorizator(data interface{}, c *gin.Context) bool {
	if data == nil {
		return false
	}
	fmt.Println("data", data)
	return true
}

// UnAuth 密码登录失败时候调用的函数
func UnAuth(c *gin.Context, code int, message string) {
	desc := "登录失败"
	msg := message
	if message == "Token is expired" {
		desc = ""
		msg = "您的登录已过期，请重新登录"
	}

	if message == "query token is empty" {
		desc = ""
		msg = "未携带身份凭证"
	}

	if message == "you don't have permission to access this resource" {
		msg = "携带身份凭证不正确，认证失败"
	}

	c.JSON(http.StatusUnauthorized, gin.H{
		"msg":  msg,
		"desc": desc,
	})
}

// LoginResponse 密码登录成功时调用的函数
func LoginResponse(c *gin.Context, code int, token string, expire time.Time) {
	u, _ := c.Get(currentUserKey)
	user := u.(*m.User)

	info := user.GetInfo()

	info["token"] = token

	now := time.Now()
	user.LastLoginTime = &now
	db.DB.Save(user)
	c.JSON(http.StatusOK, info)
}

// LogoutResponse 退出登录
func LogoutResponse(c *gin.Context, code int) {
	user := CurrentUser(c)
	c.Set(operateKey, uv.OP(I_OP_LOGOUT, user.Username))
	c.Status(201)
}

func CurrentUser(c *gin.Context) *m.User {
	u := IdHandler(c)
	currentUser, ok := u.(*m.User)
	if !ok {
		return &m.User{}
	}
	return currentUser
}
