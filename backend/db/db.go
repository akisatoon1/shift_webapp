package db

import (
	"errors"
	"time"
)

var ErrUserNotFound = errors.New("user not found")

type User struct {
	ID        int
	LoginID   string
	Password  string
	Name      string
	Role      int
	CreatedAt time.Time
}

type Request struct {
	ID        int
	CreatorID int
	StartDate time.Time
	EndDate   time.Time
	Deadline  time.Time
	CreatedAt time.Time
}

type Entry struct {
	ID        int
	RequestID int
	UserID    int
	Date      time.Time
	Hour      int
}

type DB interface {
	GetUserByID(id int) (User, error)
	GetUserByLoginID(loginID string) (User, error)
	GetRequests() ([]Request, error)
	GetRequestByID(id int) (Request, error)
	GetEntriesByRequestID(requestID int) ([]Entry, error)
	CreateRequest(creatorID int, startDate time.Time, endDate time.Time, deadline time.Time) (int, error)
	CreateEntries(entries []Entry) ([]int, error)
}
