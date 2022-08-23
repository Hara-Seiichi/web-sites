package controller

import (
	"net/http"
	"strconv"
	"web-site-go/models/repository"
	SM "web-site-go/sessions"

	"github.com/gin-gonic/gin"
)

var sm SM.SessionManager = &SM.LoginSession{}

// Userに関する操作をするクラス
type UserController struct{}

// //////////////////////////////////////////////
// サインアップ
// //////////////////////////////////////////////
func (pc UserController) Singup(c *gin.Context) {

	if c.Request.Method == http.MethodGet {
		c.HTML(http.StatusOK, "signup.html", gin.H{
			"message":  "",
			"userid":   "",
			"password": "",
		})
		return
	}

	// 入力チェック（Validationありそうだけどとりあえず）
	var userid, pw string
	userid = c.Request.FormValue("userid")
	pw = c.Request.FormValue("password")
	if userid == "" || pw == "" {
		c.HTML(http.StatusBadRequest, "signup.html", gin.H{
			"message":  "Please enter both your userid and password.",
			"userid":   userid,
			"password": pw,
		})
		return
	}

	// データを作成
	var u repository.UserRepository
	if _, err := u.SignupUser(userid, pw); err != nil {
		c.HTML(http.StatusBadRequest, "signup.html", gin.H{
			"message":  "The user is already registered.",
			"userid":   userid,
			"password": pw,
		})
		return
	}

	c.Redirect(http.StatusFound, "/signin")
}

// //////////////////////////////////////////////
// サインイン
// //////////////////////////////////////////////
func (pc UserController) Signin(c *gin.Context) {

	if c.Request.Method == http.MethodGet {
		c.HTML(http.StatusOK, "signin.html", gin.H{
			"message":  "",
			"userid":   "",
			"password": "",
		})
		return
	}

	// 入力チェック（Validationありそうだけどとりあえず）
	var userid, pw string
	userid = c.Request.FormValue("userid")
	pw = c.Request.FormValue("password")
	if userid == "" || pw == "" {
		c.HTML(http.StatusBadRequest, "signin.html", gin.H{
			"message":  "Please enter both your userid and password.",
			"userid":   userid,
			"password": pw,
		})
		return
	}

	// 登録されているか確認
	var ur repository.UserRepository
	if _, err := ur.GetUserAuthority(userid, pw); err != nil {
		c.HTML(http.StatusBadRequest, "signin.html", gin.H{
			"message":  "That user does not exist.",
			"userid":   userid,
			"password": pw,
		})
		return
	}

	// セッションにログイン情報を残して一覧画面へ遷移
	c.Set("userid", userid)
	c.Set("userName", userid)
	c.Set("isAuthenticated", true)
	sm.Set(c)
	c.Redirect(http.StatusMovedPermanently, "/app/list")
}

// //////////////////////////////////////////////
// サインアウト
// //////////////////////////////////////////////
func (pc UserController) Signout(c *gin.Context) {

	//セッションからデータを破棄する
	sm.Destroy(c)
	sm = &SM.LoginSession{}
	c.HTML(http.StatusBadRequest, "signin.html", gin.H{
		"message":  "",
		"userid":   "",
		"password": "",
	})
}

// //////////////////////////////////////////////
// 一覧表示・検索
// //////////////////////////////////////////////
func (pc UserController) List(c *gin.Context) {

	sm.GetLoginSession(c)
	var isAuthenticated, _ = c.Get("isAuthenticated")
	var userName, _ = c.Get("userName")
	// GETリクエストの場合は一覧表示
	if c.Request.Method == http.MethodGet {

		var ur repository.UserRepository
		users, err := ur.GetAll()
		if err != nil {
			c.HTML(http.StatusInternalServerError, "list.html", gin.H{})
			return
		}
		c.HTML(http.StatusOK, "list.html", gin.H{
			"isAuthenticated": isAuthenticated,
			"userName":        userName,
			"userid":          "",
			"name":            "",
			"age":             "",
			"sex":             nil,
			"items":           users,
		})
		return
	}

	// formの値取得
	// c.Request.FormValue("passwd")
	// c.HTML(http.StatusOK, "list.html", gin.H{
	// 	"userid": "",
	// 	"name":   "",
	// 	"age":    "",
	// 	"sex":    0,
	// 	"items":  users,
	// })
}

