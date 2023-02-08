package validation

type HttpStatus int

type ValidationError struct {
	ValidationError string
	HttpStatus
}

func (validationError *ValidationError) Error() string {
	return validationError.ValidationError
}
