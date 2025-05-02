package auth

import (
	"backend/context"
	"net/http"
)

// get user id from session.
// return false if user is not logged in or invalid cookie.
func GetUserID(ctx *context.AppContext, r *http.Request) (bool, int) {
	session, _ := ctx.Cookie.Get(r, "login_session")
	if session == nil || session.IsNew {
		return false, -1
	}

	userID, ok := session.Values["user_id"]
	if !ok {
		return false, -1
	}

	return true, userID.(int)
}

// 権限ビット定数
const (
	RoleEmployee = 1 << iota
	RoleManager
)

// check if user is employee
func IsEmployee(ctx *context.AppContext, userID int) (bool, error) {
	user, err := ctx.DB.GetUserByID(userID)
	if err != nil {
		return false, err
	}
	return (user.Role & RoleEmployee) != 0, nil
}

// check if user is manager
func IsManager(ctx *context.AppContext, userID int) (bool, error) {
	user, err := ctx.DB.GetUserByID(userID)
	if err != nil {
		return false, err
	}
	return (user.Role & RoleManager) != 0, nil
}
