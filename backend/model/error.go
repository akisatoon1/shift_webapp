package model

import "errors"

var (
	ErrForbidden    = errors.New("権限がありません")
	ErrInvalidInput = errors.New("入力値が間違っています")
)
