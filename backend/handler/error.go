package handler

import (
	"runtime"
)

type AppError struct {
	err     error
	message string
	code    int
	file    string // 追加: エラー発生ファイル
	line    int    // 追加: エラー発生行
}

func NewAppError(err error, message string, code int) *AppError {
	file, line := "", 0
	if _, f, l, ok := runtime.Caller(1); ok {
		file = f
		line = l
	}
	return &AppError{
		err:     err,
		message: message,
		code:    code,
		file:    file,
		line:    line,
	}
}
