package types

import (
	snettypes "soloos/common/snet/types"
)

const (
	SWALMemberRoleLeader   = 1
	SWALMemberRoleFollower = 0
)

type SWALMember struct {
	PeerID   snettypes.PeerID
	IsLeader int
}
