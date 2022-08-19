package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"path/filepath"

	"github.com/gin-contrib/multitemplate"
	"github.com/gin-gonic/gin"
)

func loadTemplates(templatesDir string) multitemplate.Renderer {
	r := multitemplate.NewRenderer()

	layout, err := filepath.Glob(templatesDir + "/layout/*.html")
	if err != nil {
		log.Fatal(err)
	}

	includes, err := filepath.Glob(templatesDir + "/includes/*.html")
	if err != nil {
		log.Fatal(err)
	}

	singes, err := filepath.Glob(templatesDir + "/*.html")
	if err != nil {
		panic(err.Error())
	}

	// Generate our templates map from our layout/ and includes/ directories
	for _, include := range includes {
		files := append(layout, include)
		r.AddFromFiles(filepath.Base(include), files...)
	}

	for _, sing := range singes {
		r.AddFromFiles(filepath.Base(sing), sing)
	}

	return r
}

func main() {
	router := gin.Default()

	// templateファイルの読込
	router.HTMLRender = loadTemplates("template")

	// 静的ファイルの配置場所(エイリアス,実際の置き場所)
	router.Static("/css", "./assets/css")
	router.Static("/js", "./assets/js")
	router.Static("/assets", "./assets")

	// ルーティングの設定
	router.GET("", singup)
	router.POST("", singup)
	router.GET("signin", signin)
	router.POST("signin", signin)
	router.GET("list", list)
	router.POST("search", list)
	router.GET("search", list)

	router.GET("signout", signout)
	// router.PUT("/somePut", putting)
	// router.DELETE("/someDelete", deleting)
	// router.PATCH("/somePatch", patching)
	// router.HEAD("/someHead", head)
	// router.OPTIONS("/someOptions", options)

	// サーバの起動
	router.Run()

	// シャッドダウンリスナーの設定
	// ※make(chan os.Signal, 1)これがCTRL + C の事らしい
	//   強制終了を安全に止めるぜ！って事か。。。
	var srv http.Server

	idleConnsClosed := make(chan struct{})
	go func() {
		sigint := make(chan os.Signal, 1)
		signal.Notify(sigint, os.Interrupt)
		<-sigint

		// We received an interrupt signal, shut down.
		if err := srv.Shutdown(context.Background()); err != nil {
			// Error from closing listeners, or context timeout:
			log.Printf("HTTP server Shutdown: %v", err)
		}
		close(idleConnsClosed)
	}()

	if err := srv.ListenAndServe(); err != http.ErrServerClosed {
		// Error starting or closing listener:
		log.Fatalf("HTTP server ListenAndServe: %v", err)
	}

	<-idleConnsClosed

}

// ///////////////////////////////////////////////////////////
// コンテキストを引数に取り、コンテキストの中に渡されるパラメータを
// 利用して処理を実行する。
// ///////////////////////////////////////////////////////////
// [gin.Context]の中身
// type Context struct {
// 	writermem responseWriter
// 	Request   *http.Request
// 	Writer    ResponseWriter

// 	Params   Params
// 	handlers HandlersChain
// 	index    int8
// 	fullPath string

// 	engine *Engine

// 	// Keys is a key/value pair exclusively for the context of each request.
//  ブラウザとやり取りする値を格納する。
// 	Keys map[string]interface{}

// 	// Errors is a list of errors attached to all the handlers/middlewares who used this context.
//  エラーの内容を格納する。
// 	Errors errorMsgs

// 	// Accepted defines a list of manually accepted formats for content negotiation.
//  手動で受け入れられる形式のリスト？
// 	Accepted []string

// 	// queryCache use url.ParseQuery cached the param query result from c.Request.URL.Query()
//  多分querystring？
// 	queryCache url.Values

//		// formCache use url.ParseQuery cached PostForm contains the parsed form data from POST, PATCH,
//		// or PUT body parameters.
//	 formの入力値。（name属性と値のペアかな？）
//		formCache url.Values
//	}
func singup(c *gin.Context) {
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

func signin(c *gin.Context) {

	if c.Request.Method == http.MethodGet {
		c.HTML(http.StatusOK, "signin.html", gin.H{
			"message": "",
		})
		return
	}

	// POSTならDBの確認とかする
	c.Redirect(http.StatusFound, "/list")

}

func list(c *gin.Context) {
	c.HTML(http.StatusOK, "list.html", gin.H{
		"message":          "",
		"is_authenticated": true,
		"user":             "hoge",
		"userid":           "",
		"name":             "",
		"age":              "",
		"sex":              0,
		"items":            "",
	})
}

func signout(c *gin.Context) {

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
