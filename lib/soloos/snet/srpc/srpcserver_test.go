package srpc

import (
	"fmt"
	"runtime"
	"soloos/snet/protocol"
	"soloos/snet/types"
	"soloos/util"
	"soloos/util/offheap"
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
	// serviceReadBuf = make([]byte, requestContentLen)
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

	srpcServer.RegisterService("/test", func(requestID uint64, requestContentLen uint32, conn *types.Connection) error {
		var err error
		{
			// read
			// var serviceReadBuf = make([]byte, requestContentLen)
			err = conn.ReadAll(serviceReadBuf)
			if err != nil {
				return err
			}

			var o protocol.MessageTest0
			o.Init(serviceReadBuf, flatbuffers.GetUOffsetT(serviceReadBuf))
			if !assert.ObjectsAreEqualValues(rpcMessageBytes, o.Data0()) {
				panic(string(o.Data0()))
			}
		}

		{
			// write
			err = conn.Response(requestID, rpcMessageBytes)
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

	assert.NoError(t, clientDriver.Init(&offheap.DefaultOffheapDriver))

	assert.NoError(t, offheap.DefaultOffheapDriver.InitRawObjectPool(&peerPool, int(types.PeerStructSize), -1, nil, nil))
	uPeer = types.PeerUintptr(peerPool.AllocRawObject())
	uPeer.Ptr().SetAddress(addr)
	assert.NoError(t, clientDriver.RegisterClient(uPeer))

	go func() {
		for i := 0; i < callTimes; i += 1 {
			if true {
				go func() {
					var (
						request         types.Request
						response        types.Response
						protocolBuilder flatbuffers.Builder
					)

					protocolBuilder.Reset()
					data0 := protocolBuilder.CreateString(rpcMessage)
					protocol.MessageTest2Start(&protocolBuilder)
					protocol.MessageTest2AddData0(&protocolBuilder, data0)
					protocol.MessageTest2AddData1(&protocolBuilder, 322)
					protocolBuilder.Finish(protocol.MessageTest0End(&protocolBuilder))
					request.Body = protocolBuilder.Bytes[protocolBuilder.Head():]

					assert.NoError(t, clientDriver.Call(uPeer, "/test", &request, &response))
					var resp = make([]byte, response.BodySize)
					util.AssertErrIsNil(clientDriver.ReadResponse(uPeer, &request, &response, &resp))
					assert.Equal(t, rpcMessageBytes, resp)
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
	// cleanRpcGOServer()
	runtime.GC()
	addr, err := runSRPCServer()
	util.AssertErrIsNil(err)

	var (
		clientDriver ClientDriver
		peerPool     offheap.RawObjectPool
		uPeer        types.PeerUintptr
	)

	util.AssertErrIsNil(clientDriver.Init(&offheap.DefaultOffheapDriver))

	util.AssertErrIsNil(offheap.DefaultOffheapDriver.InitRawObjectPool(&peerPool, int(types.PeerStructSize), -1, nil, nil))
	uPeer = types.PeerUintptr(peerPool.AllocRawObject())
	copy(uPeer.Ptr().Address[:], []byte(addr))
	util.AssertErrIsNil(clientDriver.RegisterClient(uPeer))

	var protocolBuilder flatbuffers.Builder
	protocolBuilder.Reset()
	data0 := protocolBuilder.CreateString(rpcMessage)
	protocol.MessageTest2Start(&protocolBuilder)
	protocol.MessageTest2AddData0(&protocolBuilder, data0)
	protocol.MessageTest2AddData1(&protocolBuilder, 322)
	protocolBuilder.Finish(protocol.MessageTest0End(&protocolBuilder))

	// var resp = make([]byte, response.BodySize)
	var resp = make([]byte, len(rpcMessageBytes))
	var serviceSig sync.WaitGroup
	serviceSig.Add(b.N)

	b.ResetTimer()
	requestBody := protocolBuilder.Bytes[protocolBuilder.Head():]
	for n := 0; n < b.N; n++ {
		// if false {
		if true {
			go func() {
				var (
					request = types.Request{
						Body: requestBody,
					}
					response types.Response
				)

				// util.AssertErrIsNil(clientDriver.AsyncCall(uPeer, "/test", &request, &response))
				// util.AssertErrIsNil(clientDriver.WaitResponse(uPeer, &request, &response))
				// util.AssertErrIsNil(clientDriver.ReadResponse(uPeer, &request, &response, &resp))

				util.AssertErrIsNil(clientDriver.Call(uPeer, "/test", &request, &response))
				util.AssertErrIsNil(clientDriver.ReadResponse(uPeer, &request, &response, &resp))

				// if !assert.ObjectsAreEqualValues(rpcMessageBytes, resp) {
				// panic("not equal")
				// }
				serviceSig.Done()
			}()
		} else {
			serviceSig.Done()
		}
	}

	serviceSig.Wait()

	util.AssertErrIsNil(clientDriver.CloseClient(uPeer))
}
