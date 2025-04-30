package model

import (
	"backend/context"
	"backend/db"
	"reflect"
	"testing"
)

// 新たなテスト用コンテキストを作成
func newTestContext() *context.AppContext {
	return &context.AppContext{
		DB: db.InitMockDB(),
	}
}

func TestGetRequests(t *testing.T) {
	ctx := newTestContext()

	got, err := GetRequests(ctx)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	want := GetRequestsResponse{
		Request{ID: 1, Creator: User{ID: 2, Name: "テストマネージャー"}, StartDate: "2024-06-01", EndDate: "2024-06-01", Deadline: "2024-06-01", CreatedAt: "2024-06-01 00:00:00"},
		Request{ID: 2, Creator: User{ID: 2, Name: "テストマネージャー"}, StartDate: "2024-06-01", EndDate: "2024-06-01", Deadline: "2024-06-01", CreatedAt: "2024-06-01 00:00:00"},
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

func TestCreateRequest(t *testing.T) {
	ctx := newTestContext()

	got, err := CreateRequest(ctx, NewRequest{CreatorID: 2, StartDate: "2024-06-01", EndDate: "2024-06-01", Deadline: "2024-06-01"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	want := PostRequestsResponse{ID: 3}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("got %v, want %v", got, want)
	}
}

func TestCreateEntries(t *testing.T) {
	ctx := newTestContext()

	got, err := CreateEntries(ctx, NewEntries{ID: 1, Entries: []NewEntry{{UserID: 1, Date: "2024-06-01", Hour: 8}, {UserID: 1, Date: "2024-06-01", Hour: 9}}})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	want := PostEntriesResponse{ID: 1, Entries: []PostEntriesResponseEntry{{ID: 5}, {ID: 6}}}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("got %v, want %v", got, want)
	}
}
