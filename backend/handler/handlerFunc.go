package handler

import (
	"backend/context"
	"backend/model"
	"encoding/json"
	"net/http"
	"strconv"
)

/*
	APIエンドポイントに対応するハンドラー関数
*/

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
