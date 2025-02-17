package user

import "errors"

var (
	ErrAccountNotExist    = errors.New("user not exist")
	ErrInvalidCredentials = errors.New("wrong password")
	ErrAccountStorer      = errors.New("error account storer")
	ErrAccountExist       = errors.New("account exist")

	ErrTokenNotExist = errors.New("token not exist")
	ErrTokenExpired  = errors.New("token is expired")
)
