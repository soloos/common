package snettypes

import (
	"encoding/json"
)

type HeartBeatServerOptions struct {
	PeerID     PeerID
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

	p.PeerID = StrToPeerID(ret.PeerID)
	p.DurationMS = ret.DurationMS

	return nil
}
