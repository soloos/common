package sdfsapitypes

import "encoding/gob"

func init() {
	gob.Register(FsINodeXAttr{})
}
