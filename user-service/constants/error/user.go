package error

import "errors"

var (
	ErrUserNotFound         = errors.New("user not found")
	ErrPasswordIncorrect    = errors.New("incorrect password")
	ErrUsernameExists       = errors.New("username already exist")
	ErrPasswordDoesNotMatch = errors.New("password does not match")
)

var UserErrors = []error{
	ErrUserNotFound,
	ErrPasswordDoesNotMatch,
	ErrPasswordIncorrect,
	ErrUsernameExists,
}
