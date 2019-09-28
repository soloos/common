package solofsapi

import (
	"soloos/common/fsapi"
)

type Client interface {
	Close() error
	GetPosixFs() fsapi.PosixFs
	SetSolomqClient(itSolomqClient interface{}) error
}
