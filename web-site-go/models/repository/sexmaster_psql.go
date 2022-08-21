package repository

import (
	"web-site-go/db"
	"web-site-go/models"
)

// Userテーブルの操作を受け持つ
type SexRepository struct{}

// Userテーブルの構造体
type Sex models.Sex

type SexInfo struct {
	id   int
	Name string
}

// GetByID is get a Sex by ID
func (_ SexRepository) GetAll() ([]Sex, error) {
	db := db.GetDB()
	var sm []Sex
	if err := db.Table("sex_masters").Select("id, name").Scan(&sm).Error; err != nil {
		return nil, err
	}

	return sm, nil
}
