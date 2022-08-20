package models

import "gorm.io/gorm"

// SexMasterテーブル定義
type SexMaster struct {
	gorm.Model
	SexName string
}
