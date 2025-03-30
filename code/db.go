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

func (app *App) createUser(u user) error {
	_, err := app.db.Exec("INSERT INTO users (id, password, role) VALUES (?, ?, ?)", u.ID, u.Password, RoleEmployee)
	if err != nil {
		if sqlite3Err, ok := err.(sqlite3.Error); ok && sqlite3Err.ExtendedCode == sqlite3.ErrConstraintPrimaryKey {
			return errInvalidUserID
		}
		return logErr(err)
	}
	return nil
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
func (app *App) getUserIDFromSession(sessionID string) (string, error) {
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

func generateSessionID() string {
	b := make([]byte, 32)
	rand.Read(b)
	return base64.URLEncoding.EncodeToString(b)
}

func (app *App) createRequest(manager_id string, start_date string, end_date string) error {
	id, _ := uuid.NewRandom()
	_, err := app.db.Exec("INSERT INTO shift_requests (id, manager_id, start_date, end_date) VALUES (?, ?, ?, ?)", id, manager_id, start_date, end_date)
	if err != nil {
		return logErr(err)
	}
	return nil
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
