package solomqtypes

import (
	"soloos/common/snet"
)

const (
	SolomqMemberRoleUnknown  = -1
	SolomqMemberRoleLeader   = 1
	SolomqMemberRoleFollower = 0
)

type SolomqMember struct {
	PeerID snet.PeerID
	Role   int
}
