package offheap

import "errors"

var (
	ErrAllocChunkOurOfLimit = errors.New("alloc chunk out of limit")
)
