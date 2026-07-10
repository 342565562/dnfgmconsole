//go:build windows

package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strconv"

	_ "embed"

	"dnf/tools/goldimport"

	"github.com/lxn/walk"
	. "github.com/lxn/walk/declarative"
)

//go:embed DNF.ico
var iconBytes []byte

func loadIcon() *walk.Icon {
	p := filepath.Join(os.TempDir(), "dnf_gold_import.ico")
	if err := os.WriteFile(p, iconBytes, 0644); err != nil {
		return nil
	}
	ic, err := walk.NewIconFromFile(p)
	if err != nil {
		return nil
	}
	return ic
}

func main() {
	var (
		hostLE, portLE, userLE, passLE, dbLE *walk.LineEdit
		equipLE, itemLE                      *walk.LineEdit
		configLE                             *walk.LineEdit
		modeCB                               *walk.ComboBox
		logTE                                *walk.TextEdit
		importPB                             *walk.PushButton
		mw                                   *walk.MainWindow
	)

	appendLog := func(s string) { logTE.AppendText(s + "\r\n") }

	browseFile := func(target *walk.LineEdit, filter string) {
		dlg := new(walk.FileDialog)
		dlg.Filter = filter
		dlg.Title = "选择文件"
		if ok, _ := dlg.ShowOpen(mw); ok {
			target.SetText(dlg.FilePath)
		}
	}

	// 从 server.json 一键填充
	fillFromConfig := func() {
		path := configLE.Text()
		if path == "" {
			walk.MsgBox(mw, "提示", "请先选择 server.json 路径", walk.MsgBoxIconWarning)
			return
		}
		db := dbLE.Text()
		if db == "" {
			db = "taiwan_cain_2nd"
		}
		h, pt, u, pw, cfg, ok := goldimport.LoadServerConfig(path, db)
		if !ok {
			walk.MsgBox(mw, "失败", "读取 server.json 失败", walk.MsgBoxIconError)
			return
		}
		if h != "" {
			hostLE.SetText(h)
		}
		if pt != 0 {
			portLE.SetText(strconv.Itoa(pt))
		}
		if u != "" {
			userLE.SetText(u)
		}
		if pw != "" {
			passLE.SetText(pw)
		}
		appendLog(fmt.Sprintf("已从配置读取 => 宠物:%v 时装:%v 关键字:%q", cfg.PetTypes, cfg.FashionTypes, cfg.FashionRarityKeyword))
	}

	doImport := func() {
		port, _ := strconv.Atoi(portLE.Text())
		mode := "truncate"
		if modeCB.CurrentIndex() == 1 {
			mode = "append"
		}
		p := goldimport.Params{
			Host: hostLE.Text(), Port: port, User: userLE.Text(),
			Pass: passLE.Text(), DB: dbLE.Text(),
			EquipCSV: equipLE.Text(), ItemCSV: itemLE.Text(), Mode: mode,
		}
		// 若填写了 config，用其中的分类规则
		if configLE.Text() != "" {
			if _, _, _, _, cfg, ok := goldimport.LoadServerConfig(configLE.Text(), p.DB); ok {
				p.Classify = cfg
			}
		}
		if mode == "truncate" {
			if walk.MsgBox(mw, "确认", "将【清空】gold 表并重新导入，确定?", walk.MsgBoxYesNo|walk.MsgBoxIconQuestion) != walk.DlgCmdYes {
				return
			}
		}
		importPB.SetEnabled(false)
		logTE.SetText("")
		go func() {
			_, err := goldimport.Run(p, func(format string, args ...interface{}) {
				msg := fmt.Sprintf(format, args...)
				mw.Synchronize(func() { appendLog(msg) })
			})
			mw.Synchronize(func() {
				if err != nil {
					appendLog("❌ 失败: " + err.Error())
					walk.MsgBox(mw, "失败", err.Error(), walk.MsgBoxIconError)
				} else {
					walk.MsgBox(mw, "完成", "导入完成", walk.MsgBoxIconInformation)
				}
				importPB.SetEnabled(true)
			})
		}()
	}

	_, _ = MainWindow{
		AssignTo: &mw,
		Title:    "DNF 物品导入工具",
		Icon:     loadIcon(),
		MinSize:  Size{Width: 740, Height: 620},
		Layout:   VBox{},
		Children: []Widget{
			GroupBox{
				Title:  "server.json（可选，一键读取连接与分类规则）",
				Layout: HBox{},
				Children: []Widget{
					LineEdit{AssignTo: &configLE},
					PushButton{Text: "浏览...", OnClicked: func() { browseFile(configLE, "JSON (*.json)|*.json|所有文件 (*.*)|*.*") }},
					PushButton{Text: "读取配置", OnClicked: fillFromConfig},
				},
			},
			GroupBox{
				Title:  "数据库配置",
				Layout: Grid{Columns: 4},
				Children: []Widget{
					Label{Text: "地址:"}, LineEdit{AssignTo: &hostLE, Text: "127.0.0.1"},
					Label{Text: "端口:"}, LineEdit{AssignTo: &portLE, Text: "3306"},
					Label{Text: "账号:"}, LineEdit{AssignTo: &userLE, Text: "game"},
					Label{Text: "密码:"}, LineEdit{AssignTo: &passLE, PasswordMode: true},
					Label{Text: "数据库:"}, LineEdit{AssignTo: &dbLE, Text: "taiwan_cain_2nd", ColumnSpan: 3},
				},
			},
			GroupBox{
				Title:  "CSV 文件（装备/道具至少填一个）",
				Layout: Grid{Columns: 3},
				Children: []Widget{
					Label{Text: "装备CSV:"}, LineEdit{AssignTo: &equipLE}, PushButton{Text: "浏览...", OnClicked: func() { browseFile(equipLE, "CSV (*.csv)|*.csv|所有文件 (*.*)|*.*") }},
					Label{Text: "道具CSV:"}, LineEdit{AssignTo: &itemLE}, PushButton{Text: "浏览...", OnClicked: func() { browseFile(itemLE, "CSV (*.csv)|*.csv|所有文件 (*.*)|*.*") }},
				},
			},
			Composite{
				Layout: HBox{},
				Children: []Widget{
					Label{Text: "导入模式:"},
					ComboBox{AssignTo: &modeCB, Model: []string{"清空重导 (truncate)", "追加导入 (append)"}, CurrentIndex: 0, MaxSize: Size{Width: 220}},
					HSpacer{},
					PushButton{AssignTo: &importPB, Text: "开始导入", MinSize: Size{Width: 130, Height: 38}, OnClicked: doImport},
				},
			},
			Label{Text: "运行日志:"},
			TextEdit{AssignTo: &logTE, ReadOnly: true, VScroll: true, MinSize: Size{Height: 320}},
		},
	}.Run()
}
