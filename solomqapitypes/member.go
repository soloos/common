package solomqapitypes

import (
	"soloos/common/snettypes"
)

const (
	SOLOMQMemberRoleUnknown  = -1
	SOLOMQMemberRoleLeader   = 1
	SOLOMQMemberRoleFollower = 0
)

type SOLOMQMember struct {
	PeerID snettypes.PeerID
	Role   int
}
