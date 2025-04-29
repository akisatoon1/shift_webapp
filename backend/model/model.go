package model

/*
	apiの仕様に沿ったレスポンスを適切な型(構造体やスライス)で返す
*/

import (
	"backend/context"
	"time"
)

// APIレスポンスのcreatorやuserフィールドで利用される
type User struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

// シフトリクエスト情報を表す構造体
// APIレスポンスの1件分のシフトリクエストデータ
type Request struct {
	ID        int    `json:"id"`
	Creator   User   `json:"creator"`
	StartDate string `json:"start_date"`
	EndDate   string `json:"end_date"`
	Deadline  string `json:"deadline"`
	CreatedAt string `json:"created_at"`
}

// シフトリクエスト一覧APIのレスポンス型
type GetRequestsResponse []Request

// シフトリクエスト一覧を取得し、APIレスポンス用の構造体スライスに変換して返す
func GetRequests(ctx *context.AppContext) (GetRequestsResponse, error) {
	var response GetRequestsResponse

	// DBからシフトリクエスト一覧を取得
	requests, err := ctx.DB.GetRequests()
	if err != nil {
		return nil, err
	}

	// 各シフトリクエストごとに、作成者情報を取得し、レスポンス用構造体に詰める
	for _, request := range requests {
		user, err := ctx.DB.GetUserByID(request.CreatorID)
		if err != nil {
			return nil, err
		}
		response = append(response, Request{
			ID:        request.ID,
			Creator:   User{ID: user.ID, Name: user.Name},
			StartDate: request.StartDate.Format(time.DateOnly),
			EndDate:   request.EndDate.Format(time.DateOnly),
			Deadline:  request.Deadline.Format(time.DateOnly),
			CreatedAt: request.CreatedAt.Format(time.DateTime),
		})
	}

	return response, nil
}

// APIレスポンスのentriesフィールドで利用される
type Entry struct {
	ID   int    `json:"id"`
	User User   `json:"user"`
	Date string `json:"date"`
	Hour int    `json:"hour"`
}

// エントリー一覧APIのレスポンス型
// 1つのシフトリクエストIDに紐づくエントリーのリスト
type GetEntriesResponse struct {
	ID      int     `json:"id"`
	Entries []Entry `json:"entries"`
}

// 指定したシフトリクエストIDに紐づくエントリー一覧を取得し、APIレスポンス用の構造体に変換して返す
func GetEntries(ctx *context.AppContext, requestID int) (GetEntriesResponse, error) {
	// APIレスポンスのための、エントリーに紐づいたシフトリクエストID
	response := GetEntriesResponse{
		ID: requestID,
	}

	// DBからエントリー一覧を取得
	entries, err := ctx.DB.GetEntriesByRequestID(requestID)
	if err != nil {
		return GetEntriesResponse{}, err
	}

	// 各エントリーごとに、ユーザー情報を取得し、レスポンス用構造体に詰める
	for _, entry := range entries {
		user, err := ctx.DB.GetUserByID(entry.UserID)
		if err != nil {
			return GetEntriesResponse{}, err
		}
		response.Entries = append(response.Entries, Entry{
			ID:   entry.ID,
			User: User{ID: user.ID, Name: user.Name},
			Date: entry.Date.Format(time.DateOnly),
			Hour: entry.Hour,
		})
	}

	return response, nil
}

type PostRequestsBody struct {
	CreatorID int    `json:"creator_id"`
	StartDate string `json:"start_date"`
	EndDate   string `json:"end_date"`
	Deadline  string `json:"deadline"`
}

type PostRequestsResponse struct {
	ID int `json:"id"`
}

// シフトリクエストを作成する
func CreateRequest(ctx *context.AppContext, request PostRequestsBody) (PostRequestsResponse, error) {
	// 入力値のバリデーション
	// モデルに渡す
	// レスポンス
	return PostRequestsResponse{ID: 3}, nil
}
