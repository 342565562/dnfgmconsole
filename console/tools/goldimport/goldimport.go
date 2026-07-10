// Package goldimport 提供 gold 表导入的公共逻辑，供 CLI 与 GUI 复用。
package goldimport

import (
	"bufio"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"

	"dnf/biz/gm/model"

	"golang.org/x/text/encoding/simplifiedchinese"
	"golang.org/x/text/transform"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
)

// Params 导入参数
type Params struct {
	Host     string
	Port     int
	User     string
	Pass     string
	DB       string
	EquipCSV string // 装备 CSV(可空)
	ItemCSV  string // 道具 CSV(可空)
	Mode     string // truncate | append
	Classify model.ItemClassifyConfig
}

// Result 导入结果统计
type Result struct {
	EquipImported int
	EquipFiltered int
	ItemImported  int
	ItemFiltered  int
	Normal        int
	Avata         int
	Creature      int
	TotalInTable  int64
}

// Logf 日志回调
type Logf func(format string, args ...interface{})

// server.json 最小解析
type mysqlEntry struct {
	User     string `json:"user"`
	Password string `json:"password"`
	Host     string `json:"host"`
	Port     int    `json:"port"`
	DB       string `json:"db"`
	Key      string `json:"key"`
}

type serverConfig struct {
	GameDb struct {
		Mysql []mysqlEntry `json:"mysql"`
	} `json:"game_db"`
	ItemClassify *model.ItemClassifyConfig `json:"item_classify"`
}

// LoadServerConfig 从 server.json 读取目标库连接与分类规则(找不到规则则用默认)
func LoadServerConfig(path, dbName string) (host string, port int, user, pass string, cfg model.ItemClassifyConfig, ok bool) {
	cfg = model.DefaultItemClassifyConfig()
	b, err := os.ReadFile(path)
	if err != nil {
		return "", 0, "", "", cfg, false
	}
	var sc serverConfig
	if err := json.Unmarshal(b, &sc); err != nil {
		return "", 0, "", "", cfg, false
	}
	if sc.ItemClassify != nil {
		cfg = *sc.ItemClassify
	}
	pick := func() (mysqlEntry, bool) {
		for _, e := range sc.GameDb.Mysql {
			if e.DB == dbName || e.Key == dbName {
				return e, true
			}
		}
		if len(sc.GameDb.Mysql) > 0 {
			return sc.GameDb.Mysql[0], true
		}
		return mysqlEntry{}, false
	}
	if e, found := pick(); found {
		return e.Host, e.Port, e.User, e.Password, cfg, true
	}
	return "", 0, "", "", cfg, true
}

