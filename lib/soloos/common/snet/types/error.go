package types

import "golang.org/x/xerrors"

var (
	ErrMessageTooLong  = xerrors.New("message is too long")
	ErrWrongVersion    = xerrors.New("wrong version")
	ErrServiceNotFound = xerrors.New("service not found")
	ErrClosedByUser    = xerrors.New("closed by user")
)
