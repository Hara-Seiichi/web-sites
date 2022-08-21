package models

import "gorm.io/gorm"

// Userテーブルの定義
type User struct {
	gorm.Model
	Userid   string `gorm:"uniqueIndex"`
	Name     string
	Age      int
	Password string
	SexID    uint `gorm:"default:1;"`
	Sex      Sex
}