// //////////////////////////////////////////////
// User作成
// //////////////////////////////////////////////
func (pc UserController) Create(c *gin.Context) {

	sm.GetLoginSession(c)
	var isAuthenticated, _ = c.Get("isAuthenticated")
	var userName, _ = c.Get("userName")

	if c.Request.Method == http.MethodGet {
		c.HTML(http.StatusOK, "create.html", gin.H{
			"message":         "",
			"isAuthenticated": isAuthenticated,
			"userName":        userName,
			"userid":          "",
			"name":            "",
			"age":             "",
			"sex":             nil,
		})
		return
	}

	// 入力チェック（Validationありそうだけどとりあえず）
	valid, userid, name, age, sex, pw := isValidInput(c)
	if !valid {
		c.HTML(http.StatusBadRequest, "create.html", gin.H{
			"message":         "Please enter all of the items.",
			"isAuthenticated": isAuthenticated,
			"userName":        userName,
			"userid":          userid,
			"name":            name,
			"age":             age,
			"sex":             sex,
		})
		return
	}

	// 年齢を数値に変換
	ageInt, err := strconv.ParseInt(age, 10, 32)
	if err != nil || ageInt < 0 {
		c.HTML(http.StatusBadRequest, "create.html", gin.H{
			"message":         "Enter your age as an integer greater than or equal to 0.",
			"isAuthenticated": isAuthenticated,
			"userName":        userName,
			"userid":          userid,
			"name":            name,
			"age":             age,
			"sex":             sex,
		})
		return
	}
	var ur repository.UserRepository
	var user = repository.UserProfile{
		Userid:   userid,
		Name:     name,
		Age:      int(ageInt),
		Sex:      sex,
		Password: pw,
	}

	if _, err := ur.CreateUser(&user); err != nil {
		c.HTML(http.StatusBadRequest, "create.html", gin.H{
			"message":         "Sorry. Failed to create.",
			"isAuthenticated": isAuthenticated,
			"userName":        userName,
			"userid":          userid,
			"name":            name,
			"age":             age,
			"sex":             sex,
		})
		return
	}

	// 登録成功時は一覧画面へ
	c.Redirect(http.StatusMovedPermanently, "/app/list")
}

func isValidInput(c *gin.Context) (bool, string, string, string, string, string) {
	var userid, name, age, sex, pw string
	userid = c.Request.FormValue("userid")
	name = c.Request.FormValue("name")
	age = c.Request.FormValue("age")
	sex = c.Request.FormValue("sex")
	pw = c.Request.FormValue("userid") // TODO とりあえず初期passwordはidと同じで登録。。。

	// どれか一つでも入力がない場合はエラー
	if userid == "" || name == "" || age == "" || sex == "" || pw == "" {
		return false, userid, name, age, sex, pw
	}
	return true, userid, name, age, sex, pw
}

// //////////////////////////////////////////////
// User詳細
// //////////////////////////////////////////////
func (pc UserController) Detail(c *gin.Context) {

	sm.GetLoginSession(c)
	var isAuthenticated, _ = c.Get("isAuthenticated")
	var userName, _ = c.Get("userName")
	var pk = c.Param("PK") // QueryStringにする場合は「c.Query("PK")」

	var ur repository.UserRepository
	users, err := ur.GetByID(pk)
	if err != nil {
		c.HTML(http.StatusBadRequest, "detail.html", gin.H{
			"isAuthenticated": isAuthenticated,
			"userName":        userName,
			"userid":          "",
			"name":            "",
			"age":             "",
			"sex":             "",
		})
		return
	}

	c.HTML(http.StatusOK, "detail.html", gin.H{
		"isAuthenticated": isAuthenticated,
		"userName":        userName,
		"userid":          users.Userid,
		"name":            users.Name,
		"age":             users.Age,
		"sex":             users.Sex.Name,
	})
}

