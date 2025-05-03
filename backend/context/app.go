package context

import (
	"backend/db"

	"github.com/gorilla/sessions"
)

// TODO: 外部からアクセスできないようにする

// アプリケーション全体で利用されるデータを管理
type AppContext struct {
	db           db.DB
	sessionStore *sessions.CookieStore
}

// create new AppContext and set db and sessionStore
func NewAppContext(db db.DB, sessionStore *sessions.CookieStore) *AppContext {
	return &AppContext{db: db, sessionStore: sessionStore}
}

func (ctx *AppContext) GetDB() db.DB {
	return ctx.db
}

func (ctx *AppContext) GetSessionStore() *sessions.CookieStore {
	return ctx.sessionStore
}
