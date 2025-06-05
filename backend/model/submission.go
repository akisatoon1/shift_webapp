package model

import (
	"backend/auth"
	"backend/context"
	"errors"
)

type Submission struct {
	ID          int
	RequestID   int
	SubmitterID int
	Submitter   User
	Entries     []entry
	CreatedAt   DateTime
	UpdatedAt   DateTime
}

func (*Submission) FindByRequestID(ctx *context.AppContext, requestID int) ([]Submission, error) {
	// シフトリクエストIDが存在するかチェック
	var request Request
	_, err := request.FindByID(ctx, requestID)
	if err != nil {
		return nil, err
	}

	// DBから提出一覧を取得
	submissionRecs, err := ctx.GetDB().GetSubmissionsByRequestID(requestID)
	if err != nil {
		return nil, err
	}

	// 提出一覧を構築
	var submissions []Submission
	for _, submissionRec := range submissionRecs {
		createdAt, err := NewDateTime(submissionRec.CreatedAt)
		if err != nil {
			return nil, err
		}
		updatedAt, err := NewDateTime(submissionRec.UpdatedAt)
		if err != nil {
			return nil, err
		}

		var user User
		user, err = user.FindByID(ctx, submissionRec.SubmitterID)
		if err != nil {
			return nil, err
		}

		// entryは内部構造体なのでそのまま関数を使用
		var e entry
		entries, err := e.findBySubmissionID(ctx, submissionRec.ID)
		if err != nil {
			return nil, err
		}

		submissions = append(submissions, Submission{
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

type NewSubmission struct {
	RequestID   int
	SubmitterID int
	NewEntries  []NewEntry
}

func (*Submission) Create(ctx *context.AppContext, newSubmission NewSubmission) (int, error) {
	// 提出者が従業員であるか確認する
	var user User
	foundUser, err := user.FindByID(ctx, newSubmission.SubmitterID)
	if err != nil {
		return 0, err
	}
	if foundUser.Role != auth.RoleEmployee {
		return 0, ErrForbidden
	}

	// シフトリクエストIDが存在するか確認する
	var request Request
	foundRequest, err := request.FindByID(ctx, newSubmission.RequestID)
	if err != nil {
		return 0, err
	}

	// 提出済みの場合はエラー
	alreadySubmitted, err := ctx.GetDB().AlreadySubmitted(newSubmission.RequestID, newSubmission.SubmitterID)
	if err != nil {
		return 0, err
	}
	if alreadySubmitted {
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
	var e entry
	_, err = e.create(ctx, submissionID, newSubmission.NewEntries)
	if err != nil {
		return 0, err
	}

	return submissionID, nil
}
