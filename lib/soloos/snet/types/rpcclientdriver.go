package types

import (
	"soloos/util/offheap"
)

type RpcClientDriver interface {
	Init(offheapDriver *offheap.OffheapDriver) error
	RegisterClient(uPeer PeerUintptr) error
	CloseClient(uPeer PeerUintptr) error
	Call(uPeer PeerUintptr,
		serviceID string,
		request *Request,
		response *Response) error
	AsyncCall(uPeer PeerUintptr,
		serviceID string,
		request *Request,
		response *Response) error
	WaitResponse(uPeer PeerUintptr,
		request *Request,
		response *Response) error
}
