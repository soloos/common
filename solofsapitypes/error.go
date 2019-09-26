package solofsapitypes

import (
	"soloos/common/xerrors"
)

var (
	ErrRemoteService      = xerrors.New("remote service error")
	ErrServiceNotExists   = xerrors.New("service not exists")
	ErrObjectExists       = xerrors.New("object exists")
	ErrObjectHasChildren  = xerrors.New("object has children")
	ErrNetINodeNotExists  = xerrors.New("netinode not exists")
	ErrObjectNotExists    = xerrors.New("object not exists")
	ErrObjectNotPrepared  = xerrors.New("object not prepared")
	ErrNetBlockPWrite     = xerrors.New("net block pwrite error")
	ErrNetBlockPRead      = xerrors.New("net block pread error")
	ErrBackendListIsEmpty = xerrors.New("backend list is empty")
	ErrRetryTooManyTimes  = xerrors.New("retry too many times")
	ErrRLockFailed        = xerrors.New("rlock failed")
	ErrLockFailed         = xerrors.New("lock failed")
	ErrInvalidArgs        = xerrors.New("invalid args")
	ErrHasNotPermission   = xerrors.New("has not permission")
)
