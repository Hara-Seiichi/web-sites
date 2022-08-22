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

func (pc UserController) List(c *gin.Context) {

	if c.Request.Method == http.MethodGet {
		sm.GetLoginSession(c)

		var ur repository.UserRepository
		users, err := ur.GetAll()
		if err != nil {
			c.HTML(http.StatusInternalServerError, "list.html", gin.H{})
			return
		}
		var isAuthenticated, _ = c.Get("isAuthenticated")
		var userName, _ = c.Get("userName")
		c.HTML(http.StatusOK, "list.html", gin.H{
			"isAuthenticated": isAuthenticated,
			"userName":        userName,
			"userid":          "",
			"name":            "",
			"age":             "",
			"sex":             0,
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
