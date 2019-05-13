package srpc

import (
	"soloos/common/log"
	"soloos/common/snet/types"
	"soloos/sdbone/offheap"
	"strings"
	"sync"
	"sync/atomic"
)

type Client struct {
	doingNetQueryChan chan types.NetQuery
	doingNetQueryConn types.Connection

	MaxMessageLength uint32

	clientDriver *ClientDriver
	remoteAddr   string
	maxRequestID uint64

	reqSigMapMutex sync.Mutex
	reqSigMap      map[uint64]offheap.MutexUintptr // map RequestID to netConnReadSigsIndex
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
	p.doingNetQueryChan = make(chan types.NetQuery, 1)
	err = p.doingNetQueryConn.Connect(p.remoteAddr)
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

func (p *Client) Close(closeResonErr error) error {
	var err error
	err = p.doingNetQueryConn.Close(closeResonErr)
	if err != nil {
		return err
	}

	return nil
}

func (p *Client) AllocRequestID() uint64 {
	return atomic.AddUint64(&p.maxRequestID, 1)
}
