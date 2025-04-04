package main

import (
	"crypto/rand"
	"database/sql"
	"encoding/base64"
	"log"

	"github.com/google/uuid"
	"github.com/mattn/go-sqlite3"
)

type user struct {
	ID       string
	Password string
	Role     int
}

type request struct {
	ID        string
	ManagerID string
	StartDate string
	EndDate   string
	CreatedAt string
}

type submission struct {
	ID             string
	RequestID      string
	StaffID        string
	SubmissionDate string
	SubmittedAt    string
}

type entry struct {
	ID           string
	SubmissionID string
	ShiftHour    int
	CreatedAt    string
}

/*
	database api の一覧

	規則: createは必ずidを返す

	// user
	createUser(user) (string, error)
	getUser(string) (user, error)
	getAllUsers() ([]user, error)
	getUserIDFromSessionID(string) (string, error)
	deleteUser(string) error

	// session
	createSession(string) (string, error)
	deleteSession(string) error

	// request
	createRequest(string, string, string) (string, error)
	getRequest(string) (request, error)
	getAllRequests() ([]request, error)

	// submission
	createSubmission(string, string, string) (string, error)
	getSubmissionsByRequestID(string) ([]submission, error)
	getSubmissionsByRequestIDAndDate(string string) ([]submission, error)
	getSubmissionsByRequestAndUserID(string, string) ([]submission, error)

	// entry
	createEntry(string, int) (string, error)
	getEntriesBySubmissionID(string) ([]entry, error)
*/

func logErr(err error) error {
	log.Println("DB ERROR:", err)
	return err
}

func (app *App) getUser(userID string) (user, error) {
	u := user{ID: userID}
	err := app.db.QueryRow("SELECT password, role FROM users WHERE id = ?", u.ID).Scan(&u.Password, &u.Role)
	if err != nil {
		if err != sql.ErrNoRows {
			return user{}, logErr(err)
		}
		return user{}, errUserNotFound
	}
	return u, nil
}

func (app *App) getAllUsers() ([]user, error) {
	rows, err := app.db.Query("SELECT id, password, role FROM users")
	if err != nil {
		return nil, logErr(err)
	}
	defer rows.Close()

	var users []user
	for rows.Next() {
		var usr user
		if err := rows.Scan(&usr.ID, &usr.Password, &usr.Role); err != nil {
			return nil, logErr(err)
		}
		users = append(users, usr)
	}
	if err := rows.Err(); err != nil {
		return nil, logErr(err)
	}
	return users, nil
}

func (app *App) createUser(u user) (string, error) {
	_, err := app.db.Exec("INSERT INTO users (id, password, role) VALUES (?, ?, ?)", u.ID, u.Password, RoleEmployee)
	if err != nil {
		if sqlite3Err, ok := err.(sqlite3.Error); ok && sqlite3Err.ExtendedCode == sqlite3.ErrConstraintPrimaryKey {
			return "", errInvalidUserID
		}
		return "", logErr(err)
	}
	return u.ID, nil
}

func (app *App) deleteUser(userID string) error {
	_, err := app.db.Exec("DELETE FROM users WHERE id = ?", userID)
	if err != nil {
		return logErr(err)
	}
	return nil
}

// sessionID, error
func (app *App) createSession(userID string) (string, error) {
	generateSessionID := func() string {
		b := make([]byte, 32)
		rand.Read(b)
		return base64.URLEncoding.EncodeToString(b)
	}

	sessionID := generateSessionID()
	_, err := app.db.Exec("INSERT INTO sessions (id, user_id) VALUES (?, ?)", sessionID, userID)
	if err != nil {
		return "", logErr(err)
	}
	return sessionID, nil
}

func (app *App) deleteSession(sessionID string) error {
	_, err := app.db.Exec("DELETE FROM sessions WHERE id = ?", sessionID)
	if err != nil {
		return logErr(err)
	}
	return nil
}

// sessionIDに紐づいたuserIDを探す
// userID, error
func (app *App) getUserIDFromSessionID(sessionID string) (string, error) {
	var userID string
	err := app.db.QueryRow("SELECT user_id FROM sessions WHERE id = ?", sessionID).Scan(&userID)
	if err != nil {
		if err != sql.ErrNoRows {
			return "", logErr(err)
		}
		return "", errSessionNotFound
	}
	return userID, nil
}

func (app *App) createRequest(manager_id string, start_date string, end_date string) (string, error) {
	id, _ := uuid.NewRandom()
	_, err := app.db.Exec("INSERT INTO shift_requests (id, manager_id, start_date, end_date) VALUES (?, ?, ?, ?)", id, manager_id, start_date, end_date)
	if err != nil {
		return "", logErr(err)
	}
	return id.String(), nil
}

