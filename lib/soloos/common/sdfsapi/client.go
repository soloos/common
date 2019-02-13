package sdfsapi

import (
	"soloos/common/fsapi"
)

type Client interface {
	Close() error
	GetRawFileSystem() fsapi.RawFileSystem
}
