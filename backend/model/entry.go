package model

import (
	"backend/context"
	"backend/db"
)

type entry struct {
	ID           int
	SubmissionID int
	Date         DateOnly
	Hour         int
}

func (*entry) findBySubmissionID(ctx *context.AppContext, submissionID int) ([]entry, error) {
	// DBからエントリー一覧を取得
	entryRecs, err := ctx.GetDB().GetEntriesBySubmissionID(submissionID)
	if err != nil {
		return nil, err
	}

	// エントリー一覧を構築
	var entries []entry
	for _, entryRec := range entryRecs {
		// 時間型に変換
		date, err := NewDateOnly(entryRec.Date)
		if err != nil {
			return nil, err
		}

		entries = append(entries, entry{
			ID:           entryRec.ID,
			SubmissionID: entryRec.SubmissionID,
			Date:         date,
			Hour:         entryRec.Hour,
		})
	}

	return entries, nil
}

type NewEntry struct {
	Date DateOnly
	Hour int
}

// エントリーを作成する
func (*entry) create(ctx *context.AppContext, submissionID int, newEntries []NewEntry) ([]int, error) {
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
