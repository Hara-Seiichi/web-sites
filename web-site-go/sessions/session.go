package sessions

import (
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
)

// ReadMiddleWare sessions.Defaultを呼び出す
func readMiddleWare(c *gin.Context) sessions.Session {
	session := sessions.Default(c)
	return session
}

// Manager セッションマネージャを定義するインターフェース
type SessionManager interface {

	// Start セッションの開始
	Start(r *gin.Engine)

	// Get セッションから値を取得 => 構造体に格納
	Get(c *gin.Context) SessionManager

	// GetLoginSession 構造体 => コンテキストに格納
	GetLoginSession(c *gin.Context)

	// Set 構造体を受け取る => セッションに各値を格納
	Set(c *gin.Context) error

	// Destroy セッションを削除 削除対象のセッションキーは構造体ごとに決まるため関数内で定義する
	Destroy(c *gin.Context) error

	// Certified 認証を確認する
	Certified(c *gin.Context) bool
}

// セッション保存
type LoginSession struct {
	userid          string
	userName        string
	isAuthenticated bool
}

// Start セッションの開始
func (l *LoginSession) Start(r *gin.Engine) {
	store := cookie.NewStore([]byte("secret"))
	r.Use(sessions.Sessions("LoginSession", store))
}

// Get セッションから値を取得 => 構造体に格納
func (l *LoginSession) Get(c *gin.Context) SessionManager {
	session := readMiddleWare(c)
	userid := session.Get("userid")
	if userid != "" && userid != nil {
		l.userid = userid.(string)
	} else {
		l.userid = ""
	}

	userName := session.Get("userName")
	if userName != "" && userName != nil {
		l.userName = userName.(string)
	} else {
		l.userName = ""
	}

	// idとnameが揃っている時は構造体に入れる（認証が正しく終わっている）
	if userid != "" && userid != nil && userName != "" && userName != nil {
		l.isAuthenticated = session.Get("isAuthenticated").(bool)
	} else {
		l.isAuthenticated = false
	}
	return l
}

// GetLoginSession 構造体 => コンテキストに格納
func (l *LoginSession) GetLoginSession(c *gin.Context) {
	l.Get(c)
	c.Set("userName", l.userName)
	c.Set("isAuthenticated", l.isAuthenticated)
}

// Set 構造体を受け取る => セッションに各値を格納
func (l *LoginSession) Set(c *gin.Context) error {
	session := readMiddleWare(c)
	session.Set("userid", c.GetString("userid"))
	session.Set("userName", c.GetString("userName"))
	session.Set("isAuthenticated", c.GetBool("isAuthenticated"))
	// Setしたセッション情報を保存
	if err := session.Save(); err != nil {
		return err
	}
	return nil
}

// Destroy セッションを削除 削除対象のセッションキーは構造体ごとに決まるため関数内で定義する
func (l *LoginSession) Destroy(c *gin.Context) error {
	session := readMiddleWare(c)
	keyList := [...]string{"userid", "userName", "isAuthenticated"}
	for _, v := range keyList {
		session.Delete(v)
	}
	// セッション情報の変更を保存
	if err := session.Save(); err != nil {
		return err
	}
	l.userid = ""
	l.userName = ""
	l.isAuthenticated = false
	c.Set("userid", "")
	c.Set("userName", "")
	c.Set("isAuthenticated", false)
	return nil
}

// Certified 認証を確認する
func (l *LoginSession) Certified(c *gin.Context) bool {
	l.Get(c)
	return l.isAuthenticated
}
