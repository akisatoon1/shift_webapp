package handler

import (
	"backend/auth"
	"backend/context"
	"backend/db"
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
	"time"

	"github.com/gorilla/sessions"
	"golang.org/x/crypto/bcrypt"
)

// JSON形式のレスポンスを評価するヘルパー関数
func AssertRes(t *testing.T, got []byte, wantJSON string) {
	t.Helper()
	var gotInterface interface{}
	var wantInterface interface{}

	if err := json.Unmarshal(got, &gotInterface); err != nil {
		t.Fatalf("got json decode error: %v\ngot response body: %s", err, string(got))
	}

	if err := json.Unmarshal([]byte(wantJSON), &wantInterface); err != nil {
		t.Fatalf("want json decode error: %v", err)
	}

	if !reflect.DeepEqual(gotInterface, wantInterface) {
		t.Errorf("\ngot  %#v\nwant %#v\ngot response body: %s", gotInterface, wantInterface, string(got))
	}
}

// HTTPステータスコードを評価するヘルパー関数
func AssertCode(t *testing.T, got, want int, body []byte) {
	t.Helper()
	if got != want {
		t.Fatalf("want status code %d, got %d\nresponse body: %s", want, got, string(body))
	}
}

// ログインしてCookieを取得するヘルパー関数
func getLoginCookies(appCtx *context.AppContext, loginID, password string) []*http.Cookie {
	loginBody := map[string]string{
		"login_id": loginID,
		"password": password,
	}
	jsonBody, _ := json.Marshal(loginBody)
	loginReq := httptest.NewRequest("POST", "/login", bytes.NewBuffer(jsonBody))
	loginReq.Header.Set("Content-Type", "application/json")
	loginW := httptest.NewRecorder()
	loginMux := setHandlerToEndpoint(appCtx, "POST /login", LoginRequest)
	loginMux.ServeHTTP(loginW, loginReq)
	return loginW.Result().Cookies()
}

// リクエストにCookieをセットするヘルパー関数
func addCookiesToRequest(req *http.Request, cookies []*http.Cookie) {
	for _, cookie := range cookies {
		req.AddCookie(cookie)
	}
}

func newTestContext(requests []db.Request, users []db.User, entries []db.Entry) *context.AppContext {
	return context.NewAppContext(db.NewMockDB(requests, users, entries), sessions.NewCookieStore([]byte("test-secret")))
}

// 1つのAPIエンドポイントに、1つのハンドラーをセットする
func setHandlerToEndpoint(appCtx *context.AppContext, endpoint string, handlerFn HandlerFuncWithContext) *http.ServeMux {
	mux := http.NewServeMux()
	mux.Handle(endpoint, NewHandler(appCtx, handlerFn))
	return mux
}

func TestGetRequestsHandler(t *testing.T) {
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte("password"), bcrypt.DefaultCost)
	appCtx := newTestContext(
		[]db.Request{
			{ID: 1, CreatorID: 2, StartDate: db.DateOnly(time.Date(2024, 6, 1, 0, 0, 0, 0, time.UTC)), EndDate: db.DateOnly(time.Date(2024, 6, 1, 0, 0, 0, 0, time.UTC)), Deadline: db.DateTime(time.Date(2024, 6, 1, 0, 0, 0, 0, time.UTC)), CreatedAt: db.DateTime(time.Date(2024, 6, 1, 0, 0, 0, 0, time.UTC))},
			{ID: 2, CreatorID: 2, StartDate: db.DateOnly(time.Date(2024, 6, 1, 0, 0, 0, 0, time.UTC)), EndDate: db.DateOnly(time.Date(2024, 6, 1, 0, 0, 0, 0, time.UTC)), Deadline: db.DateTime(time.Date(2024, 6, 1, 0, 0, 0, 0, time.UTC)), CreatedAt: db.DateTime(time.Date(2024, 6, 1, 0, 0, 0, 0, time.UTC))},
		},
		[]db.User{
			{ID: 2, LoginID: "test_user", Password: string(hashedPassword), Name: "テストマネージャー", Role: auth.RoleManager},
		},
		[]db.Entry{},
	)
	mux := setHandlerToEndpoint(appCtx, "GET /requests", GetRequestsRequest)

	// ログイン用のCookieを取得
	cookies := getLoginCookies(appCtx, "test_user", "password")

	req := httptest.NewRequest("GET", "/requests", nil)
	addCookiesToRequest(req, cookies)
	w := httptest.NewRecorder()
	mux.ServeHTTP(w, req)

	AssertCode(t, w.Code, http.StatusOK, w.Body.Bytes())

	wantJSON := `
	[
		{
			"id": 1,
			"creator": {"id": 2, "name": "テストマネージャー"},
			"start_date": "2024-06-01",
			"end_date": "2024-06-01",
			"deadline": "2024-06-01 00:00:00",
			"created_at": "2024-06-01 00:00:00"
		},
		{
			"id": 2,
			"creator": {"id": 2, "name": "テストマネージャー"},
			"start_date": "2024-06-01",
			"end_date": "2024-06-01",
			"deadline": "2024-06-01 00:00:00",
			"created_at": "2024-06-01 00:00:00"
		}
	]
	`
	AssertRes(t, w.Body.Bytes(), wantJSON)
}

