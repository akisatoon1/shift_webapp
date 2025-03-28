package main

import (
	"crypto/sha256"
	"encoding/base64"
	"errors"
	"fmt"
	"html/template"
	"net/http"
	"net/url"
)

var (
	errUserNotFound    = errors.New("user not found")
	errSessionNotFound = errors.New("session not found")
	errInvalidUserID   = errors.New("invalid user id")
)

// /home
// ユーザーのホーム画面を表示する
func (app *App) homeHandler(w http.ResponseWriter, r *http.Request) {
	// check session
	userID, err := app.getUserIDFromCookie(r)
	if err != nil {
		if err == errSessionNotFound {
			redirectWithError(w, r, "/login", "ログインしてください")
		} else {
			responseServerError(w)
		}
		return
	}

	// check role
	// adminならば、adminのホーム画面にリダイレクトする
	usr, err := app.getUser(userID)
	if err != nil {
		responseServerError(w)
		return
	}
	if isAdmin(usr.Role) {
		http.Redirect(w, r, "/admin/home", http.StatusFound)
		return
	}

	tmpl, _ := template.ParseFiles("./html/home.html")
	tmpl.Execute(w, usr)
}

// /login
// 2つの役割がある
// 1: loginページの表示をする
// 1: login処理をする
func (app *App) loginHandler(w http.ResponseWriter, r *http.Request) {
	// http methodがPOSTならばlogin処理をする
	if r.Method == http.MethodPost {
		userID := r.FormValue("user_id")
		password := r.FormValue("password")
		hash := sha256.Sum256([]byte(password))

		usr, err := app.getUser(userID)
		if err != nil {
			if err == errUserNotFound {
				redirectWithError(w, r, "/login", "userIDまたはpasswordが間違っています")
			} else {
				responseServerError(w)
			}
			return
		}
		// password check
		if usr.Password == base64.URLEncoding.EncodeToString(hash[:]) {
			// passwordが一致するので、sessionを始める
			// session_idをCookieでクライアントに送信
			// ホーム画面にリダイレクト
			sessionID, err := app.createSession(usr.ID)
			if err != nil {
				responseServerError(w)
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
			// paswordが一致しない時
			redirectWithError(w, r, "/login", "userIDまたはpasswordが間違っています")
			return
		}
	} else if r.Method == http.MethodGet {
		// http methodがGETなので、loginページを表示
		errorMessage := r.URL.Query().Get("error")
		tmpl, _ := template.ParseFiles("./html/login.html")
		tmpl.Execute(w, map[string]string{"Error": errorMessage})
	}
}

// /logout
// logout処理のみ
func (app *App) logoutHandler(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("session_id")
	if err != nil {
		http.Redirect(w, r, "/login", http.StatusFound)
		return
	}

	err = app.deleteSession(cookie.Value)
	if err != nil {
		responseServerError(w)
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
			if err == errSessionNotFound {
				responseForbidden(w)
			} else {
				responseServerError(w)
			}
			return
		}

		usr, err := app.getUser(userID)
		if err != nil {
			if err == errUserNotFound {
				responseForbidden(w)
			} else {
				responseServerError(w)
			}
			return
		}
		if !isAdmin(usr.Role) {
			responseForbidden(w)
			return
		}
		handler(w, r, usr)
	}
}

// /admin/home
// adminのホーム画面を表示する
func (app *App) adminHomeHandler(w http.ResponseWriter, r *http.Request, usr user) {
	tmpl, _ := template.ParseFiles("./html/admin/home.html")
	tmpl.Execute(w, usr)
}

// /admin/register
// 2つ役割がある
// 1: ユーザ作成画面を表示
// 2: adminではないユーザを新しく作成する
func (app *App) adminRegisterHandler(w http.ResponseWriter, r *http.Request, usr user) {
	// http methodがPOSTならばユーザを作成する
	if r.Method == http.MethodPost {
		userID := r.FormValue("user_id")
		password := r.FormValue("password")
		hash := sha256.Sum256([]byte(password))
		hashedPassword := base64.URLEncoding.EncodeToString(hash[:])

		err := app.createUser(user{ID: userID, Password: hashedPassword})
		if err != nil {
			if err == errInvalidUserID {
				http.Error(w, "無効なuser idです", http.StatusBadRequest)
			} else {
				responseServerError(w)
			}
			return
		}
		http.Redirect(w, r, "/admin/users", http.StatusFound)
		return

	} else if r.Method == http.MethodGet {
		// http methodがGETなのでユーザ作成画面を表示する
		tmpl, _ := template.ParseFiles("./html/admin/register.html")
		tmpl.Execute(w, nil)
	}
}

// /admin/users
// ユーザ管理画面を表示する
func (app *App) adminUsersHandler(w http.ResponseWriter, r *http.Request, usr user) {
	users, err := app.getAllUsers()
	if err != nil {
		responseServerError(w)
		return
	}

	tmpl, _ := template.ParseFiles("./html/admin/users.html")
	tmpl.Execute(w, users)
}

// /admin/delete
// ユーザを削除する
func (app *App) adminDeleteHandler(w http.ResponseWriter, r *http.Request, usr user) {
	deletedUserID := r.FormValue("user_id")
	err := app.deleteUser(deletedUserID)
	if err != nil {
		responseServerError(w)
		return
	}
	http.Redirect(w, r, "/admin/users", http.StatusFound)
}

// クライアントが認証されていることを確認する
func (app *App) getUserIDFromCookie(r *http.Request) (string, error) {
	cookie, err := r.Cookie("session_id")
	if err != nil {
		return "", errSessionNotFound
	}

	userID, err := app.getUserIDFromSession(cookie.Value)
	if err != nil {
		return "", err
	}
	return userID, nil
}

func responseServerError(w http.ResponseWriter) {
	http.Error(w, "サーバーでデータベースエラー", http.StatusInternalServerError)
}

func redirectWithError(w http.ResponseWriter, r *http.Request, path string, errMessage string) {
	redirectURL := fmt.Sprintf("%s?error=%s", path, url.QueryEscape(errMessage))
	http.Redirect(w, r, redirectURL, http.StatusFound)
	return
}
