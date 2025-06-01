package model

import (
	"backend/context"
)

type Submission struct {
	ID          int
	RequestID   int
	SubmitterID int
	Submitter   User
	CreatedAt   DateTime
	UpdatedAt   DateTime
}

func GetSubmissionsByRequestID(ctx *context.AppContext, requestID int) ([]Submission, error) {
	// シフトリクエストIDが存在するかチェック
	_, err := GetRequestByID(ctx, requestID)
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

		user, err := GetUserByID(ctx, submissionRec.SubmitterID)
		if err != nil {
			return nil, err
		}
		submissions = append(submissions, Submission{
			ID:          submissionRec.ID,
			RequestID:   submissionRec.RequestID,
			SubmitterID: submissionRec.SubmitterID,
			Submitter:   user,
			CreatedAt:   createdAt,
			UpdatedAt:   updatedAt,
		})
	}

	return submissions, nil
}