func TestGetRequestHandler(t *testing.T) {
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte("password"), bcrypt.DefaultCost)
	appCtx := newTestContext(
		[]db.Request{
			{ID: 1, CreatorID: 3, StartDate: db.DateOnly(time.Date(2024, 6, 1, 0, 0, 0, 0, time.UTC)), EndDate: db.DateOnly(time.Date(2024, 6, 1, 0, 0, 0, 0, time.UTC)), Deadline: db.DateTime(time.Date(2024, 6, 1, 0, 0, 0, 0, time.UTC)), CreatedAt: db.DateTime(time.Date(2024, 6, 1, 0, 0, 0, 0, time.UTC))},
		},
		[]db.User{
			{ID: 1, LoginID: "test_user1", Password: string(hashedPassword), Name: "テストユーザー1", Role: auth.RoleEmployee},
			{ID: 2, LoginID: "test_user2", Password: string(hashedPassword), Name: "テストユーザー2", Role: auth.RoleEmployee},
			{ID: 3, LoginID: "test_user3", Password: string(hashedPassword), Name: "テストマネージャー", Role: auth.RoleManager},
		},
		[]db.Entry{
			{ID: 1, RequestID: 1, UserID: 1, Date: db.DateOnly(time.Date(2024, 6, 1, 0, 0, 0, 0, time.UTC)), Hour: 8},
			{ID: 2, RequestID: 1, UserID: 2, Date: db.DateOnly(time.Date(2024, 6, 1, 0, 0, 0, 0, time.UTC)), Hour: 8},
		},
	)
	mux := setHandlerToEndpoint(appCtx, "GET /requests/{id}", GetRequestRequest)

	// ログイン用のCookieを取得
	cookies := getLoginCookies(appCtx, "test_user1", "password")

	req := httptest.NewRequest("GET", "/requests/1", nil)
	addCookiesToRequest(req, cookies)
	w := httptest.NewRecorder()
	mux.ServeHTTP(w, req)

	AssertCode(t, w.Code, http.StatusOK, w.Body.Bytes())

	wantJSON := `
	{
		"id": 1,
		"creator": {"id": 3, "name": "テストマネージャー"},
		"start_date": "2024-06-01",
		"end_date": "2024-06-01",
		"deadline": "2024-06-01 00:00:00",
		"created_at": "2024-06-01 00:00:00",
		"entries": [
			{
				"id": 1,
				"user": {"id": 1, "name": "テストユーザー1"},
				"date": "2024-06-01",
				"hour": 8
			},
			{
				"id": 2,
				"user": {"id": 2, "name": "テストユーザー2"},
				"date": "2024-06-01",
				"hour": 8
			}
		]
	}
	`
	AssertRes(t, w.Body.Bytes(), wantJSON)
}

func TestPostRequestsHandler(t *testing.T) {
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte("password"), bcrypt.DefaultCost)
	appCtx := newTestContext(
		[]db.Request{},
		[]db.User{
			{ID: 1, LoginID: "test_user", Password: string(hashedPassword), Name: "テストマネージャー", Role: auth.RoleManager},
		},
		[]db.Entry{},
	)
	mux := setHandlerToEndpoint(appCtx, "POST /requests", PostRequestsRequest)

	// ログイン用のCookieを取得
	cookies := getLoginCookies(appCtx, "test_user", "password")

	requestBody := map[string]string{
		"start_date": "2024-06-01",
		"end_date":   "2024-06-30",
		"deadline":   "2024-05-25 00:00:00",
	}
	body, _ := json.Marshal(requestBody)

	req := httptest.NewRequest("POST", "/requests", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	addCookiesToRequest(req, cookies)
	w := httptest.NewRecorder()
	mux.ServeHTTP(w, req)

	AssertCode(t, w.Code, http.StatusCreated, w.Body.Bytes())

	wantJSON := `
	{
		"id": 1
	}
	`
	AssertRes(t, w.Body.Bytes(), wantJSON)
}

