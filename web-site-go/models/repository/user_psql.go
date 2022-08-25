package repository

import (
	"strconv"
	"web-site-go/db"
	"web-site-go/models"
)

// Userテーブルの操作を受け持つ
type UserRepository struct{}

// Userテーブルの構造体
type User models.User

type UserProfile struct {
	PK       uint
	Userid   string
	Name     string
	Age      int
	Sex      string
	Password string
}

// Userの一覧を取得
func (ur UserRepository) GetAll() ([]UserProfile, error) {
	db := db.GetDB()
	var u []models.User
	var r []UserProfile
	var p *UserProfile

	if err := db.Model(&User{}).Preload("Sex").Find(&u).Error; err != nil {
		return nil, err
	}

	for _, v := range u {
		p = new(UserProfile)
		p.PK = v.ID
		p.Userid = v.Userid
		p.Name = v.Name
		p.Age = v.Age
		p.Sex = v.Sex.Name
		r = append(r, *p)
	}

	return r, nil
}

// idで絞り込んでユーザを1人取得
func (ur UserRepository) GetByID(id string) (models.User, error) {

	idInt, _ := strconv.Atoi(id)

	db := db.GetDB()
	var u models.User
	if err := db.Model(&User{}).Where("id = ?", idInt).Preload("Sex").Find(&u).Error; err != nil {
		return u, err
	}

	return u, nil
}

// userid,passwordで絞り込んでユーザを1人取得
func (ur UserRepository) GetUserAuthority(userid string, password string) (models.User, error) {

	db := db.GetDB()
	var u models.User
	// .Preload("Sex")の中に書くカラムは、参照先のモデル型を定義しているカラム名
	if err := db.Model(&User{}).Where("userid = ? and password = ?", userid, password).Preload("Sex").First(&u).Error; err != nil {
		return u, err
	}

	return u, nil
}

// idで絞り込んでユーザを1人取得
func (ur UserRepository) UserSearch(userid string, name string, age string, sex string) ([]UserProfile, error) {

	db := db.GetDB()
	var u []models.User
	var r []UserProfile
	var p *UserProfile
	ageInt, _ := strconv.ParseInt(age, 10, 32)

	chain := db.Model(&User{}).Where("")
	if userid != "" {
		chain = chain.Where("userid like ?", "%"+userid+"%")
	}
	if name != "" {
		chain = chain.Where("name like ?", "%"+name+"%")
	}
	if age != "" {
		chain = chain.Where("age = ?", ageInt)
	}
	if sex != "" {
		chain = chain.Where("sex_id = ?", sex)
	}
	chain.LogMode(true)
	if err := chain.Preload("Sex").Find(&u).Error; err != nil {
		chain.LogMode(false)
		return nil, err
	}
	chain.LogMode(false)

	for _, v := range u {
		p = new(UserProfile)
		p.PK = v.ID
		p.Userid = v.Userid
		p.Name = v.Name
		p.Age = v.Age
		p.Sex = v.Sex.Name
		r = append(r, *p)
	}
	return r, nil
}

// サインアップでデータを登録する
func (ur UserRepository) SignupUser(userid string, password string) (User, error) {
	var u User
	u.Userid = userid
	u.Name = userid
	u.Password = password
	if u, err := createModel(&u); err != nil {
		return u, err
	}
	return u, nil
}

// Userを作成するデータを登録する
func (ur UserRepository) CreateUser(form *UserProfile) (User, error) {

	// 性別を数値に変換。未選択""はチェック済みなので変換エラーは見ない
	sexInt, _ := strconv.ParseInt(form.Sex, 10, 32)

	var u User
	u.Userid = form.Userid
	u.Name = form.Name
	u.Age = form.Age
	u.SexID = uint(sexInt)
	u.Password = form.Password

	if u, err := createModel(&u); err != nil {
		return u, err
	}
	return u, nil
}

// CreateModel is create User model
func createModel(u *User) (User, error) {
	db := db.GetDB()
	if err := db.Create(&u).Error; err != nil {
		return *u, err
	}
	return *u, nil
}

// Userを更新する
func (ur UserRepository) UpdateByID(id string, form *UserProfile) (models.User, error) {

	var u models.User
	if u, err := ur.GetByID(id); err != nil {
		return u, err
	}

	// 性別を数値に変換。未選択""はチェック済みなので変換エラーは見ない
	sexInt, _ := strconv.ParseInt(form.Sex, 10, 32)
	u.Userid = form.Userid
	u.Name = form.Name
	u.Age = form.Age
	u.SexID = uint(sexInt)

	idInt, _ := strconv.ParseInt(id, 10, 32)
	u.ID = uint(idInt)
	db := db.GetDB()
	db.Save(u)

	return u, nil
}

// Userを削除する
func (ur UserRepository) DeleteByID(id string) error {

	var u models.User
	if err := db.GetDB().Where("id = ?", id).Delete(&u).Error; err != nil {
		return err
	}

	return nil
}
