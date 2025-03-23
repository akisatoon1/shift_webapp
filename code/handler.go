package main

import (
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"net/http"
)

func (app *App) rootHandler(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, "/home", http.StatusFound)
}

func (app *App) registerHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		userID := r.FormValue("user_id")
		password := r.FormValue("password")
		hash := sha256.Sum256([]byte(password))
		hashedPassword := base64.URLEncoding.EncodeToString(hash[:])

		err := app.createUser(userID, hashedPassword)
		if err != nil {
			http.Redirect(w, r, "/register", http.StatusFound)
			return
		}
		http.Redirect(w, r, "/login", http.StatusFound)
		return

	} else {
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		fmt.Fprintln(w, `
			<form method="POST">
				<input type="text" name="user_id" placeholder="User ID"><br>
				<input type="password" name="password" placeholder="Password"><br>
				<button type="submit">Register</button>
            </form>
		`)
	}
}

func (app *App) homeHandler(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("session_id")
	if err != nil {
		http.Redirect(w, r, "/login", http.StatusFound)
		return
	}

	userID, err := app.getUserIDFromSession(cookie.Value)
	if err != nil {
		http.Redirect(w, r, "/login", http.StatusFound)
		return
	}

	_, role, err := app.getUser(userID)
	var roleName string
	if isAdmin(role) {
		roleName = "Admin"
	} else {
		roleName = "Employee"
	}
	fmt.Fprintf(w, `
		<h1>Your user ID is '%v'</h1>
		<h2>You are %v<h2>
		<a href="/logout">Logout</a>
	`, userID, roleName)
}

func (app *App) loginHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		userID := r.FormValue("user_id")
		password := r.FormValue("password")
		hash := sha256.Sum256([]byte(password))

		hashedPassword, _, err := app.getUser(userID)
		if err != nil {
			http.Redirect(w, r, "/login", http.StatusFound)
			return
		}
		if hashedPassword == base64.URLEncoding.EncodeToString(hash[:]) {
			sessionID, err := app.createSession(userID)
			if err != nil {
				http.Redirect(w, r, "/login", http.StatusFound)
				return
			}
			http.SetCookie(w, &http.Cookie{
				Name:     "session_id",
				Value:    sessionID,
				HttpOnly: true,
				Secure:   false,
				Path:     "/",
			})
			http.Redirect(w, r, "/home", http.StatusFound)
			return
		} else {
			http.Redirect(w, r, "/login", http.StatusFound)
			return
		}

	} else {
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		fmt.Fprintln(w, `
                <form method="POST">
                        <input type="text" name="user_id" placeholder="User ID"><br>
                        <input type="password" name="password" placeholder="Password"><br>
                        <button type="submit">Login</button>
                </form>
				<a href="/register">create new account</a>
        `)
	}

}

func (app *App) logoutHandler(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("session_id")
	if err != nil {
		http.Redirect(w, r, "/login", http.StatusFound)
		return
	}

	err = app.deleteSession(cookie.Value)
	if err != nil {
		http.Redirect(w, r, "/login", http.StatusFound)
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:     "session_id",
		Value:    "",
		HttpOnly: true,
		Secure:   false,
		Path:     "/",
	})
	http.Redirect(w, r, "/login", http.StatusFound)
	return
}
