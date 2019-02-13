package types

import (
	"soloos/common/util/offheap"
)

type RpcClientDriver interface {
	Init(offheapDriver *offheap.OffheapDriver) error
	RegisterClient(uPeer PeerUintptr, client interface{}) error
	CloseClient(uPeer PeerUintptr) error
	Call(uPeer PeerUintptr,
		serviceID string,
		req *Request,
		resp *Response) error
	AsyncCall(uPeer PeerUintptr,
		serviceID string,
		req *Request,
		resp *Response) error
	WaitResponse(uPeer PeerUintptr,
		req *Request,
		resp *Response) error
	ReadResponse(uPeer PeerUintptr,
		req *Request,
		resp *Response,
		respBody []byte) error
}
