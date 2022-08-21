package db

import (
	"web-site-go/configs"
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

	db, err = gorm.Open("postgres", "host="+configs.DB_HOST+" port="+configs.DB_PORT+" user="+configs.DB_USER+" dbname="+configs.DB_NAME+" password="+configs.DB_PASSWORD+" sslmode=disable")
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
	db.AutoMigrate(&models.Sex{})
	db.AutoMigrate(&models.User{}).AddForeignKey("sex_id", "sexes(id)", "SET DEFAULT", "SET DEFAULT")
}
