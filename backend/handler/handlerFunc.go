package handler

import (
	"backend/auth"
	"backend/context"
	"backend/model"
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
)

/*
	APIエンドポイントに対応するハンドラー関数
*/

var ErrNotLoggedIn = errors.New("user not logged in")

func LoginRequest(ctx *context.AppContext, w http.ResponseWriter, r *http.Request) *AppError {
	type Body struct {
		LoginID  string `json:"login_id"`
		Password string `json:"password"`
	}

	var body Body
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		return NewAppError(err, "リクエストボディのデコードに失敗しました", http.StatusBadRequest)
	}

	err := auth.Login(ctx, w, r, body.LoginID, body.Password)
	if err != nil {
		if errors.Is(err, auth.ErrIncorrectAuth) {
			return NewAppError(err, "ログインIDまたはパスワードが間違っています", http.StatusUnauthorized)
		}
		return NewAppError(err, "ログインに失敗しました", http.StatusInternalServerError)
	}
	return nil
}

func LogoutRequest(ctx *context.AppContext, w http.ResponseWriter, r *http.Request) *AppError {
	err := auth.Logout(ctx, w, r)
	if err != nil {
		return NewAppError(err, "ログアウトに失敗しました", http.StatusInternalServerError)
	}
	return nil
}

func GetRequestsRequest(ctx *context.AppContext, w http.ResponseWriter, r *http.Request) *AppError {
	// ログインユーザのみ認可
	if _, isLoggedIn := auth.GetUserID(ctx, r); !isLoggedIn {
		return NewAppError(ErrNotLoggedIn, "ログインしていません", http.StatusUnauthorized)
	}

	requests, err := model.GetRequests(ctx)
	if err != nil {
		return NewAppError(err, "シフトリクエストの取得に失敗しました", http.StatusInternalServerError)
	}
	json.NewEncoder(w).Encode(requests)
	return nil
}

func GetRequestRequest(ctx *context.AppContext, w http.ResponseWriter, r *http.Request) *AppError {
	// ログインユーザのみ認可
	if _, isLoggedIn := auth.GetUserID(ctx, r); !isLoggedIn {
		return NewAppError(ErrNotLoggedIn, "ログインしていません", http.StatusUnauthorized)
	}

	requestId := r.PathValue("id")
	requestIdInt, err := strconv.Atoi(requestId)
	if err != nil {
		return NewAppError(err, "requestIdが整数ではありません", http.StatusBadRequest)
	}

	response, err := model.GetRequest(ctx, requestIdInt)
	if err != nil {
		return NewAppError(err, "リクエストの取得に失敗しました", http.StatusInternalServerError)
	}
	json.NewEncoder(w).Encode(response)
	return nil
}

func PostRequestsRequest(ctx *context.AppContext, w http.ResponseWriter, r *http.Request) *AppError {
	// ログインしているユーザーのIDを取得する
	userID, isLoggedIn := auth.GetUserID(ctx, r)
	if !isLoggedIn {
		return NewAppError(ErrNotLoggedIn, "ログインしていません", http.StatusUnauthorized)
	}

	// リクエストボディの形式を定義する
	type Body struct {
		StartDate string `json:"start_date"`
		EndDate   string `json:"end_date"`
		Deadline  string `json:"deadline"`
	}

	// リクエストボディのバリデーション
	var body Body
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		return NewAppError(err, "リクエストボディのデコードに失敗しました", http.StatusBadRequest)
	}

	// 新しいシフトリクエストを作成する
	response, err := model.CreateRequest(ctx, model.NewRequest{
		CreatorID: userID,
		StartDate: body.StartDate,
		EndDate:   body.EndDate,
		Deadline:  body.Deadline,
	})
	if err != nil {
		if errors.Is(err, model.ErrForbidden) {
			return NewAppError(err, "権限がありません", http.StatusForbidden)
		}

		var inputErr model.InputError
		if errors.As(err, &inputErr) {
			return NewAppError(inputErr, inputErr.Message(), http.StatusBadRequest)
		}

		return NewAppError(err, "シフトリクエストの作成に失敗しました", http.StatusInternalServerError)
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(response)
	return nil
}

func PostEntriesRequest(ctx *context.AppContext, w http.ResponseWriter, r *http.Request) *AppError {
	// ログインしているユーザーのIDを取得する
	userID, isLoggedIn := auth.GetUserID(ctx, r)
	if !isLoggedIn {
		return NewAppError(ErrNotLoggedIn, "ログインしていません", http.StatusUnauthorized)
	}

	// シフトリクエストのIDを取得する
	// 整数ではない場合はエラーを返す
	requestId := r.PathValue("id")
	requestIdInt, err := strconv.Atoi(requestId)
	if err != nil {
		return NewAppError(err, "requestidが整数ではありません", http.StatusBadRequest)
	}

	// リクエストボディの形式を定義する
	type Body []struct {
		Date string `json:"date"`
		Hour int    `json:"hour"`
	}

	// リクエストボディのバリデーション
	var body Body
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		return NewAppError(err, "リクエストボディのデコードに失敗しました", http.StatusBadRequest)
	}

	// モデルに渡す形に変換する
	entries := model.NewEntries{
		ID:      requestIdInt,
		UserID:  userID,
		Entries: []model.NewEntry{},
	}
	for _, entry := range body {
		entries.Entries = append(entries.Entries, model.NewEntry{
			Date: entry.Date,
			Hour: entry.Hour,
		})
	}

	// エントリーを作成する
	response, err := model.CreateEntries(ctx, entries)
	if err != nil {
		if errors.Is(err, model.ErrForbidden) {
			return NewAppError(err, "権限がありません", http.StatusForbidden)
		}

		var inputErr model.InputError
		if errors.As(err, &inputErr) {
			return NewAppError(inputErr, inputErr.Message(), http.StatusBadRequest)
		}

		return NewAppError(err, "エントリーの作成に失敗しました", http.StatusInternalServerError)
	}

	// レスポンスを返す
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(response)
	return nil
}
