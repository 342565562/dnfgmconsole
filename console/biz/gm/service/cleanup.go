package service

import (
	"console/biz/gm/model"
	"console/mods/game_db"
	"errors"

	"gorm.io/gorm"
)

// ClearAllMailsByCharac 删除指定角色的所有邮件（letter 与 postal）
func ClearAllMailsByCharac(characNo int) error {
	// 校验角色是否存在（taiwan_cain）
	dbxTaiwanCain := game_db.DBPools.Get(model.TaiwanCain)
	var characInfo model.CharacInfo
	if err := dbxTaiwanCain.Table("charac_info").Where("charac_no = ?", characNo).First(&characInfo).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("角色不存在！")
		}
		return err
	}

	// 删除邮件（taiwan_cain_2nd）
	dbxTaiwanCain2nd := game_db.DBPools.Get(model.TaiwanCain2nd)
	if err := dbxTaiwanCain2nd.Table("letter").Where("charac_no = ?", characNo).Delete(&model.Letter{}).Error; err != nil {
		return err
	}
	if err := dbxTaiwanCain2nd.Table("postal").Where("receive_charac_no = ?", characNo).Delete(&model.Postal{}).Error; err != nil {
		return err
	}
	return nil
}

// ClearCreaturesNotEquipped 清除该角色宠物栏（非穿戴）中的宠物：creature_items.slot != 238
func ClearCreaturesNotEquipped(characNo int) error {
	// 先清邮件
	if err := ClearAllMailsByCharac(characNo); err != nil {
		return err
	}

	// 删除 creature_items 中 slot != 238 的记录（taiwan_cain_2nd）
	dbxTaiwanCain2nd := game_db.DBPools.Get(model.TaiwanCain2nd)
	if err := dbxTaiwanCain2nd.Table("creature_items").Where("charac_no = ? AND slot <> ?", characNo, 238).Delete(nil).Error; err != nil {
		return err
	}
	return nil
}

// ClearAvatarsInBag 清除该角色时装栏（非穿戴）的时装：user_items.slot > 10
func ClearAvatarsInBag(characNo int) error {
	// 先清邮件
	if err := ClearAllMailsByCharac(characNo); err != nil {
		return err
	}

	// 删除 user_items 中 slot > 9 的记录（taiwan_cain_2nd）
	dbxTaiwanCain2nd := game_db.DBPools.Get(model.TaiwanCain2nd)
	if err := dbxTaiwanCain2nd.Table("user_items").Where("charac_no = ? AND slot > ?", characNo, 9).Delete(nil).Error; err != nil {
		return err
	}
	return nil
}

// RestoreAccount 一键恢复功能：同时执行删除邮件、删除宠物、删除时装
// 用于解决无法登录游戏和网络中断问题
func RestoreAccount(characNo int) error {
	// 校验角色是否存在（taiwan_cain）
	dbxTaiwanCain := game_db.DBPools.Get(model.TaiwanCain)
	var characInfo model.CharacInfo
	if err := dbxTaiwanCain.Table("charac_info").Where("charac_no = ?", characNo).First(&characInfo).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("角色不存在！")
		}
		return err
	}

	dbxTaiwanCain2nd := game_db.DBPools.Get(model.TaiwanCain2nd)

	// 1. 删除邮件
	if err := dbxTaiwanCain2nd.Table("letter").Where("charac_no = ?", characNo).Delete(&model.Letter{}).Error; err != nil {
		return err
	}
	if err := dbxTaiwanCain2nd.Table("postal").Where("receive_charac_no = ?", characNo).Delete(&model.Postal{}).Error; err != nil {
		return err
	}

	// 2. 删除宠物（非穿戴）：creature_items.slot != 238
	if err := dbxTaiwanCain2nd.Table("creature_items").Where("charac_no = ? AND slot <> ?", characNo, 238).Delete(nil).Error; err != nil {
		return err
	}

	// 3. 删除时装（非穿戴）：user_items.slot > 9
	if err := dbxTaiwanCain2nd.Table("user_items").Where("charac_no = ? AND slot > ?", characNo, 9).Delete(nil).Error; err != nil {
		return err
	}

	return nil
}
