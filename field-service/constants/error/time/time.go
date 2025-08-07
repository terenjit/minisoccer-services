package error

import "errors"

var (
	ErrTimeNotFound = errors.New("Time not found")
	ErrTimeExists   = errors.New("Time already exist")
)

var TimeErrors = []error{
	ErrTimeNotFound,
	ErrTimeExists,
}
