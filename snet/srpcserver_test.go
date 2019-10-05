package snet

import (
	"fmt"
	"runtime"
	"soloos/common/iron"
	"soloos/common/snetprotocol"
	"soloos/common/snettypes"
	"soloos/common/util"
	"soloos/solodb/offheap"
	"sync"
	"sync/atomic"
	"testing"
	"time"
	"unsafe"

	"github.com/stretchr/testify/assert"
)

type MessageTest0 = snetprotocol.MessageTest0
type MessageTest1 = snetprotocol.MessageTest1
type MessageTest2 = snetprotocol.MessageTest2

const (
	rpcMessageLen = 1024 * 64
)

var (
	addrPort           int32 = 10200
	rawRpcMessageBytes       = [rpcMessageLen]byte{}
	rpcMessage         string
	rpcMessageBytes    []byte

	srpcServer     *SrpcServer
	srpcServerAddr string
)

func init() {
	for i := 0; i < len(rawRpcMessageBytes); i++ {
		rawRpcMessageBytes[i] = 'a'
	}
	rpcMessage = string(rawRpcMessageBytes[:])
	rpcMessageBytes = []byte(rpcMessage)
}

func allocAddr() string {
	return fmt.Sprintf("127.0.0.1:%d", atomic.AddInt32(&addrPort, 1))
}

func runSrpcServer() (string, error) {
	if srpcServer != nil {
		return srpcServerAddr, nil
	}

	srpcServer = new(SrpcServer)
	srpcServerAddr = allocAddr()
	util.AssertErrIsNil(srpcServer.Init("tcp", srpcServerAddr))

	srpcServer.RegisterService("/test", func(reqCtx *snettypes.SNetReqContext, msg string) (string, error) {
		return msg, nil
	})

	srpcServer.RegisterService("/testoffheap", func(reqCtx *snettypes.SNetReqContext) (string, error) {
		util.AssertTrue(int(reqCtx.BodySize) == len(rpcMessageBytes))
		return "fuck", nil
	})

	go func() {
		util.AssertErrIsNil(srpcServer.Serve())
	}()

	time.Sleep(1 * time.Second)

	return srpcServerAddr, nil
}

func assertIsCmdNotFound(err error) {
	util.AssertTrue(err.Error() == iron.ErrCmdNotFound.Error())
}

