package srpc

import (
	"soloos/snet/types"
)

func (p *ClientDriver) Call(uPeer types.PeerUintptr,
	serviceID string,
	request *types.Request,
	response *types.Response) error {
	var (
		err error
	)

	err = p.AsyncCall(uPeer, serviceID, request, response)
	if err != nil {
		return err
	}

	err = p.WaitResponse(uPeer, request, response)
	if err != nil {
		return err
	}

	return nil
}

func (p *ClientDriver) AsyncCall(uPeer types.PeerUintptr,
	serviceID string,
	request *types.Request,
	response *types.Response) error {
	var (
		client = p.clients[uPeer]
		err    error
	)

	requestID := client.AllocRequestID()
	err = client.PrepareWaitResponse(requestID, response)
	if err != nil {
		return err
	}

	err = client.Write(requestID, serviceID, request)
	if err != nil {
		return err
	}

	return nil
}

func (p *ClientDriver) WaitResponse(uPeer types.PeerUintptr,
	request *types.Request,
	response *types.Response) error {
	var (
		client = p.clients[uPeer]
		err    error
	)

	err = client.WaitResponse(request, response)
	if err != nil {
		return err
	}

	return nil
}
