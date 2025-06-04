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
		[]db.Entry{
			{ID: 1, SubmissionID: 1, Date: "2024-06-01", Hour: 8},
			{ID: 2, SubmissionID: 2, Date: "2024-06-01", Hour: 9},
		},
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
			Entries: []entry{
				{ID: 1, SubmissionID: 1, Date: mustNewDateOnly("2024-06-01"), Hour: 8},
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
			Entries: []entry{
				{ID: 2, SubmissionID: 2, Date: mustNewDateOnly("2024-06-01"), Hour: 9},
			},
			CreatedAt: mustNewDateTime("2023-01-03 00:00:00"),
			UpdatedAt: mustNewDateTime("2023-01-04 00:00:00"),
		},
	}
	assert(t, submissions, want)
}

func TestCreateSubmission(t *testing.T) {
	ctx := newTestContext(
		[]db.User{
			{ID: 1, LoginID: "test_manager", Password: "password", Name: "テストマネージャー", Role: auth.RoleManager, CreatedAt: "2024-06-01 00:00:00"},
			{ID: 2, LoginID: "test_employee", Password: "password", Name: "テスト従業員", Role: auth.RoleEmployee, CreatedAt: "2023-01-01 00:00:00"},
			{ID: 3, LoginID: "test_client", Password: "password", Name: "テストクライアント", Role: auth.RoleManager, CreatedAt: "2023-01-02 00:00:00"},
		},
		[]db.Request{
			{ID: 1, CreatorID: 1, StartDate: "2024-06-01", EndDate: "2024-06-07", Deadline: "2024-05-30 00:00:00", CreatedAt: "2024-05-20 00:00:00"},
			{ID: 999, CreatorID: 1, StartDate: "2024-07-01", EndDate: "2024-07-07", Deadline: "2024-06-30 00:00:00", CreatedAt: "2024-06-20 00:00:00"},
		},
		[]db.Entry{},
		[]db.Submission{},
	)

	// 正常系: 有効なリクエストに対して従業員が提出する
	submissionID, err := CreateSubmission(ctx, NewSubmission{
		RequestID:   1,
		SubmitterID: 2,
		NewEntries: []NewEntry{
			{Date: mustNewDateOnly("2024-06-01"), Hour: 9},
			{Date: mustNewDateOnly("2024-06-02"), Hour: 10},
		},
	})
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	if submissionID != 1 {
		t.Errorf("Expected valid submission ID, got %d", submissionID)
	}

	// 異常系: 存在しないリクエストID
	_, err = CreateSubmission(ctx, NewSubmission{
		RequestID:   9999,
		SubmitterID: 2,
		NewEntries: []NewEntry{
			{Date: mustNewDateOnly("2024-06-01"), Hour: 9},
		},
	})
	if err == nil {
		t.Errorf("Expected error for non-existent request ID, got nil")
	}

	// 異常系: 従業員でないユーザーが提出しようとする
	_, err = CreateSubmission(ctx, NewSubmission{
		RequestID:   1,
		SubmitterID: 3, // manager
		NewEntries: []NewEntry{
			{Date: mustNewDateOnly("2024-06-01"), Hour: 9},
		},
	})
	if err != ErrForbidden {
		t.Errorf("Expected ErrForbidden for non-employee user, got %v", err)
	}

	// 異常系: リクエスト期間外の日付
	_, err = CreateSubmission(ctx, NewSubmission{
		RequestID:   1,
		SubmitterID: 2,
		NewEntries: []NewEntry{
			{Date: mustNewDateOnly("2024-05-31"), Hour: 9}, // リクエスト開始日より前
		},
	})
	if _, ok := err.(InputError); !ok {
		t.Errorf("Expected InputError for date outside request range, got %v", err)
	}

	// 異常系: 無効な時間
	_, err = CreateSubmission(ctx, NewSubmission{
		RequestID:   1,
		SubmitterID: 2,
		NewEntries: []NewEntry{
			{Date: mustNewDateOnly("2024-06-01"), Hour: 24}, // 24時は無効
		},
	})
	if _, ok := err.(InputError); !ok {
		t.Errorf("Expected InputError for invalid hour, got %v", err)
	}
}
