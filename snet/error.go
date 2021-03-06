package snet

import "soloos/common/xerrors"

var (
	ErrObjectNotExists = xerrors.New("object is not exists")
	ErrMessageTooLong  = xerrors.New("message is too long")
	ErrWrongVersion    = xerrors.New("wrong version")
	ErrServiceNotFound = xerrors.New("service not found")
	ErrClosedByUser    = xerrors.New("closed by user")
	Err502             = xerrors.New("502")
)
