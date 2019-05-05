package types

import (
	snettypes "soloos/common/snet/types"
)

type SWALMember struct {
	PeerID   snettypes.PeerID
	IsLeader int
}
