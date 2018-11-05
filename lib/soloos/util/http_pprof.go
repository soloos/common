package util

import (
	"soloos/log"
	"net/http"
	_ "net/http/pprof"
)

func PProfServe(pprofListenAddr string) {
	log.Info("pprof start", pprofListenAddr)
	err := http.ListenAndServe(pprofListenAddr, nil)
	AssertErrIsNil(err)
}
