package snet

import (
	"soloos/common/snettypes"
)

func (p *SRPCClientDriver) ReadResponse(peerID snettypes.PeerID,
	req *snettypes.Request, resp *snettypes.Response, respBody []byte) error {
	var (
		client *SRPCClient
		err    error
	)

	client, err = p.getClient(peerID)
	if err != nil {
		return err
	}

	err = client.ReadResponse(resp, respBody)
	if err != nil {
		return err
	}

	return err
}
