package sdbapitypes

import "soloos/common/xerrors"

var (
	ErrObjectNotExists    = xerrors.New("object not exists")
	ErrCreateObjectFailed = xerrors.New("create object failed")
)
