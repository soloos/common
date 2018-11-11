package srpc

import (
	"soloos/snet/types"
)

func (p *ClientDriver) ReadResponse(uPeer types.PeerUintptr,
	request *types.Request,
	response *types.Response,
	resp *[]byte) error {
	var (
		client = p.clients[uPeer]
		err    error
	)

	err = client.ReadResponse(resp)
	if err != nil {
		return err
	}

	return err
}
