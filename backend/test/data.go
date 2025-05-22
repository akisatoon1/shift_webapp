package test

import (
	"backend/auth"
	"backend/db"

	"golang.org/x/crypto/bcrypt"
)

// サーバをテストモードで起動するときに使うデータ

func mustHashPassword(password string) string {
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(hashedPassword)
}

var MockUsers = []db.User{
	{ID: 1, LoginID: "employee1", Password: mustHashPassword("pass_employee1"), Name: "テストエンプロイ1", Role: auth.RoleEmployee, CreatedAt: "2023-12-31 12:00:00"},
	{ID: 3, LoginID: "employee2", Password: mustHashPassword("pass_employee2"), Name: "テストエンプロイ2", Role: auth.RoleEmployee, CreatedAt: "2024-01-31 12:00:00"},
	{ID: 2, LoginID: "manager1", Password: mustHashPassword("pass_manager1"), Name: "テストマネージャー1", Role: auth.RoleManager, CreatedAt: "2024-01-31 12:00:00"},
}

var MockRequests = []db.Request{
	{ID: 1, CreatorID: 2, StartDate: "2024-01-01", EndDate: "2024-01-07", Deadline: "2024-01-01 12:00:00", CreatedAt: "2023-12-31 12:00:00"},
	{ID: 2, CreatorID: 2, StartDate: "2024-02-01", EndDate: "2024-02-07", Deadline: "2024-02-01 12:00:00", CreatedAt: "2024-01-31 12:00:00"},
}

var MockEntries = []db.Entry{
	{ID: 1, RequestID: 1, UserID: 1, Date: "2024-01-01", Hour: 8},
	{ID: 2, RequestID: 1, UserID: 1, Date: "2024-01-02", Hour: 6},
	{ID: 3, RequestID: 2, UserID: 1, Date: "2024-02-01", Hour: 7},
	{ID: 4, RequestID: 2, UserID: 3, Date: "2024-02-02", Hour: 8},
	{ID: 5, RequestID: 2, UserID: 3, Date: "2024-02-03", Hour: 10},
	{ID: 6, RequestID: 2, UserID: 3, Date: "2024-02-04", Hour: 12},
}
