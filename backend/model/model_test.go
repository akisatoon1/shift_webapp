package model

import (
	"backend/auth"
	"backend/context"
	"backend/db"
	"reflect"
	"testing"
	"time"
)

// 新たなテスト用コンテキストを作成
func newTestContext(users []db.User, requests []db.Request, entries []db.Entry) *context.AppContext {
	return context.NewAppContext(db.NewMockDB(requests, users, entries), nil)
}

func TestGetRequests(t *testing.T) {
	ctx := newTestContext(
		[]db.User{
			{ID: 2, Name: "テストマネージャー"},
		},
		[]db.Request{
			{ID: 1, CreatorID: 2, StartDate: db.DateOnly(time.Date(2024, 6, 1, 0, 0, 0, 0, time.UTC)), EndDate: db.DateOnly(time.Date(2024, 6, 1, 0, 0, 0, 0, time.UTC)), Deadline: db.DateTime(time.Date(2024, 6, 1, 0, 0, 0, 0, time.UTC)), CreatedAt: db.DateTime(time.Date(2024, 6, 1, 0, 0, 0, 0, time.UTC))},
			{ID: 2, CreatorID: 2, StartDate: db.DateOnly(time.Date(2024, 6, 1, 0, 0, 0, 0, time.UTC)), EndDate: db.DateOnly(time.Date(2024, 6, 1, 0, 0, 0, 0, time.UTC)), Deadline: db.DateTime(time.Date(2024, 6, 1, 0, 0, 0, 0, time.UTC)), CreatedAt: db.DateTime(time.Date(2024, 6, 1, 0, 0, 0, 0, time.UTC))},
		},
		[]db.Entry{},
	)

	got, err := GetRequests(ctx)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	want := GetRequestsResponse{
		Request{ID: 1, Creator: User{ID: 2, Name: "テストマネージャー"}, StartDate: "2024-06-01", EndDate: "2024-06-01", Deadline: "2024-06-01 00:00:00", CreatedAt: "2024-06-01 00:00:00"},
		Request{ID: 2, Creator: User{ID: 2, Name: "テストマネージャー"}, StartDate: "2024-06-01", EndDate: "2024-06-01", Deadline: "2024-06-01 00:00:00", CreatedAt: "2024-06-01 00:00:00"},
	}

	if !reflect.DeepEqual(got, want) {
		t.Errorf("got %v, want %v", got, want)
	}
}

func TestGetEntries(t *testing.T) {
	ctx := newTestContext(
		[]db.User{
			{ID: 1, Name: "テストユーザー1"},
			{ID: 2, Name: "テストユーザー2"},
		},
		[]db.Request{
			{ID: 1, CreatorID: 3, StartDate: db.DateOnly(time.Date(2024, 6, 1, 0, 0, 0, 0, time.UTC)), EndDate: db.DateOnly(time.Date(2024, 6, 1, 0, 0, 0, 0, time.UTC)), Deadline: db.DateTime(time.Date(2024, 6, 1, 0, 0, 0, 0, time.UTC)), CreatedAt: db.DateTime(time.Date(2024, 6, 1, 0, 0, 0, 0, time.UTC))},
		},
		[]db.Entry{
			{ID: 1, UserID: 1, RequestID: 1, Date: db.DateOnly(time.Date(2024, 6, 1, 0, 0, 0, 0, time.UTC)), Hour: 8},
			{ID: 2, UserID: 2, RequestID: 1, Date: db.DateOnly(time.Date(2024, 6, 1, 0, 0, 0, 0, time.UTC)), Hour: 8},
		},
	)
	got, err := GetEntries(ctx, 1)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	want := GetEntriesResponse{
		ID: 1,
		Entries: []Entry{
			{ID: 1, User: User{ID: 1, Name: "テストユーザー1"}, Date: "2024-06-01", Hour: 8},
			{ID: 2, User: User{ID: 2, Name: "テストユーザー2"}, Date: "2024-06-01", Hour: 8},
		},
	}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("got %v, want %v", got, want)
	}
}

func TestCreateRequest(t *testing.T) {
	ctx := newTestContext(
		[]db.User{
			{ID: 2, Name: "テストマネージャー", Role: auth.RoleManager},
		},
		[]db.Request{},
		[]db.Entry{},
	)

	got, err := CreateRequest(ctx, NewRequest{CreatorID: 2, StartDate: "2024-06-01", EndDate: "2024-06-01", Deadline: "2024-06-01 00:00:00"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	want := PostRequestsResponse{ID: 1}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("got %v, want %v", got, want)
	}
}

func TestCreateEntries(t *testing.T) {
	ctx := newTestContext(
		[]db.User{
			{ID: 1, Name: "テストユーザー1", Role: auth.RoleEmployee},
		},
		[]db.Request{
			{ID: 1, CreatorID: 2, StartDate: db.DateOnly(time.Date(2024, 6, 1, 0, 0, 0, 0, time.UTC)), EndDate: db.DateOnly(time.Date(2024, 6, 1, 0, 0, 0, 0, time.UTC)), Deadline: db.DateTime(time.Date(2024, 6, 1, 0, 0, 0, 0, time.UTC)), CreatedAt: db.DateTime(time.Date(2024, 6, 1, 0, 0, 0, 0, time.UTC))},
		},
		[]db.Entry{},
	)
	got, err := CreateEntries(ctx, NewEntries{ID: 1, UserID: 1, Entries: []NewEntry{{Date: "2024-06-01", Hour: 8}, {Date: "2024-06-01", Hour: 9}}})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	want := PostEntriesResponse{ID: 1, Entries: []PostEntriesResponseEntry{{ID: 1}, {ID: 2}}}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("got %v, want %v", got, want)
	}
}

// test utility

func TestIsBeforeOrEqual(t *testing.T) {
	tests := []struct {
		a    db.DateOnly
		b    db.DateOnly
		want bool
	}{
		{db.DateOnly(time.Date(2024, 6, 1, 0, 0, 0, 0, time.UTC)), db.DateOnly(time.Date(2024, 6, 1, 0, 0, 0, 0, time.UTC)), true},
		{db.DateOnly(time.Date(2024, 6, 1, 0, 0, 0, 0, time.UTC)), db.DateOnly(time.Date(2024, 6, 2, 0, 0, 0, 0, time.UTC)), true},
		{db.DateOnly(time.Date(2024, 6, 1, 0, 0, 0, 0, time.UTC)), db.DateOnly(time.Date(2024, 5, 31, 0, 0, 0, 0, time.UTC)), false},
	}
	for _, test := range tests {
		got := isBeforeOrEqual(test.a, test.b)
		if got != test.want {
			t.Errorf("got %v, want %v", got, test.want)
		}
	}
}
