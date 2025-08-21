package error

import "errors"

var (
	ErrOrderNotFound = errors.New("order not found")
	ErrOrderExists   = errors.New("order already exist")
	ErrAlreadyBooked = errors.New("field already booked")
)

var OrderErrors = []error{
	ErrOrderNotFound,
	ErrOrderExists,
	ErrAlreadyBooked,
}
