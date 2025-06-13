package usecase

import "errors"

var (
	ErrForbidden = errors.New("forbidden access")
)

type InputError struct {
	err     error
	message string
}

func NewInputError(err error, message string) InputError {
	return InputError{
		err:     err,
		message: message,
	}
}

func (e InputError) Error() string {
	return e.err.Error()
}

func (e InputError) Message() string {
	return e.message
}
