package model

import (
	"strings"

	"dnf/mods/game_db"

	"github.com/localhostjason/webserver/server/config"
)

// 物品分类：用于邮件自动判定时装/宠物
const (
	CategoryNormal   = 0 // 普通
	CategoryAvata    = 1 // 时装
	CategoryCreature = 2 // 宠物
)

type Gold struct {
	Id   int    `json:"id" gorm:"primaryKey"`
	Code int    `json:"code" gorm:"index"`
	Name string `json:"name" gorm:"type:varchar(128) CHARACTER SET utf8mb4"`
	// Rarity 稀有度，web 展示用
	Rarity string `json:"rarity" gorm:"type:varchar(16) CHARACTER SET utf8mb4"`
	// Category 分类(0普通/1时装/2宠物)，仅后端邮件判定用，不下发给前端
	Category int `json:"-" gorm:"type:tinyint;default:0;index"`
}

// ItemClassifyConfig 物品分类规则，写在 server.json 的 item_classify 段，可手工增删改。
// 分类在“导入工具”写入 gold.category 时生效；修改后需重新导入。
type ItemClassifyConfig struct {
	// PetTypes 归为宠物的“原始种类”(道具表英文原始类型)
	PetTypes []string `json:"pet_types"`
	// FashionTypes 归为时装的“原始种类”(道具表英文原始类型)
	FashionTypes []string `json:"fashion_types"`
	// FashionRarityKeyword 装备表“稀有度”含此关键字即判为时装(如 时装/稀有时装/史诗时装)
	FashionRarityKeyword string `json:"fashion_rarity_keyword"`
}

// DefaultItemClassifyConfig 默认分类规则
func DefaultItemClassifyConfig() ItemClassifyConfig {
	return ItemClassifyConfig{
		PetTypes:             []string{"creature", "creature expitem", "feed"},
		FashionTypes:         []string{"disguise", "disguise random", "dye"},
		FashionRarityKeyword: "时装",
	}
}

// Classify 依据规则判定分类：
// 1) 原始种类命中宠物集合 → 宠物；
// 2) 原始种类命中时装集合，或稀有度含时装关键字 → 时装；
// 3) 其余(含徽章、礼包、装备、未知/新类型) → 普通。
func (c ItemClassifyConfig) Classify(rawType, rarity string) int {
	rt := strings.ToLower(strings.TrimSpace(rawType))
	for _, t := range c.PetTypes {
		if rt != "" && rt == strings.ToLower(strings.TrimSpace(t)) {
			return CategoryCreature
		}
	}
	for _, t := range c.FashionTypes {
		if rt != "" && rt == strings.ToLower(strings.TrimSpace(t)) {
			return CategoryAvata
		}
	}
	if c.FashionRarityKeyword != "" && strings.Contains(rarity, c.FashionRarityKeyword) {
		return CategoryAvata
	}
	return CategoryNormal
}

// EnsureGoldTable 程序启动时调用：gold 表不存在才创建，已存在则跳过忽略。
// name/rarity 列显式 utf8mb4，即使所在库 taiwan_cain_2nd 为 latin1 也不会中文乱码。
func EnsureGoldTable() error {
	dbx := game_db.DBPools.Get(TaiwanCain2nd)
	if dbx == nil {
		return nil
	}
	if dbx.Migrator().HasTable(&Gold{}) {
		// 表已存在，跳过，不做任何结构改动
		return nil
	}
	return dbx.AutoMigrate(&Gold{})
}

func init() {
	// 注意：gold 表不走通用 AutoMigrate(避免对已存在表做结构改动)，
	// 改由 EnsureGoldTable 在启动时“不存在才创建”。
	_ = config.RegConfig("item_classify", DefaultItemClassifyConfig())
}
