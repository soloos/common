package swalapitypes

import (
	"soloos/common/snettypes"
)

const (
	SWALMemberRoleUnknown  = -1
	SWALMemberRoleLeader   = 1
	SWALMemberRoleFollower = 0
)

type SWALMember struct {
	PeerID snettypes.PeerID
	Role   int
}
