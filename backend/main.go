package main

import (
	"backend/context"
	"backend/db"
	"backend/router"
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"
)

func main() {
	// .envファイルの読み込み
	err := godotenv.Load()
	if err != nil {
		log.Println(".envファイルが見つかりませんでした")
	}

	// 環境変数から設定値を取得
	dbPath := os.Getenv("DB_PATH")
	if dbPath == "" {
		log.Fatal("DB_PATHが設定されていません")
	}
	port := os.Getenv("PORT")
	if port == "" {
		log.Fatal("PORTが設定されていません")
	}

	// DBの初期化
	db, err := db.NewSqlite3DB(dbPath)
	if err != nil {
		log.Fatal("DBの接続に失敗しました: " + err.Error())
	}
	log.Println("DBの接続に成功しました")
	defer db.Close()

	// アプリケーション全体で使うデータを管理するコンテキストを作成
	appCtx := context.NewAppContext(db)

	// ルーティングの設定
	router.Routes(appCtx)

	// サーバーの起動
	log.Println("サーバーを起動します: http://localhost:" + port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
