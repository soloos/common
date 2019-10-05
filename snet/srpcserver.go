package snet

import (
	"net"
	"reflect"
	"soloos/common/iron"
	"soloos/common/log"
	"soloos/common/snettypes"
	"time"
)

type SrpcServer struct {
	MaxMessageLength uint32
	ln               net.Listener
	network          string
	address          string
	iron.Proxy
}

func (p *SrpcServer) Init(network, address string) error {
	var err error

	p.MaxMessageLength = 1024 * 1024 * 512
	p.network = network
	p.address = address

	err = p.Proxy.Init()
	if err != nil {
		return err
	}

	p.RegisterService("/Close", func(reqCtx *snettypes.SNetReqContext) error {
		go func() {
			time.Sleep(time.Second * 3)
			reqCtx.ConnClose(snettypes.ErrClosedByUser)
		}()
		return nil
	})

	return nil
}

func (p *SrpcServer) Serve() error {
	var err error
	p.ln, err = makeListener(p.network, p.address)
	if err != nil {
		return err
	}

	p.serveListener(p.ln)

	return nil
}

func (p *SrpcServer) serveListener(ln net.Listener) error {
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

func (p *SrpcServer) serveConn(netConn net.Conn) {
	var (
		conn             snettypes.Connection
		closeConnErrChan = make(chan error)
		err              error
	)

	conn.SetNetConn(netConn)

	for {
		var reqCtx snettypes.SNetReqContext
		reqCtx.Init(&conn)

		// read reqHeader
		err = reqCtx.ReadHeader(p.MaxMessageLength)
		if err != nil {
			goto CONN_END
		}

		go p.serveService(closeConnErrChan, reqCtx)

		err = <-closeConnErrChan
		if err != nil {
			goto CONN_END
		}

		if reqCtx.BodySize > 0 {
			conn.WaitReadDone()
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

func (p *SrpcServer) serveService(closeConnErrChan chan<- error,
	reqCtx snettypes.SNetReqContext) {
	var path = reqCtx.Header.Url
	var err error
	if !p.IsServiceExists(path) {
		err = reqCtx.SkipReadRemaining()
		if err != nil {
			closeConnErrChan <- err
			return
		}

		err = reqCtx.SimpleResponse(reqCtx.ReqID, iron.MustSpecMarshalResponseErr(nil, iron.ErrCmdNotFound))
		if err != nil {
			closeConnErrChan <- err
			return
		}

		closeConnErrChan <- nil
		return
	}
	closeConnErrChan <- nil

	var resp snettypes.IRespData

	var reqArgElems []interface{}
	var reqArgSize uint32 = reqCtx.ParamSize
	if reqArgSize > 0 {
		var service = p.ServiceTable[path]
		var reqArgBytes = make([]byte, reqArgSize)
		var reqArgValues []reflect.Value
		var reqArgInterfaces []interface{}

		err = reqCtx.ReadAll(reqArgBytes)
		if err != nil {
			goto PARSE_ARGS_DONE
		}

		for i, _ := range service.Params {
			var serviceParam = service.Params[i]
			var reqArgValue = reflect.New(serviceParam)
			reqArgValues = append(reqArgValues, reqArgValue)
			reqArgInterfaces = append(reqArgInterfaces, reqArgValue.Interface())
		}

		if len(reqArgInterfaces) == 1 {
			err = iron.SpecUnmarshalRequest(reqArgBytes, reqArgInterfaces[0])
		} else {
			err = iron.SpecUnmarshalRequest(reqArgBytes, reqArgInterfaces)
		}
		if err != nil {
			goto PARSE_ARGS_DONE
		}

		for i, _ := range reqArgValues {
			reqArgElems = append(reqArgElems, reqArgValues[i].Elem())
		}

	PARSE_ARGS_DONE:
	}

	resp = p.Proxy.Dispatch(path, &reqCtx, reqArgElems...)
	reqCtx.SkipReadRemaining()

	if !reqCtx.IsResponseInService {
		err = reqCtx.SimpleResponse(reqCtx.ReqID, iron.MustSpecMarshalResponse(resp))
		if err != nil {
			log.Debug("SrpcServer serveService SimpleRespons error, err:", err)
		}
	}

	return
}

func (p *SrpcServer) Close() error {
	return p.ln.Close()
}
