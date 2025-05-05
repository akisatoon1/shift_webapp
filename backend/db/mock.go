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

func (m *mockDB) GetRequestByID(id int) (Request, error) {
	for _, request := range m.Requests {
		if request.ID == id {
			return request, nil
		}
	}
	return Request{}, errors.New("request not found")
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

// テスト用データを入れたモックDBを生成
func NewMockDB(requests []Request, users []User, entries []Entry) *mockDB {
	return &mockDB{
		Requests: requests,
		Users:    users,
		Entries:  entries,
	}
}
