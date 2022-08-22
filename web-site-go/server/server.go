package server

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"path/filepath"
	"web-site-go/configs"
	"web-site-go/controller"
	SM "web-site-go/sessions"

	"github.com/gin-contrib/multitemplate"
	"github.com/gin-gonic/gin"
)

var sm SM.SessionManager = &SM.LoginSession{}

// サーバ初期化
func Init() {
	r := router()
	r.HTMLRender = loadTemplates()
	setStatic(r)

	r.Run()

	//server起動より後に呼ぶ（listenerだからね）
	shutdownlistener()
}

// ルーティング
func router() *gin.Engine {
	r := gin.Default()

	sm.Start(r)

	ctrl := controller.UserController{}
	// ルーティングの設定
	r.GET("", ctrl.Singup)
	r.POST("", ctrl.Singup)
	r.GET("signin", ctrl.Signin)
	r.POST("signin", ctrl.Signin)

	// ログインチェックをする画面は「app」を使う
	app := r.Group("/app")
	app.Use(sessionCheck())
	{
		app.GET("list", ctrl.List)
		app.POST("search", ctrl.List)
		app.GET("search", ctrl.List)
	}
	r.GET("signout", ctrl.Signout)

	// router.PUT("/somePut", putting)
	// router.DELETE("/someDelete", deleting)
	// router.PATCH("/somePatch", patching)
	// router.HEAD("/someHead", head)
	// router.OPTIONS("/someOptions", options)

	return r
}

// テンプレートファイルの読込
func loadTemplates() multitemplate.Renderer {
	r := multitemplate.NewRenderer()

	// baseのtemplate
	layout, err := filepath.Glob(configs.BASE_FILE_PATH + "/*.html")
	if err != nil {
		log.Fatal(err)
	}

	// 各画面のコンテンツ
	includes, err := filepath.Glob(configs.CONTENTS_FILE_PATH + "/*.html")
	if err != nil {
		log.Fatal(err)
	}

	// サインアップ、サインイン
	singes, err := filepath.Glob(configs.TEMPLATE_BASE + "/*.html")
	if err != nil {
		panic(err.Error())
	}

	// baseファイルと各画面のcontentsをセットでレンダラーに登録
	for _, include := range includes {
		files := append(layout, include)
		r.AddFromFiles(filepath.Base(include), files...)
	}

	// サインアップ、サインインは単体でそれぞれ登録
	for _, sing := range singes {
		r.AddFromFiles(filepath.Base(sing), sing)
	}

	return r
}

func setStatic(r *gin.Engine) {
	// 静的ファイルの配置場所(エイリアス,実際の置き場所)
	r.Static("/css", configs.CSS_FILE_PATH)
	r.Static("/js", configs.JS_FILE_PATH)
	r.Static("/assets", configs.ASSETS_BASE)
}

// ログインセッションの確認
func sessionCheck() gin.HandlerFunc {
	return func(c *gin.Context) {
		sm.Get(c)
		// 未承認の場合は終了
		if !sm.Certified(c) {
			sm.Destroy(c)
			c.Redirect(http.StatusMovedPermanently, "/signin")
			c.Abort() // これがないと続けて処理されてしまう
			return
		}
		// sm.Set(c)
		c.Next()
	}
}

// シャッドダウンリスナーの設定
func shutdownlistener() {
	// ※make(chan os.Signal, 1)これがCTRL + C の事らしい
	//   強制終了を安全に止めるぜ！って事か。。。
	var srv http.Server

	idleConnsClosed := make(chan struct{})
	go func() {
		sigint := make(chan os.Signal, 1)
		signal.Notify(sigint, os.Interrupt)
		<-sigint

		// 割り込み信号を受信しました。シャットダウンします
		if err := srv.Shutdown(context.Background()); err != nil {
			// リスナーを閉じる際のエラー、またはコンテキスト タイムアウト:
			log.Printf("HTTP server Shutdown: %v", err)
		}
		close(idleConnsClosed)
	}()

	if err := srv.ListenAndServe(); err != http.ErrServerClosed {
		// リスナーの開始または終了中にエラーが発生しました:
		log.Fatalf("HTTP server ListenAndServe: %v", err)
	}

	<-idleConnsClosed
}
