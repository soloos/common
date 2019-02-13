package srpc

import (
	"soloos/common/snet/types"
)

func (p *Client) Write(reqID uint64, serviceID string, req *types.Request) error {
	var (
		err error
	)

	req.ID = reqID

	// post data
	err = p.Conn.WriteRequestHeader(req.ID,
		serviceID,
		uint32(len(req.Param)+(req.OffheapBody.BodySize())),
		uint32(len(req.Param)))
	if err != nil {
		goto POST_DATA_DONE
	}

	err = p.Conn.WriteAll(req.Param)
	if err != nil {
		goto POST_DATA_DONE
	}

	err = req.OffheapBody.Copy(&p.Conn)
	if err != nil {
		goto POST_DATA_DONE
	}

POST_DATA_DONE:
	if err != nil {
		err = p.Conn.Close(err)
	}
	return err
}
