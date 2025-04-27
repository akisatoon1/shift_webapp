package model

import (
	"backend/context"
	"backend/db"
	"errors"
	"reflect"
	"testing"
	"time"
)

// モック用のDB構造体
type mockDB struct {
	Requests []db.Request
	Users    []db.User
	Entries  []db.Entry
}

func (m *mockDB) GetRequests() ([]db.Request, error) {
	return m.Requests, nil
}

func (m *mockDB) GetUserByID(id int) (db.User, error) {
	for _, user := range m.Users {
		if user.ID == id {
			return user, nil
		}
	}
	return db.User{}, errors.New("user not found")
}

func (m *mockDB) GetEntriesByRequestID(requestID int) ([]db.Entry, error) {
	entries := []db.Entry{}
	for _, entry := range m.Entries {
		if entry.RequestID == requestID {
			entries = append(entries, entry)
		}
	}
	return entries, nil
}

// モックDBを初期化
func initMockDB() *mockDB {
	testTime := time.Date(2024, 6, 1, 0, 0, 0, 0, time.UTC)
	return &mockDB{
		Requests: []db.Request{
			{ID: 1, CreatorID: 2, StartDate: testTime, EndDate: testTime, Deadline: testTime, CreatedAt: testTime},
			{ID: 2, CreatorID: 2, StartDate: testTime, EndDate: testTime, Deadline: testTime, CreatedAt: testTime},
		},
		Users: []db.User{
			{ID: 1, LoginID: "test_user", Password: "password", Name: "テストユーザー", Role: 0, CreatedAt: testTime},
			{ID: 2, LoginID: "test_manager", Password: "password2", Name: "テストマネージャー", Role: 1, CreatedAt: testTime},
		},
		Entries: []db.Entry{
			{ID: 1, RequestID: 1, UserID: 1, Date: testTime, Hour: 8},
			{ID: 2, RequestID: 1, UserID: 2, Date: testTime, Hour: 8},
			{ID: 3, RequestID: 2, UserID: 1, Date: testTime, Hour: 8},
			{ID: 4, RequestID: 2, UserID: 2, Date: testTime, Hour: 8},
		},
	}
}

// 新たなテスト用コンテキストを作成
func newTestContext() *context.AppContext {
	return &context.AppContext{
		DB: initMockDB(),
	}
}

func TestGetRequests(t *testing.T) {
	ctx := newTestContext()

	got, err := GetRequests(ctx)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	want := GetRequestsResponse{
		Request{ID: 1, Creator: User{ID: 2, Name: "テストマネージャー"}, StartDate: "2024-06-01", EndDate: "2024-06-01", Deadline: "2024-06-01", CreatedAt: "2024-06-01"},
		Request{ID: 2, Creator: User{ID: 2, Name: "テストマネージャー"}, StartDate: "2024-06-01", EndDate: "2024-06-01", Deadline: "2024-06-01", CreatedAt: "2024-06-01"},
	}

	if !reflect.DeepEqual(got, want) {
		t.Errorf("got %v, want %v", got, want)
	}
}

func TestGetEntries(t *testing.T) {
	ctx := newTestContext()

	got, err := GetEntries(ctx, 1)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	want := GetEntriesResponse{
		ID: 1,
		Entries: []Entry{
			{ID: 1, User: User{ID: 1, Name: "テストユーザー"}, Date: "2024-06-01", Hour: 8},
			{ID: 2, User: User{ID: 2, Name: "テストマネージャー"}, Date: "2024-06-01", Hour: 8},
		},
	}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("got %v, want %v", got, want)
	}
}
