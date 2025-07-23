package error

import "errors"

var (
	ErrFieldNotFound = errors.New("field not found")
	ErrFieldExists   = errors.New("field already exist")
)

var FieldErrors = []error{
	ErrFieldNotFound,
	ErrFieldExists,
}
