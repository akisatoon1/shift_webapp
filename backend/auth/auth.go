package auth

import (
	"backend/context"
	"backend/db"
	"errors"
	"net/http"
	"time"

	"golang.org/x/crypto/bcrypt"
)

func Login(ctx *context.AppContext, w http.ResponseWriter, r *http.Request, loginID string, password string) error {
	// login_idとpasswordを比較
	user, err := ctx.GetDB().GetUserByLoginID(loginID)
	if err != nil {
		if err == db.ErrUserNotFound {
			return errors.New("invalid login_id or password")
		}
		return err
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		return errors.New("invalid login_id or password")
	}

	// 成功時はセッションを作成し、Cookieに保存
	session, _ := ctx.GetSessionStore().Get(r, "login_session")
	if session == nil {
		return errors.New("session is nil")
	}

	session.Values["user_id"] = user.ID
	session.Options.MaxAge = int((time.Hour * 3).Seconds())
	return session.Save(r, w)
}

func Logout(ctx *context.AppContext, w http.ResponseWriter, r *http.Request) error {
	session, _ := ctx.GetSessionStore().Get(r, "login_session")

	// セッションが存在しない場合でもempty sessionが返される仕様なので、
	// session==nilの場合は想定されていない
	if session == nil {
		return errors.New("session is nil")
	}
	session.Options.MaxAge = -1
	return session.Save(r, w)
}

// get user id from session.
// return false if user is not logged in or invalid cookie.
func GetUserID(ctx *context.AppContext, r *http.Request) (int, bool) {
	session, _ := ctx.GetSessionStore().Get(r, "login_session")
	if session == nil || session.IsNew {
		return -1, false
	}

	userID, ok := session.Values["user_id"]
	if !ok {
		return -1, false
	}

	return userID.(int), true
}

// 権限ビット定数
const (
	RoleEmployee = 1 << iota
	RoleManager
)

// check if user is employee
func IsEmployee(ctx *context.AppContext, userID int) (bool, error) {
	user, err := ctx.GetDB().GetUserByID(userID)
	if err != nil {
		if err == db.ErrUserNotFound {
			return false, nil
		}
		return false, err
	}
	return (user.Role & RoleEmployee) != 0, nil
}

// check if user is manager
func IsManager(ctx *context.AppContext, userID int) (bool, error) {
	user, err := ctx.GetDB().GetUserByID(userID)
	if err != nil {
		if err == db.ErrUserNotFound {
			return false, nil
		}
		return false, err
	}
	return (user.Role & RoleManager) != 0, nil
}
