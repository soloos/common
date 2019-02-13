package util

import (
	"bytes"
	"encoding/gob"
	"runtime/debug"
	"soloos/common/log"
)

func MustGobEncode(target interface{}) []byte {
	var buf bytes.Buffer
	enc := gob.NewEncoder(&buf)
	AssertErrIsNil(enc.Encode(target))
	return buf.Bytes()
}

func MustGobDecode(bufBytes []byte, target interface{}) {
	buf := bytes.NewBuffer(bufBytes)
	dec := gob.NewDecoder(buf)
	err := dec.Decode(target)
	if err != nil {
		log.Error(bufBytes)
		log.Error(err)
		debug.PrintStack()
	}

	AssertErrIsNil(err)
}
