package snet

import (
	"encoding/gob"
	"soloos/common/util"
	"soloos/common/xerrors"
)

func (p *Connection) SimpleUnmarshalResponse(snetResp *SNetResp, resp IResponse) error {
	// func (p *Connection) SimpleUnmarshalResponse(snetResp *SNetResp, resp *Response) error {
	var err error
	var decBytes = snetCodecBytesPool.Get().(util.Buffer)
	defer snetCodecBytesPoolPut(decBytes)

	decBytes.Resize(int(snetResp.NetQuery.ParamSize))
	err = snetResp.ReadAll(decBytes.Bytes())
	if err != nil {
		return err
	}

	if resp == nil {
		resp = &Response{}
	}
	var dec = gob.NewDecoder(&decBytes)
	err = dec.Decode(resp)
	if err != nil {
		return err
	}

	if resp.GetErrorStr() != "" {
		return xerrors.New(resp.GetErrorStr())
	}

	return nil
}
