package handler

type AppError struct {
	err     error
	message string
	code    int
}

func NewAppError(err error, message string, code int) *AppError {
	return &AppError{
		err:     err,
		message: message,
		code:    code,
	}
}
