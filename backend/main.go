package main

import (
	"backend/context"
	"backend/db"
	"backend/router"
	"backend/test"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/sessions"
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

	mode := os.Getenv("MODE")
	var database db.DB

	if mode == "test" {
		log.Println("テストモードで起動します(mock DB使用)")
		database = db.NewMockDB(test.MockRequests, test.MockUsers, test.MockEntries)
	} else {
		log.Println("本番モードで起動します(SQLite3使用)")
		// DBの初期化
		sqliteDB, err := db.NewSqlite3DB(dbPath)
		if err != nil {
			log.Fatal("DBの接続に失敗しました: " + err.Error())
		}
		log.Println("DBの接続に成功しました")
		database = sqliteDB
		defer sqliteDB.Close()
	}

	// セッションの初期化
	cookie := sessions.NewCookieStore([]byte(os.Getenv("SESSION_KEY")))

	// アプリケーション全体で使うデータを管理するコンテキストを作成
	appCtx := context.NewAppContext(database, cookie)

	// ルーティングの設定
	mux := http.NewServeMux()
	router.Routes(mux, appCtx)

	// サーバーの起動
	log.Println("サーバーを起動します: http://localhost:" + port)
	log.Fatal(http.ListenAndServe(":"+port, mux))
}
