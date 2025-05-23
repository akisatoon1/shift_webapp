package model

import (
	"backend/auth"
	"backend/db"
	"testing"
)

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
	assert(t, got, want)
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

	// 正常系
	got, err := CreateEntries(ctx, NewEntries{RequestID: 1, CreatorID: 1, Entries: []NewEntry{{Date: mustNewDateOnly("2024-06-01"), Hour: 8}, {Date: mustNewDateOnly("2024-06-01"), Hour: 9}}})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	want := []int{1, 2}
	assert(t, got, want)

	// 異常系

	// 存在しないリクエストID
	_, err = CreateEntries(ctx, NewEntries{RequestID: 999, CreatorID: 1, Entries: []NewEntry{{Date: mustNewDateOnly("2024-06-01"), Hour: 8}}})
	if err == nil {
		t.Fatalf("expected error for non-existent request ID")
	}

	// 作成者が存在しない
	_, err = CreateEntries(ctx, NewEntries{RequestID: 1, CreatorID: 999, Entries: []NewEntry{{Date: mustNewDateOnly("2024-06-01"), Hour: 8}}})
	if err == nil {
		t.Fatalf("expected error for non-existent creator ID")
	}

	// 作成者がemployeeでない場合
	_, err = CreateEntries(ctx, NewEntries{RequestID: 1, CreatorID: 2, Entries: []NewEntry{{Date: mustNewDateOnly("2024-06-01"), Hour: 8}}})
	if err == nil {
		t.Fatalf("expected error for non-employee creator")
	}

	// 日付がリクエストの範囲外
	_, err = CreateEntries(ctx, NewEntries{RequestID: 1, CreatorID: 1, Entries: []NewEntry{{Date: mustNewDateOnly("2024-06-02"), Hour: 8}}})
	if err == nil {
		t.Fatalf("expected error for date outside request range")
	}

	// 時間が0未満または24以上
	_, err = CreateEntries(ctx, NewEntries{RequestID: 1, CreatorID: 1, Entries: []NewEntry{{Date: mustNewDateOnly("2024-06-01"), Hour: -1}}})
	if err == nil {
		t.Fatalf("expected error for hour less than 0")
	}
}
