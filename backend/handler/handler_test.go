package handler

import (
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
)

// TODO
// テスト失敗時のレスポンスボディのエラーメッセージを表示する

// JSON形式のレスポンスを評価するヘルパー関数
func AssertRes(t *testing.T, got []byte, wantJSON string) {
	t.Helper()
	var gotInterface interface{}
	var wantInterface interface{}

	if err := json.Unmarshal(got, &gotInterface); err != nil {
		t.Fatalf("got json decode error: %v", err)
	}

	if err := json.Unmarshal([]byte(wantJSON), &wantInterface); err != nil {
		t.Fatalf("want json decode error: %v", err)
	}

	if !reflect.DeepEqual(gotInterface, wantInterface) {
		t.Errorf("\ngot  %#v\nwant %#v", gotInterface, wantInterface)
	}
}

// HTTPステータスコードを評価するヘルパー関数
func AssertCode(t *testing.T, got, want int) {
	t.Helper()
	if got != want {
		t.Fatalf("want status code %d, got %d", want, got)
	}
}

func newTestContext(requests []db.Request, users []db.User, entries []db.Entry) *context.AppContext {
	return &context.AppContext{
		DB:     db.NewMockDB(requests, users, entries),
		Cookie: sessions.NewCookieStore([]byte("test-secret")),
	}
}

// 1つのAPIエンドポイントに、1つのハンドラーをセットする
func setHandlerToEndpoint(appCtx *context.AppContext, endpoint string, handlerFn HandlerFuncWithContext) *http.ServeMux {
	mux := http.NewServeMux()
	mux.Handle(endpoint, NewHandler(appCtx, handlerFn))
	return mux
}

func TestGetRequestsHandler(t *testing.T) {
	appCtx := newTestContext(
		[]db.Request{
			{ID: 1, CreatorID: 2, StartDate: time.Date(2024, 6, 1, 0, 0, 0, 0, time.UTC), EndDate: time.Date(2024, 6, 1, 0, 0, 0, 0, time.UTC), Deadline: time.Date(2024, 6, 1, 0, 0, 0, 0, time.UTC), CreatedAt: time.Date(2024, 6, 1, 0, 0, 0, 0, time.UTC)},
			{ID: 2, CreatorID: 2, StartDate: time.Date(2024, 6, 1, 0, 0, 0, 0, time.UTC), EndDate: time.Date(2024, 6, 1, 0, 0, 0, 0, time.UTC), Deadline: time.Date(2024, 6, 1, 0, 0, 0, 0, time.UTC), CreatedAt: time.Date(2024, 6, 1, 0, 0, 0, 0, time.UTC)},
		},
		[]db.User{
			{ID: 2, Name: "テストマネージャー"},
		},
		[]db.Entry{},
	)
	mux := setHandlerToEndpoint(appCtx, "GET /requests", GetRequestsRequest)

	req := httptest.NewRequest("GET", "/requests", nil)
	w := httptest.NewRecorder()
	mux.ServeHTTP(w, req)

	AssertCode(t, w.Code, http.StatusOK)

	wantJSON := `
	[
		{
			"id": 1,
			"creator": {"id": 2, "name": "テストマネージャー"},
			"start_date": "2024-06-01",
			"end_date": "2024-06-01",
			"deadline": "2024-06-01",
			"created_at": "2024-06-01 00:00:00"
		},
		{
			"id": 2,
			"creator": {"id": 2, "name": "テストマネージャー"},
			"start_date": "2024-06-01",
			"end_date": "2024-06-01",
			"deadline": "2024-06-01",
			"created_at": "2024-06-01 00:00:00"
		}
	]
	`
	AssertRes(t, w.Body.Bytes(), wantJSON)
}

