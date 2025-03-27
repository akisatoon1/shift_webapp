package main

import (
	"crypto/rand"
	"database/sql"
	"encoding/base64"
	"errors"
	"log"
)

type user struct {
	ID       string
	Password string
	Role     int
}

var errDB = errors.New("db error")

func handleErrDB(err error) error {
	log.Println("DB ERROR:", err)
	return errDB
}

func (app *App) getUser(userID string) (user, error) {
	u := user{ID: userID}
	err := app.db.QueryRow("SELECT password, role FROM users WHERE id = ?", u.ID).Scan(&u.Password, &u.Role)
	if err != nil {
		if err != sql.ErrNoRows {
			return user{}, handleErrDB(err)
		}
		return user{}, sql.ErrNoRows
	}
	return u, nil
}

func (app *App) getAllUsers() ([]user, error) {
	rows, err := app.db.Query("SELECT id, password, role FROM users")
	if err != nil {
		return nil, handleErrDB(err)
	}
	defer rows.Close()

	var users []user
	for rows.Next() {
		var usr user
		if err := rows.Scan(&usr.ID, &usr.Password, &usr.Role); err != nil {
			return nil, handleErrDB(err)
		}
		users = append(users, usr)
	}
	if err := rows.Err(); err != nil {
		return nil, handleErrDB(err)
	}
	return users, nil
}

func (app *App) createUser(u user) error {
	_, err := app.db.Exec("INSERT INTO users (id, password, role) VALUES (?, ?, ?)", u.ID, u.Password, RoleEmployee)
	if err != nil {
		return handleErrDB(err)
	}
	return nil
}

func (app *App) deleteUser(userID string) error {
	_, err := app.db.Exec("DELETE FROM users WHERE id = ?", userID)
	if err != nil {
		return handleErrDB(err)
	}
	return nil
}

// sessionID, error
func (app *App) createSession(userID string) (string, error) {
	sessionID := generateSessionID()
	_, err := app.db.Exec("INSERT INTO sessions (id, user_id) VALUES (?, ?)", sessionID, userID)
	if err != nil {
		return "", handleErrDB(err)
	}
	return sessionID, nil
}

func (app *App) deleteSession(sessionID string) error {
	_, err := app.db.Exec("DELETE FROM sessions WHERE id = ?", sessionID)
	if err != nil {
		return handleErrDB(err)
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
			return "", handleErrDB(err)
		}
		return "", sql.ErrNoRows
	}
	return userID, nil
}

func generateSessionID() string {
	b := make([]byte, 32)
	rand.Read(b)
	return base64.URLEncoding.EncodeToString(b)
}
