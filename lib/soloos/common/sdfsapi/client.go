package sdfsapi

import (
	"soloos/common/fsapi"
)

type Client interface {
	Close() error
	GetPosixFS() fsapi.PosixFS
}
