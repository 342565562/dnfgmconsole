package model

import "time"

// Task定义的数据包传递到前端vue，想传递什么值，直接在Task结构体中修改
type Task struct {
	CharacNo  int  `json:"charac_no"`
	Letter_id int  `json:"play_1" gorm:"column:letter_id;"`
	Stat      bool `json:"play_1_trigger" gorm:"column:stat;"`
}

type Letter struct {
	Letter_id      int       `json:"letter_id" gorm:"column:letter_id;"`
	CharacNo       int       `json:"charac_no" gorm:"column:charac_no;"`
	SendCharacNo   int       `json:"send_charac_no" gorm:"column:send_charac_no;"`
	SendCharacName string    `json:"send_charac_name" gorm:"column:send_charac_name;"`
	LettText       string    `json:"letter_text" gorm:"column:letter_text;"`
	RegDate        time.Time `json:"reg_date" gorm:"column:reg_date;"`
	Stat           int8      `json:"stat" gorm:"column:stat;"`
}
type Postal struct {
	ReceiveCharacNo int   `json:"receive_charac_no" gorm:"column:receive_charac_no;"`
	ItemID          uint  `json:"item_id" gorm:"column:item_id;"`
	AddInfo         int   `json:"add_info" gorm:"column:add_info;"`
	Upgrade         uint8 `json:"upgrade" gorm:"column:upgrade;"`
	LetterID        int   `json:"letter_id" gorm:"column:letter_id;"`
}

/*
	type Task struct {
		CharacNo     int `json:"charac_no"`
		Play1        int `json:"play_1" gorm:"column:play_1;"`
		Play1Trigger int `json:"play_1_trigger" gorm:"column:play_1_trigger;"`
	}

	type Letter struct {
		CharacNo       int `gorm:"column:charac_no"`
		SendCharacNo   int `json:"send_charac_no" gorm:"column:send_charac_no;"`
		SendCharacName int `json:"send_charac_name" gorm:"column:send_charac_name;"`
		// 其他字段...
	}
*/
