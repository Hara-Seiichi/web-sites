package db

import (
	"web-site-go/models"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres" // Use PostgreSQL in gorm
)

var (
	db  *gorm.DB
	err error
)

// DB接続の初期化
func Init() {
	db, err = gorm.Open("postgres", "host=127.0.0.1 port=5432 user=postgres dbname=golang_gin_db password=sei1013sei sslmode=disable")
	if err != nil {
		panic(err)
	}

	autoMaigration()

}

// DB接続の取得
func GetDB() *gorm.DB {
	return db
}

// DB切断
func Close() {
	if err := db.Close(); err != nil {
		panic(err)
	}
}

// エンティティーの移行。不足フィールドの追加のみ実行する。データは変更しない。
func autoMaigration() {
	db.AutoMigrate(&models.User{})
	db.AutoMigrate(&models.SexMaster{})
}
