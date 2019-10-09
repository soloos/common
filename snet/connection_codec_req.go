package snet

import (
	"encoding/gob"
	"io"
	"soloos/common/util"
)

func (p *Connection) SimpleUnmarshalRequest(reqCtx *SNetReqContext, reqArgValue interface{}) error {
	var err error
	var decBytes = snetCodecBytesPool.Get().(util.Buffer)
	defer snetCodecBytesPoolPut(decBytes)

	decBytes.Resize(int(reqCtx.ParamSize))
	err = reqCtx.ReadAll(decBytes.Bytes())
	if err != nil {
		return err
	}

	var dec = gob.NewDecoder(&decBytes)
	err = dec.Decode(reqArgValue)
	if err != nil && err != io.ErrUnexpectedEOF && err != io.EOF {
		return err
	}

	return nil
}