func TestPostEntriesHandler(t *testing.T) {
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte("password"), bcrypt.DefaultCost)
	appCtx := newTestContext(
		[]db.Request{
			{ID: 1, CreatorID: 2, StartDate: db.DateOnly(time.Date(2024, 6, 1, 0, 0, 0, 0, time.UTC)), EndDate: db.DateOnly(time.Date(2024, 6, 1, 0, 0, 0, 0, time.UTC)), Deadline: db.DateTime(time.Date(2024, 6, 1, 0, 0, 0, 0, time.UTC)), CreatedAt: db.DateTime(time.Date(2024, 6, 1, 0, 0, 0, 0, time.UTC))},
		},
		[]db.User{
			{ID: 1, LoginID: "test_user", Password: string(hashedPassword), Name: "テストユーザー", Role: auth.RoleEmployee},
		},
		[]db.Entry{},
	)
	mux := setHandlerToEndpoint(appCtx, "POST /requests/{id}/entries", PostEntriesRequest)

	// ログイン用のCookieを取得
	cookies := getLoginCookies(appCtx, "test_user", "password")

	requestBody := []map[string]interface{}{
		{
			"date": "2024-06-01",
			"hour": 8,
		},
		{
			"date": "2024-06-01",
			"hour": 9,
		},
	}

	body, _ := json.Marshal(requestBody)

	req := httptest.NewRequest("POST", "/requests/1/entries", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	addCookiesToRequest(req, cookies)
	w := httptest.NewRecorder()
	mux.ServeHTTP(w, req)

	AssertCode(t, w.Code, http.StatusCreated, w.Body.Bytes())

	wantJSON := `
	{
		"id": 1,
		"entries": [{"id": 1}, {"id": 2}]
	}
	`
	AssertRes(t, w.Body.Bytes(), wantJSON)
}

func TestLoginHandler(t *testing.T) {
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte("password"), bcrypt.DefaultCost)
	appCtx := newTestContext(
		[]db.Request{},
		[]db.User{
			{ID: 1, LoginID: "test_user", Password: string(hashedPassword), Name: "テストユーザー"},
		},
		[]db.Entry{},
	)
	mux := setHandlerToEndpoint(appCtx, "POST /login", LoginRequest)

	// --- 正常系 ---
	body := map[string]string{
		"login_id": "test_user",
		"password": "password",
	}
	jsonBody, _ := json.Marshal(body)
	req := httptest.NewRequest("POST", "/login", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	mux.ServeHTTP(w, req)

	AssertCode(t, w.Code, http.StatusOK, w.Body.Bytes())

	// セッションCookieがセットされているか
	cookies := w.Result().Cookies()
	found := false
	for _, c := range cookies {
		if c.Name == "login_session" && c.Value != "" {
			found = true
		}
	}
	if !found {
		t.Errorf("login_session cookieがセットされていません")
	}

	// --- 異常系: パスワード間違い ---
	body2 := map[string]string{
		"login_id": "test_user",
		"password": "wrong",
	}
	jsonBody2, _ := json.Marshal(body2)
	req2 := httptest.NewRequest("POST", "/login", bytes.NewBuffer(jsonBody2))
	req2.Header.Set("Content-Type", "application/json")
	w2 := httptest.NewRecorder()
	mux.ServeHTTP(w2, req2)
	AssertCode(t, w2.Code, http.StatusUnauthorized, w2.Body.Bytes())
}

func TestLogoutHandler(t *testing.T) {
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte("password"), bcrypt.DefaultCost)
	appCtx := newTestContext(
		[]db.Request{},
		[]db.User{
			{ID: 1, LoginID: "test_user", Password: string(hashedPassword), Name: "テストユーザー"},
		},
		[]db.Entry{},
	)
	mux := setHandlerToEndpoint(appCtx, "DELETE /session", LogoutRequest)

	// まずloginしてCookie取得
	cookies := getLoginCookies(appCtx, "test_user", "password")

	// loginで得たCookieをlogoutリクエストに付与
	logoutReq := httptest.NewRequest("DELETE", "/session", nil)
	addCookiesToRequest(logoutReq, cookies)
	logoutW := httptest.NewRecorder()
	mux.ServeHTTP(logoutW, logoutReq)
	AssertCode(t, logoutW.Code, http.StatusOK, logoutW.Body.Bytes())

	// セッションが無効化されているか
	// MaxAge=-1のCookieが返る
	found := false
	for _, c := range logoutW.Result().Cookies() {
		if c.Name == "login_session" && c.MaxAge == -1 {
			found = true
		}
	}
	if !found {
		t.Errorf("logout後、login_session Cookieが無効化されていません")
	}
}
