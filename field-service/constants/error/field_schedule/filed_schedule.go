package error

import "errors"

var (
	ErrFieldScheduleNotFound = errors.New("field schedule not found")
	ErrFieldScheduleExists   = errors.New("field schedule already exist")
)

var FieldScheduleErrors = []error{
	ErrFieldScheduleNotFound,
	ErrFieldScheduleExists,
}
