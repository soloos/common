package srpc

import (
	"runtime"
	"soloos/snet/protocol"
	"soloos/snet/types"
	"sync"
	"testing"
	"time"

	flatbuffers "github.com/google/flatbuffers/go"
	"github.com/stretchr/testify/assert"
)

func TestSRPCServer(t *testing.T) {
	var (
		server       Server
		clientDriver ClientDriver
		addr         = "127.0.0.1:12052"
	)
	assert.NoError(t, server.Init())
	assert.NoError(t, clientDriver.Init())

	runtime.GOMAXPROCS(4)

	var serviceSig sync.WaitGroup
	server.RegisterService("/test", func(requestID uint64, conn *types.Connection) {
		var d = make([]byte, conn.LastRequestReadLimit)
		err := conn.ReadAll(d)
		assert.NoError(t, err)

		var o protocol.MessageTest0
		o.Init(d, flatbuffers.GetUOffsetT(d))
		assert.Equal(t, addr, string(o.Data0()))

		{
			resp := []byte("success")
			respLen := uint32(len(resp))
			conn.WriteAcquire()
			conn.WriteHeader(requestID, respLen)
			conn.WriteAll(resp)
			conn.WriteRelease()
		}

		serviceSig.Done()
	})
	go func() {
		assert.NoError(t, server.Serve("tcp", addr))
	}()

	serviceSig.Add(2)
	go func() {
		time.Sleep(1 * time.Second)
		var (
			client         Client
			clientRequest  types.ClientRequest
			clientResponse types.ClientResponse
		)

		assert.NoError(t, client.Init(addr))

		{
			var protocolBuilder flatbuffers.Builder
			protocolBuilder.Reset()
			data0 := protocolBuilder.CreateString(addr)
			protocol.MessageTest2Start(&protocolBuilder)
			protocol.MessageTest2AddData0(&protocolBuilder, data0)
			protocol.MessageTest2AddData1(&protocolBuilder, 322)
			protocolBuilder.Finish(protocol.MessageTest0End(&protocolBuilder))
			clientRequest.Body = protocolBuilder.Bytes[protocolBuilder.Head():]
		}

		{
			client.Conn.ReadAcquire()
			clientDriver.Call(&client, "/test", &clientRequest, &clientResponse)
			var resp = make([]byte, clientResponse.BodySize)
			assert.NoError(t, client.Conn.ReadAll(resp))
			assert.Equal(t, "success", string(resp))
			client.Conn.ReadRelease()
			serviceSig.Done()
		}
	}()

	serviceSig.Wait()
}
