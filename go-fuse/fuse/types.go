package fuse

import (
	. "soloos/common/fsapi"
)

const (
	_DEFAULT_BACKGROUND_TASKS = 12
)

type _BmapIn struct {
	InHeader
	Block     uint64
	Blocksize uint32
	Padding   uint32
}

type _BmapOut struct {
	Block uint64
}

type _IoctlIn struct {
	InHeader
	Fh      uint64
	Flags   uint32
	Cmd     uint32
	Arg     uint64
	InSize  uint32
	OutSize uint32
}

type _IoctlOut struct {
	Result  int32
	Flags   uint32
	InIovs  uint32
	OutIovs uint32
}

type _PollIn struct {
	InHeader
	Fh      uint64
	Kh      uint64
	Flags   uint32
	Padding uint32
}

type _PollOut struct {
	Revents uint32
	Padding uint32
}

type _NotifyPollWakeupOut struct {
	Kh uint64
}

// batch forget is handled internally.
type _ForgetOne struct {
	NodeId  uint64
	Nlookup uint64
}

// batch forget is handled internally.
type _BatchForgetIn struct {
	InHeader
	Count uint32
	Dummy uint32
}

type _CuseInitIn struct {
	InHeader
	Major  uint32
	Minor  uint32
	Unused uint32
	Flags  uint32
}

type _CuseInitOut struct {
	Major    uint32
	Minor    uint32
	Unused   uint32
	Flags    uint32
	MaxRead  uint32
	MaxWrite uint32
	DevMajor uint32
	DevMinor uint32
	Spare    [10]uint32
}
