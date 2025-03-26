package main

import (
	"crypto/rand"
	"encoding/base64"
)

type user struct {
	ID       string
	Password string
	Role     int
}

func (app *App) getUser(userID string) (user, error) {
	u := user{ID: userID}
	err := app.db.QueryRow("SELECT password, role FROM users WHERE id = ?", u.ID).Scan(&u.Password, &u.Role)
	if err != nil {
		return user{}, err
	}
	return u, nil
}

func (app *App) getAllUsers() ([]user, error) {
	rows, err := app.db.Query("SELECT id, password, role FROM users")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []user
	for rows.Next() {
		var usr user
		if err := rows.Scan(&usr.ID, &usr.Password, &usr.Role); err != nil {
			return nil, err
		}
		users = append(users, usr)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return users, nil
}

func (app *App) createUser(u user) error {
	_, err := app.db.Exec("INSERT INTO users (id, password, role) VALUES (?, ?, ?)", u.ID, u.Password, RoleEmployee)
	return err
}

func (app *App) deleteUser(userID string) error {
	_, err := app.db.Exec("DELETE FROM users WHERE id = ?", userID)
	return err
}

// sessionID, error
func (app *App) createSession(userID string) (string, error) {
	sessionID := generateSessionID()
	_, err := app.db.Exec("INSERT INTO sessions (id, user_id) VALUES (?, ?)", sessionID, userID)
	if err != nil {
		return "", err
	}
	return sessionID, nil
}

func (app *App) deleteSession(sessionID string) error {
	_, err := app.db.Exec("DELETE FROM sessions WHERE id = ?", sessionID)
	if err != nil {
		return err
	}
	return nil
}

// sessionIDに紐づいたuserIDを探す
// userID, error
func (app *App) getUserIDFromSession(sessionID string) (string, error) {
	var userID string
	err := app.db.QueryRow("SELECT user_id FROM sessions WHERE id = ?", sessionID).Scan(&userID)
	if err != nil {
		return "", err
	}
	return userID, nil
}

func generateSessionID() string {
	b := make([]byte, 32)
	rand.Read(b)
	return base64.URLEncoding.EncodeToString(b)
}