// //////////////////////////////////////////////
// User更新
// //////////////////////////////////////////////
func (pc UserController) Update(c *gin.Context) {

	sm.GetLoginSession(c)
	var isAuthenticated, _ = c.Get("isAuthenticated")
	var userName, _ = c.Get("userName")
	var pk = c.Param("PK") // QueryStringにする場合は「c.Query("PK")」

	var ur repository.UserRepository
	if c.Request.Method == http.MethodGet {
		users, err := ur.GetByID(pk)
		if err != nil {
			c.HTML(http.StatusBadRequest, "update.html", gin.H{
				"message":         "That user does not exist.",
				"isAuthenticated": isAuthenticated,
				"userName":        userName,
				"PK":              pk,
				"userid":          "",
				"name":            "",
				"age":             "",
				"sex":             nil,
			})
			return
		}

		c.HTML(http.StatusOK, "update.html", gin.H{
			"message":         "",
			"isAuthenticated": isAuthenticated,
			"userName":        userName,
			"PK":              pk,
			"userid":          users.Userid,
			"name":            users.Name,
			"age":             users.Age,
			"sex":             users.Sex.ID,
		})
		return
	}

	// 入力チェック（Validationありそうだけどとりあえず）
	valid, userid, name, age, sex, pw := isValidInput(c)
	if !valid {
		c.HTML(http.StatusBadRequest, "update.html", gin.H{
			"message":         "Please enter all of the items.",
			"isAuthenticated": isAuthenticated,
			"userName":        userName,
			"PK":              pk,
			"userid":          userid,
			"name":            name,
			"age":             age,
			"sex":             sex,
		})
		return
	}

	// 年齢を数値に変換
	ageInt, err := strconv.ParseInt(age, 10, 32)
	if err != nil || ageInt < 0 {
		c.HTML(http.StatusBadRequest, "update.html", gin.H{
			"message":         "Enter your age as an integer greater than or equal to 0.",
			"isAuthenticated": isAuthenticated,
			"userName":        userName,
			"PK":              pk,
			"userid":          userid,
			"name":            name,
			"age":             age,
			"sex":             sex,
		})
		return
	}

	var user = repository.UserProfile{
		Userid:   userid,
		Name:     name,
		Age:      int(ageInt),
		Sex:      sex,
		Password: pw,
	}

	if _, err := ur.UpdateByID(pk, &user); err != nil {
		c.HTML(http.StatusBadRequest, "update.html", gin.H{
			"message":         "Sorry. Failed to update.",
			"isAuthenticated": isAuthenticated,
			"userName":        userName,
			"PK":              pk,
			"userid":          userid,
			"name":            name,
			"age":             age,
			"sex":             sex,
		})
		return
	}

	// 更新成功時は一覧画面へ
	c.Redirect(http.StatusMovedPermanently, "/app/list")

}

// //////////////////////////////////////////////
// User削除
// //////////////////////////////////////////////
func (pc UserController) Delete(c *gin.Context) {

	sm.GetLoginSession(c)
	var isAuthenticated, _ = c.Get("isAuthenticated")
	var userName, _ = c.Get("userName")
	var pk = c.Param("PK") // QueryStringにする場合は「c.Query("PK")」

	var ur repository.UserRepository
	if c.Request.Method == http.MethodGet {
		users, err := ur.GetByID(pk)
		if err != nil {
			c.HTML(http.StatusBadRequest, "delete.html", gin.H{
				"message":         "That user does not exist.",
				"isAuthenticated": isAuthenticated,
				"userName":        userName,
				"PK":              pk,
				"userid":          "",
				"name":            "",
				"age":             "",
				"sex":             nil,
			})
			return
		}

		c.HTML(http.StatusOK, "delete.html", gin.H{
			"message":         "",
			"isAuthenticated": isAuthenticated,
			"userName":        userName,
			"PK":              pk,
			"userid":          users.Userid,
			"name":            users.Name,
			"age":             users.Age,
			"sex":             users.Sex.Name,
		})
		return
	}

	if err := ur.DeleteByID(pk); err != nil {
		c.HTML(http.StatusBadRequest, "update.html", gin.H{
			"message":         "Sorry. Failed to delete.",
			"isAuthenticated": isAuthenticated,
			"userName":        userName,
			"PK":              pk,
			"userid":          "",
			"name":            "",
			"age":             "",
			"sex":             "",
		})
		return
	}

	// 削除成功時は一覧画面へ
	c.Redirect(http.StatusMovedPermanently, "/app/list")
}
