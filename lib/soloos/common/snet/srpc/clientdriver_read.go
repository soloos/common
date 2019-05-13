package srpc

import (
	"soloos/common/snet/types"
)

func (p *ClientDriver) ReadResponse(uPeer types.PeerUintptr,
	req *types.Request,
	resp *types.Response,
	respBody []byte) error {
	var (
		client *Client
		err    error
	)

	client, err = p.getClient(uPeer)
	if err != nil {
		return err
	}

	err = client.ReadResponse(resp, respBody)
	if err != nil {
		return err
	}

	return err
}
