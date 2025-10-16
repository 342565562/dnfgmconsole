package service

import (
	"console/biz/gm/model"
	"console/mods/game_db"
	"time"

	"github.com/localhostjason/webserver/server/util/uv"
)

// SendEmail 插入 2条记录 letter postal
func SendEmail(characNo int, email *Email) error {
	dbx := game_db.DBPools.Get(model.TaiwanCain2nd)

	err := dbx.Table("letter").Create(map[string]interface{}{
		"charac_no":        characNo,
		"send_charac_no":   0,
		"send_charac_name": "GM",
		"letter_text":      "Thanks!",
		"stat":             1,
	}).Error

	if err != nil {
		return err
	}

	type _R struct {
		LetterId int `json:"letter_id"`
	}
	var data _R
	dbx.Table("letter").Where("charac_no = ?", characNo).Order("letter_id desc").Take(&data)

	//fmt.Println(123, data.LetterId)
	amplifyOption := 0
	amplifyValue := 0
	if email.IsAmplify {
		amplifyOption = email.AmplifyOption
		amplifyValue = email.AmplifyValue
	}

	var addInfo int = 0
	if email.AvataFlag {
		// 时装邮件，插入user_items表
		userItem := map[string]interface{}{
			"charac_no":   characNo,
			"it_id":       email.Code,
			"expire_date": "9999-12-31 23:59:59",
			"obtain_from": 1,
			"reg_date":    time.Now().Format("2006-01-02 15:04:05"),
			"stat":        2,
		}
		dbx.Table("user_items").Create(userItem)
		var ui struct {
			UiId int `json:"ui_id"`
		}
		dbx.Table("user_items").Where("charac_no = ? AND it_id = ? AND reg_date = ?", characNo, email.Code, userItem["reg_date"]).Order("ui_id desc").Take(&ui)
		addInfo = ui.UiId
	} else if email.CreatureFlag {
		// 宠物邮件，插入creature_items表
		creatureType := 1
		creatureItem := map[string]interface{}{
			"charac_no":     characNo,
			"it_id":         email.Code,
			"expire_date":   "9999-12-31 23:59:59",
			"reg_date":      time.Now().Format("2006-01-02 15:04:05"),
			"stat":          0,
			"item_lock_key": 0,
			"creature_type": creatureType,
			"stomach":       100,
		}
		dbx.Table("creature_items").Create(creatureItem)
		var ci struct {
			UiId int `json:"ui_id"`
		}
		dbx.Table("creature_items").Where("charac_no = ? AND it_id = ? AND reg_date = ?", characNo, email.Code, creatureItem["reg_date"]).Order("ui_id desc").Take(&ci)
		addInfo = ci.UiId
	} else {
		addInfo = email.Number
	}

	now := time.Now()
	return dbx.Debug().Table("postal").Create(map[string]interface{}{
		"occ_time":          now.Format("2006-01-02 15:04:05"),
		"send_charac_no":    0,
		"send_charac_name":  "GAME MASTER",
		"receive_charac_no": characNo,
		"item_id":           email.Code,
		"add_info":          addInfo,
		"letter_id":         data.LetterId,
		"seperate_upgrade":  email.SeperateUpgrade,
		"upgrade":           email.Upgrade,
		"amplify_option":    amplifyOption,
		"amplify_value":     amplifyValue,
		"gold":              email.Gold,
		"seal_flag":         email.SealFlag,
		"avata_flag":        email.AvataFlag,
		"creature_flag":     email.CreatureFlag,
	}).Error

}

func GetGoldList(q *GoldQ, pi *uv.PagingIn, order *uv.Order) ([]model.Gold, *uv.PagingOut, error) {
	dbx := game_db.DBPools.Get(model.TaiwanCain2nd)
	tx := q.FilterQuery(dbx)

	var lst = make([]model.Gold, 0)
	po, err := uv.PagingFind(tx, &lst, pi, order)

	return lst, po, err
}
