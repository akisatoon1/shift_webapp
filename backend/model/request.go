package model

import (
	"backend/auth"
	"backend/context"
	"errors"
)

// シフトリクエスト
type Request struct {
	ID        int
	Creator   User
	StartDate DateOnly
	EndDate   DateOnly
	Deadline  DateTime
	CreatedAt DateTime
}

func (*Request) FindByID(ctx *context.AppContext, requestID int) (Request, error) {
	// シフトリクエストを取得
	requestRec, err := ctx.GetDB().GetRequestByID(requestID)
	if err != nil {
		return Request{}, err
	}

	// 作成者ユーザーを取得
	userRec, err := ctx.GetDB().GetUserByID(requestRec.CreatorID)
	if err != nil {
		return Request{}, err
	}

	// 日付を適切な型に変換
	start_date, err := NewDateOnly(requestRec.StartDate)
	end_date, err := NewDateOnly(requestRec.EndDate)
	deadline, err := NewDateTime(requestRec.Deadline)
	reqCreated_at, err := NewDateTime(requestRec.CreatedAt)
	userCreated_at, err := NewDateTime(userRec.CreatedAt)
	if err != nil {
		return Request{}, err
	}

	return Request{
		ID: requestRec.ID,
		Creator: User{
			ID:        userRec.ID,
			LoginID:   userRec.LoginID,
			Password:  userRec.Password,
			Name:      userRec.Name,
			Role:      userRec.Role,
			CreatedAt: userCreated_at,
		},
		StartDate: start_date,
		EndDate:   end_date,
		Deadline:  deadline,
		CreatedAt: reqCreated_at,
	}, nil
}

func (*Request) FindAll(ctx *context.AppContext) ([]Request, error) {
	// すべてのシフトリクエストを取得
	requestRecs, err := ctx.GetDB().GetRequests()
	if err != nil {
		return nil, err
	}

	var requests []Request
	for _, rec := range requestRecs {
		var request Request
		foundRequest, err := request.FindByID(ctx, rec.ID)
		if err != nil {
			return nil, err
		}
		requests = append(requests, foundRequest)
	}

	return requests, nil
}

// リクエスト作成用のコマンド構造体
type NewRequest struct {
	CreatorID int
	StartDate DateOnly
	EndDate   DateOnly
	Deadline  DateTime
}

func (*Request) Create(ctx *context.AppContext, newRequest NewRequest) (int, error) {
	// 作成するユーザーがマネージャーであるか確認する
	isUserManager, err := auth.IsManager(ctx, newRequest.CreatorID)
	if err != nil {
		return -1, err
	}
	if !isUserManager {
		return -1, ErrForbidden
	}

	// 日付の整合性チェック
	// 期限 <= 開始日 <= 終了日 でなければいけない
	if !((isBeforeOrEqual(newRequest.Deadline, newRequest.StartDate)) && (isBeforeOrEqual(newRequest.StartDate, newRequest.EndDate))) {
		return -1, NewInputError(
			errors.New("must be deadline <= start_date <= end_date"),
			"期限 <= 開始日 <= 終了日 でなければいけない",
		)
	}

	// dbに作成
	requestID, err := ctx.GetDB().CreateRequest(newRequest.CreatorID, newRequest.StartDate.Format(), newRequest.EndDate.Format(), newRequest.Deadline.Format())
	if err != nil {
		return -1, err
	}

	return requestID, nil
}
