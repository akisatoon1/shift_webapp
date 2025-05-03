package model

/*
	apiの仕様に沿ったレスポンスを適切な型(構造体やスライス)で返す
*/

import (
	"backend/context"
	"errors"
	"time"
)

// TODO
// validation
// time.Timeの形式がそれぞれ違うのが面倒

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
	requests, err := ctx.GetDB().GetRequests()
	if err != nil {
		return nil, err
	}

	// 各シフトリクエストごとに、作成者情報を取得し、レスポンス用構造体に詰める
	for _, request := range requests {
		user, err := ctx.GetDB().GetUserByID(request.CreatorID)
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
	entries, err := ctx.GetDB().GetEntriesByRequestID(requestID)
	if err != nil {
		return GetEntriesResponse{}, err
	}

	// 各エントリーごとに、ユーザー情報を取得し、レスポンス用構造体に詰める
	for _, entry := range entries {
		user, err := ctx.GetDB().GetUserByID(entry.UserID)
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

// 新しいシフトリクエストを作成するため
// 新しいシフトリクエストの内容を表す
type NewRequest struct {
	CreatorID int    `json:"creator_id"`
	StartDate string `json:"start_date"`
	EndDate   string `json:"end_date"`
	Deadline  string `json:"deadline"`
}

// 作成したシフトリクエストのidのみを返す
type PostRequestsResponse struct {
	ID int `json:"id"`
}

// 新しいシフトリクエストを作成する
func CreateRequest(ctx *context.AppContext, request NewRequest) (PostRequestsResponse, error) {
	// 日付の整合性チェック
	// 期限 <= 開始日 <= 終了日 でなければいけない
	deadline, _ := time.Parse(time.DateOnly, request.Deadline)
	startDate, _ := time.Parse(time.DateOnly, request.StartDate)
	endDate, _ := time.Parse(time.DateOnly, request.EndDate)
	if !((deadline.Before(startDate) || deadline.Equal(startDate)) && (startDate.Before(endDate) || startDate.Equal(endDate))) {
		return PostRequestsResponse{}, errors.New("input date must be deadline <= start_date <= end_date")
	}

	// DBに保存
	requestID, err := ctx.GetDB().CreateRequest(request.CreatorID, startDate, endDate, deadline)
	if err != nil {
		return PostRequestsResponse{}, err
	}

	// レスポンス
	return PostRequestsResponse{ID: requestID}, nil
}

// 新しく作成する1回分のエントリーの内容を表す
type NewEntry struct {
	UserID int    `json:"user_id"`
	Date   string `json:"date"`
	Hour   int    `json:"hour"`
}

// 全ての新しく作成するエントリーの内容を表す
// IDはシフトリクエストのID
type NewEntries struct {
	ID      int        `json:"id"`
	Entries []NewEntry `json:"entries"`
}

// 作成した1回分のエントリーのidのみを返す
type PostEntriesResponseEntry struct {
	ID int `json:"id"`
}

// 作成した全てのエントリーのidのみを返す
// IDはシフトリクエストのID
type PostEntriesResponse struct {
	ID      int                        `json:"id"`
	Entries []PostEntriesResponseEntry `json:"entries"`
}

// エントリーを作成する
func CreateEntries(ctx *context.AppContext, entries NewEntries) (PostEntriesResponse, error) {
	// TODO: start_date <= date <= end_date のバリデーション

	// エントリーを提出するシフトリクエストのID
	requestID := entries.ID

	response := PostEntriesResponse{
		ID:      requestID,
		Entries: []PostEntriesResponseEntry{},
	}

	// 全てのエントリーを作成
	for _, entry := range entries.Entries {
		date, _ := time.Parse(time.DateOnly, entry.Date)
		entryID, err := ctx.GetDB().CreateEntry(requestID, entry.UserID, date, entry.Hour)
		if err != nil {
			return PostEntriesResponse{}, err
		}
		response.Entries = append(response.Entries, PostEntriesResponseEntry{ID: entryID})
	}

	return response, nil
}
