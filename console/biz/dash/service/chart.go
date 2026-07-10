package service

import (
	"dnf/biz/gm/model"
	logModel "dnf/biz/log/model"
	"dnf/mods/game_db"
	"fmt"
	"time"
)

func GetDashChart() ChartResult {
	return ChartResult{
		Cera:      getChartByType("cera"),
		CeraPoint: getChartByType("cera_point"),
	}
}

func getChartByType(typ string) ChartTranInfo {
	dbx := game_db.DBPools.Get(model.WebServer)
	if dbx == nil {
		return ChartTranInfo{
			Date:  make([]string, 0),
			Total: make([]int, 0),
		}
	}

	var data []ChartInfo

	now := time.Now()
	nowYear := now.Format("2006")

	tx := dbx.Model(&logModel.RechargeLog{}).
		Select("sum(number) as total, DATE_FORMAT(time, '%Y') as year, DATE_FORMAT(time, '%m') AS month").
		Group("year, month").
		Having("year = ?", nowYear)

	if typ == "cera" {
		tx = tx.Where("action = ?", logModel.RechargeCera)
	} else {
		tx = tx.Where("action = ?", logModel.RechargeCeraPoint)
	}
	tx.Find(&data)

	date := make([]string, 0)
	for i := 1; i <= 12; i++ {
		date = append(date, fmt.Sprintf("%d月", i))
	}

	number := make([]int, 0)
	for i := 1; i <= 12; i++ {
		isAppend := false
		for _, info := range data {
			if info.Month == i {
				number = append(number, info.Total)
				isAppend = true
				break
			}
		}

		if !isAppend {
			number = append(number, 0)
		}

	}

	return ChartTranInfo{
		Date:  date,
		Total: number,
	}
}
