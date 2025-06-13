package usecase

import (
	"backend/context"
	"backend/db"
	"backend/domain"
	"encoding/json"
	"reflect"
	"testing"
)

// 新たなテスト用コンテキストを作成
func newTestContext(users []db.User, requests []db.Request, entries []db.Entry, submissions []db.Submission) *context.AppContext {
	return context.NewAppContext(db.NewMockDB(requests, users, entries, submissions), nil)
}

func mustNewDateTime(s string) domain.DateTime {
	t, err := domain.NewDateTime(s)
	if err != nil {
		panic(err)
	}
	return t
}

func mustNewDateOnly(s string) domain.DateOnly {
	t, err := domain.NewDateOnly(s)
	if err != nil {
		panic(err)
	}
	return t
}

func assert(t *testing.T, got, want interface{}) {
	if !reflect.DeepEqual(got, want) {
		// JSON形式で構造体を出力するためのエンコーディング
		gotJSON, _ := json.MarshalIndent(got, "", "  ")
		wantJSON, _ := json.MarshalIndent(want, "", "  ")

		t.Errorf("\nGOT:\n%s\n\nWANT:\n%s", gotJSON, wantJSON)
	}
}
