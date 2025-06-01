package model

import (
	"backend/auth"
	"backend/context"
	"backend/db"
	"errors"
)

type Entry struct {
	ID        int
	RequestID int
	User      User
	Date      DateOnly
	Hour      int
}

func GetEntriesByRequestID(ctx *context.AppContext, requestID int) ([]Entry, error) {
	// シフトリクエストIDが存在するかチェック
	_, err := GetRequestByID(ctx, requestID)
	if err != nil {
		return nil, err
	}

	// DBからエントリー一覧を取得
	entryRecs, err := ctx.GetDB().GetEntriesByRequestID(requestID)
	if err != nil {
		return nil, err
	}

	// エントリー一覧を構築
	var entries []Entry
	for _, entryRec := range entryRecs {
		// TODO: 最適化(他の箇所もあるかも)
		user, err := GetUserByID(ctx, entryRec.UserID)
		if err != nil {
			return nil, err
		}

		// 時間型に変換
		date, err := NewDateOnly(entryRec.Date)
		if err != nil {
			return nil, err
		}

		entries = append(entries, Entry{
			ID:        entryRec.ID,
			RequestID: entryRec.RequestID,
			User: User{
				ID:        user.ID,
				LoginID:   user.LoginID,
				Password:  user.Password,
				Name:      user.Name,
				Role:      user.Role,
				CreatedAt: user.CreatedAt,
			},
			Date: date,
			Hour: entryRec.Hour,
		})
	}

	return entries, nil
}

type NewEntry struct {
	Date DateOnly
	Hour int
}

type NewEntries struct {
	RequestID int
	CreatorID int
	Entries   []NewEntry
}

// エントリーを作成する
func CreateEntries(ctx *context.AppContext, newEntries NewEntries) ([]int, error) {
	// 作成者が従業員であるか確認する
	isUserEmployee, err := auth.IsEmployee(ctx, newEntries.CreatorID)
	if err != nil {
		return nil, err
	}
	if !isUserEmployee {
		return nil, ErrForbidden
	}

	// シフトリクエストIDが存在するか確認する
	// start_date <= date <= end_date のバリデーションのため
	request, err := GetRequestByID(ctx, newEntries.RequestID)
	if err != nil {
		return nil, err
	}

	// dbに作成するためにrecord型を作成
	var entryRecs []db.Entry
	// 全てのエントリーについて、バリデーションを行う
	for _, entry := range newEntries.Entries {
		// start_date <= date <= end_date のバリデーション
		if !((isBeforeOrEqual(request.StartDate, entry.Date)) && (isBeforeOrEqual(entry.Date, request.EndDate))) {
			return nil, NewInputError(
				errors.New("must be start_date <= date <= end_date"),
				"開始日 <= 日付 <= 終了日 でなければいけない",
			)
		}

		// hourのvalidation
		// 0 <= hour <= 23 でなければいけない
		if !(0 <= entry.Hour && entry.Hour <= 23) {
			return nil, NewInputError(
				errors.New("must be 0 <= hour <= 23"),
				"0 <= 時間 <= 23 でなければいけない",
			)
		}

		entryRecs = append(entryRecs, db.Entry{
			RequestID: newEntries.RequestID,
			UserID:    newEntries.CreatorID,
			Date:      entry.Date.Format(),
			Hour:      entry.Hour,
		})
	}

	entryIDs, err := ctx.GetDB().CreateEntries(entryRecs)
	if err != nil {
		return nil, err
	}

	return entryIDs, nil
}
