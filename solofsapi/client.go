package solofsapi

import (
	"soloos/common/fsapi"
)

type Client interface {
	Close() error
	GetPosixFS() fsapi.PosixFS
	SetSolomqClient(itSolomqClient interface{}) error
}
