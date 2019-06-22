package swalapi

import (
	"soloos/common/soloosbase"
	"time"
)

type SWALAgentClient struct {
	*soloosbase.SoloOSEnv
	normalCallRetryTimes        int
	waitAliveEveryRetryWaitTime time.Duration
}

func (p *SWALAgentClient) Init(soloOSEnv *soloosbase.SoloOSEnv) error {
	p.SoloOSEnv = soloOSEnv
	p.normalCallRetryTimes = 3
	p.waitAliveEveryRetryWaitTime = time.Second * 3
	return nil
}
