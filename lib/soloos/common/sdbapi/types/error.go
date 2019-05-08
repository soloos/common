package types

import "golang.org/x/xerrors"

var (
	ErrObjectNotExists = xerrors.New("object not exists")
)
