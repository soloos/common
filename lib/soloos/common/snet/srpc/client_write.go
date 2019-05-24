package srpc

import (
	"soloos/common/snettypes"
)

func (p *Client) Write(req *snettypes.Request) error {
	var (
		err error
	)

	// post data
	err = req.NetQuery.WriteRequestHeader(req.ReqID,
		req.ServiceID,
		uint32(len(req.Param)+(req.OffheapBody.BodySize())),
		uint32(len(req.Param)))
	if err != nil {
		goto POST_DATA_DONE
	}

	err = req.NetQuery.WriteAll(req.Param)
	if err != nil {
		goto POST_DATA_DONE
	}

	err = req.OffheapBody.Copy(&req.NetQuery)
	if err != nil {
		goto POST_DATA_DONE
	}

POST_DATA_DONE:
	if err != nil {
		err = req.NetQuery.ConnClose(err)
	}
	return err
}
