package error

import "errors"

var (
	ErrInternalServerError = errors.New("internal server error")
	ErrSQLError            = errors.New("database server failed to execute query")
	ErrToManyRequests      = errors.New("too many requests")
	ErrUnauthorized        = errors.New("unauthorized")
	ErrInvalidToken        = errors.New("invalid token")
	ErrInvalidUploadFile   = errors.New("invalid upload file")
	ErrForbidden           = errors.New("forbidden")
	ErrSizetooBig          = errors.New("upload file size too big")
)

var GeneralErrors = []error{
	ErrInternalServerError,
	ErrSQLError,
	ErrToManyRequests,
	ErrUnauthorized, ErrInvalidToken, ErrForbidden,
}
