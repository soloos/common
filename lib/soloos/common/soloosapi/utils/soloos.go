package base

import (
	"soloos/common/fsapi"
	soloosbase "soloos/common/soloosapi/base"
	"soloos/sdfs/libsdfs"
	"soloos/swal/libswal"
)

var (
	SoloOSInstance        SoloOS
	isDefaultSoloOSInited bool = false
)

type SoloOS struct {
	options Options
	soloosbase.SoloOSEnv

	SDFSClientDriver libsdfs.ClientDriver
	SDFSClient       libsdfs.Client
	PosixFS          fsapi.PosixFS
	SWALClientDriver libswal.ClientDriver
	SWALClient       libswal.Client
}

func InitSoloOSInstance(options Options) error {
	if isDefaultSoloOSInited {
		return nil
	}
	isDefaultSoloOSInited = true
	return SoloOSInstance.Init(options)
}

func (p *SoloOS) Init(options Options) error {
	var err error

	p.options = options
	err = p.SoloOSEnv.Init()
	if err != nil {
		return err
	}

	err = p.initSDFS()
	if err != nil {
		return err
	}

	err = p.initSWAL()
	if err != nil {
		return err
	}

	return nil
}

func (p *SoloOS) Serve() error {
	var err error
	err = p.SWALClientDriver.Serve()
	if err != nil {
		return err
	}

	return nil
}

func (p *SoloOS) Close() error {
	var err error
	err = p.SWALClientDriver.Close()
	if err != nil {
		return err
	}

	return nil
}
