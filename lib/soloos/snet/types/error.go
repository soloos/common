package types

import "errors"

var (
	ErrMessageTooLong  = errors.New("message is too long")
	ErrWrongVersion    = errors.New("wrong version")
	ErrServiceNotFound = errors.New("service not found")
)
