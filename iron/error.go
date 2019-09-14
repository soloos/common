package iron

import "errors"

var (
	ErrCmdNotFound     = errors.New("command not found.")
	ErrCmdParamInvalid = errors.New("command params invalid.")
	ErrCmdParamEmpty   = errors.New("command params empty.")
)
