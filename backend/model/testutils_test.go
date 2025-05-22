package model

import (
	"backend/context"
	"backend/db"
)

// 新たなテスト用コンテキストを作成
func newTestContext(users []db.User, requests []db.Request, entries []db.Entry) *context.AppContext {
	return context.NewAppContext(db.NewMockDB(requests, users, entries), nil)
}

func mustNewDateTime(s string) DateTime {
	t, err := NewDateTime(s)
	if err != nil {
		panic(err)
	}
	return t
}

func mustNewDateOnly(s string) DateOnly {
	t, err := NewDateOnly(s)
	if err != nil {
		panic(err)
	}
	return t
}
