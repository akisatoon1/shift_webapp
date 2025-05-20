package model

/*
	apiの仕様に沿ったレスポンスを適切な型(構造体やスライス)で返す
*/

// TODO: modelという命名は間違い。リファクタリングする
// TODO: sessionインターフェイス

import (
	"backend/auth"
	"backend/context"
	"backend/db"
	"errors"
	"time"
)

/*
	注意!
	structのtagはhttp responseのjsonのkeyに利用される
	(modelがhandlerに依存している、悪い設計です)
*/

type SessionUser struct {
	ID        int      `json:"id"`
	Name      string   `json:"name"`
	Roles     []string `json:"roles"`
	CreatedAt string   `json:"created_at"`
}

type Session struct {
	User SessionUser `json:"user"`
}

// ユーザのセッション情報を返す
func GetSession(ctx *context.AppContext, userID int) (Session, error) {
	user, err := ctx.GetDB().GetUserByID(userID)
	if err != nil {
		return Session{}, err
	}

	// TODO: 抽象化できてない
	var roles []string
	if user.Role&auth.RoleEmployee != 0 {
		roles = append(roles, "employee")
	}
	if user.Role&auth.RoleManager != 0 {
		roles = append(roles, "manager")
	}

	// APIレスポンス用の構造体に詰める
	response := Session{
		User: SessionUser{
			ID:        user.ID,
			Name:      user.Name,
			Roles:     roles,
			CreatedAt: user.CreatedAt.Format(),
		},
	}

	return response, nil
}

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
			StartDate: request.StartDate.Format(),
			EndDate:   request.EndDate.Format(),
			Deadline:  request.Deadline.Format(),
			CreatedAt: request.CreatedAt.Format(),
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
	ID        int     `json:"id"`
	Creator   User    `json:"creator"`
	StartDate string  `json:"start_date"`
	EndDate   string  `json:"end_date"`
	Deadline  string  `json:"deadline"`
	CreatedAt string  `json:"created_at"`
	Entries   []Entry `json:"entries"`
}

// 指定したシフトリクエストIDに紐づくエントリー一覧と、シフトリクエスト詳細データを取得し、
// APIレスポンス用の構造体に変換して返す
func GetRequest(ctx *context.AppContext, requestID int) (GetEntriesResponse, error) {
	// シフトリクエストIDが存在するかチェック
	request, err := ctx.GetDB().GetRequestByID(requestID)
	if err != nil {
		return GetEntriesResponse{}, err
	}

	// 作成者の名前を取得
	user, err := ctx.GetDB().GetUserByID(request.CreatorID)
	if err != nil {
		return GetEntriesResponse{}, err
	}

	// APIレスポンスのための、エントリーに紐づいたシフトリクエストID
	response := GetEntriesResponse{
		ID:        requestID,
		Creator:   User{ID: request.CreatorID, Name: user.Name},
		StartDate: request.StartDate.Format(),
		EndDate:   request.EndDate.Format(),
		Deadline:  request.Deadline.Format(),
		CreatedAt: request.CreatedAt.Format(),
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
			Date: entry.Date.Format(),
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
	// ログインしているユーザーが従業員であるか確認する
	isUserManager, err := auth.IsManager(ctx, request.CreatorID)
	if err != nil {
		return PostRequestsResponse{}, err
	}
	if !isUserManager {
		return PostRequestsResponse{}, ErrForbidden
	}

	// 日付の整合性チェック
	// 期限 <= 開始日 <= 終了日 でなければいけない
	deadline, err := db.NewDateTime(request.Deadline)
	startDate, err := db.NewDateOnly(request.StartDate)
	endDate, err := db.NewDateOnly(request.EndDate)
	if err != nil {
		// 時間に関するデータが定められたフォーマットではないとき
		return PostRequestsResponse{}, err
	}
	if !((isBeforeOrEqual(deadline, startDate)) && (isBeforeOrEqual(startDate, endDate))) {
		return PostRequestsResponse{}, NewInputError(
			errors.New("must be deadline <= start_date <= end_date"),
			"期限 <= 開始日 <= 終了日 でなければいけない",
		)
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
	Date string `json:"date"`
	Hour int    `json:"hour"`
}

// 全ての新しく作成するエントリーの内容を表す
// IDはシフトリクエストのID
type NewEntries struct {
	ID      int        `json:"id"`
	UserID  int        `json:"user_id"`
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
	// ログインしているユーザーが従業員であるか確認する
	isUserEmployee, err := auth.IsEmployee(ctx, entries.UserID)
	if err != nil {
		return PostEntriesResponse{}, err
	}
	if !isUserEmployee {
		return PostEntriesResponse{}, ErrForbidden
	}

	// シフトリクエストIDが存在するかチェック
	if _, err := ctx.GetDB().GetRequestByID(entries.ID); err != nil {
		return PostEntriesResponse{}, err
	}

	// エントリーを提出するシフトリクエストのID
	requestID := entries.ID

	response := PostEntriesResponse{
		ID:      requestID,
		Entries: []PostEntriesResponseEntry{},
	}

	// start_date <= date <= end_date のバリデーションのため
	request, err := ctx.GetDB().GetRequestByID(entries.ID)
	if err != nil {
		return PostEntriesResponse{}, err
	}
	startDate := request.StartDate
	endDate := request.EndDate

	// 全てのエントリーについて、バリデーションを行う
	for _, entry := range entries.Entries {
		// start_date <= date <= end_date のバリデーション
		date, err := db.NewDateOnly(entry.Date)
		if err != nil {
			// 日付に関するデータが定められたフォーマットではないとき
			return PostEntriesResponse{}, err
		}
		if !((isBeforeOrEqual(startDate, date)) && (isBeforeOrEqual(date, endDate))) {
			return PostEntriesResponse{}, NewInputError(
				errors.New("must be start_date <= date <= end_date"),
				"開始日 <= 日付 <= 終了日 でなければいけない",
			)
		}

		// hourのvalidation
		// 0 <= hour <= 23 でなければいけない
		if !(0 <= entry.Hour && entry.Hour <= 23) {
			return PostEntriesResponse{}, NewInputError(
				errors.New("must be 0 <= hour <= 23"),
				"0 <= 時間 <= 23 でなければいけない",
			)
		}
	}

	// エントリーをdbに保存するために
	// db.Entry型に変換
	dbEntries := []db.Entry{}
	for _, entry := range entries.Entries {
		date, err := db.NewDateOnly(entry.Date)
		if err != nil {
			// 日付に関するデータが定められたフォーマットではないとき
			return PostEntriesResponse{}, err
		}
		dbEntries = append(dbEntries, db.Entry{
			RequestID: entries.ID,
			UserID:    entries.UserID,
			Date:      date,
			Hour:      entry.Hour,
		})
	}

	// DBに保存
	entryIDs, err := ctx.GetDB().CreateEntries(dbEntries)
	if err != nil {
		return PostEntriesResponse{}, err
	}

	// レスポンスを作成
	for _, entryID := range entryIDs {
		response.Entries = append(response.Entries, PostEntriesResponseEntry{ID: entryID})
	}

	return response, nil
}

// utility

func isBeforeOrEqual[T1 db.DateOnly | db.DateTime, T2 db.DateOnly | db.DateTime](a T1, b T2) bool {
	return time.Time(a).Before(time.Time(b)) || time.Time(a).Equal(time.Time(b))
}
