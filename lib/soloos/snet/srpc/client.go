package srpc

import (
	"soloos/log"
	"soloos/snet/types"
	"soloos/util/offheap"
	"strings"
	"sync"
	"sync/atomic"
)

type Client struct {
	Conn             types.Connection
	MaxMessageLength uint32

	clientDriver       *ClientDriver
	remoteAddr         string
	maxRequestID       uint64
	lastResponseHeader types.ResponseHeader
	reqSigMapMutex     sync.Mutex
	reqSigMap          map[uint64]offheap.MutexUintptr // map RequestID to netConnReadSigsIndex
}

func (p *Client) Init(clientDriver *ClientDriver, address string) error {
	p.MaxMessageLength = 1024 * 1024 * 512

	p.clientDriver = clientDriver
	p.remoteAddr = address
	p.reqSigMap = make(map[uint64]offheap.MutexUintptr, 128)

	return nil
}

func (p *Client) Start() error {
	var err error
	err = p.Conn.Connect(p.remoteAddr)
	if err != nil {
		return err
	}

	go func() {
		err = p.cronReadResponse()
		if err != nil &&
			strings.Contains(err.Error(), "use of closed network connection") == false {
			log.Warn(err)
		}
	}()

	return nil
}

func (p *Client) Close() error {
	var err error
	err = p.Conn.Close()
	if err != nil {
		return err
	}

	return nil
}

func (p *Client) AllocRequestID() uint64 {
	return atomic.AddUint64(&p.maxRequestID, 1)
}
