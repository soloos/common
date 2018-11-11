package srpc

import (
	"net"
	"soloos/log"
	"soloos/snet/types"
)

type Server struct {
	MaxMessageLength uint32
	network          string
	address          string
	services         map[types.ServiceID]types.Service
}

func (p *Server) Init(network, address string) error {
	p.MaxMessageLength = 1024 * 1024 * 512
	p.network = network
	p.address = address
	p.services = make(map[types.ServiceID]types.Service)
	p.RegisterService("/Close", func(requestID uint64, requestContentLen uint32, conn *types.Connection) error {
		return conn.Close()
	})
	return nil
}

func (p *Server) RegisterService(serviceIDStr string, service types.Service) {
	var serviceID types.ServiceID
	copy(serviceID[:], []byte(serviceIDStr))
	p.services[serviceID] = service
}

func (p *Server) Serve() error {
	var (
		ln  net.Listener
		err error
	)
	ln, err = makeListener(p.network, p.address)
	if err != nil {
		return err
	}

	p.serveListener(ln)

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
		conn      types.Connection
		header    types.RequestHeader
		serviceID types.ServiceID
		service   types.Service
		exists    bool
		err       error
	)

	conn.SetNetConn(netConn)

	for {
		// read header
		err = conn.ReadRequestHeader(p.MaxMessageLength, &header)
		if err != nil {
			goto CONN_END
		}

		header.ServiceID(&serviceID)
		service, exists = p.services[serviceID]
		if !exists {
			conn.AfterReadHeaderError()
			err = types.ErrServiceNotFound
			goto CONN_END
		}

		// call service
		go func() {
			requestContentLen := conn.LastReadLimit
			localService := service
			err = conn.AfterReadHeaderSuccess()
			if err != nil {
				return
			}

			err = localService(header.ID(), requestContentLen, &conn)
			if err != nil {
				return
			}
		}()

		if header.ContentLen() > 0 {
			conn.ContinueReadSig.Wait()
		}

		if err != nil {
			goto CONN_END
		}
	}

CONN_END:
	if err != nil {
		log.Debug("serveConn err ", netConn.RemoteAddr().Network(), err)
	}

	err = conn.Close()
	if err != nil {
		log.Debug("serveConn err ", netConn.RemoteAddr().Network(), err)
	}
}
