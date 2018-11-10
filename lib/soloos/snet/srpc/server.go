package srpc

import (
	"net"
	"soloos/log"
	"soloos/snet/types"
)

type Server struct {
	unsafeBuf        [512]byte
	MaxMessageLength uint32
	services         map[types.ServiceID]types.Service
}

func (p *Server) Init() error {
	p.MaxMessageLength = 1024 * 1024 * 512
	p.services = make(map[types.ServiceID]types.Service)
	return nil
}

func (p *Server) RegisterService(serviceIDStr string, service types.Service) {
	var serviceID types.ServiceID
	copy(serviceID[:], []byte(serviceIDStr))
	p.services[serviceID] = service
}

func (p *Server) Serve(network, address string) error {
	var (
		ln  net.Listener
		err error
	)
	ln, err = makeListener(network, address)
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
		conn   types.Connection
		header types.RequestHeader
		err    error
	)

	conn.Init(netConn)

	for {
		// read header
		err = conn.ReadRequestHeader(&header)
		if err != nil {
			log.Debug("serveConn err ", err)
			goto CONN_END
		}

		// prepare & check request
		conn.LastRequestReadLimit = header.ContentLen()
		if header.Version() != types.SNetVersion {
			err = types.ErrWrongVersion
			goto CONN_END
		}
		if conn.LastRequestReadLimit > p.MaxMessageLength {
			err = types.ErrMessageTooLong
			goto CONN_END
		}

		conn.ContinueReadSig.Add(1)

		// call service
		go func() {
			p.services[header.ServiceID()](header.ID(), &conn)
		}()

		if header.ContentLen() > 0 {
			conn.ContinueReadSig.Wait()
		}

		// read the rest
		if conn.LastRequestReadLimit > 0 {
			err = conn.SkipReadAllRest(&p.unsafeBuf)
		}
	}

CONN_END:
	err = conn.Close()
	if err != nil {
		log.Debug("serveConn err ", netConn.RemoteAddr().Network(), err)
	}
}
