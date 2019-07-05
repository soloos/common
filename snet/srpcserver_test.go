package snet

import (
	"fmt"
	"runtime"
	"soloos/common/snetprotocol"
	"soloos/common/snettypes"
	"soloos/common/util"
	"soloos/sdbone/offheap"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	flatbuffers "github.com/google/flatbuffers/go"
	"github.com/stretchr/testify/assert"
)

const (
	rpcMessageLen = 1024 * 64
)

var (
	addrPort int32 = 10200
	// serviceReadBuf = make([]byte, reqBodySize)
	serviceReadBuf     = make([]byte, rpcMessageLen+32)
	rawRpcMessageBytes = [rpcMessageLen]byte{}
	rpcMessage         = string(rawRpcMessageBytes[:])
	rpcMessageBytes    = []byte(rpcMessage)

	srpcServer     *SRPCServer
	srpcServerAddr string
)

func allocAddr() string {
	return fmt.Sprintf("127.0.0.1:%d", atomic.AddInt32(&addrPort, 1))
}

func runSRPCServer() (string, error) {
	if srpcServer != nil {
		return srpcServerAddr, nil
	}

	srpcServer = new(SRPCServer)
	srpcServerAddr = allocAddr()
	util.AssertErrIsNil(srpcServer.Init("tcp", srpcServerAddr))

	srpcServer.RegisterService("/test", func(serviceReq *snettypes.NetQuery) error {
		var err error
		{
			// read
			// var serviceReadBuf = make([]byte, reqBodySize)
			if len(serviceReadBuf) != int(serviceReq.BodySize) {
				panic("error")
			}
			err = serviceReq.ReadAll(serviceReadBuf)
			if err != nil {
				return err
			}

			var o snetprotocol.MessageTest0
			o.Init(serviceReadBuf, flatbuffers.GetUOffsetT(serviceReadBuf))
			if assert.ObjectsAreEqualValues(rpcMessageBytes, o.Data0()) == false {
				panic(string(o.Data0()))
			}
		}

		{
			// write
			err = serviceReq.SimpleResponse(serviceReq.ReqID, rpcMessageBytes)
			if err != nil {
				return err
			}
		}

		return nil
	})

	go func() {
		util.AssertErrIsNil(srpcServer.Serve())
	}()

	time.Sleep(1 * time.Second)

	return srpcServerAddr, nil
}

func TestSRPCServer(t *testing.T) {
	var defaultOffheapDriver offheap.OffheapDriver
	var defaultSNetDriver NetDriver
	util.AssertErrIsNil(defaultOffheapDriver.Init())
	util.AssertErrIsNil(defaultSNetDriver.Init(&defaultOffheapDriver))

	addr, err := runSRPCServer()
	assert.NoError(t, err)

	var (
		callTimes    = 128
		serviceSig   sync.WaitGroup
		clientDriver SRPCClientDriver
		peerPool     offheap.RawObjectPool
		peer         snettypes.Peer
	)
	serviceSig.Add(callTimes)

	assert.NoError(t, clientDriver.Init(&defaultOffheapDriver, &defaultSNetDriver))

	assert.NoError(t, peerPool.Init(int(snettypes.PeerStructSize), -1, nil, nil))
	clientDriver.netDriver.InitPeerID((*snettypes.PeerID)(&peer.ID))
	peer.SetAddress(addr)
	peer.ServiceProtocol = snettypes.ProtocolSRPC
	clientDriver.netDriver.RegisterPeer(peer)

	go func() {
		for i := 0; i < callTimes; i += 1 {
			if true {
				go func() {
					var (
						req             [7]snettypes.Request
						resp            [7]snettypes.Response
						protocolBuilder flatbuffers.Builder
					)

					protocolBuilder.Reset()
					data0 := protocolBuilder.CreateString(rpcMessage)
					snetprotocol.MessageTest2Start(&protocolBuilder)
					snetprotocol.MessageTest2AddData0(&protocolBuilder, data0)
					snetprotocol.MessageTest2AddData1(&protocolBuilder, 322)
					protocolBuilder.Finish(snetprotocol.MessageTest0End(&protocolBuilder))

					for i := 0; i < len(req); i++ {
						req[i].Param = protocolBuilder.Bytes[protocolBuilder.Head():]
					}

					util.AssertErrIsNil(clientDriver.Call(peer.ID, "/notexist", &req[0], &resp[0]))
					util.AssertErrIsNil(clientDriver.ReadResponse(peer.ID, &req[0], &resp[0], nil))

					{
						util.AssertErrIsNil(clientDriver.AsyncCall(peer.ID, "/notexist", &req[1], &resp[1]))
						util.AssertErrIsNil(clientDriver.AsyncCall(peer.ID, "/notexist", &req[2], &resp[2]))
						var wg sync.WaitGroup
						wg.Add(1)
						go func() {
							util.AssertErrIsNil(clientDriver.WaitResponse(peer.ID, &req[1], &resp[1]))
							util.AssertErrIsNil(clientDriver.ReadResponse(peer.ID, &req[1], &resp[1], nil))
							util.AssertErrIsNil(clientDriver.WaitResponse(peer.ID, &req[2], &resp[2]))
							util.AssertErrIsNil(clientDriver.ReadResponse(peer.ID, &req[2], &resp[2], nil))
							wg.Done()
						}()
						wg.Wait()
					}

					util.AssertErrIsNil(clientDriver.Call(peer.ID, "/notexist", &req[3], &resp[3]))
					util.AssertErrIsNil(clientDriver.ReadResponse(peer.ID, &req[3], &resp[3], nil))

					util.AssertErrIsNil(clientDriver.AsyncCall(peer.ID, "/test", &req[4], &resp[4]))
					util.AssertErrIsNil(clientDriver.WaitResponse(peer.ID, &req[4], &resp[4]))
					var respBody = make([]byte, resp[4].BodySize)
					util.AssertErrIsNil(clientDriver.ReadResponse(peer.ID, &req[4], &resp[4], respBody))
					assert.Equal(t, rpcMessageBytes, respBody)

					util.AssertErrIsNil(clientDriver.AsyncCall(peer.ID, "/notexist", &req[5], &resp[5]))
					util.AssertErrIsNil(clientDriver.WaitResponse(peer.ID, &req[5], &resp[5]))
					util.AssertErrIsNil(clientDriver.ReadResponse(peer.ID, &req[5], &resp[5], nil))

					util.AssertErrIsNil(clientDriver.Call(peer.ID, "/notexist", &req[6], &resp[6]))
					util.AssertErrIsNil(clientDriver.ReadResponse(peer.ID, &req[6], &resp[6], nil))

					serviceSig.Done()
				}()
			} else {
				serviceSig.Done()
			}
		}
	}()

	serviceSig.Wait()

	assert.NoError(t, clientDriver.CloseClient(peer.ID))

	time.Sleep(1 * time.Second)
}

