package service

import (
	"console/biz/gm/model"
	"console/mods/game_db"
	"errors"

	"gorm.io/gorm"
)

func GetTaskByRole(characNo int) ([]model.Task, error) {
	// 获取 TaiwanCain 数据库连接
	dbxTaiwanCain := game_db.DBPools.Get(model.TaiwanCain)

	// 检查角色是否存在（在 TaiwanCain 数据库）
	var characInfo model.CharacInfo
	err := dbxTaiwanCain.Table("charac_info").Where("charac_no = ?", characNo).First(&characInfo).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, errors.New("角色不存在！")
	}
	// 初始化用于存储查询结果的切片
	var data = make([]model.Task, 0)
	// 获取 taiwan_cain_2nd 数据库连接
	dbxTaiwanCain2nd := game_db.DBPools.Get(model.TaiwanCain2nd)

	// 查询 letter 表
	var letters []model.Letter
	dbxTaiwanCain2nd.Table("letter").Where("charac_no = ?", characNo).Find(&letters)
	for _, l := range letters {
		task := model.Task{
			CharacNo:  l.CharacNo,
			Letter_id: l.Letter_id,
			Stat:      l.Stat != 0, // 原有逻辑，stat非0为true
		}
		data = append(data, task)
	}

	// 查询 postal 表
	var postals []model.Postal
	dbxTaiwanCain2nd.Table("postal").Where("receive_charac_no = ?", characNo).Find(&postals)
	for _, p := range postals {
		task := model.Task{
			CharacNo:  p.ReceiveCharacNo,
			Letter_id: p.LetterID,
			Stat:      false, // play_1_trigger全部为false
		}
		data = append(data, task)
	}

	return data, nil
}

// UpdateTaskByRole 修改为根据 characNo 删除 taiwan_cain_2nd 库中的letter和postal表的记录
func UpdateTaskByRole(characNo int, ids []int) error {
	// 获取 TaiwanCain 数据库连接
	dbxTaiwanCain := game_db.DBPools.Get(model.TaiwanCain)

	// 检查角色是否存在（在 TaiwanCain 数据库）
	var characInfo model.CharacInfo
	err := dbxTaiwanCain.Table("charac_info").Where("charac_no = ?", characNo).First(&characInfo).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return errors.New("角色不存在！")
	}

	// 获取 taiwan_cain_2nd 数据库连接
	dbxTaiwanCain2nd := game_db.DBPools.Get(model.TaiwanCain2nd)

	if len(ids) == 0 {
		return errors.New("请选择邮件！")
	}
	// 只要ids长度大于0，删除该角色所有邮件（letter和postal表）
	dbxTaiwanCain2nd.Debug().Table("letter").Where("charac_no = ?", characNo).Delete(&model.Letter{})
	dbxTaiwanCain2nd.Debug().Table("postal").Where("receive_charac_no = ?", characNo).Delete(&model.Postal{})
	return nil
}
