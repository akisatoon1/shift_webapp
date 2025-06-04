package model

import (
	"backend/auth"
	"backend/db"
	"testing"
)

func TestGetEntriesBySubmissionID(t *testing.T) {
	ctx := newTestContext(
		[]db.User{
			{ID: 1, LoginID: "test_user1", Password: "password", Name: "テストユーザー1", Role: auth.RoleEmployee, CreatedAt: "2024-06-01 00:00:00"},
			{ID: 3, LoginID: "test_manager", Password: "password", Name: "テストマネージャー", Role: auth.RoleManager, CreatedAt: "2024-06-01 00:00:00"},
		},
		[]db.Request{
			{ID: 1, CreatorID: 3, StartDate: "2024-06-01", EndDate: "2024-06-01", Deadline: "2024-06-01 00:00:00", CreatedAt: "2024-06-01 00:00:00"},
		},
		[]db.Entry{
			{ID: 1, SubmissionID: 1, Date: "2024-06-01", Hour: 8},
			{ID: 2, SubmissionID: 1, Date: "2024-06-01", Hour: 9},
		},
		[]db.Submission{
			{ID: 1, RequestID: 1, SubmitterID: 1, CreatedAt: "2024-06-01 00:00:00", UpdatedAt: "2024-06-01 00:00:00"},
		},
	)
	got, err := getEntriesBySubmissionID(ctx, 1)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	want := []entry{
		{ID: 1, SubmissionID: 1, Date: mustNewDateOnly("2024-06-01"), Hour: 8},
		{ID: 2, SubmissionID: 1, Date: mustNewDateOnly("2024-06-01"), Hour: 9},
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
		[]db.Submission{},
	)

	// 正常系
	got, err := createEntries(ctx, 1, []NewEntry{{Date: mustNewDateOnly("2024-06-01"), Hour: 8}, {Date: mustNewDateOnly("2024-06-01"), Hour: 9}})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	want := []int{1, 2}
	assert(t, got, want)
}
