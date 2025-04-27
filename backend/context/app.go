package context

import "backend/db"

// アプリケーション全体で利用されるデータを管理
type AppContext struct {
	DB db.DB
}

func NewAppContext(db db.DB) *AppContext {
	return &AppContext{DB: db}
}
