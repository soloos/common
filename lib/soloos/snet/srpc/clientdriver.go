package srpc

import (
	"io"
	"soloos/snet/types"
)

type ClientDriver struct {
	maxRequestID uint64
}

func (p *ClientDriver) Init() error {
	return nil
}

func (p *ClientDriver) Call(client *Client,
	serviceID string,
	request *types.ClientRequest,
	response *types.ClientResponse) error {
	var (
		err error
	)

	// post data
	var requestHeader types.RequestHeader
	requestHeader.SetID(p.maxRequestID)
	requestHeader.SetVersion(types.SNetVersion)
	requestHeader.SetContentLen(uint32(len(request.Body) +
		(request.OffheapBody.CopyEnd - request.OffheapBody.CopyOffset)))
	requestHeader.SetServiceID(serviceID)
	p.maxRequestID++

	err = client.Conn.WriteAll(requestHeader[:])
	if err != nil {
		return err
	}

	if request.Body != nil {
		err = client.Conn.WriteAll(request.Body)
		if err != nil {
			return err
		}
	}

	for {
		err = request.OffheapBody.Copy(&client.Conn)
		if err == io.EOF {
			break
		}
	}

	// fetch data
	client.Conn.ContinueReadSig.Add(1)
	var responseHeader types.ResponseHeader
	err = client.Conn.ReadResponseHeader(&responseHeader)
	if err != nil {
		return err
	}

	response.BodySize = responseHeader.ContentLen()
	client.Conn.LastRequestReadLimit = response.BodySize
	return nil
}