func TestGetEntriesHandler(t *testing.T) {
	appCtx := newTestContext(
		[]db.Request{
			{ID: 1, CreatorID: 2, StartDate: time.Date(2024, 6, 1, 0, 0, 0, 0, time.UTC), EndDate: time.Date(2024, 6, 1, 0, 0, 0, 0, time.UTC), Deadline: time.Date(2024, 6, 1, 0, 0, 0, 0, time.UTC), CreatedAt: time.Date(2024, 6, 1, 0, 0, 0, 0, time.UTC)},
		},
		[]db.User{
			{ID: 1, Name: "テストユーザー1"},
			{ID: 2, Name: "テストユーザー2"},
		},
		[]db.Entry{
			{ID: 1, RequestID: 1, UserID: 1, Date: time.Date(2024, 6, 1, 0, 0, 0, 0, time.UTC), Hour: 8},
			{ID: 2, RequestID: 1, UserID: 2, Date: time.Date(2024, 6, 1, 0, 0, 0, 0, time.UTC), Hour: 8},
		},
	)
	mux := setHandlerToEndpoint(appCtx, "GET /requests/{id}/entries", GetEntriesRequest)

	req := httptest.NewRequest("GET", "/requests/1/entries", nil)
	w := httptest.NewRecorder()
	mux.ServeHTTP(w, req)

	AssertCode(t, w.Code, http.StatusOK)

	wantJSON := `
	{
		"id": 1,
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
	appCtx := newTestContext(
		[]db.Request{},
		[]db.User{},
		[]db.Entry{},
	)
	mux := setHandlerToEndpoint(appCtx, "POST /requests", PostRequestsRequest)

	requestBody := map[string]string{
		"start_date": "2024-06-01",
		"end_date":   "2024-06-30",
		"deadline":   "2024-05-25",
	}
	body, _ := json.Marshal(requestBody)

	req := httptest.NewRequest("POST", "/requests", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	mux.ServeHTTP(w, req)

	AssertCode(t, w.Code, http.StatusCreated)

	wantJSON := `
	{
		"id": 1
	}
	`
	AssertRes(t, w.Body.Bytes(), wantJSON)
}

func TestPostEntriesHandler(t *testing.T) {
	appCtx := newTestContext(
		[]db.Request{
			{ID: 1, CreatorID: 2, StartDate: time.Date(2024, 6, 1, 0, 0, 0, 0, time.UTC), EndDate: time.Date(2024, 6, 1, 0, 0, 0, 0, time.UTC), Deadline: time.Date(2024, 6, 1, 0, 0, 0, 0, time.UTC), CreatedAt: time.Date(2024, 6, 1, 0, 0, 0, 0, time.UTC)},
		},
		[]db.User{},
		[]db.Entry{},
	)
	mux := setHandlerToEndpoint(appCtx, "POST /requests/{id}/entries", PostEntriesRequest)

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
	w := httptest.NewRecorder()
	mux.ServeHTTP(w, req)

	AssertCode(t, w.Code, http.StatusCreated)

	wantJSON := `
	{
		"id": 1,
		"entries": [{"id": 1}, {"id": 2}]
	}
	`
	AssertRes(t, w.Body.Bytes(), wantJSON)
}

func TestLoginHandler(t *testing.T) {
	appCtx := newTestContext(
		[]db.Request{},
		[]db.User{
			{ID: 1, LoginID: "test_user", Password: "password", Name: "テストユーザー"},
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

	AssertCode(t, w.Code, http.StatusOK)

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
	AssertCode(t, w2.Code, http.StatusUnauthorized)
}

func TestLogoutHandler(t *testing.T) {
	appCtx := newTestContext(
		[]db.Request{},
		[]db.User{
			{ID: 1, LoginID: "test_user", Password: "password", Name: "テストユーザー"},
		},
		[]db.Entry{},
	)
	mux := setHandlerToEndpoint(appCtx, "DELETE /session", LogoutRequest)

	// まずloginしてCookie取得
	loginBody := map[string]string{
		"login_id": "test_user",
		"password": "password",
	}
	jsonBody, _ := json.Marshal(loginBody)
	loginReq := httptest.NewRequest("POST", "/login", bytes.NewBuffer(jsonBody))
	loginReq.Header.Set("Content-Type", "application/json")
	loginMux := setHandlerToEndpoint(appCtx, "POST /login", LoginRequest)
	loginW := httptest.NewRecorder()
	loginMux.ServeHTTP(loginW, loginReq)

	// loginで得たCookieをlogoutリクエストに付与
	logoutReq := httptest.NewRequest("DELETE", "/session", nil)
	for _, c := range loginW.Result().Cookies() {
		logoutReq.AddCookie(c)
	}
	logoutW := httptest.NewRecorder()
	mux.ServeHTTP(logoutW, logoutReq)
	AssertCode(t, logoutW.Code, http.StatusOK)

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
