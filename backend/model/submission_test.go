package model

import (
	"backend/auth"
	"backend/db"
	"testing"
)

func TestGetSubmissionsByRequestID(t *testing.T) {
	ctx := newTestContext(
		[]db.User{
			{ID: 1, LoginID: "test_manager", Password: "password", Name: "テストマネージャー", Role: auth.RoleManager, CreatedAt: "2024-06-01 00:00:00"},
			{ID: 2, LoginID: "test_user_2", Password: "password2", Name: "テストユーザー2", Role: auth.RoleEmployee, CreatedAt: "2023-01-01 00:00:00"},
			{ID: 3, LoginID: "test_user_3", Password: "password3", Name: "テストユーザー3", Role: auth.RoleEmployee, CreatedAt: "2023-01-03 00:00:00"},
		},
		[]db.Request{
			{ID: 1, CreatorID: 1, StartDate: "2024-06-01", EndDate: "2024-06-01", Deadline: "2024-06-01 00:00:00", CreatedAt: "2024-06-01 00:00:00"},
		},
		[]db.Entry{},
		[]db.Submission{
			{ID: 1, RequestID: 1, SubmitterID: 2, CreatedAt: "2023-01-01 00:00:00", UpdatedAt: "2023-01-02 00:00:00"},
			{ID: 2, RequestID: 1, SubmitterID: 3, CreatedAt: "2023-01-03 00:00:00", UpdatedAt: "2023-01-04 00:00:00"},
		},
	)

	// シフトリクエストIDが存在しない場合のテスト
	_, err := GetSubmissionsByRequestID(ctx, 123)
	if err == nil {
		t.Errorf("Expected error for non-existent request ID, got nil")
	}

	// シフトリクエストIDが存在する場合のテスト
	submissions, err := GetSubmissionsByRequestID(ctx, 1)
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}

	want := []Submission{
		{
			ID:          1,
			RequestID:   1,
			SubmitterID: 2,
			Submitter: User{
				ID:        2,
				LoginID:   "test_user_2",
				Password:  "password2",
				Name:      "テストユーザー2",
				Role:      auth.RoleEmployee,
				CreatedAt: mustNewDateTime("2023-01-01 00:00:00"),
			},
			CreatedAt: mustNewDateTime("2023-01-01 00:00:00"),
			UpdatedAt: mustNewDateTime("2023-01-02 00:00:00"),
		},
		{
			ID:          2,
			RequestID:   1,
			SubmitterID: 3,
			Submitter: User{
				ID:        3,
				LoginID:   "test_user_3",
				Password:  "password3",
				Name:      "テストユーザー3",
				Role:      auth.RoleEmployee,
				CreatedAt: mustNewDateTime("2023-01-03 00:00:00"),
			},
			CreatedAt: mustNewDateTime("2023-01-03 00:00:00"),
			UpdatedAt: mustNewDateTime("2023-01-04 00:00:00"),
		},
	}
	assert(t, submissions, want)
}
