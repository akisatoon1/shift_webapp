package model

import (
	"backend/context"
	"backend/db"
	"errors"
	"testing"
	"time"
)

// モック用のDB構造体
type mockDB struct {
	Requests []db.Request
	Users    []db.User
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
	return []db.Entry{}, nil
}

func TestGetRequests(t *testing.T) {
	// テストデータ
	testTime := time.Date(2024, 6, 1, 0, 0, 0, 0, time.UTC)
	mock := &mockDB{
		Requests: []db.Request{
			{
				ID:        1,
				CreatorID: 1,
				StartDate: testTime,
				EndDate:   testTime,
				Deadline:  testTime,
				CreatedAt: testTime,
			},
		},
		Users: []db.User{
			{
				ID:        1,
				LoginID:   "test_user",
				Password:  "password",
				Name:      "テストユーザー",
				Role:      0,
				CreatedAt: testTime,
			},
		},
	}

	ctx := &context.AppContext{DB: mock}

	got, err := GetRequests(ctx)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	want := `[{"id":1,"creator":{"id":1,"name":"テストユーザー"},"start_date":"2024-06-01","end_date":"2024-06-01","deadline":"2024-06-01","created_at":"2024-06-01"}]`
	if got != want {
		t.Errorf("got %s, want %s", got, want)
	}
}
