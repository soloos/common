package types

import (
	"fmt"
	"syscall"
)

func (code Status) String() string {
	if code <= 0 {
		return []string{
			"OK",
			"NOTIFY_POLL",
			"NOTIFY_INVAL_INODE",
			"NOTIFY_INVAL_ENTRY",
			"NOTIFY_STORE_CACHE",
			"NOTIFY_RETRIEVE_CACHE",
			"NOTIFY_DELETE",
		}[-code]
	}
	return fmt.Sprintf("%d=%v", int(code), syscall.Errno(code))
}

func (code Status) Ok() bool {
	return code == OK
}
