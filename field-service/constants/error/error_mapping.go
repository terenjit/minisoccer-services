package error

import (
	errField "field-service/constants/error/field"
	errFieldSch "field-service/constants/error/field_schedule"
	errTime "field-service/constants/error/time"
)

func ErrMapping(err error) bool {
	var (
		GeneralErrors       = GeneralErrors
		FieldErrors         = errField.FieldErrors
		FieldScheduleErrors = errFieldSch.FieldScheduleErrors
		TimeErrors          = errTime.TimeErrors
	)
	allErrors := make([]error, 0)
	allErrors = append(allErrors, GeneralErrors...)
	allErrors = append(allErrors, FieldErrors...)
	allErrors = append(allErrors, FieldScheduleErrors...)
	allErrors = append(allErrors, TimeErrors...)

	for _, item := range allErrors {
		if err.Error() == item.Error() {
			return true
		}
	}
	return false
}
