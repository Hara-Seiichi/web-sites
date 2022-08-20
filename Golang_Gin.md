# ２．ユーザ管理WEBアプリ-GOlang & Gin

間違いや古い情報も含む可能性があるので作業を進める時は都度、確認する事。

### 利用したアーキテクチャの構成
リンクとバージョンは必ずしも一致しない
* [PostgreSQL 10.21](https://www.postgresql.org/download/)
* [GOlang 1.19](https://go.dev/doc/)
  * [ドキュメント](https://pkg.go.dev/std)
* [Gin](https://pkg.go.dev/github.com/gin-gonic/gin)
* [UIkit 3.14.3](https://getuikit.com/)
* [VSCode 1.70.0](https://azure.microsoft.com/ja-jp/products/visual-studio-code/)

### 環境構築の簡易手順
参考になるサイトが沢山あるので順序のまとめだけ。
1. PostgreSQLの導入[(参考サイト)](https://marunaka-blog.com/postgresql-download-install/3704/)

2. GOlangの導入[(参考サイト)](https://go.dev/dl/)
   
   特に難しいことは無く、installerを落として実行。

   ※Path環境変数にbinフォルダ自動設定。

   ※GOPATH環境変数は自動作成。`C:\Users\<username>\go`

   バージョン確認  
   `go version`

   環境変数確認  
   `go env GOROOT`→インストールフォルダ表示  
   `go env GOPATH`→ `C:\Users\<username>\go` 表示

3. VSCodeの導入
   * 「GO」で検索して拡張機能インストール
   * VScode再起動
   * `Ctrl + Shift + p` > コマンドパレット > gotools  
      基本機能のtoolが出てくるので全て導入する  
      何度かインストールの確認通知が出てるので全てOKで進める  
      完了するとにっこりしてくれる  
     `All tools successfully installed. You are ready to Go. :)`

     インストールされたtoolは`GOPATH\bin`に作成される

   * 設定ファイルの修正  
     `歯車（設定）＞拡張機能＞GO＞setting.json`を編集する

         "go.alternateTools": {
            "editor.tabSize": 4, // 公式のサイズは不明。2か4  
            "editor.insertSpaces": false, // Go公式はスペースでなくタブ  
            "editor.formatOnSave": true, // ファイル保存時にフォーマット  
            "editor.defaultFormatter": "golang.go" // 変える  
          }

4. プロジェクトのフォルダを作成
   
   プロジェクトフォルダを用意してPowerShellで移動  
   アプリの初期化コマンド実行(go.modが作成される)  
   `go mod init web-site-go`

5. Ginのインストール[(参考サイト)](https://github.com/gin-gonic/gin)  
   * プロジェクトで使用するライブラリをインストールするときは`go get`コマンド
   * Goをインストールした環境全体で使うライブラリは`go install`コマンド  
     ※インストール先は`$%GOBIN%`の指定による。  
     ※存在しなければ`$%GOPATH%\bin`になる。
   
   Ginのインストールコマンドを実行（go.modが書きかわる）  
   `go get -u github.com/gin-gonic/gin`

6. UIkitの導入
   
   ダウンロードした３ファイルをプロジェクトの`assets`フォルダに入れる。

   配置先フォルダはcode参照。  
   設定はsetting.py参照。

### 【未編集です】フレームワーク動作機序の簡易図
今回の独学で大まかに理解した内容
```mermaid
sequenceDiagram
    participant Chrome
    Chrome->>urls.py of PJ: URL
    urls.py of PJ->>urls.py of APP: URL
    participant Template.html
    participant forms.py
    urls.py of APP->>views.py: class or function
    Note right of forms.py: POST param
    Note right of forms.py: Validation
    views.py->>models.py: DB操作を依頼
    models.py->>DB: データ操作実行
    DB-->>models.py: 結果を返却
    models.py-->>views.py: 結果を返却
    views.py-->>Chrome: 結果を返却
    Note left of views.py: POST param
    participant DB
    Note left of Template.html: 結果をHTMLに埋込
```

### 【未編集です】作成機能の概要
```mermaid
sequenceDiagram
    participant Singup
    Singup->>Singin: Singup success
    Singup->>Singin: Mutual link
    Singin->>Singup: Mutual link
    Singin->>List: Signin success
    List->>List: Search
    List->>Create: Forward
    Create-->>List: Create success
    Create-->>List: Back
    List->>Update: Forward query string pk
    Update-->>List: Update success
    Update-->>List: Back
    List->>Delete: Forward query string pk
    Delete-->>List: Delete success
    Delete-->>List: Back
    List->>Detail: Forward query string pk
    Detail-->>List: Back
    List-->>Singin: Singout
    Create-->>Singin: Singout
    Update-->>Singin: Singout
    Delete-->>Singin: Singout
    Detail-->>Singin: Singout
```

## メモ
### Goコマンド
* `go run ./cmd/jisyo/main.go` サーバの起動
  
[Goコマンド参考](https://pkg.go.dev/cmd/go)
