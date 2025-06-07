package db

import (
	"errors"
	"time"
)

// モック用のDB構造体
type mockDB struct {
	Requests    []Request
	Users       []User
	Entries     []Entry
	Submissions []Submission
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
	return User{}, ErrUserNotFound
}

func (m *mockDB) GetUserByLoginID(loginID string) (User, error) {
	for _, user := range m.Users {
		if user.LoginID == loginID {
			return user, nil
		}
	}
	return User{}, ErrUserNotFound
}

func (m *mockDB) GetEntriesBySubmissionID(submissionID int) ([]Entry, error) {
	entries := []Entry{}
	for _, entry := range m.Entries {
		if entry.SubmissionID == submissionID {
			entries = append(entries, entry)
		}
	}
	return entries, nil
}

func (m *mockDB) GetSubmissionsByRequestID(requestID int) ([]Submission, error) {
	submissions := []Submission{}
	for _, submission := range m.Submissions {
		if submission.RequestID == requestID {
			submissions = append(submissions, submission)
		}
	}
	return submissions, nil
}

func (m *mockDB) GetSubmissionByRequestIDAndSubmitterID(requestID int, submitterID int) (*Submission, error) {
	for _, submission := range m.Submissions {
		if submission.RequestID == requestID && submission.SubmitterID == submitterID {
			return &submission, nil
		}
	}
	return nil, nil
}

func (m *mockDB) CreateRequest(creatorID int, startDate string, endDate string, deadline string) (int, error) {
	m.Requests = append(m.Requests, Request{ID: len(m.Requests) + 1, CreatorID: creatorID, StartDate: startDate, EndDate: endDate, Deadline: deadline, CreatedAt: time.Now().Format(time.DateTime)})
	return len(m.Requests), nil
}

func (m *mockDB) CreateEntries(entries []Entry) ([]int, error) {
	lastID := len(m.Entries)
	ids := []int{}
	for i := range entries {
		id := lastID + i + 1
		entries[i].ID = id
		m.Entries = append(m.Entries, entries[i])
		ids = append(ids, id)
	}
	return ids, nil
}

func (m *mockDB) CreateSubmission(submitterID int, requestID int) (int, error) {
	submission := Submission{
		ID:          len(m.Submissions) + 1,
		RequestID:   requestID,
		SubmitterID: submitterID,
		CreatedAt:   time.Now().Format(time.DateTime),
		UpdatedAt:   time.Now().Format(time.DateTime),
	}
	m.Submissions = append(m.Submissions, submission)
	return submission.ID, nil
}

// テスト用データを入れたモックDBを生成
func NewMockDB(requests []Request, users []User, entries []Entry, submissions []Submission) *mockDB {
	return &mockDB{
		Requests:    requests,
		Users:       users,
		Entries:     entries,
		Submissions: submissions,
	}
}
