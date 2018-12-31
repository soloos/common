package util

import (
	"net/http"
	_ "net/http/pprof"
	"soloos/log"
)

func PProfServe(pprofListenAddr string) {
	log.Info("pprof start", pprofListenAddr)
	err := http.ListenAndServe(pprofListenAddr, nil)
	AssertErrIsNil(err)
}
