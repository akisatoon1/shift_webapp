package main

import (
	"crypto/sha256"
	"encoding/base64"
	"errors"
	"fmt"
	"html/template"
	"net/http"
	"net/url"
	"time"
)

var (
	errUserNotFound    = errors.New("user not found")
	errSessionNotFound = errors.New("session not found")
	errInvalidUserID   = errors.New("invalid user id")
	errRequestNotFound = errors.New("request not found")
)

/*
	user management handlers
	ユーザ管理やログインに関する処理をするハンドラーたち。
*/

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

		_, err := app.createUser(user{ID: userID, Password: hashedPassword})
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

	userID, err := app.getUserIDFromSessionID(cookie.Value)
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

/*
	shift request handlers
	シフト要請に関する処理をするハンドラーたち。
*/

func (app *App) adminShowRequestsHandler(w http.ResponseWriter, r *http.Request, usr user) {
	requests, err := app.getAllRequests()
	if err != nil {
		responseServerError(w)
	}
	tmpl, _ := template.ParseFiles("./html/admin/requests.html")
	tmpl.Execute(w, requests)
}

func (app *App) requestCreatePageHandler(w http.ResponseWriter, r *http.Request, usr user) {
	tmpl, _ := template.ParseFiles("./html/admin/requests/create.html")
	tmpl.Execute(w, nil)
}

func (app *App) createRequestHandler(w http.ResponseWriter, r *http.Request, usr user) {
	start_date := r.FormValue("start_date")
	end_date := r.FormValue("end_date")
	_, err := app.createRequest(usr.ID, start_date, end_date)
	if err != nil {
		responseServerError(w)
		return
	}
	http.Redirect(w, r, "/admin/requests", http.StatusFound)
}

/*
	admin submission handler
	管理者アカウントの提出されたシフトに対する処理をするハンドラーたち。
*/

func (app *App) adminShowSubmissionsHandler(w http.ResponseWriter, r *http.Request, usr user) {}

func (app *App) showUserSubmissionHandler(w http.ResponseWriter, r *http.Request, usr user) {

}

/*
	user requests handler
	ユーザのシフト提出に関する処理をするハンドラーたち。
*/

// userがシフト要請の一覧を見る
func (app *App) showRequestsHandler(w http.ResponseWriter, r *http.Request) {
	requests, err := app.getAllRequests()
	if err != nil {
		responseServerError(w)
	}
	tmpl, err := template.ParseFiles("./html/requests.html")
	tmpl.Execute(w, requests)
}

func (app *App) submissionPageHandler(w http.ResponseWriter, r *http.Request) {
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

	// get request
	requestID := r.PathValue("request_id")
	req, err := app.getRequest(requestID)
	if err != nil {
		if err == errRequestNotFound {
			http.Error(w, "invalid request id", http.StatusBadRequest)
		} else {
			responseServerError(w)
		}
		return
	}

	// 提出済みかチェック
	subs, err := app.getSubmissionsByRequestAndUserID(requestID, userID)
	if len(subs) > 0 {
		http.Error(w, "この要請は既に提出済みです。", http.StatusBadRequest)
		return
	}

	// シフトの日付を得る
	var dates []string
	startDate, _ := time.Parse("2006-01-02", req.StartDate)
	endDate, _ := time.Parse("2006-01-02", req.EndDate)
	for date := startDate; !date.After(endDate); date = date.AddDate(0, 0, 1) {
		dates = append(dates, date.Format("2006-01-02"))
	}

	// シフトの選択できる時間帯を作成
	var hours []string
	for i := 9; i <= 22; i++ {
		hours = append(hours, fmt.Sprintf("%v", i))
	}

	tmpl, _ := template.ParseFiles("./html/requests/request/submit.html")
	tmpl.Execute(w, struct {
		RequestID string
		Dates     []string
		Hours     []string
	}{
		RequestID: requestID,
		Dates:     dates,
		Hours:     hours,
	})
}

func (app *App) submitShiftHandler(w http.ResponseWriter, r *http.Request) {
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

	requestID := r.PathValue("request_id")
	r.ParseForm()
	for submissionDate, hours := range r.Form {
		submissionID, err := app.createSubmission(requestID, userID, submissionDate)
		if err != nil {
			responseServerError(w)
			return
		}
		for _, hour := range hours {
			_, err = app.createEntry(submissionID, hour)
			if err != nil {
				responseServerError(w)
				return
			}
		}
	}
	http.Redirect(w, r, fmt.Sprintf("/requests/%v/submissions", requestID), http.StatusFound)
}
func (app *App) showSubmissionsHandler(w http.ResponseWriter, r *http.Request) {
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

	requestID := r.PathValue("request_id")
	shift := make(map[string][]int)
	submissions, err := app.getSubmissionsByRequestAndUserID(requestID, userID)
	if err != nil {
		responseServerError(w)
		return
	}
	for _, sub := range submissions {
		entries, err := app.getEntriesBySubmissionID(sub.ID)
		if err != nil {
			responseServerError(w)
			return
		}
		var hours []int
		for _, ent := range entries {
			hours = append(hours, ent.ShiftHour)
		}
		shift[sub.SubmissionDate] = hours
	}

	tmpl, err := template.ParseFiles("./html/requests/request/submissions.html")
	tmpl.Execute(w, shift)
}
