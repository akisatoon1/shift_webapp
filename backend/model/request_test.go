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
		Request{ID: 1, Creator: User{ID: 2, LoginID: "test_manager", Password: "password", Name: "テストマネージャー", Role: auth.RoleManager, CreatedAt: mustNewDateTime("2024-06-01 00:00:00")}, StartDate: mustNewDateOnly("2024-06-01"), EndDate: mustNewDateOnly("2024-06-01"), Deadline: mustNewDateTime("2024-06-01 00:00:00"), CreatedAt: mustNewDateTime("2024-06-01 00:00:00")},
		Request{ID: 2, Creator: User{ID: 2, LoginID: "test_manager", Password: "password", Name: "テストマネージャー", Role: auth.RoleManager, CreatedAt: mustNewDateTime("2024-06-01 00:00:00")}, StartDate: mustNewDateOnly("2024-06-01"), EndDate: mustNewDateOnly("2024-06-01"), Deadline: mustNewDateTime("2024-06-01 00:00:00"), CreatedAt: mustNewDateTime("2024-06-01 00:00:00")},
	}

	assert(t, got, want)
}

func TestCreateRequest(t *testing.T) {
	ctx := newTestContext(
		[]db.User{
			{ID: 2, Name: "テストマネージャー", Role: auth.RoleManager},
		},
		[]db.Request{},
		[]db.Entry{},
	)

	got, err := CreateRequest(ctx, NewRequest{CreatorID: 2, StartDate: mustNewDateOnly("2024-06-01"), EndDate: mustNewDateOnly("2024-06-01"), Deadline: mustNewDateTime("2024-06-01 00:00:00")})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	want := 1
	assert(t, got, want)
}
