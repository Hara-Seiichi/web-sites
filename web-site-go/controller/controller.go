package controller

import (
	"net/http"
	"web-site-go/models/repository"
	SM "web-site-go/sessions"

	"github.com/gin-gonic/gin"
)

var sm SM.SessionManager = &SM.LoginSession{}

// Userに関する操作をするクラス
type UserController struct{}

func (pc UserController) Singup(c *gin.Context) {

	if c.Request.Method == http.MethodGet {
		c.HTML(http.StatusOK, "signup.html", gin.H{
			"message":  "",
			"userid":   "",
			"password": "",
		})
		return
	}

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

func (pc UserController) Signin(c *gin.Context) {

	if c.Request.Method == http.MethodGet {
		c.HTML(http.StatusOK, "signin.html", gin.H{
			"message":  "",
			"userid":   "",
			"password": "",
		})
		return
	}

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

	var ur repository.UserRepository
	if _, err := ur.GetUserAuthority(userid, pw); err != nil {
		c.HTML(http.StatusBadRequest, "signin.html", gin.H{
			"message":  "That user does not exist.",
			"userid":   userid,
			"password": pw,
		})
		return
	}

	c.Set("userid", userid)
	c.Set("userName", userid)
	c.Set("isAuthenticated", true)
	sm.Set(c)
	c.Redirect(http.StatusMovedPermanently, "/app/list")
}

type M_USER struct {
	Userid string
	Name   string
	Age    int
	Sex    string
}
type UserList []M_USER

func (pc UserController) List(c *gin.Context) {
	var users UserList
	var user M_USER
	user.Userid = "hoge"
	user.Name = "hoge"
	user.Age = 35
	user.Sex = "-"
	users = append(users, user)
	// formの値取得
	// c.Request.FormValue("passwd")
	c.HTML(http.StatusOK, "list.html", gin.H{
		"message":          "",
		"is_authenticated": true,
		"user":             "hoge",
		"userid":           "",
		"name":             "",
		"age":              "",
		"sex":              0,
		"items":            users,
	})
}

func (pc UserController) Signout(c *gin.Context) {

	//セッションからデータを破棄する
	sm.Destroy(c)
	c.HTML(http.StatusBadRequest, "signin.html", gin.H{
		"message":  "",
		"userid":   "",
		"password": "",
	})
}
