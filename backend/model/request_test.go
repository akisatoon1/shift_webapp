package model

import (
	"backend/auth"
	"backend/db"
	"testing"
)

func TestGetRequestByID(t *testing.T) {
	ctx := newTestContext(
		[]db.User{
			{ID: 2, LoginID: "test_manager", Password: "password", Name: "テストマネージャー", Role: auth.RoleManager, CreatedAt: "2024-06-01 00:00:00"},
		},
		[]db.Request{
			{ID: 1, CreatorID: 2, StartDate: "2024-06-01", EndDate: "2024-06-01", Deadline: "2024-06-01 00:00:00", CreatedAt: "2024-06-01 00:00:00"},
		},
		[]db.Entry{},
	)

	// 正常系
	got, err := GetRequestByID(ctx, 1)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	want := Request{
		ID:        1,
		Creator:   User{ID: 2, LoginID: "test_manager", Password: "password", Name: "テストマネージャー", Role: auth.RoleManager, CreatedAt: mustNewDateTime("2024-06-01 00:00:00")},
		StartDate: mustNewDateOnly("2024-06-01"),
		EndDate:   mustNewDateOnly("2024-06-01"),
		Deadline:  mustNewDateTime("2024-06-01 00:00:00"),
		CreatedAt: mustNewDateTime("2024-06-01 00:00:00"),
	}

	assert(t, got, want)

	// 異常系
	_, err = GetRequestByID(ctx, 999)
	if err == nil {
		t.Fatalf("expected error")
	}
}

func TestGetRequests(t *testing.T) {
	ctx := newTestContext(
		[]db.User{
			{ID: 2, LoginID: "test_manager", Password: "password", Name: "テストマネージャー", Role: auth.RoleManager, CreatedAt: "2024-06-01 00:00:00"},
		},
		[]db.Request{
			{ID: 1, CreatorID: 2, StartDate: "2024-06-01", EndDate: "2024-06-01", Deadline: "2024-06-01 00:00:00", CreatedAt: "2024-06-01 00:00:00"},
			{ID: 2, CreatorID: 2, StartDate: "2024-06-01", EndDate: "2024-06-01", Deadline: "2024-06-01 00:00:00", CreatedAt: "2024-06-01 00:00:00"},
		},
		[]db.Entry{},
	)

	got, err := GetRequests(ctx)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	want := []Request{
		{ID: 1, Creator: User{ID: 2, LoginID: "test_manager", Password: "password", Name: "テストマネージャー", Role: auth.RoleManager, CreatedAt: mustNewDateTime("2024-06-01 00:00:00")}, StartDate: mustNewDateOnly("2024-06-01"), EndDate: mustNewDateOnly("2024-06-01"), Deadline: mustNewDateTime("2024-06-01 00:00:00"), CreatedAt: mustNewDateTime("2024-06-01 00:00:00")},
		{ID: 2, Creator: User{ID: 2, LoginID: "test_manager", Password: "password", Name: "テストマネージャー", Role: auth.RoleManager, CreatedAt: mustNewDateTime("2024-06-01 00:00:00")}, StartDate: mustNewDateOnly("2024-06-01"), EndDate: mustNewDateOnly("2024-06-01"), Deadline: mustNewDateTime("2024-06-01 00:00:00"), CreatedAt: mustNewDateTime("2024-06-01 00:00:00")},
	}

	assert(t, got, want)
}

func TestCreateRequest(t *testing.T) {
	ctx := newTestContext(
		[]db.User{
			{ID: 1, LoginID: "test_user", Password: "password", Name: "テストユーザー", Role: auth.RoleEmployee},
			{ID: 2, LoginID: "test_manager", Password: "password", Name: "テストマネージャー", Role: auth.RoleManager},
		},
		[]db.Request{},
		[]db.Entry{},
	)

	// 正常系
	got, err := CreateRequest(ctx, NewRequest{CreatorID: 2, StartDate: mustNewDateOnly("2024-06-01"), EndDate: mustNewDateOnly("2024-06-01"), Deadline: mustNewDateTime("2024-06-01 00:00:00")})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	want := 1
	assert(t, got, want)

	// 異常系

	// 作成者が存在しない場合
	_, err = CreateRequest(ctx, NewRequest{CreatorID: 999, StartDate: mustNewDateOnly("2024-06-01"), EndDate: mustNewDateOnly("2024-06-01"), Deadline: mustNewDateTime("2024-06-01 00:00:00")})
	if err == nil {
		t.Fatalf("expected error")
	}

	// 作成者がマネージャーでない場合
	_, err = CreateRequest(ctx, NewRequest{CreatorID: 1, StartDate: mustNewDateOnly("2024-06-01"), EndDate: mustNewDateOnly("2024-06-01"), Deadline: mustNewDateTime("2024-06-01 00:00:00")})
	if err == nil {
		t.Fatalf("expected error")
	}

	// 開始日が終了日より後の場合
	_, err = CreateRequest(ctx, NewRequest{CreatorID: 2, StartDate: mustNewDateOnly("2024-06-02"), EndDate: mustNewDateOnly("2024-06-01"), Deadline: mustNewDateTime("2024-06-01 00:00:00")})
	if err == nil {
		t.Fatalf("expected error")
	}

	// 締切日が開始日より後の場合
	_, err = CreateRequest(ctx, NewRequest{CreatorID: 2, StartDate: mustNewDateOnly("2024-06-01"), EndDate: mustNewDateOnly("2024-06-01"), Deadline: mustNewDateTime("2024-06-02 00:00:00")})
	if err == nil {
		t.Fatalf("expected error")
	}
}
