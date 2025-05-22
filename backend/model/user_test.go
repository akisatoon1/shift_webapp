package model

import (
	"backend/auth"
	"backend/db"
	"reflect"
	"testing"
)

func TestGetUserByID(t *testing.T) {
	ctx := newTestContext(
		[]db.User{
			{ID: 1, LoginID: "testuser", Password: "password", Name: "テストユーザー", Role: auth.RoleEmployee, CreatedAt: "2024-06-01 00:00:00"},
		},
		[]db.Request{},
		[]db.Entry{},
	)

	got, err := GetUserByID(ctx, 1)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	want := User{
		ID:        1,
		LoginID:   "testuser",
		Password:  "password",
		Name:      "テストユーザー",
		Role:      auth.RoleEmployee,
		CreatedAt: mustNewDateTime("2024-06-01 00:00:00"),
	}

	if !reflect.DeepEqual(got, want) {
		t.Errorf("got %v, want %v", got, want)
	}
}
