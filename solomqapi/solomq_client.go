package solomqapi

import (
	"soloos/common/soloosbase"
	"time"
)

type SolomqClient struct {
	*soloosbase.SoloosEnv
	normalCallRetryTimes        int
	waitAliveEveryRetryWaitTime time.Duration
}

func (p *SolomqClient) Init(soloosEnv *soloosbase.SoloosEnv) error {
	p.SoloosEnv = soloosEnv
	p.normalCallRetryTimes = 3
	p.waitAliveEveryRetryWaitTime = time.Second * 3
	return nil
}