// Run 执行导入
func Run(p Params, logf Logf) (*Result, error) {
	if logf == nil {
		logf = func(string, ...interface{}) {}
	}
	if p.Mode != "truncate" && p.Mode != "append" {
		return nil, fmt.Errorf("无效模式: %s (只能 truncate/append)", p.Mode)
	}
	if p.EquipCSV == "" && p.ItemCSV == "" {
		return nil, fmt.Errorf("至少指定一个装备或道具 CSV 文件")
	}
	if p.Port == 0 {
		p.Port = 3306
	}
	if p.Host == "" || p.User == "" || p.DB == "" {
		return nil, fmt.Errorf("缺少数据库连接参数(host/user/db)")
	}
	// 分类规则为空则用默认
	if len(p.Classify.PetTypes) == 0 && len(p.Classify.FashionTypes) == 0 && p.Classify.FashionRarityKeyword == "" {
		p.Classify = model.DefaultItemClassifyConfig()
	}
	logf("分类规则 => 宠物:%v 时装:%v 时装稀有度关键字:%q", p.Classify.PetTypes, p.Classify.FashionTypes, p.Classify.FashionRarityKeyword)

	// 连接(utf8mb4，避免 latin1 库中文乱码)
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local&timeout=10s",
		p.User, p.Pass, p.Host, p.Port, p.DB)
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		SkipDefaultTransaction: true,
		NamingStrategy:         schema.NamingStrategy{SingularTable: true},
		Logger:                 logger.Default.LogMode(logger.Silent),
	})
	if err != nil {
		return nil, fmt.Errorf("连接数据库失败: %v", err)
	}
	logf("已连接 %s@%s:%d/%s", p.User, p.Host, p.Port, p.DB)

	// gold 表不存在才创建(name/rarity utf8mb4)
	if !db.Migrator().HasTable(&model.Gold{}) {
		if err := db.AutoMigrate(&model.Gold{}); err != nil {
			return nil, fmt.Errorf("建表失败: %v", err)
		}
		logf("gold 表不存在，已创建")
	} else {
		logf("gold 表已存在")
	}

	res := &Result{}
	var rows []model.Gold
	if p.EquipCSV != "" {
		part, filtered := parseCSV(p.EquipCSV, 0, 1, 6, -1, p.Classify)
		rows = append(rows, part...)
		res.EquipImported, res.EquipFiltered = len(part), filtered
		logf("装备: 有效 %d 条，过滤 %d 条(空/乱码/问号)", len(part), filtered)
	}
	if p.ItemCSV != "" {
		part, filtered := parseCSV(p.ItemCSV, 0, 1, 5, 2, p.Classify)
		rows = append(rows, part...)
		res.ItemImported, res.ItemFiltered = len(part), filtered
		logf("道具: 有效 %d 条，过滤 %d 条(空/乱码/问号)", len(part), filtered)
	}
	for _, r := range rows {
		switch r.Category {
		case model.CategoryAvata:
			res.Avata++
		case model.CategoryCreature:
			res.Creature++
		default:
			res.Normal++
		}
	}
	logf("合计有效 %d 条 —— 普通:%d 时装:%d 宠物:%d", len(rows), res.Normal, res.Avata, res.Creature)
	if len(rows) == 0 {
		return res, fmt.Errorf("没有可导入的数据")
	}

	if p.Mode == "truncate" {
		if err := db.Exec("TRUNCATE TABLE gold").Error; err != nil {
			if err2 := db.Where("1=1").Delete(&model.Gold{}).Error; err2 != nil {
				return res, fmt.Errorf("清空表失败: %v / %v", err, err2)
			}
		}
		logf("已清空 gold 表")
	}

	if err := db.CreateInBatches(&rows, 2000).Error; err != nil {
		return res, fmt.Errorf("导入失败: %v", err)
	}
	db.Model(&model.Gold{}).Count(&res.TotalInTable)
	logf("导入完成 ✅ 本次写入 %d 条，当前表内共 %d 条", len(rows), res.TotalInTable)
	return res, nil
}

// parseCSV 读取 GBK 编码 CSV，分类并过滤脏数据。rawTypeIdx<0 表示无原始种类列(装备表)
func parseCSV(path string, nameIdx, codeIdx, rarityIdx, rawTypeIdx int, cfg model.ItemClassifyConfig) ([]model.Gold, int) {
	f, err := os.Open(path)
	if err != nil {
		return nil, 0
	}
	defer f.Close()

	r := csv.NewReader(transform.NewReader(bufio.NewReader(f), simplifiedchinese.GBK.NewDecoder()))
	r.FieldsPerRecord = -1
	r.LazyQuotes = true

	var out []model.Gold
	filtered := 0
	first := true
	for {
		rec, err := r.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			continue
		}
		if first {
			first = false
			continue
		}
		if len(rec) <= codeIdx {
			continue
		}
		code, err := strconv.Atoi(strings.TrimSpace(rec[codeIdx]))
		if err != nil {
			continue
		}
		name := safeGet(rec, nameIdx)
		if IsBadName(name) {
			filtered++
			continue
		}
		rarity := safeGet(rec, rarityIdx)
		rawType := ""
		if rawTypeIdx >= 0 {
			rawType = safeGet(rec, rawTypeIdx)
		}
		out = append(out, model.Gold{
			Code:     code,
			Name:     name,
			Rarity:   rarity,
			Category: cfg.Classify(rawType, rarity),
		})
	}
	return out, filtered
}

// IsBadName 判定脏数据：空、纯问号、含解码替换符或西里尔字母(乱码)
func IsBadName(name string) bool {
	n := strings.TrimSpace(name)
	if n == "" {
		return true
	}
	if strings.Trim(n, "?？ 　") == "" {
		return true
	}
	for _, r := range n {
		if r == '�' {
			return true
		}
		if r >= 0x0400 && r <= 0x04FF {
			return true
		}
	}
	return false
}

func safeGet(rec []string, idx int) string {
	if idx >= 0 && idx < len(rec) {
		return strings.TrimSpace(rec[idx])
	}
	return ""
}
