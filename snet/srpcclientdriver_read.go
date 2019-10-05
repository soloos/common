package snet

import (
	"soloos/common/iron"
	"soloos/common/snettypes"
)

func (p *SrpcClientDriver) ReadResponse(peerID snettypes.PeerID,
	snetReq *snettypes.SNetReq, snetResp *snettypes.SNetResp,
	respBody []byte,
	ret interface{},
) error {
	var err error

	err = p.ReadRawResponse(peerID, snetReq, snetResp, respBody)
	if err != nil {
		return err
	}

	err = iron.SpecUnmarshalResponse(respBody, ret)
	if err != nil {
		return err
	}

	return err
}

func (p *SrpcClientDriver) ReadRawResponse(peerID snettypes.PeerID,
	snetReq *snettypes.SNetReq, snetResp *snettypes.SNetResp,
	respBody []byte,
) error {
	var (
		client *SrpcClient
		err    error
	)

	client, err = p.getClient(peerID)
	if err != nil {
		return err
	}

	err = client.ReadResponse(snetResp, respBody)
	if err != nil {
		return err
	}

	return err
}
