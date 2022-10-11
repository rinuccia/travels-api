package apperrors

import "errors"

var (
	ErrRecordNotFound = errors.New("record not found")
	ErrIncorrectQuery = errors.New("incorrect query")
)
