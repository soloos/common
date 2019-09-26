package solofsapi

import (
	"soloos/common/fsapi"
)

type Client interface {
	Close() error
	GetPosixFS() fsapi.PosixFS
	SetSOLOMQClient(itSOLOMQClient interface{}) error
}
