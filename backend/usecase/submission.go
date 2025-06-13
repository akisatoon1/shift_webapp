package usecase

import (
	"backend/auth"
	"backend/context"
	"backend/domain"
	"errors"
)

type ISubmissionUsecase interface {
	FindByRequestID(ctx *context.AppContext, requestID int) ([]domain.Submission, error)
	FindByRequestIDAndSubmitterID(ctx *context.AppContext, requestID, submitterID int) (*domain.Submission, error)
	Create(ctx *context.AppContext, newSubmission NewSubmission) (int, error)
}

type submissionUsecase struct{}

func NewSubmissionUsecase() ISubmissionUsecase {
	return &submissionUsecase{}
}

func (*submissionUsecase) FindByRequestID(ctx *context.AppContext, requestID int) ([]domain.Submission, error) {
	// シフトリクエストIDが存在するかチェック
	var r IRequestUsecase = &requestUsecase{}
	_, err := r.FindByID(ctx, requestID)
	if err != nil {
		return nil, err
	}

	// DBから提出一覧を取得
	submissionRecs, err := ctx.GetDB().GetSubmissionsByRequestID(requestID)
	if err != nil {
		return nil, err
	}

	// 提出一覧を構築
	var submissions []domain.Submission
	for _, submissionRec := range submissionRecs {
		createdAt, err := domain.NewDateTime(submissionRec.CreatedAt)
		if err != nil {
			return nil, err
		}
		updatedAt, err := domain.NewDateTime(submissionRec.UpdatedAt)
		if err != nil {
			return nil, err
		}

		var u IUserUsecase = &userUsecase{}
		user, err := u.FindByID(ctx, submissionRec.SubmitterID)
		if err != nil {
			return nil, err
		}

		// entryは内部構造体なのでそのまま関数を使用
		entries, err := findEntriesBySubmissionID(ctx, submissionRec.ID)
		if err != nil {
			return nil, err
		}

		submissions = append(submissions, domain.Submission{
			ID:          submissionRec.ID,
			RequestID:   submissionRec.RequestID,
			SubmitterID: submissionRec.SubmitterID,
			Submitter:   user,
			Entries:     entries,
			CreatedAt:   createdAt,
			UpdatedAt:   updatedAt,
		})
	}

	return submissions, nil
}

func (*submissionUsecase) FindByRequestIDAndSubmitterID(ctx *context.AppContext, requestID, submitterID int) (*domain.Submission, error) {
	// シフトリクエストIDが存在するかチェック
	var r IRequestUsecase = &requestUsecase{}
	_, err := r.FindByID(ctx, requestID)
	if err != nil {
		return nil, err
	}

	// 提出者が従業員であるか確認する
	var u IUserUsecase = &userUsecase{}
	user, err := u.FindByID(ctx, submitterID)
	if err != nil {
		return nil, err
	}
	if user.Role != auth.RoleEmployee {
		return nil, ErrForbidden
	}

	// DBから提出を取得
	submissionRec, err := ctx.GetDB().GetSubmissionByRequestIDAndSubmitterID(requestID, submitterID)
	if err != nil {
		return nil, err
	}

	if submissionRec == nil {
		return nil, nil
	}

	entries, err := findEntriesBySubmissionID(ctx, submissionRec.ID)
	if err != nil {
		return nil, err
	}

	createdAt, err := domain.NewDateTime(submissionRec.CreatedAt)
	if err != nil {
		return nil, err
	}
	updatedAt, err := domain.NewDateTime(submissionRec.UpdatedAt)
	if err != nil {
		return nil, err
	}

	return &domain.Submission{
		ID:          submissionRec.ID,
		RequestID:   submissionRec.RequestID,
		SubmitterID: submissionRec.SubmitterID,
		Submitter:   user,
		Entries:     entries,
		CreatedAt:   createdAt,
		UpdatedAt:   updatedAt,
	}, nil
}

type NewSubmission struct {
	RequestID   int
	SubmitterID int
	NewEntries  []NewEntry
}

func (*submissionUsecase) Create(ctx *context.AppContext, newSubmission NewSubmission) (int, error) {
	// 提出者が従業員であるか確認する
	var u IUserUsecase = &userUsecase{}
	foundUser, err := u.FindByID(ctx, newSubmission.SubmitterID)
	if err != nil {
		return 0, err
	}
	if foundUser.Role != auth.RoleEmployee {
		return 0, ErrForbidden
	}

	// シフトリクエストIDが存在するか確認する
	var r IRequestUsecase = &requestUsecase{}
	foundRequest, err := r.FindByID(ctx, newSubmission.RequestID)
	if err != nil {
		return 0, err
	}

	// 提出済みの場合はエラー
	subRec, err := ctx.GetDB().GetSubmissionByRequestIDAndSubmitterID(newSubmission.RequestID, newSubmission.SubmitterID)
	if err != nil {
		return 0, err
	}
	if subRec != nil {
		return 0, errors.New("already submitted")
	}

	// エントリーのvalidation
	for _, entry := range newSubmission.NewEntries {
		// 日付のvalidation
		if !isBeforeOrEqual(foundRequest.StartDate, entry.Date) || !isBeforeOrEqual(entry.Date, foundRequest.EndDate) {
			return 0, NewInputError(
				errors.New("date must be within request range"),
				"日付はリクエストの範囲内でなければいけない",
			)
		}

		// 0 <= hour <= 23 でなければいけない
		if !(0 <= entry.Hour && entry.Hour <= 23) {
			return 0, NewInputError(
				errors.New("must be 0 <= hour <= 23"),
				"0 <= 時間 <= 23 でなければいけない",
			)
		}
	}

	// DBに提出を作成
	submissionID, err := ctx.GetDB().CreateSubmission(newSubmission.SubmitterID, newSubmission.RequestID)
	if err != nil {
		return 0, err
	}

	// エントリーを作成
	_, err = createEntries(ctx, submissionID, newSubmission.NewEntries)
	if err != nil {
		return 0, err
	}

	return submissionID, nil
}
