package model

import (
	"backend/auth"
	"backend/db"
	"reflect"
	"testing"
)

// TODO: 全てのテストで間違い入力時のテストを追加する

func TestGetEntriesByRequestID(t *testing.T) {
	ctx := newTestContext(
		[]db.User{
			{ID: 1, LoginID: "test_user1", Password: "password", Name: "テストユーザー1", Role: auth.RoleEmployee, CreatedAt: "2024-06-01 00:00:00"},
			{ID: 2, LoginID: "test_user2", Password: "password", Name: "テストユーザー2", Role: auth.RoleEmployee, CreatedAt: "2024-06-01 00:00:00"},
			{ID: 3, LoginID: "test_manager", Password: "password", Name: "テストマネージャー", Role: auth.RoleManager, CreatedAt: "2024-06-01 00:00:00"},
		},
		[]db.Request{
			{ID: 1, CreatorID: 3, StartDate: "2024-06-01", EndDate: "2024-06-01", Deadline: "2024-06-01 00:00:00", CreatedAt: "2024-06-01 00:00:00"},
		},
		[]db.Entry{
			{ID: 1, UserID: 1, RequestID: 1, Date: "2024-06-01", Hour: 8},
			{ID: 2, UserID: 2, RequestID: 1, Date: "2024-06-01", Hour: 8},
		},
	)
	got, err := GetEntriesByRequestID(ctx, 1)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	want := []Entry{
		{ID: 1, RequestID: 1, User: User{ID: 1, LoginID: "test_user1", Password: "password", Name: "テストユーザー1", Role: auth.RoleEmployee, CreatedAt: mustNewDateTime("2024-06-01 00:00:00")}, Date: mustNewDateOnly("2024-06-01"), Hour: 8},
		{ID: 2, RequestID: 1, User: User{ID: 2, LoginID: "test_user2", Password: "password", Name: "テストユーザー2", Role: auth.RoleEmployee, CreatedAt: mustNewDateTime("2024-06-01 00:00:00")}, Date: mustNewDateOnly("2024-06-01"), Hour: 8},
	}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("got %v, want %v", got, want)
	}
}

func TestCreateEntries(t *testing.T) {
	ctx := newTestContext(
		[]db.User{
			{ID: 1, Name: "テストユーザー1", Role: auth.RoleEmployee, CreatedAt: "2024-06-01 00:00:00"},
			{ID: 2, Name: "テストマネージャ", Role: auth.RoleManager, CreatedAt: "2024-06-01 00:00:00"},
		},
		[]db.Request{
			{ID: 1, CreatorID: 2, StartDate: "2024-06-01", EndDate: "2024-06-01", Deadline: "2024-06-01 00:00:00", CreatedAt: "2024-06-01 00:00:00"},
		},
		[]db.Entry{},
	)
	got, err := CreateEntries(ctx, NewEntries{RequestID: 1, CreatorID: 1, Entries: []NewEntry{{Date: mustNewDateOnly("2024-06-01"), Hour: 8}, {Date: mustNewDateOnly("2024-06-01"), Hour: 9}}})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	want := []int{1, 2}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("got %v, want %v", got, want)
	}
}