func TestSrpcServer(t *testing.T) {
	var defaultOffheapDriver offheap.OffheapDriver
	var defaultSNetDriver NetDriver
	util.AssertErrIsNil(defaultOffheapDriver.Init())
	util.AssertErrIsNil(defaultSNetDriver.Init(&defaultOffheapDriver))

	addr, err := runSrpcServer()
	assert.NoError(t, err)

	var (
		callTimes    = 128
		serviceSig   sync.WaitGroup
		clientDriver SrpcClientDriver
		peerPool     offheap.RawObjectPool
		peer         snettypes.Peer
	)
	serviceSig.Add(callTimes)

	assert.NoError(t, clientDriver.Init(&defaultOffheapDriver, &defaultSNetDriver))

	assert.NoError(t, peerPool.Init(int(snettypes.PeerStructSize), -1, nil, nil))
	clientDriver.netDriver.InitPeerID((*snettypes.PeerID)(&peer.ID))
	peer.SetAddress(addr)
	peer.ServiceProtocol = snettypes.ProtocolSrpc
	clientDriver.netDriver.RegisterPeer(peer)

	for i := 0; i < callTimes; i += 1 {
		go func(t *testing.T,
			serviceSig *sync.WaitGroup,
			peer *snettypes.Peer,
			clientDriver *SrpcClientDriver) {
			var (
				req      [7]snettypes.SNetReq
				resp     [7]snettypes.SNetResp
				respBody []byte
			)

			data0 := MessageTest2{
				Data0: rpcMessage,
				Data1: 322,
			}

			util.AssertErrIsNil(clientDriver.Call(peer.ID, "/notexist", &req[0], &resp[0], data0))
			respBody = make([]byte, resp[0].BodySize)
			assertIsCmdNotFound(clientDriver.ReadResponse(peer.ID, &req[0], &resp[0], respBody, nil))

			{
				util.AssertErrIsNil(clientDriver.AsyncCall(peer.ID, "/notexist", &req[1], &resp[1], data0))

				util.AssertErrIsNil(clientDriver.AsyncCall(peer.ID, "/notexist", &req[2], &resp[2], data0))

				var wg sync.WaitGroup
				wg.Add(1)
				go func() {
					var respBody []byte
					util.AssertErrIsNil(clientDriver.WaitResponse(peer.ID, &req[1], &resp[1]))
					respBody = make([]byte, resp[1].BodySize)
					assertIsCmdNotFound(clientDriver.ReadResponse(peer.ID, &req[1], &resp[1], respBody, nil))

					util.AssertErrIsNil(clientDriver.WaitResponse(peer.ID, &req[2], &resp[2]))
					respBody = make([]byte, resp[2].BodySize)
					assertIsCmdNotFound(clientDriver.ReadResponse(peer.ID, &req[2], &resp[2], respBody, nil))
					wg.Done()
				}()
				wg.Wait()
			}

			util.AssertErrIsNil(clientDriver.Call(peer.ID, "/notexist", &req[3], &resp[3], data0))
			respBody = make([]byte, resp[3].BodySize)
			assertIsCmdNotFound(clientDriver.ReadResponse(peer.ID, &req[3], &resp[3], respBody, nil))

			util.AssertErrIsNil(clientDriver.AsyncCall(peer.ID, "/test", &req[4], &resp[4], rpcMessage))
			util.AssertErrIsNil(clientDriver.WaitResponse(peer.ID, &req[4], &resp[4]))
			respBody = make([]byte, resp[4].BodySize)
			var msg = snettypes.Response{RespData: ""}
			util.AssertErrIsNil(clientDriver.ReadResponse(peer.ID, &req[4], &resp[4], respBody, &msg))
			assert.Equal(t, rpcMessage, msg.RespData)

			util.AssertErrIsNil(clientDriver.AsyncCall(peer.ID, "/notexist", &req[5], &resp[5], "test"))
			util.AssertErrIsNil(clientDriver.WaitResponse(peer.ID, &req[5], &resp[5]))
			respBody = make([]byte, resp[5].BodySize)
			assertIsCmdNotFound(clientDriver.ReadResponse(peer.ID, &req[5], &resp[5], respBody, nil))

			util.AssertErrIsNil(clientDriver.Call(peer.ID, "/notexist", &req[6], &resp[6], "test"))
			respBody = make([]byte, resp[6].BodySize)
			assertIsCmdNotFound(clientDriver.ReadResponse(peer.ID, &req[6], &resp[6], respBody, nil))

			serviceSig.Done()
		}(t, &serviceSig, &peer, &clientDriver)
	}

	serviceSig.Wait()

	assert.NoError(t, clientDriver.CloseClient(peer.ID))

	time.Sleep(1 * time.Second)
}

func BenchmarkSrpcServer(b *testing.B) {
	var defaultOffheapDriver offheap.OffheapDriver
	var defaultSNetDriver NetDriver
	util.AssertErrIsNil(defaultOffheapDriver.Init())
	util.AssertErrIsNil(defaultSNetDriver.Init(&defaultOffheapDriver))

	// cleanRpcGOServer()
	runtime.GC()
	addr, err := runSrpcServer()
	util.AssertErrIsNil(err)

	var (
		clientDriver SrpcClientDriver
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

	// var resp = make([]byte, resp.BodySize)
	var respBody = make([]byte, len(rpcMessageBytes)+1024)
	var serviceSig sync.WaitGroup

	b.ResetTimer()
	serviceSig.Add(b.N)
	for n := 0; n < b.N; n++ {
		// if false {
		if true {
			go func() {
				var (
					req  snettypes.SNetReq
					resp snettypes.SNetResp
				)

				// util.AssertErrIsNil(clientDriver.AsyncCall(peerID, "/test", &req, &resp))
				// util.AssertErrIsNil(clientDriver.WaitResponse(peerID, &req, &resp))
				// assertIsCmdNotFound(clientDriver.ReadResponse(peerID, &req, &resp, &resp))

				req.OffheapBody.OffheapBytes = uintptr(unsafe.Pointer(&rpcMessageBytes))
				req.OffheapBody.CopyOffset = 0
				req.OffheapBody.CopyEnd = len(rpcMessageBytes)
				util.AssertErrIsNil(clientDriver.Call(peerID, "/testoffheap", &req, &resp, nil))
				util.ChangeBytesArraySize(&respBody, int(resp.BodySize))
				util.AssertErrIsNil(clientDriver.ReadRawResponse(peerID, &req, &resp, respBody))
				util.AssertErrIsNil(resp.SkipReadRemaining())
				util.Ignore(respBody)
				util.Ignore(req)
				util.Ignore(resp)

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
