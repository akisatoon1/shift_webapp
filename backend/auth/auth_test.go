package auth

import (
	"backend/context"
	"backend/db"
	"net/http/httptest"
	"testing"

	"github.com/gorilla/sessions"
	"golang.org/x/crypto/bcrypt"
)

func TestGetUserID(t *testing.T) {
	store := sessions.NewCookieStore([]byte("test-secret"))
	ctx := context.NewAppContext(nil, store)

	// --- 正常系: セッションにuser_idが入っている場合 ---
	req := httptest.NewRequest("GET", "/", nil)

	// セッションにuser_idをセット
	session, _ := store.Get(req, "login_session")
	session.Values["user_id"] = 42

	// セッションをレスポンスに保存する
	rr := httptest.NewRecorder()
	session.Save(req, rr)

	// レスポンスのCookieをリクエストに追加
	req2 := httptest.NewRequest("GET", "/", nil)
	for _, cookie := range rr.Result().Cookies() {
		req2.AddCookie(cookie)
	}

	// テスト実行
	userID, ok := GetUserID(ctx, req2)
	if !ok || userID != 42 {
		t.Errorf("want ok=true, userID=42, got ok=%v, userID=%v", ok, userID)
	}

	// --- 異常系: セッションが存在しない場合 ---
	req3 := httptest.NewRequest("GET", "/", nil)
	userID2, ok2 := GetUserID(ctx, req3)
	if ok2 || userID2 != -1 {
		t.Errorf("want ok=false, userID=-1, got ok=%v, userID=%v", ok2, userID2)
	}
}

func newTestContext(user db.User, store *sessions.CookieStore) *context.AppContext {
	return context.NewAppContext(db.NewMockDB(nil, []db.User{user}, nil, nil), store)
}

func TestIsEmployee(t *testing.T) {
	// 従業員権限を持つユーザー
	user := db.User{ID: 10, Role: RoleEmployee}
	ctx := newTestContext(user, nil)
	ok, err := IsEmployee(ctx, 10)
	if err != nil || !ok {
		t.Errorf("want employee, got ok=%v, err=%v", ok, err)
	}

	// 権限なしユーザー
	user2 := db.User{ID: 11, Role: 0}
	ctx2 := newTestContext(user2, nil)
	ok2, err2 := IsEmployee(ctx2, 11)
	if err2 != nil || ok2 {
		t.Errorf("want not employee, got ok=%v, err=%v", ok2, err2)
	}
}

func TestIsManager(t *testing.T) {
	// マネージャー権限を持つユーザー
	user := db.User{ID: 20, Role: RoleManager}
	ctx := newTestContext(user, nil)
	ok, err := IsManager(ctx, 20)
	if err != nil || !ok {
		t.Errorf("want manager, got ok=%v, err=%v", ok, err)
	}

	// 権限なしユーザー
	user2 := db.User{ID: 21, Role: 0}
	ctx2 := newTestContext(user2, nil)
	ok2, err2 := IsManager(ctx2, 21)
	if err2 != nil || ok2 {
		t.Errorf("want not manager, got ok=%v, err=%v", ok2, err2)
	}
}

func TestLogout(t *testing.T) {
	store := sessions.NewCookieStore([]byte("test-secret"))
	ctx := context.NewAppContext(nil, store)

	// --- 正常系: セッションが存在する場合 ---
	req := httptest.NewRequest("GET", "/", nil)
	rr := httptest.NewRecorder()

	// セッションを作成し保存
	session, _ := store.Get(req, "login_session")
	session.Values["user_id"] = 123
	session.Save(req, rr)

	// Cookieを新しいリクエストにセット
	req2 := httptest.NewRequest("GET", "/", nil)
	for _, cookie := range rr.Result().Cookies() {
		req2.AddCookie(cookie)
	}
	rr2 := httptest.NewRecorder()

	err := Logout(ctx, rr2, req2)
	if err != nil {
		t.Errorf("Logout returned error: %v", err)
	}

	// セッションが無効化されているか確認
	session2, _ := store.Get(req2, "login_session")
	if session2.Options.MaxAge != -1 {
		t.Errorf("session MaxAge should be -1 after logout, got %d", session2.Options.MaxAge)
	}
}

func TestLogin(t *testing.T) {
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte("pass123"), bcrypt.DefaultCost)
	user := db.User{ID: 42, LoginID: "testuser", Password: string(hashedPassword), Role: 0}
	store := sessions.NewCookieStore([]byte("test-secret"))
	ctx := newTestContext(user, store)

	// --- 正常系: 正しいloginIDとpassword ---
	req := httptest.NewRequest("POST", "/login", nil)
	rr := httptest.NewRecorder()
	err := Login(ctx, rr, req, "testuser", "pass123")
	if err != nil {
		t.Errorf("Login returned error: %v", err)
	}

	// セッションにuser_idがセットされているか確認
	req2 := httptest.NewRequest("GET", "/", nil)
	for _, cookie := range rr.Result().Cookies() {
		req2.AddCookie(cookie)
	}
	session, _ := store.Get(req2, "login_session")
	userID, ok := session.Values["user_id"]
	if !ok || userID != 42 {
		t.Errorf("user_id should be 42 in session, got %v", userID)
	}

	// --- 異常系: loginIDが存在しない ---
	rr2 := httptest.NewRecorder()
	err2 := Login(ctx, rr2, req, "notfound", "pass123")
	if err2 == nil {
		t.Errorf("want error for invalid loginID, got nil")
	}

	// --- 異常系: パスワードが間違っている ---
	rr3 := httptest.NewRecorder()
	err3 := Login(ctx, rr3, req, "testuser", "wrongpass")
	if err3 == nil {
		t.Errorf("want error for wrong password, got nil")
	}
}
