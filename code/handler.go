package main

import (
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"html/template"
	"net/http"
	"net/url"
)

func (app *App) homeHandler(w http.ResponseWriter, r *http.Request) {
	userID, err := app.getUserIDFromCookie(r)
	if err != nil {
		if err == errDB {
			responseError(w)
			return
		}
		redirectWithError(w, r, "/login", "ログインしてください")
		return
	}

	usr, err := app.getUser(userID)
	if err != nil {
		if err == errDB {
			responseError(w)
			return
		}
		redirectWithError(w, r, "/login", "ログインしてください")
		return
	}
	if isAdmin(usr.Role) {
		http.Redirect(w, r, "/admin/home", http.StatusFound)
		return
	}

	tmpl, _ := template.ParseFiles("./html/home.html")
	tmpl.Execute(w, usr)
}

func (app *App) loginHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		userID := r.FormValue("user_id")
		password := r.FormValue("password")
		hash := sha256.Sum256([]byte(password))

		usr, err := app.getUser(userID)
		if err != nil {
			if err == errDB {
				responseError(w)
				return
			}
			redirectWithError(w, r, "/login", "userIDまたはpasswordが間違っています")
			return
		}
		if usr.Password == base64.URLEncoding.EncodeToString(hash[:]) {
			sessionID, err := app.createSession(usr.ID)
			if err != nil {
				if err == errDB {
					responseError(w)
					return
				}
				http.Error(w, "サーバーエラーによりログイン失敗", http.StatusInternalServerError)
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
			redirectWithError(w, r, "/login", "userIDまたはpasswordが間違っています")
			return
		}

	} else {
		errorMessage := r.URL.Query().Get("error")
		tmpl, _ := template.ParseFiles("./html/login.html")
		tmpl.Execute(w, map[string]string{"Error": errorMessage})
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
		if err == errDB {
			responseError(w)
			return
		}
		http.Error(w, "サーバーエラーによりログアウト失敗", http.StatusInternalServerError)
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

// /admin/*へのアクセスは、adminアカウントのみ許可する。
func (app *App) adminMiddleware(handler func(http.ResponseWriter, *http.Request, user)) http.HandlerFunc {
	responseForbidden := func(w http.ResponseWriter) {
		http.Error(w, "Forbidden", http.StatusForbidden)
	}

	return func(w http.ResponseWriter, r *http.Request) {
		userID, err := app.getUserIDFromCookie(r)
		if err != nil {
			if err == errDB {
				responseError(w)
				return
			}
			responseForbidden(w)
			return
		}

		usr, err := app.getUser(userID)
		if err != nil {
			if err == errDB {
				responseError(w)
				return
			}
			responseForbidden(w)
			return
		}
		if !isAdmin(usr.Role) {
			responseForbidden(w)
			return
		}
		handler(w, r, usr)
	}
}

func (app *App) adminHomeHandler(w http.ResponseWriter, r *http.Request, usr user) {
	tmpl, _ := template.ParseFiles("./html/admin/home.html")
	tmpl.Execute(w, usr)
}

func (app *App) adminRegisterHandler(w http.ResponseWriter, r *http.Request, usr user) {
	if r.Method == http.MethodPost {
		userID := r.FormValue("user_id")
		password := r.FormValue("password")
		hash := sha256.Sum256([]byte(password))
		hashedPassword := base64.URLEncoding.EncodeToString(hash[:])

		err := app.createUser(user{ID: userID, Password: hashedPassword})
		if err != nil {
			if err == errDB {
				responseError(w)
				return
			}
			// useridが既に存在している場合、/admin/registerにリダイレクトして
			// エラーメッセージを表示。
			// 現在はぜんぶdbエラーでresponseする
			return
		}
		http.Redirect(w, r, "/admin/users", http.StatusFound)
		return

	} else {
		tmpl, _ := template.ParseFiles("./html/admin/register.html")
		tmpl.Execute(w, nil)
	}
}

func (app *App) adminUsersHandler(w http.ResponseWriter, r *http.Request, usr user) {
	users, err := app.getAllUsers()
	if err != nil {
		if err == errDB {
			responseError(w)
			return
		}
		// 現在はぜんぶdbエラーでresponseする
		return
	}

	tmpl, _ := template.ParseFiles("./html/admin/users.html")
	tmpl.Execute(w, users)
}

func (app *App) adminDeleteHandler(w http.ResponseWriter, r *http.Request, usr user) {
	if r.Method == http.MethodPost {
		deletedUserID := r.FormValue("user_id")
		err := app.deleteUser(deletedUserID)
		if err != nil {
			if err == errDB {
				responseError(w)
				return
			}
			// 現在はぜんぶdbエラーでresponseする
			return
		}
		http.Redirect(w, r, "/admin/users", http.StatusFound)
	} else {
		http.Redirect(w, r, "/admin/users", http.StatusFound)
	}
}

func (app *App) getUserIDFromCookie(r *http.Request) (string, error) {
	cookie, err := r.Cookie("session_id")
	if err != nil {
		return "", err
	}

	userID, err := app.getUserIDFromSession(cookie.Value)
	if err != nil {
		return "", err
	}
	return userID, nil
}

func responseError(w http.ResponseWriter) {
	http.Error(w, "サーバーでデータベースエラー", http.StatusInternalServerError)
}

func redirectWithError(w http.ResponseWriter, r *http.Request, path string, errMessage string) {
	redirectURL := fmt.Sprintf("%s?error=%s", path, url.QueryEscape(errMessage))
	http.Redirect(w, r, redirectURL, http.StatusFound)
	return
}
