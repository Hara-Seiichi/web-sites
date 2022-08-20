package controller

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// Userに関する操作をするクラス
type UserController struct{}

func (pc UserController) Singup(c *gin.Context) {
	if c.Request.Method == http.MethodGet {
		c.HTML(http.StatusOK, "signup.html", gin.H{
			"message": "",
		})
		return
	}

	c.JSONP(http.StatusOK, gin.H{
		"message": "ok",
		"data":    "singupのPOSTの時は登録する予定",
	})
}

func (pc UserController) Signin(c *gin.Context) {

	if c.Request.Method == http.MethodGet {
		c.HTML(http.StatusOK, "signin.html", gin.H{
			"message": "",
		})
		return
	}

	// POSTならDBの確認とかする
	c.Redirect(http.StatusFound, "/list")

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

	/*
		DB操作など
	*/

	// err := ""
	// if err != "" {
	// 	c.String(http.StatusInternalServerError, "Server Error")
	// 	return
	// }
	c.HTML(http.StatusOK, "signin.html", gin.H{
		"message": "",
	})
}