func (app *App) getRequest(requestID string) (request, error) {
	req := request{ID: requestID}
	err := app.db.QueryRow("SELECT manager_id, start_date, end_date, created_at FROM shift_requests WHERE id = ?", req.ID).Scan(&req.ManagerID, &req.StartDate, &req.EndDate, &req.CreatedAt)
	if err != nil {
		if err != sql.ErrNoRows {
			return request{}, logErr(err)
		}
		return request{}, errRequestNotFound
	}
	return req, nil
}

func (app *App) getAllRequests() ([]request, error) {
	rows, err := app.db.Query("SELECT id, manager_id, start_date, end_date, created_at FROM shift_requests")
	if err != nil {
		return nil, logErr(err)
	}
	defer rows.Close()

	var requests []request
	for rows.Next() {
		var req request
		if err := rows.Scan(&req.ID, &req.ManagerID, &req.StartDate, &req.EndDate, &req.CreatedAt); err != nil {
			return nil, logErr(err)
		}
		requests = append(requests, req)
	}
	if err := rows.Err(); err != nil {
		return nil, logErr(err)
	}
	return requests, nil
}

func (app *App) getSubmissionsByRequestID(requestID string) ([]submission, error) {
	rows, err := app.db.Query("SELECT id, request_id, staff_id, submission_date, submitted_at FROM shift_submissions WHERE request_id = ?", requestID)
	if err != nil {
		return nil, logErr(err)
	}
	defer rows.Close()

	var submissions []submission
	for rows.Next() {
		var sub submission
		if err := rows.Scan(&sub.ID, &sub.RequestID, &sub.StaffID, &sub.SubmissionDate, &sub.SubmittedAt); err != nil {
			return nil, logErr(err)
		}
		submissions = append(submissions, sub)
	}
	if err := rows.Err(); err != nil {
		return nil, logErr(err)
	}
	return submissions, nil
}

func (app *App) getSubmissionsByRequestIDAndDate(requestID string, date string) ([]submission, error) {
	rows, err := app.db.Query("SELECT id, request_id, staff_id, submission_date, submitted_at FROM shift_submissions WHERE request_id = ? AND submission_date = ?", requestID, date)
	if err != nil {
		return nil, logErr(err)
	}
	defer rows.Close()

	var submissions []submission
	for rows.Next() {
		var sub submission
		if err := rows.Scan(&sub.ID, &sub.RequestID, &sub.StaffID, &sub.SubmissionDate, &sub.SubmittedAt); err != nil {
			return nil, logErr(err)
		}
		submissions = append(submissions, sub)
	}
	if err := rows.Err(); err != nil {
		return nil, logErr(err)
	}
	return submissions, nil
}

func (app *App) getSubmissionsByRequestAndUserID(requestID string, staffID string) ([]submission, error) {
	rows, err := app.db.Query("SELECT id, request_id, staff_id, submission_date, submitted_at FROM shift_submissions WHERE request_id = ? AND staff_id = ?", requestID, staffID)
	if err != nil {
		return nil, logErr(err)
	}
	defer rows.Close()

	var submissions []submission
	for rows.Next() {
		var sub submission
		if err := rows.Scan(&sub.ID, &sub.RequestID, &sub.StaffID, &sub.SubmissionDate, &sub.SubmittedAt); err != nil {
			return nil, logErr(err)
		}
		submissions = append(submissions, sub)
	}
	if err := rows.Err(); err != nil {
		return nil, logErr(err)
	}
	return submissions, nil
}

func (app *App) createSubmission(requestID string, staffID string, submissionDate string) (string, error) {
	id, _ := uuid.NewRandom()
	_, err := app.db.Exec("INSERT INTO shift_submissions (id, request_id, staff_id, submission_date) VALUES (?, ?, ?, ?)", id, requestID, staffID, submissionDate)
	if err != nil {
		return "", logErr(err)
	}
	return id.String(), nil
}

func (app *App) createEntry(submissionID string, shiftHour string) (string, error) {
	id, _ := uuid.NewRandom()
	_, err := app.db.Exec("INSERT INTO shift_entries (id, submission_id, shift_hour) VALUES (?, ?, ?)", id, submissionID, shiftHour)
	if err != nil {
		return "", logErr(err)
	}
	return id.String(), nil
}

func (app *App) getEntriesBySubmissionID(submissionID string) ([]entry, error) {
	rows, err := app.db.Query("SELECT id, submission_id, shift_hour, created_at FROM shift_entries WHERE submission_id = ?", submissionID)
	if err != nil {
		return nil, logErr(err)
	}
	defer rows.Close()

	var entries []entry
	for rows.Next() {
		var ent entry
		if err := rows.Scan(&ent.ID, &ent.SubmissionID, &ent.ShiftHour, &ent.CreatedAt); err != nil {
			return nil, logErr(err)
		}
		entries = append(entries, ent)
	}
	if err := rows.Err(); err != nil {
		return nil, logErr(err)
	}
	return entries, nil
}
