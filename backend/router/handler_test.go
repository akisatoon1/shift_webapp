package router

import (
	"backend/context"
	"backend/db"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
)

/*
	jsonの比較がめんどうくさい...
*/

func newTestContext() *context.AppContext {
	return &context.AppContext{
		DB: db.InitMockDB(),
	}
}

func setupTestMux(appCtx *context.AppContext) *http.ServeMux {
	mux := http.NewServeMux()
	Routes(mux, appCtx)
	return mux
}

func TestGetRequestsHandler(t *testing.T) {
	appCtx := newTestContext()
	mux := setupTestMux(appCtx)

	req := httptest.NewRequest("GET", "/requests", nil)
	w := httptest.NewRecorder()
	mux.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("want 200, got %d", w.Code)
	}

	var got interface{}
	if err := json.Unmarshal(w.Body.Bytes(), &got); err != nil {
		t.Fatalf("json decode error: %v", err)
	}

	wantJSON := `
	[
		{
			"id": 1,
			"creator": {"id": 2, "name": "テストマネージャー"},
			"start_date": "2024-06-01",
			"end_date": "2024-06-01",
			"deadline": "2024-06-01",
			"created_at": "2024-06-01"
		},
		{
			"id": 2,
			"creator": {"id": 2, "name": "テストマネージャー"},
			"start_date": "2024-06-01",
			"end_date": "2024-06-01",
			"deadline": "2024-06-01",
			"created_at": "2024-06-01"
		}
	]
	`
	var want interface{}
	if err := json.Unmarshal([]byte(wantJSON), &want); err != nil {
		t.Fatalf("want json decode error: %v", err)
	}

	if !reflect.DeepEqual(got, want) {
		t.Errorf("\ngot  %#v\nwant %#v", got, want)
	}
}

func TestGetEntriesHandler(t *testing.T) {
	appCtx := newTestContext()
	mux := setupTestMux(appCtx)

	// requestID=1のエントリーを取得
	req := httptest.NewRequest("GET", "/requests/1/entries", nil)
	w := httptest.NewRecorder()
	mux.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("want 200, got %d", w.Code)
	}

	var got map[string]interface{}
	if err := json.Unmarshal(w.Body.Bytes(), &got); err != nil {
		t.Fatalf("json decode error: %v", err)
	}

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
	var want interface{}
	if err := json.Unmarshal([]byte(wantJSON), &want); err != nil {
		t.Fatalf("want json decode error: %v", err)
	}

	if !reflect.DeepEqual(got, want) {
		t.Errorf("\ngot  %#v\nwant %#v", got, want)
	}
}
