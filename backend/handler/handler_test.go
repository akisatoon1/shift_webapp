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
)

// TODO
// テスト失敗時のレスポンスボディのエラーメッセージを表示する

/*
	jsonの比較がめんどうくさい...
*/

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

func newTestContext() *context.AppContext {
	return &context.AppContext{
		DB: db.InitMockDB(),
	}
}

// 1つのAPIエンドポイントに、1つのハンドラーをセットする
func setHandlerToEndpoint(appCtx *context.AppContext, endpoint string, handlerFn HandlerFuncWithContext) *http.ServeMux {
	mux := http.NewServeMux()
	mux.Handle(endpoint, NewHandler(appCtx, handlerFn))
	return mux
}

func TestGetRequestsHandler(t *testing.T) {
	appCtx := newTestContext()
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
	appCtx := newTestContext()
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
				"user": {"id": 1, "name": "テストユーザー"},
				"date": "2024-06-01",
				"hour": 8
			},
			{
				"id": 2,
				"user": {"id": 2, "name": "テストマネージャー"},
				"date": "2024-06-01",
				"hour": 8
			}
		]
	}
	`
	AssertRes(t, w.Body.Bytes(), wantJSON)
}

func TestPostRequestsHandler(t *testing.T) {
	appCtx := newTestContext()
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
		"id": 3
	}
	`
	AssertRes(t, w.Body.Bytes(), wantJSON)
}

func TestPostEntriesHandler(t *testing.T) {
	appCtx := newTestContext()
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
		"entries": [{"id": 5}, {"id": 6}]
	}
	`
	AssertRes(t, w.Body.Bytes(), wantJSON)
}
