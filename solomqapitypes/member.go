package solomqapitypes

import (
	"soloos/common/snettypes"
)

const (
	SolomqMemberRoleUnknown  = -1
	SolomqMemberRoleLeader   = 1
	SolomqMemberRoleFollower = 0
)

type SolomqMember struct {
	PeerID snettypes.PeerID
	Role   int
}
