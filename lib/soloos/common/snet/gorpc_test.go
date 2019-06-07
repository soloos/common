package snet

import (
	"net"
	"net/http"
	"net/rpc"
	"runtime"
	"soloos/common/util"
	"sync"
	"testing"
	"time"
)

var (
	goRPCServer     *rpc.Server
	goRPCServerLn   net.Listener
	goRPCServerAddr string
)

type Hello int
type HelloReply [rpcMessageLen]byte

func (t *Hello) Say(args *string, reply *HelloReply) error {
	copy((*reply)[:], []byte(*args))
	return nil
}

func cleanGORpcServer() {
	if goRPCServer != nil {
		goRPCServerLn.Close()
		goRPCServer = nil
		runtime.GC()
	}
}

func runGORPCServer() (string, error) {
	var err error
	if goRPCServer != nil {
		return goRPCServerAddr, nil
	}

	goRPCServerAddr = allocAddr()
	goRPCServer = rpc.NewServer()
	goRPCServer.Register(new(Hello))
	goRPCServer.HandleHTTP(rpc.DefaultRPCPath, rpc.DefaultDebugPath)
	goRPCServerLn, err = net.Listen("tcp", goRPCServerAddr)
	util.AssertErrIsNil(err)
	go func() {
		http.Serve(goRPCServerLn, nil)
	}()

	time.Sleep(1 * time.Second)

	return goRPCServerAddr, nil
}

func TestGORPCServer(t *testing.T) {
	addr, err := runGORPCServer()
	util.AssertErrIsNil(err)

	var serviceSig sync.WaitGroup
	serviceSig.Add(1)

	client, err := rpc.DialHTTP("tcp", addr)
	util.AssertErrIsNil(err)
	var reply HelloReply

	go func() {
		err = client.Call("Hello.Say", rpcMessage, &reply)
		util.AssertErrIsNil(err)
		serviceSig.Done()
	}()

	serviceSig.Wait()
}

func BenchmarkGORPCServer(b *testing.B) {
	runtime.GC()
	addr, err := runGORPCServer()
	util.AssertErrIsNil(err)

	var serviceSig sync.WaitGroup
	serviceSig.Add(b.N)

	client, err := rpc.DialHTTP("tcp", addr)
	util.AssertErrIsNil(err)
	var reply string

	b.ResetTimer()

	for n := 0; n < b.N; n++ {
		// if false {
		if true {
			go func() {
				client.Call("Hello.Say", rpcMessage, &reply)

				serviceSig.Done()
			}()
		} else {
			serviceSig.Done()
		}
	}

	serviceSig.Wait()
}
