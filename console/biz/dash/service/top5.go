package service

import (
	"dnf/biz/gm/model"
	logModel "dnf/biz/log/model"
	"dnf/mods/game_db"
)

func GetRechargeTop5() StatTop5Result {
	return StatTop5Result{
		Cera:           getRechargeTop5ByType("cera"),
		CeraPoint:      getRechargeTop5ByType("cera_point"),
		CeraTotal:      GetRechargeTotal("cera"),
		CeraPointTotal: GetRechargeTotal("cera_point"),
	}
}

func getRechargeTop5ByType(typ string) []TopInfo {
	dbx := game_db.DBPools.Get(model.WebServer)
	if dbx == nil {
		return make([]TopInfo, 0)
	}

	var data []TopInfo
	tx := dbx.Model(&logModel.RechargeLog{}).Select("sum(number) as total, uid").Group("action, uid").Order("total desc").Limit(5)

	if typ == "cera" {
		tx.Having("action = 1")
	} else {
		tx.Having("action = 2")
	}
	tx.Find(&data)

	result := make([]TopInfo, 0)
	for _, info := range data {
		info.AccountName = getAccountNameByUid(info.Uid)
		result = append(result, info)
	}
	return result
}

func getAccountNameByUid(uid int) string {
	dbx := game_db.DBPools.Get(model.DTaiwan)
	var data model.Accounts
	dbx.Where("UID = ?", uid).First(&data)
	return data.AccountName
}
