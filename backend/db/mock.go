package db

import (
	"errors"
	"time"
)

// モック用のDB構造体
type mockDB struct {
	Requests []Request
	Users    []User
	Entries  []Entry
}

func (m *mockDB) GetRequests() ([]Request, error) {
	return m.Requests, nil
}

func (m *mockDB) GetUserByID(id int) (User, error) {
	for _, user := range m.Users {
		if user.ID == id {
			return user, nil
		}
	}
	return User{}, errors.New("user not found")
}

func (m *mockDB) GetUserByLoginID(loginID string) (User, error) {
	for _, user := range m.Users {
		if user.LoginID == loginID {
			return user, nil
		}
	}
	return User{}, errors.New("user not found")
}

func (m *mockDB) GetEntriesByRequestID(requestID int) ([]Entry, error) {
	entries := []Entry{}
	for _, entry := range m.Entries {
		if entry.RequestID == requestID {
			entries = append(entries, entry)
		}
	}
	return entries, nil
}

func (m *mockDB) CreateRequest(creatorID int, startDate time.Time, endDate time.Time, deadline time.Time) (int, error) {
	m.Requests = append(m.Requests, Request{ID: len(m.Requests) + 1, CreatorID: creatorID, StartDate: startDate, EndDate: endDate, Deadline: deadline, CreatedAt: time.Now()})
	return len(m.Requests), nil
}

func (m *mockDB) CreateEntry(requestID int, userID int, date time.Time, hour int) (int, error) {
	m.Entries = append(m.Entries, Entry{ID: len(m.Entries) + 1, RequestID: requestID, UserID: userID, Date: date, Hour: hour})
	return len(m.Entries), nil
}

// モックDBを初期化
func InitMockDB() *mockDB {
	testTime := time.Date(2024, 6, 1, 0, 0, 0, 0, time.UTC)
	return &mockDB{
		Requests: []Request{
			{ID: 1, CreatorID: 2, StartDate: testTime, EndDate: testTime, Deadline: testTime, CreatedAt: testTime},
			{ID: 2, CreatorID: 2, StartDate: testTime, EndDate: testTime, Deadline: testTime, CreatedAt: testTime},
		},
		Users: []User{
			{ID: 1, LoginID: "test_user", Password: "password", Name: "テストユーザー", Role: 0, CreatedAt: testTime},
			{ID: 2, LoginID: "test_manager", Password: "password2", Name: "テストマネージャー", Role: 1, CreatedAt: testTime},
		},
		Entries: []Entry{
			{ID: 1, RequestID: 1, UserID: 1, Date: testTime, Hour: 8},
			{ID: 2, RequestID: 1, UserID: 2, Date: testTime, Hour: 8},
			{ID: 3, RequestID: 2, UserID: 1, Date: testTime, Hour: 8},
			{ID: 4, RequestID: 2, UserID: 2, Date: testTime, Hour: 8},
		},
	}
}
