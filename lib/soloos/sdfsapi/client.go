package sdfsapi

import (
	"soloos/fsapi"
)

type Client interface {
	Close() error
	GetRawFileSystem() fsapi.RawFileSystem
}
