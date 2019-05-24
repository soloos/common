package srpc

import (
	"soloos/common/snettypes"
)

func (p *ClientDriver) ReadResponse(uPeer snettypes.PeerUintptr,
	req *snettypes.Request,
	resp *snettypes.Response,
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
