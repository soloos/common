package srpc

import (
	"soloos/snet/types"
)

func (p *Client) Write(requestID uint64, serviceID string, request *types.Request) error {
	var (
		err error
	)

	request.ID = requestID

	// post data
	err = p.Conn.WriteRequestHeader(request.ID,
		serviceID,
		uint32(len(request.Body)+(request.OffheapBody.ContentLen())))
	if err != nil {
		goto POST_DATA_DONE
	}

	if len(request.Body) > 0 {
		err = p.Conn.WriteAll(request.Body)
		if err != nil {
			goto POST_DATA_DONE
		}
	}

	err = request.OffheapBody.Copy(&p.Conn)
	if err != nil {
		goto POST_DATA_DONE
	}

POST_DATA_DONE:
	return err
}
