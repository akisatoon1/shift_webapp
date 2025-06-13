package usecase

import (
	"backend/context"
	"backend/domain"
)

type IUserUsecase interface {
	FindByID(ctx *context.AppContext, userID int) (domain.User, error)
}

type userUsecase struct{}

func NewUserUsecase() IUserUsecase {
	return &userUsecase{}
}

func (u *userUsecase) FindByID(ctx *context.AppContext, userID int) (domain.User, error) {
	// ユーザーIDが見つからない時はエラーを返す
	userRec, err := ctx.GetDB().GetUserByID(userID)
	if err != nil {
		return domain.User{}, err
	}

	createdAt, err := domain.NewDateTime(userRec.CreatedAt)
	if err != nil {
		return domain.User{}, err
	}

	return domain.User{
		ID:        userRec.ID,
		LoginID:   userRec.LoginID,
		Password:  userRec.Password,
		Name:      userRec.Name,
		Role:      userRec.Role,
		CreatedAt: createdAt,
	}, nil
}
