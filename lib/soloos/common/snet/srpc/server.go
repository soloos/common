package srpc

import (
	"net"
	"soloos/common/log"
	"soloos/common/snet/types"
)

type Server struct {
	MaxMessageLength uint32
	ln               net.Listener
	network          string
	address          string
	services         map[types.ServiceID]types.Service
}

func (p *Server) Init(network, address string) error {
	p.MaxMessageLength = 1024 * 1024 * 512
	p.network = network
	p.address = address
	p.services = make(map[types.ServiceID]types.Service)
	p.RegisterService("/Close", func(serviceReq types.ServiceRequest) error {
		return serviceReq.Conn.Close(types.ErrClosedByUser)
	})
	return nil
}

func (p *Server) RegisterService(serviceIDStr string, service types.Service) {
	var serviceID types.ServiceID
	copy(serviceID[:], []byte(serviceIDStr))
	p.services[serviceID] = service
}

func (p *Server) Serve() error {
	var err error
	p.ln, err = makeListener(p.network, p.address)
	if err != nil {
		return err
	}

	p.serveListener(p.ln)

	return nil
}

func (p *Server) serveListener(ln net.Listener) error {
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

func (p *Server) serveConn(netConn net.Conn) {
	var (
		conn          types.Connection
		header        types.RequestHeader
		serviceID     types.ServiceID
		service       types.Service
		reqID         uint64
		localService  types.Service
		reqBodySize   uint32
		reqParamSize  uint32
		serviceExists bool
		err           error
	)

	conn.SetNetConn(netConn)

	for {
		// read header
		err = conn.ReadRequestHeader(p.MaxMessageLength, &header)
		if err != nil {
			goto CONN_END
		}

		header.ServiceID(&serviceID)
		service, serviceExists = p.services[serviceID]

		// call service
		go func() {
			reqID = header.ID()
			localService = service
			reqBodySize = conn.LastReadLimit
			reqParamSize = header.ParamSize()

			if serviceExists == false {
				conn.SkipReadRemaining()
				conn.SimpleResponse(reqID, nil)
				return
			} else {
				err = conn.AfterReadHeaderSuccess()
				if err != nil {
					return
				}
			}

			var serviceReq = types.ServiceRequest{
				ReqID:        header.ID(),
				ReqBodySize:  reqBodySize,
				ReqParamSize: reqParamSize,
				Conn:         &conn,
			}
			err = localService(serviceReq)
			if err != nil {
				return
			}
		}()

		if header.BodySize() > 0 {
			conn.ContinueReadSig.Lock()
			conn.ContinueReadSig.Unlock()
		}

		if err != nil {
			goto CONN_END
		}
	}

CONN_END:
	if err != nil {
		log.Debug("serveConn err ", netConn.RemoteAddr().Network(), err)
	}

	err = conn.Close(err)
	if err != nil {
		log.Debug("serveConn err ", netConn.RemoteAddr().Network(), err)
	}
}

func (p *Server) Close() error {
	return p.ln.Close()
}
