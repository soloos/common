package snet

func (p *SrpcClient) Write(req *SNetReq) error {
	return req.Request(req)
}
