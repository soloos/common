package xerrors

import "errors"

func New(msg string) error {
	return errors.New(msg)
}
