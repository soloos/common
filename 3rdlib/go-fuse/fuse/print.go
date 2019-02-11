package fuse

import "fmt"

func (me *_BatchForgetIn) string() string {
	return fmt.Sprintf("{Count=%d}", me.Count)
}
