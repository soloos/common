package solofstypes

import "encoding/gob"

func init() {
	gob.Register(FsINodeMeta{})
	gob.Register([]FsINodeMeta{})
	gob.Register(NetINodeMeta{})
	gob.Register([]NetINodeMeta{})
	gob.Register(NetBlockMeta{})
	gob.Register([]NetBlockMeta{})
}
