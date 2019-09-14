package iron

import (
	"encoding/json"
)

type ResponseJSON struct {
	RData interface{} `json:"RData"`
	Err   string      `json:"Err"`
}

type Response struct {
	RData interface{}
	Err   error
}

func (p Response) Resolve() (interface{}, error) {
	return p.RData, p.Err
}

func (p Response) Must() interface{} {
	if p.Err != nil {
		panic(p.Err)
	}
	return p.RData
}

func MarshalResponse(resp Response) string {
	var ret ResponseJSON
	ret.RData = resp.RData
	if resp.Err != nil {
		ret.Err = resp.Err.Error()
	}

	b, err := json.Marshal(ret)

	if err != nil {
		return err.Error()
	}

	return string(b)
}
