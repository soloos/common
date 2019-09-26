package solofsapitypes

import "encoding/gob"

func init() {
	gob.Register(FsINodeXAttr{})
}
