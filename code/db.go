package main

import (
	"crypto/rand"
	"encoding/base64"
)

// passwordHash, role, error
func (app *App) getUser(userID string) (string, int, error) {
	var password string
	var role int
	err := app.db.QueryRow("SELECT password, role FROM users WHERE id = ?", userID).Scan(&password, &role)
	if err != nil {
		return "", 0, err
	}
	return password, role, nil
}

func (app *App) createUser(userID string, hashedPassword string) error {
	_, err := app.db.Exec("INSERT INTO users (id, password, role) VALUES (?, ?, ?)", userID, hashedPassword, RoleEmployee)
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
