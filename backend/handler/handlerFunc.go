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

func GetRequestsRequest(ctx *context.AppContext, w http.ResponseWriter, r *http.Request) {
	requests, err := model.GetRequests(ctx)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(requests)
}

func GetEntriesRequest(ctx *context.AppContext, w http.ResponseWriter, r *http.Request) {
	requestId := r.PathValue("id")
	requestIdInt, err := strconv.Atoi(requestId)
	if err != nil {
		http.Error(w, "requestidが整数ではありません", http.StatusBadRequest)
		return
	}

	entries, err := model.GetEntries(ctx, requestIdInt)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(entries)
}

func PostRequestsRequest(ctx *context.AppContext, w http.ResponseWriter, r *http.Request) {
	// リクエストボディの形式を定義する
	type Body struct {
		StartDate string `json:"start_date"`
		EndDate   string `json:"end_date"`
		Deadline  string `json:"deadline"`
	}

	// リクエストボディのバリデーション
	var body Body
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// 新しいシフトリクエストを作成する
	response, err := model.CreateRequest(ctx, model.NewRequest{
		CreatorID: 2, // TODO: ログインしているユーザーのIDを取得する
		StartDate: body.StartDate,
		EndDate:   body.EndDate,
		Deadline:  body.Deadline,
	})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(response)
}
