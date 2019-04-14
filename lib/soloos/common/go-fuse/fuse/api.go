// Copyright 2016 the Go-FUSE Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package fuse provides APIs to implement filesystems in
// userspace in terms of raw FUSE protocol.
//
// A filesystem is implemented by implementing its server that provides a
// PosixFS interface. Typically the server embeds
// NewDefaultPosixFS() and implements only subset of filesystem methods:
//
//	type MyFS struct {
//		fuse.PosixFS
//		...
//	}
//
//	func NewMyFS() *MyFS {
//		return &MyFS{
//			PosixFS: fuse.NewDefaultPosixFS(),
//			...
//		}
//	}
//
//	// Mkdir implements "mkdir" request handler.
//	//
//	// For other requests - not explicitly implemented by MyFS - ENOSYS
//	// will be typically returned to client.
//	func (fs *MyFS) Mkdir(...) {
//		...
//	}
//
// Then the filesystem can be mounted and served to a client (typically OS
// kernel) by creating Server:
//
//	fs := NewMyFS() // implements PosixFS
//	fssrv, err := fuse.NewServer(fs, mountpoint, &fuse.MountOptions{...})
//	if err != nil {
//		...
//	}
//
// and letting the server do its work:
//
//	// either synchronously - .Serve() blocks until the filesystem is unmounted.
//	fssrv.Serve()
//
//	// or in the background - .Serve() is spawned in another goroutine, but
//	// before interacting with fssrv from current context we have to wait
//	// until the filesystem mounting is complete.
//	go fssrv.Serve()
//	err = fssrv.WaitMount()
//	if err != nil {
//		...
//	}
//
// The server will serve clients by dispatching their requests to the
// filesystem implementation and conveying responses back. For example "mkdir"
// FUSE request dispatches to call
//
//	fs.Mkdir(*MkdirIn, ..., *EntryOut)
//
// "stat" to call
//
//	fs.GetAttr(*GetAttrIn, *AttrOut)
//
// etc. Please refer to PosixFS documentation for details.
//
// Typically, each call of the API happens in its own
// goroutine, so take care to make the file system thread-safe.
//
//
// Higher level interfaces
//
// As said above this packages provides way to implement filesystems in terms of
// raw FUSE protocol. Additionally packages nodefs and pathfs provide ways to
// implement filesystem at higher levels:
//
// Package github.com/hanwen/go-fuse/fuse/nodefs provides way to implement
// filesystems in terms of inodes. This resembles kernel's idea of what a
// filesystem looks like.
//
// Package github.com/hanwen/go-fuse/fuse/pathfs provides way to implement
// filesystems in terms of path names. Working with path names is somewhat
// easier compared to inodes, however renames can be racy. Do not use pathfs if
// you care about correctness.
package fuse

// Types for users to implement.

type MountOptions struct {
	AllowOther bool

	// Options are passed as -o string to fusermount.
	Options []string

	// Default is _DEFAULT_BACKGROUND_TASKS, 12.  This numbers
	// controls the allowed number of requests that relate to
	// async I/O.  Concurrency for synchronous I/O is not limited.
	MaxBackground int

	// Write size to use.  If 0, use default. This number is
	// capped at the kernel maximum.
	MaxWrite int

	// Max read ahead to use.  If 0, use default. This number is
	// capped at the kernel maximum.
	MaxReadAhead int

	// If IgnoreSecurityLabels is set, all security related xattr
	// requests will return NO_DATA without passing through the
	// user defined filesystem.  You should only set this if you
	// file system implements extended attributes, and you are not
	// interested in security labels.
	IgnoreSecurityLabels bool // ignoring labels should be provided as a fusermount mount option.

	// If given, use this buffer pool instead of the global one.
	Buffers BufferPool

	// If RememberInodes is set, we will never forget inodes.
	// This may be useful for NFS.
	RememberInodes bool

	// Values shown in "df -T" and friends
	// First column, "Filesystem"
	FsName string

	// Second column, "Type", will be shown as "fuse." + Name
	Name string

	// If set, return ENOSYS for Getxattr calls, so the kernel does not issue any
	// Xattr operations at all.
	DisableXAttrs bool

	// If set, print debugging information.
	Debug bool

	// If set, ask kernel to forward file locks to FUSE. If using,
	// you must implement the GetLk/SetLk/SetLkw methods.
	EnableLocks bool
}
