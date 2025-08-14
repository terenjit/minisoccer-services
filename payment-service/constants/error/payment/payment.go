package error

import "errors"

var (
	ErrPaymentNotFound = errors.New("Payment not found")
	ErrExpireAtInvalid = errors.New("expired time must be greater than current time")
	ErrPaymentExists   = errors.New("Payment already exist")
)

var PaymentErrors = []error{
	ErrPaymentNotFound,
	ErrPaymentExists,
}
