package solofstypes

import "encoding/gob"

func init() {
	gob.Register(FsINodeXAttr{})
}
