package sdbapitypes

import "golang.org/x/xerrors"

var (
	ErrObjectNotExists    = xerrors.New("object not exists")
	ErrCreateObjectFailed = xerrors.New("create object failed")
)
