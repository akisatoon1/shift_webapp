package context

import (
	"backend/db"

	"github.com/gorilla/sessions"
)

// アプリケーション全体で利用されるデータを管理
type AppContext struct {
	DB     db.DB
	Cookie *sessions.CookieStore
}

func NewAppContext(db db.DB, cookie *sessions.CookieStore) *AppContext {
	return &AppContext{DB: db, Cookie: cookie}
}
