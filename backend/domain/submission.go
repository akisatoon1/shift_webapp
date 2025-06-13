package domain

type Submission struct {
	ID          int
	RequestID   int
	SubmitterID int
	Submitter   User
	Entries     []Entry
	CreatedAt   DateTime
	UpdatedAt   DateTime
}

type Entry struct {
	ID           int
	SubmissionID int
	Date         DateOnly
	Hour         int
}
