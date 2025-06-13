package domain

type User struct {
	ID        int
	LoginID   string
	Password  string
	Name      string
	Role      int
	CreatedAt DateTime
}
