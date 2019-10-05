package snet

import (
	"soloos/common/log"
	"soloos/common/snettypes"
	"soloos/solodb/offheap"
	"strings"
	"sync"
	"sync/atomic"
)

type SrpcClient struct {
	doingNetQueryChan chan snettypes.NetQuery
	doingNetQueryConn snettypes.Connection

	MaxMessageLength uint32

	clientDriver *SrpcClientDriver
	remoteAddr   string
	maxRequestID uint64

	reqSigMapMutex sync.Mutex
	reqSigMap      map[uint64]offheap.MutexUintptr // map RequestID to netConnReadSigsIndex
}

func (p *SrpcClient) Init(clientDriver *SrpcClientDriver, address string) error {
	p.MaxMessageLength = 1024 * 1024 * 512

	p.clientDriver = clientDriver
	p.remoteAddr = address
	p.reqSigMap = make(map[uint64]offheap.MutexUintptr, 128)

	return nil
}

func (p *SrpcClient) Start() error {
	var err error
	p.doingNetQueryChan = make(chan snettypes.NetQuery, 1)
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

func (p *SrpcClient) Close(closeResonErr error) error {
	var err error
	err = p.doingNetQueryConn.Close(closeResonErr)
	if err != nil {
		return err
	}

	return nil
}

func (p *SrpcClient) AllocRequestID() uint64 {
	return atomic.AddUint64(&p.maxRequestID, 1)
}
