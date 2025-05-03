package handler

import (
	"backend/auth"
	"backend/context"
	"backend/model"
	"encoding/json"
	"net/http"
	"strconv"
)

/*
	APIエンドポイントに対応するハンドラー関数
*/

func LoginRequest(ctx *context.AppContext, w http.ResponseWriter, r *http.Request) *AppError {
	type Body struct {
		LoginID  string `json:"login_id"`
		Password string `json:"password"`
	}

	var body Body
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		return NewAppError(err, "LoginRequest: リクエストボディのデコードに失敗しました", http.StatusBadRequest)
	}

	err := auth.Login(ctx, w, r, body.LoginID, body.Password)
	if err != nil {
		return NewAppError(err, "LoginRequest: ログインに失敗しました", http.StatusUnauthorized)
	}
	return nil
}

func LogoutRequest(ctx *context.AppContext, w http.ResponseWriter, r *http.Request) *AppError {
	err := auth.Logout(ctx, w, r)
	if err != nil {
		return NewAppError(err, "LogoutRequest: ログアウトに失敗しました", http.StatusInternalServerError)
	}
	return nil
}

func GetRequestsRequest(ctx *context.AppContext, w http.ResponseWriter, r *http.Request) *AppError {
	requests, err := model.GetRequests(ctx)
	if err != nil {
		return NewAppError(err, "GetRequestsRequest: シフトリクエストの取得に失敗しました", http.StatusInternalServerError)
	}
	json.NewEncoder(w).Encode(requests)
	return nil
}

func GetEntriesRequest(ctx *context.AppContext, w http.ResponseWriter, r *http.Request) *AppError {
	requestId := r.PathValue("id")
	requestIdInt, err := strconv.Atoi(requestId)
	if err != nil {
		return NewAppError(err, "GetEntriesRequest: requestidが整数ではありません", http.StatusBadRequest)
	}

	entries, err := model.GetEntries(ctx, requestIdInt)
	if err != nil {
		return NewAppError(err, "GetEntriesRequest: エントリーの取得に失敗しました", http.StatusInternalServerError)
	}
	json.NewEncoder(w).Encode(entries)
	return nil
}

func PostRequestsRequest(ctx *context.AppContext, w http.ResponseWriter, r *http.Request) *AppError {
	// リクエストボディの形式を定義する
	type Body struct {
		StartDate string `json:"start_date"`
		EndDate   string `json:"end_date"`
		Deadline  string `json:"deadline"`
	}

	// リクエストボディのバリデーション
	var body Body
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		return NewAppError(err, "PostRequestsRequest: リクエストボディのデコードに失敗しました", http.StatusBadRequest)
	}

	// 新しいシフトリクエストを作成する
	response, err := model.CreateRequest(ctx, model.NewRequest{
		CreatorID: 2, // TODO: ログインしているユーザーのIDを取得する
		StartDate: body.StartDate,
		EndDate:   body.EndDate,
		Deadline:  body.Deadline,
	})
	if err != nil {
		return NewAppError(err, "PostRequestsRequest: シフトリクエストの作成に失敗しました", http.StatusInternalServerError)
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(response)
	return nil
}

func PostEntriesRequest(ctx *context.AppContext, w http.ResponseWriter, r *http.Request) *AppError {
	// シフトリクエストのIDを取得する
	// 整数ではない場合はエラーを返す
	requestId := r.PathValue("id")
	requestIdInt, err := strconv.Atoi(requestId)
	if err != nil {
		return NewAppError(err, "PostEntriesRequest: requestidが整数ではありません", http.StatusBadRequest)
	}

	// リクエストボディの形式を定義する
	type Body []struct {
		Date string `json:"date"`
		Hour int    `json:"hour"`
	}

	// リクエストボディのバリデーション
	var body Body
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		return NewAppError(err, "PostEntriesRequest: リクエストボディのデコードに失敗しました", http.StatusBadRequest)
	}

	// モデルに渡す形に変換する
	entries := model.NewEntries{
		ID:      requestIdInt,
		Entries: []model.NewEntry{},
	}
	for _, entry := range body {
		entries.Entries = append(entries.Entries, model.NewEntry{
			UserID: 1, // TODO: ログインしているユーザーのIDを取得する
			Date:   entry.Date,
			Hour:   entry.Hour,
		})
	}

	// エントリーを作成する
	response, err := model.CreateEntries(ctx, entries)
	if err != nil {
		return NewAppError(err, "PostEntriesRequest: エントリーの作成に失敗しました", http.StatusInternalServerError)
	}

	// レスポンスを返す
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(response)
	return nil
}