func BenchmarkSRPCServer(b *testing.B) {
	var defaultOffheapDriver offheap.OffheapDriver
	var defaultSNetDriver NetDriver
	util.AssertErrIsNil(defaultOffheapDriver.Init())
	util.AssertErrIsNil(defaultSNetDriver.Init(&defaultOffheapDriver))

	// cleanRpcGOServer()
	runtime.GC()
	addr, err := runSRPCServer()
	util.AssertErrIsNil(err)

	var (
		clientDriver SRPCClientDriver
		peerPool     offheap.RawObjectPool
		uPeer        snettypes.PeerUintptr
		peerID       snettypes.PeerID
	)

	util.AssertErrIsNil(clientDriver.Init(&defaultOffheapDriver, &defaultSNetDriver))

	util.AssertErrIsNil(peerPool.Init(int(snettypes.PeerStructSize), -1, nil, nil))
	uPeer = snettypes.PeerUintptr(peerPool.AllocRawObject())
	uPeer.Ptr().SetAddress(addr)
	peerID = uPeer.Ptr().ID
	defaultSNetDriver.RegisterPeer(*uPeer.Ptr())

	var protocolBuilder flatbuffers.Builder
	protocolBuilder.Reset()
	data0 := protocolBuilder.CreateString(rpcMessage)
	snetprotocol.MessageTest2Start(&protocolBuilder)
	snetprotocol.MessageTest2AddData0(&protocolBuilder, data0)
	snetprotocol.MessageTest2AddData1(&protocolBuilder, 322)
	protocolBuilder.Finish(snetprotocol.MessageTest0End(&protocolBuilder))

	// var resp = make([]byte, resp.BodySize)
	var respBody = make([]byte, len(rpcMessageBytes))
	var serviceSig sync.WaitGroup

	b.ResetTimer()
	serviceSig.Add(b.N)
	reqBody := protocolBuilder.Bytes[protocolBuilder.Head():]
	for n := 0; n < b.N; n++ {
		// if false {
		if true {
			go func() {
				var (
					req = snettypes.Request{
						Param: reqBody,
					}
					resp snettypes.Response
				)

				// util.AssertErrIsNil(clientDriver.AsyncCall(peerID, "/test", &req, &resp))
				// util.AssertErrIsNil(clientDriver.WaitResponse(peerID, &req, &resp))
				// util.AssertErrIsNil(clientDriver.ReadResponse(peerID, &req, &resp, &resp))

				util.AssertErrIsNil(clientDriver.Call(peerID, "/test", &req, &resp))
				util.AssertErrIsNil(clientDriver.ReadResponse(peerID, &req, &resp, respBody))

				// if assert.ObjectsAreEqualValues(rpcMessageBytes, resp) == false{
				// panic("not equal")
				// }
				serviceSig.Done()
			}()
		} else {
			serviceSig.Done()
		}
	}

	serviceSig.Wait()

	clientDriver.CloseClient(peerID)
}
