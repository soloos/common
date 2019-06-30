package snet

import (
	"net"
	"soloos/common/log"
	"soloos/common/snettypes"
)

type SRPCServer struct {
	MaxMessageLength uint32
	ln               net.Listener
	network          string
	address          string
	services         map[snettypes.ServiceID]snettypes.Service
}

func (p *SRPCServer) Init(network, address string) error {
	p.MaxMessageLength = 1024 * 1024 * 512
	p.network = network
	p.address = address
	p.services = make(map[snettypes.ServiceID]snettypes.Service)
	p.RegisterService("/Close", func(serviceReq *snettypes.NetQuery) error {
		return serviceReq.ConnClose(snettypes.ErrClosedByUser)
	})
	return nil
}

func (p *SRPCServer) RegisterService(serviceIDStr string, service snettypes.Service) {
	var serviceID snettypes.ServiceID
	copy(serviceID[:], []byte(serviceIDStr))
	p.services[serviceID] = service
}

func (p *SRPCServer) Serve() error {
	var err error
	p.ln, err = makeListener(p.network, p.address)
	if err != nil {
		return err
	}

	p.serveListener(p.ln)

	return nil
}

func (p *SRPCServer) serveListener(ln net.Listener) error {
	var (
		netConn net.Conn
		err     error
	)

	for {
		netConn, err = ln.Accept()
		if err != nil {
			return err
		}

		go p.serveConn(netConn)
	}
}

func (p *SRPCServer) serveConn(netConn net.Conn) {
	var (
		conn          snettypes.Connection
		reqHeader     snettypes.RequestHeader
		serviceID     snettypes.ServiceID
		service       snettypes.Service
		serviceReq    snettypes.NetQuery
		serviceExists bool
		err           error
	)

	conn.SetNetConn(netConn)

	serviceReq = snettypes.NetQuery{}
	serviceReq.Init(&conn)

	for {
		// read reqHeader
		err = serviceReq.ReadRequestHeader(p.MaxMessageLength, &reqHeader)
		if err != nil {
			goto CONN_END
		}

		reqHeader.ServiceID(&serviceID)
		service, serviceExists = p.services[serviceID]

		if serviceExists == false {
			serviceReq.SkipReadRemaining()
			serviceReq.SimpleResponse(serviceReq.ReqID, nil)
			goto QUERY_DONE
		}

		// call service
		go func(localService snettypes.Service, localServiceReq snettypes.NetQuery) {
			localService(&localServiceReq)
			localServiceReq.EnsureServiceReadDone()
		}(service, serviceReq)

	QUERY_DONE:
		if serviceReq.BodySize > 0 {
			conn.WaitReadDone()
		}

		if err != nil {
			goto CONN_END
		}
	}

CONN_END:
	if err != nil {
		log.Info("serveConn err ", netConn.RemoteAddr().Network(), err)
	}

	err = conn.Close(err)
	if err != nil {
		log.Info("serveConn err ", netConn.RemoteAddr().Network(), err)
	}
}

func (p *SRPCServer) Close() error {
	return p.ln.Close()
}
