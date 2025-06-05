package db

import (
	"errors"
)

var ErrUserNotFound = errors.New("user not found")

type User struct {
	ID        int
	LoginID   string
	Password  string
	Name      string
	Role      int
	CreatedAt string
}

type Request struct {
	ID        int
	CreatorID int
	StartDate string
	EndDate   string
	Deadline  string
	CreatedAt string
}

type Entry struct {
	ID           int
	SubmissionID int
	Date         string
	Hour         int
}

type Submission struct {
	ID          int
	RequestID   int
	SubmitterID int
	CreatedAt   string
	UpdatedAt   string
}

type DB interface {
	GetUserByID(id int) (User, error)
	GetUserByLoginID(loginID string) (User, error)
	GetRequests() ([]Request, error)
	GetRequestByID(id int) (Request, error)
	GetEntriesBySubmissionID(submissionID int) ([]Entry, error)
	GetSubmissionsByRequestID(requestID int) ([]Submission, error)
	AlreadySubmitted(requestID int, submitterID int) (bool, error)
	CreateRequest(creatorID int, startDate string, endDate string, deadline string) (int, error)
	CreateEntries(entries []Entry) ([]int, error)
	CreateSubmission(submitterID int, requestID int) (int, error)
}
