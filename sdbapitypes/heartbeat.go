package sdbapitypes

import (
	"encoding/json"
	"soloos/common/snettypes"
)

type HeartBeatServerOptions struct {
	PeerID     snettypes.PeerID
	DurationMS int64
}

type HeartBeatServerOptionsJson struct {
	PeerID     string
	DurationMS int64
}

func (p *HeartBeatServerOptions) UnmarshalJSON(data []byte) error {
	var ret HeartBeatServerOptionsJson
	if err := json.Unmarshal(data, &ret); err != nil {
		return err
	}

	p.PeerID = snettypes.StrToPeerID(ret.PeerID)
	p.DurationMS = ret.DurationMS

	return nil
}
