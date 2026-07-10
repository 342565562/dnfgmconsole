package main

import (
	"bufio"
	"crypto/rand"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"math/big"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
)

// 数据库配置结构
type MysqlDBConfig struct {
	Key             string `json:"key"`
	User            string `json:"user"`
	Password        string `json:"password"`
	Host            string `json:"host"`
	Port            int    `json:"port"`
	DB              string `json:"db"`
	Charset         string `json:"charset"`
	Timeout         int    `json:"timeout"`
	MultiStatements bool   `json:"multi_statements"`
	Debug           bool   `json:"debug"`
}

type DbConfig struct {
	Enable bool            `json:"enable"`
	Mysql  []MysqlDBConfig `json:"mysql"`
}

type ServerConfig struct {
	GameDb DbConfig `json:"game_db"`
}

// 激活码模型
type ActivationCode struct {
	Id        int64      `gorm:"primaryKey"`
	Code      string     `gorm:"type:string;size:64;unique;not null"`
	IsUsed    bool       `gorm:"default:false"`
	UserId    *int64     `gorm:"type:bigint"`
	UsedAt    *time.Time
	CreatedAt time.Time `gorm:"autoCreateTime"`
}

// 字符集：数字0-9，小写字母a-z，大写字母A-Z
const charset = "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

// 生成16位随机激活码（只包含数字和英文字母大小写）
func generateActivationCode() (string, error) {
	code := make([]byte, 16)
	charsetLen := big.NewInt(int64(len(charset)))

	for i := range code {
		randomIndex, err := rand.Int(rand.Reader, charsetLen)
		if err != nil {
			return "", err
		}
		code[i] = charset[randomIndex.Int64()]
	}

	return string(code), nil
}

// 读取server.json配置文件
func loadConfig(configPath string) (*ServerConfig, error) {
	file, err := os.Open(configPath)
	if err != nil {
		return nil, fmt.Errorf("无法打开配置文件: %v", err)
	}
	defer file.Close()

	var config ServerConfig
	decoder := json.NewDecoder(file)
	if err := decoder.Decode(&config); err != nil {
		return nil, fmt.Errorf("解析配置文件失败: %v", err)
	}

	return &config, nil
}

// 查找webserver数据库配置
func findWebServerConfig(config *ServerConfig) (*MysqlDBConfig, error) {
	for _, mysqlConfig := range config.GameDb.Mysql {
		if mysqlConfig.Key == "webserver" {
			return &mysqlConfig, nil
		}
	}
	return nil, fmt.Errorf("未找到webserver数据库配置")
}

// 连接数据库
func connectDB(config *MysqlDBConfig) (*gorm.DB, error) {
	multiStatements := "false"
	if config.MultiStatements {
		multiStatements = "true"
	}

	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?multiStatements=%s&charset=%s&parseTime=True&loc=Local&timeout=%ds",
		config.User, config.Password, config.Host, config.Port, config.DB,
		multiStatements, config.Charset, config.Timeout)

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		NamingStrategy: schema.NamingStrategy{SingularTable: true},
	})
	if err != nil {
		return nil, fmt.Errorf("连接数据库失败: %v", err)
	}

	return db, nil
}

// 确保激活码表存在（工具独立运行时需要自行建表）
func ensureTables(db *gorm.DB) error {
	return db.AutoMigrate(&ActivationCode{})
}

