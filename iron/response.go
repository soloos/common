package iron

import "encoding/gob"

type ResponseJSON struct {
	Data  interface{} `json:"Data"`
	Error string      `json:"Error"`
	Code  int         `json:"Code"`
}

func MakeResp(data interface{}, err error) Response {
	if err != nil {
		return Response{data, RespDataCommon{err.Error()}}
	}
	return Response{data, RespDataCommon{""}}
}

type IRespData interface {
	GetError() string
}

type RespData = interface{}

type Response struct {
	RespData
	RespDataCommon
}

type RespDataCommon struct {
	Error string
}

func (p RespDataCommon) GetError() string {
	return p.Error
}

func init() {
	gob.Register(Response{})
	gob.Register(RespDataCommon{})
}
