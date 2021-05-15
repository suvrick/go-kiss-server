package errors

import "errors"

var (
	ErrInvalidParam             = errors.New("invalid param")
	ErrIncorrectEmailOrPassword = errors.New("incorrect email or password")
	ErrNotAuthenticated         = errors.New("not authenticated")
	ErrRecordNotFound           = errors.New("not found")
	ErrLimitBot                 = errors.New("limit bot")
	ErrNeedTime                 = errors.New("need time")
	ErrProxyListEmpty           = errors.New("proxy list empty")
	ErrParseBotInvalid          = errors.New("parse bot error")
)
