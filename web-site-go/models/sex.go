package models

import "gorm.io/gorm"

// SexMasterテーブル定義
type Sex struct {
	gorm.Model
	Name string `gorm:"uniqueIndex;default:'-';"`
}
