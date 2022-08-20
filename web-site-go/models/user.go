package models

import "gorm.io/gorm"

// Userテーブルの定義
type User struct {
	gorm.Model
	Userid      string `unique:"true"`
	Name        string
	Age         int
	SexMasterID int
	SexMaster   SexMaster `gorm:"constraint:OnUpdate:CASCADE,OnDelete:SET 1;"`
	Password    string
}
