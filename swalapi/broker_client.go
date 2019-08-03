package swalapi

import (
	"soloos/common/soloosbase"
	"time"
)

type BrokerClient struct {
	*soloosbase.SoloOSEnv
	normalCallRetryTimes        int
	waitAliveEveryRetryWaitTime time.Duration
}

func (p *BrokerClient) Init(soloOSEnv *soloosbase.SoloOSEnv) error {
	p.SoloOSEnv = soloOSEnv
	p.normalCallRetryTimes = 3
	p.waitAliveEveryRetryWaitTime = time.Second * 3
	return nil
}