// 检查激活码是否已存在
func codeExists(db *gorm.DB, code string) (bool, error) {
	var count int64
	err := db.Model(&ActivationCode{}).Where("code = ?", code).Count(&count).Error
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

// 生成唯一激活码
func generateUniqueCode(db *gorm.DB) (string, error) {
	maxRetries := 100 // 最多重试100次
	for i := 0; i < maxRetries; i++ {
		code, err := generateActivationCode()
		if err != nil {
			return "", err
		}

		exists, err := codeExists(db, code)
		if err != nil {
			return "", err
		}

		if !exists {
			return code, nil
		}
	}

	return "", fmt.Errorf("无法生成唯一激活码，已重试%d次", maxRetries)
}

// 批量生成激活码
func generateActivationCodes(db *gorm.DB, count int) error {
	fmt.Printf("\n开始生成 %d 个激活码...\n", count)

	successCount := 0
	failCount := 0
	codes := make([]ActivationCode, 0, count)

	for i := 0; i < count; i++ {
		code, err := generateUniqueCode(db)
		if err != nil {
			fmt.Printf("生成第 %d 个激活码失败: %v\n", i+1, err)
			failCount++
			continue
		}

		activationCode := ActivationCode{
			Code:      code,
			IsUsed:    false,
			UserId:    nil,
			UsedAt:    nil,
			CreatedAt: time.Now(),
		}

		codes = append(codes, activationCode)
		successCount++

		// 每生成10个显示进度
		if (i+1)%10 == 0 {
			fmt.Printf("已生成 %d/%d 个激活码...\n", i+1, count)
		}
	}

	// 批量插入数据库
	if len(codes) > 0 {
		fmt.Printf("\n正在将激活码保存到数据库...\n")
		if err := db.CreateInBatches(codes, 100).Error; err != nil {
			return fmt.Errorf("保存激活码到数据库失败: %v", err)
		}
		fmt.Printf("成功保存 %d 个激活码到数据库\n", len(codes))
	}

	// 显示生成的激活码
	fmt.Printf("\n生成的激活码列表：\n")
	fmt.Println(strings.Repeat("=", 52))
	for i, code := range codes {
		fmt.Printf("%d. %s\n", i+1, code.Code)
	}
	fmt.Println(strings.Repeat("=", 52))

	if failCount > 0 {
		fmt.Printf("\n警告: 有 %d 个激活码生成失败\n", failCount)
	}

	return nil
}

func main() {
	fmt.Println("==========================================")
	fmt.Println("    激活码生成工具")
	fmt.Println("==========================================")

	// 配置文件路径（默认读取“当前工作目录”的 server.json）
	// 允许通过 -config 显式指定配置文件路径
	var configPathFlag string
	flag.StringVar(&configPathFlag, "config", "", "server.json 配置文件路径（默认读取当前目录下的 server.json）")
	flag.Parse()

	cwd, err := os.Getwd()
	if err != nil {
		log.Fatalf("获取当前工作目录失败: %v", err)
	}

	candidates := make([]string, 0, 6)
	if strings.TrimSpace(configPathFlag) != "" {
		candidates = append(candidates, configPathFlag)
	}
	candidates = append(candidates,
		filepath.Join(cwd, "server.json"),
		filepath.Join(cwd, "config", "server.json"),
	)

	// 兼容：仍尝试可执行文件同级（避免老用法直接失效）
	if exePath, e := os.Executable(); e == nil {
		exeDir := filepath.Dir(exePath)
		candidates = append(candidates,
			filepath.Join(exeDir, "server.json"),
			filepath.Join(exeDir, "config", "server.json"),
		)
	}

	var configPath string
	for _, p := range candidates {
		if p == "" {
			continue
		}
		if _, e := os.Stat(p); e == nil {
			configPath = p
			break
		}
	}
	if configPath == "" {
		log.Fatalf("未找到配置文件。请在当前目录放置 server.json，或使用 -config 指定路径")
	}

	fmt.Printf("配置文件路径: %s\n", configPath)

	// 读取配置
	config, err := loadConfig(configPath)
	if err != nil {
		log.Fatalf("加载配置失败: %v", err)
	}

	// 查找webserver数据库配置
	dbConfig, err := findWebServerConfig(config)
	if err != nil {
		log.Fatalf("查找数据库配置失败: %v", err)
	}

	fmt.Printf("数据库配置: %s@%s:%d/%s\n", dbConfig.User, dbConfig.Host, dbConfig.Port, dbConfig.DB)

	// 连接数据库
	fmt.Println("\n正在连接数据库...")
	db, err := connectDB(dbConfig)
	if err != nil {
		log.Fatalf("连接数据库失败: %v", err)
	}
	defer func() {
		sqlDB, _ := db.DB()
		if sqlDB != nil {
			sqlDB.Close()
		}
	}()
	fmt.Println("数据库连接成功！")

	// 确保表存在（避免 activation_codes 表不存在导致生成失败）
	if err := ensureTables(db); err != nil {
		log.Fatalf("初始化表结构失败（请确认账号有建表权限）: %v", err)
	}

	// 提示用户输入数量
	reader := bufio.NewReader(os.Stdin)
	fmt.Print("\n请输入要生成的激活码数量: ")
	input, err := reader.ReadString('\n')
	if err != nil {
		log.Fatalf("读取输入失败: %v", err)
	}

	// 解析输入
	input = input[:len(input)-1] // 移除换行符
	if len(input) > 0 && input[len(input)-1] == '\r' {
		input = input[:len(input)-1] // Windows下移除\r
	}

	count, err := strconv.Atoi(input)
	if err != nil || count <= 0 {
		log.Fatalf("无效的数量: 请输入一个大于0的整数")
	}

	if count > 10000 {
		fmt.Printf("警告: 生成数量较大 (%d)，可能需要较长时间，是否继续？(y/n): ", count)
		confirm, _ := reader.ReadString('\n')
		if confirm[0] != 'y' && confirm[0] != 'Y' {
			fmt.Println("已取消")
			return
		}
	}

	// 生成激活码
	if err := generateActivationCodes(db, count); err != nil {
		log.Fatalf("生成激活码失败: %v", err)
	}

	fmt.Println("\n==========================================")
	fmt.Println("激活码生成完成！")
	fmt.Println("==========================================")
}
