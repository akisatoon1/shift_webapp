package usecase

import (
	"backend/context"
	"backend/db"
	"backend/domain"
)

func findEntriesBySubmissionID(ctx *context.AppContext, submissionID int) ([]domain.Entry, error) {
	// DBからエントリー一覧を取得
	entryRecs, err := ctx.GetDB().GetEntriesBySubmissionID(submissionID)
	if err != nil {
		return nil, err
	}

	// エントリー一覧を構築
	var entries []domain.Entry
	for _, entryRec := range entryRecs {
		// 時間型に変換
		date, err := domain.NewDateOnly(entryRec.Date)
		if err != nil {
			return nil, err
		}

		entries = append(entries, domain.Entry{
			ID:           entryRec.ID,
			SubmissionID: entryRec.SubmissionID,
			Date:         date,
			Hour:         entryRec.Hour,
		})
	}

	return entries, nil
}

type NewEntry struct {
	Date domain.DateOnly
	Hour int
}

// エントリーを作成する
func createEntries(ctx *context.AppContext, submissionID int, newEntries []NewEntry) ([]int, error) {
	var entryRecs []db.Entry

	for _, newEntry := range newEntries {
		entryRecs = append(entryRecs, db.Entry{
			SubmissionID: submissionID,
			Date:         newEntry.Date.Format(),
			Hour:         newEntry.Hour,
		})
	}

	entryIDs, err := ctx.GetDB().CreateEntries(entryRecs)
	if err != nil {
		return nil, err
	}

	return entryIDs, nil
}
