package iron

import (
	"bytes"
	"encoding/gob"
	"soloos/common/util"

	"golang.org/x/xerrors"
)

func SpecMarshal(v interface{}) ([]byte, error) {
	// return json.Marshal(v)
	var buf bytes.Buffer
	var err error
	enc := gob.NewEncoder(&buf)
	err = enc.Encode(v)
	if err != nil {
		return []byte{}, err
	}
	return buf.Bytes(), nil
}

func SpecUnmarshal(data []byte, v interface{}) error {
	// return json.Unmarshal(data, v)
	var buf = bytes.NewBuffer(data)
	dec := gob.NewDecoder(buf)
	var err = dec.Decode(v)
	return err
}

func MustSpecMarshalRequest(req interface{}) []byte {
	if req == nil {
		return []byte{}
	}
	var bs, err = SpecMarshal(req)
	util.AssertErrIsNil(err)
	return bs
}

func SpecUnmarshalRequest(bs []byte, ret interface{}) error {
	return SpecUnmarshal(bs, ret)
}

func SpecUnmarshalResponse(bs []byte, paramResp interface{}) error {
	var (
		resp  IRespData
		match bool
	)
	if paramResp == nil {
		resp = &Response{}
	} else {
		resp, match = paramResp.(IRespData)
		if !match {
			return ErrRespIsNotRespData
		}
	}

	var err = SpecUnmarshal(bs, resp)
	if err != nil {
		return err
	}
	if resp.GetError() != "" {
		return xerrors.New(resp.GetError())
	}

	return nil
}

func MustSpecMarshalResponse(resp IRespData) []byte {
	var bs, err = SpecMarshal(resp)
	if err != nil {
		bs, _ = SpecMarshal(MakeResp(nil, err))
		return bs
	}
	return bs
}

func MustSpecMarshalResponseErr(paramResp interface{}, respError error) []byte {
	var resp = MakeResp(paramResp, respError)
	var iresp = &resp
	var bs, err = SpecMarshal(iresp)
	if err != nil {
		bs, _ = SpecMarshal(MakeResp(nil, err))
		return bs
	}
	return bs
}
