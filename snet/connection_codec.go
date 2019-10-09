package snet

import (
	"encoding/gob"
	"soloos/common/util"
	"sync"
)

var (
	snetCodecBytesPool sync.Pool
)

func init() {
	snetCodecBytesPool.New = func() interface{} {
		return util.Buffer{}
	}
}

func snetCodecBytesPoolPut(buf util.Buffer) {
	buf.Reset()
	snetCodecBytesPool.Put(buf)
}

func (p *Connection) prepareCodec() {
}

func (p *Connection) releaseCodec() {
}

func (p *Connection) Marshal(buf *util.Buffer, req interface{}) error {
	var enc = gob.NewEncoder(buf)
	var err error
	err = enc.Encode(req)
	return err
}
