package srpc

import (
	"fmt"
	"runtime"
	"soloos/common/snet/protocol"
	"soloos/common/snet/types"
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

	srpcServer     *Server
	srpcServerAddr string
)

func allocAddr() string {
	return fmt.Sprintf("127.0.0.1:%d", atomic.AddInt32(&addrPort, 1))
}

func runSRPCServer() (string, error) {
	if srpcServer != nil {
		return srpcServerAddr, nil
	}

	srpcServer = new(Server)
	srpcServerAddr = allocAddr()
	util.AssertErrIsNil(srpcServer.Init("tcp", srpcServerAddr))

	srpcServer.RegisterService("/test", func(serviceReq *types.NetQuery) error {
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

			var o protocol.MessageTest0
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
	util.AssertErrIsNil(defaultOffheapDriver.Init())

	addr, err := runSRPCServer()
	assert.NoError(t, err)

	var (
		callTimes    = 128
		serviceSig   sync.WaitGroup
		clientDriver ClientDriver
		peerPool     offheap.RawObjectPool
		uPeer        types.PeerUintptr
	)
	serviceSig.Add(callTimes)

	assert.NoError(t, clientDriver.Init(&defaultOffheapDriver))

	assert.NoError(t, peerPool.Init(int(types.PeerStructSize), -1, nil, nil))
	uPeer = types.PeerUintptr(peerPool.AllocRawObject())
	uPeer.Ptr().SetAddress(addr)

	go func() {
		for i := 0; i < callTimes; i += 1 {
			if true {
				go func() {
					var (
						req             [7]types.Request
						resp            [7]types.Response
						protocolBuilder flatbuffers.Builder
					)

					protocolBuilder.Reset()
					data0 := protocolBuilder.CreateString(rpcMessage)
					protocol.MessageTest2Start(&protocolBuilder)
					protocol.MessageTest2AddData0(&protocolBuilder, data0)
					protocol.MessageTest2AddData1(&protocolBuilder, 322)
					protocolBuilder.Finish(protocol.MessageTest0End(&protocolBuilder))

					for i := 0; i < len(req); i++ {
						req[i].Param = protocolBuilder.Bytes[protocolBuilder.Head():]
					}

					util.AssertErrIsNil(clientDriver.Call(uPeer, "/notexist", &req[0], &resp[0]))
					util.AssertErrIsNil(clientDriver.ReadResponse(uPeer, &req[0], &resp[0], nil))

					{
						util.AssertErrIsNil(clientDriver.AsyncCall(uPeer, "/notexist", &req[1], &resp[1]))
						util.AssertErrIsNil(clientDriver.AsyncCall(uPeer, "/notexist", &req[2], &resp[2]))
						var wg sync.WaitGroup
						wg.Add(1)
						go func() {
							util.AssertErrIsNil(clientDriver.WaitResponse(uPeer, &req[1], &resp[1]))
							util.AssertErrIsNil(clientDriver.ReadResponse(uPeer, &req[1], &resp[1], nil))
							util.AssertErrIsNil(clientDriver.WaitResponse(uPeer, &req[2], &resp[2]))
							util.AssertErrIsNil(clientDriver.ReadResponse(uPeer, &req[2], &resp[2], nil))
							wg.Done()
						}()
						wg.Wait()
					}

					util.AssertErrIsNil(clientDriver.Call(uPeer, "/notexist", &req[3], &resp[3]))
					util.AssertErrIsNil(clientDriver.ReadResponse(uPeer, &req[3], &resp[3], nil))

					util.AssertErrIsNil(clientDriver.AsyncCall(uPeer, "/test", &req[4], &resp[4]))
					util.AssertErrIsNil(clientDriver.WaitResponse(uPeer, &req[4], &resp[4]))
					var respBody = make([]byte, resp[4].BodySize)
					util.AssertErrIsNil(clientDriver.ReadResponse(uPeer, &req[4], &resp[4], respBody))
					assert.Equal(t, rpcMessageBytes, respBody)

					util.AssertErrIsNil(clientDriver.AsyncCall(uPeer, "/notexist", &req[5], &resp[5]))
					util.AssertErrIsNil(clientDriver.WaitResponse(uPeer, &req[5], &resp[5]))
					util.AssertErrIsNil(clientDriver.ReadResponse(uPeer, &req[5], &resp[5], nil))

					util.AssertErrIsNil(clientDriver.Call(uPeer, "/notexist", &req[6], &resp[6]))
					util.AssertErrIsNil(clientDriver.ReadResponse(uPeer, &req[6], &resp[6], nil))

					serviceSig.Done()
				}()
			} else {
				serviceSig.Done()
			}
		}
	}()

	serviceSig.Wait()

	assert.NoError(t, clientDriver.CloseClient(uPeer))

	time.Sleep(1 * time.Second)
}

func BenchmarkSRPCServer(b *testing.B) {
	var defaultOffheapDriver offheap.OffheapDriver
	util.AssertErrIsNil(defaultOffheapDriver.Init())

	// cleanRpcGOServer()
	runtime.GC()
	addr, err := runSRPCServer()
	util.AssertErrIsNil(err)

	var (
		clientDriver ClientDriver
		peerPool     offheap.RawObjectPool
		uPeer        types.PeerUintptr
	)

	util.AssertErrIsNil(clientDriver.Init(&defaultOffheapDriver))

	util.AssertErrIsNil(peerPool.Init(int(types.PeerStructSize), -1, nil, nil))
	uPeer = types.PeerUintptr(peerPool.AllocRawObject())
	uPeer.Ptr().SetAddress(addr)

	var protocolBuilder flatbuffers.Builder
	protocolBuilder.Reset()
	data0 := protocolBuilder.CreateString(rpcMessage)
	protocol.MessageTest2Start(&protocolBuilder)
	protocol.MessageTest2AddData0(&protocolBuilder, data0)
	protocol.MessageTest2AddData1(&protocolBuilder, 322)
	protocolBuilder.Finish(protocol.MessageTest0End(&protocolBuilder))

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
					req = types.Request{
						Param: reqBody,
					}
					resp types.Response
				)

				// util.AssertErrIsNil(clientDriver.AsyncCall(uPeer, "/test", &req, &resp))
				// util.AssertErrIsNil(clientDriver.WaitResponse(uPeer, &req, &resp))
				// util.AssertErrIsNil(clientDriver.ReadResponse(uPeer, &req, &resp, &resp))

				util.AssertErrIsNil(clientDriver.Call(uPeer, "/test", &req, &resp))
				util.AssertErrIsNil(clientDriver.ReadResponse(uPeer, &req, &resp, respBody))

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

	clientDriver.CloseClient(uPeer)
}
