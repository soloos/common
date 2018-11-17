package util

import (
	"crypto/sha256"
	"fmt"

	uuid "github.com/satori/go.uuid"
)

func InitUUID64(d *[64]byte) {
	u2, err := uuid.NewV4()
	AssertErrIsNil(err)

	h := sha256.New()
	h.Write(u2[:])
	b := h.Sum(nil)
	copy((*d)[:], []byte(fmt.Sprintf("%x", b)))
}
