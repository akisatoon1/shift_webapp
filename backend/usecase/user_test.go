package usecase

import (
	"backend/auth"
	"backend/db"
	"backend/domain"
	"testing"
)

func TestGetUserByID(t *testing.T) {
	ctx := newTestContext(
		[]db.User{
			{ID: 1, LoginID: "testuser", Password: "password", Name: "テストユーザー", Role: auth.RoleEmployee, CreatedAt: "2024-06-01 00:00:00"},
		},
		[]db.Request{},
		[]db.Entry{},
		[]db.Submission{},
	)

	var u IUserUsecase = &userUsecase{}

	// 正常系
	got, err := u.FindByID(ctx, 1)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	want := domain.User{
		ID:        1,
		LoginID:   "testuser",
		Password:  "password",
		Name:      "テストユーザー",
		Role:      auth.RoleEmployee,
		CreatedAt: mustNewDateTime("2024-06-01 00:00:00"),
	}

	assert(t, got, want)

	// 異常系
	_, err = u.FindByID(ctx, 999)
	if err == nil {
		t.Fatalf("expected error")
	}
}
