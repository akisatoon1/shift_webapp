package model

import (
	"backend/context"
	"backend/db"
	"reflect"
	"testing"
)

// 新たなテスト用コンテキストを作成
func newTestContext(users []db.User, requests []db.Request, entries []db.Entry, submissions []db.Submission) *context.AppContext {
	return context.NewAppContext(db.NewMockDB(requests, users, entries, submissions), nil)
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

func assert(t *testing.T, got, want interface{}) {
	if !reflect.DeepEqual(got, want) {
		t.Errorf("got %v, want %v", got, want)
	}
}
