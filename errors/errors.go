package errors

import "errors"

var ErrUnauthorized = errors.New("Invalid credentials!")
var ErrDuplicated = errors.New("Duplicated!")
var ErrAlreadyExisted = errors.New("Already existed!")
var ErrNotExisted = errors.New("Not existed!")
var ErrValidationCheck = errors.New("Fail validate check")
var ErrInvalidOrExpiredToken = errors.New("Invalid or expired token")
var ErrInvalidInput = errors.New("Invalid input: Please provide accurate and complete information.")
var ErrNotSupportEventType = errors.New("Not support event type")
var ErrNotSupportEventSource = errors.New("Not supported event source")
