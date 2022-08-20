package server

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"path/filepath"
	"web-site-go/controller"

	"github.com/gin-contrib/multitemplate"
	"github.com/gin-gonic/gin"
)

// サーバ初期化
func Init() {
	r := router()
	r.HTMLRender = loadTemplates("template")
	setStatic(r)

	r.Run()

	//server起動より後に呼ぶ（listenerだからね）
	shutdownlistener()
}

// ルーティング
func router() *gin.Engine {
	r := gin.Default()

	ctrl := controller.UserController{}
	// ルーティングの設定
	r.GET("", ctrl.Singup)
	r.POST("", ctrl.Singup)
	r.GET("signin", ctrl.Signin)
	r.POST("signin", ctrl.Signin)
	r.GET("list", ctrl.List)
	r.POST("search", ctrl.List)
	r.GET("search", ctrl.List)

	r.GET("signout", ctrl.Signout)
	// router.PUT("/somePut", putting)
	// router.DELETE("/someDelete", deleting)
	// router.PATCH("/somePatch", patching)
	// router.HEAD("/someHead", head)
	// router.OPTIONS("/someOptions", options)

	return r
}

// テンプレートファイルの読込
func loadTemplates(templatesDir string) multitemplate.Renderer {
	r := multitemplate.NewRenderer()

	// baseのtemplate
	layout, err := filepath.Glob(templatesDir + "/layout/*.html")
	if err != nil {
		log.Fatal(err)
	}

	// 各画面のコンテンツ
	includes, err := filepath.Glob(templatesDir + "/includes/*.html")
	if err != nil {
		log.Fatal(err)
	}

	// サインアップ、サインイン
	singes, err := filepath.Glob(templatesDir + "/*.html")
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
	r.Static("/css", "./assets/css")
	r.Static("/js", "./assets/js")
	r.Static("/assets", "./assets")
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
