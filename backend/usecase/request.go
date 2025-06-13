package usecase

import (
	"backend/auth"
	"backend/context"
	"backend/domain"
	"errors"
)

type IRequestUsecase interface {
	FindByID(ctx *context.AppContext, requestID int) (domain.Request, error)
	FindAll(ctx *context.AppContext) ([]domain.Request, error)
	Create(ctx *context.AppContext, newRequest NewRequest) (int, error)
}

type requestUsecase struct{}

func NewRequestUsecase() IRequestUsecase {
	return &requestUsecase{}
}

func (*requestUsecase) FindByID(ctx *context.AppContext, requestID int) (domain.Request, error) {
	// シフトリクエストを取得
	requestRec, err := ctx.GetDB().GetRequestByID(requestID)
	if err != nil {
		return domain.Request{}, err
	}

	// 作成者ユーザーを取得
	userRec, err := ctx.GetDB().GetUserByID(requestRec.CreatorID)
	if err != nil {
		return domain.Request{}, err
	}

	// 日付を適切な型に変換
	start_date, err := domain.NewDateOnly(requestRec.StartDate)
	end_date, err := domain.NewDateOnly(requestRec.EndDate)
	deadline, err := domain.NewDateTime(requestRec.Deadline)
	reqCreated_at, err := domain.NewDateTime(requestRec.CreatedAt)
	userCreated_at, err := domain.NewDateTime(userRec.CreatedAt)
	if err != nil {
		return domain.Request{}, err
	}

	return domain.Request{
		ID: requestRec.ID,
		Creator: domain.User{
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

func (*requestUsecase) FindAll(ctx *context.AppContext) ([]domain.Request, error) {
	// すべてのシフトリクエストを取得
	requestRecs, err := ctx.GetDB().GetRequests()
	if err != nil {
		return nil, err
	}

	var requests []domain.Request
	for _, rec := range requestRecs {
		var r IRequestUsecase = &requestUsecase{}
		foundRequest, err := r.FindByID(ctx, rec.ID)
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
	StartDate domain.DateOnly
	EndDate   domain.DateOnly
	Deadline  domain.DateTime
}

func (*requestUsecase) Create(ctx *context.AppContext, newRequest NewRequest) (int, error) {
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
