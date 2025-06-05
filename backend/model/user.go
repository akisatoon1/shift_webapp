package model

import "backend/context"

type User struct {
	ID        int
	LoginID   string
	Password  string
	Name      string
	Role      int
	CreatedAt DateTime
}

func (*User) FindByID(ctx *context.AppContext, userID int) (User, error) {
	// ユーザーIDが見つからない時はエラーを返す
	userRec, err := ctx.GetDB().GetUserByID(userID)
	if err != nil {
		return User{}, err
	}

	createdAt, err := NewDateTime(userRec.CreatedAt)
	if err != nil {
		return User{}, err
	}

	return User{
		ID:        userRec.ID,
		LoginID:   userRec.LoginID,
		Password:  userRec.Password,
		Name:      userRec.Name,
		Role:      userRec.Role,
		CreatedAt: createdAt,
	}, nil
}
