package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"os"
	"strings"

	"dnf/tools/goldimport"
)

func main() {
	var (
		configPath string
		host       string
		port       int
		user       string
		pass       string
		dbName     string
		equipCSV   string
		itemCSV    string
		mode       string
		assumeYes  bool
	)

	flag.StringVar(&configPath, "config", "", "server.json 路径(读取 game_db 与 item_classify)")
	flag.StringVar(&host, "host", "", "数据库地址(覆盖 config)")
	flag.IntVar(&port, "port", 0, "数据库端口(覆盖 config，默认 3306)")
	flag.StringVar(&user, "user", "", "数据库账号(覆盖 config)")
	flag.StringVar(&pass, "pass", "", "数据库密码(覆盖 config)")
	flag.StringVar(&dbName, "db", "taiwan_cain_2nd", "目标数据库名")
	flag.StringVar(&equipCSV, "equip", "", "装备 CSV 路径(稀有度含时装关键字判为时装)")
	flag.StringVar(&itemCSV, "item", "", "道具 CSV 路径(按原始种类分类)")
	flag.StringVar(&mode, "mode", "truncate", "导入模式：truncate(清空重导) 或 append(追加)")
	flag.BoolVar(&assumeYes, "yes", false, "跳过确认(用于 truncate)")
	flag.Parse()

	p := goldimport.Params{
		Host: host, Port: port, User: user, Pass: pass, DB: dbName,
		EquipCSV: equipCSV, ItemCSV: itemCSV, Mode: mode,
	}

	// 从 server.json 补齐连接与分类规则(显式 flag 优先)
	if configPath != "" {
		h, pt, u, pw, cfg, ok := goldimport.LoadServerConfig(configPath, dbName)
		if ok {
			if p.Host == "" {
				p.Host = h
			}
			if p.Port == 0 {
				p.Port = pt
			}
			if p.User == "" {
				p.User = u
			}
			if p.Pass == "" {
				p.Pass = pw
			}
			p.Classify = cfg
		}
	}

	if mode == "truncate" && !assumeYes {
		if !confirm(fmt.Sprintf("即将【清空】%s.gold 表并重新导入，确认? (yes/no): ", dbName)) {
			fmt.Println("已取消")
			return
		}
	}

	_, err := goldimport.Run(p, func(format string, args ...interface{}) {
		fmt.Printf(format+"\n", args...)
	})
	if err != nil {
		log.Fatalln(err)
	}
}

func confirm(prompt string) bool {
	fmt.Print(prompt)
	sc := bufio.NewScanner(os.Stdin)
	if sc.Scan() {
		ans := strings.ToLower(strings.TrimSpace(sc.Text()))
		return ans == "yes" || ans == "y"
	}
	return false
}
