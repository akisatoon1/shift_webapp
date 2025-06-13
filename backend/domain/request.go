package domain

type Request struct {
	ID        int
	Creator   User
	StartDate DateOnly
	EndDate   DateOnly
	Deadline  DateTime
	CreatedAt DateTime
}
