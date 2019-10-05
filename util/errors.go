package util

import "soloos/common/xerrors"

func NewError(errMsg string) error {
	if errMsg == "" {
		return nil
	}
	return xerrors.New(errMsg)
}
