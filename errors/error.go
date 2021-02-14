package errors

import "errors"

var (
	ErrIncorrectEmailOrPassword = errors.New("incorrect email or password")
	ErrNotAuthenticated         = errors.New("not authenticated")
	ErrRecordNotFound           = errors.New("Запись не найдена")
)
