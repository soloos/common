package util

import (
	"crypto/sha256"
	"fmt"

	uuid "github.com/satori/go.uuid"
)

func InitUUID64(d *[64]byte) {
	u2 := uuid.NewV4()

	copy((*d)[:], []byte(fmt.Sprintf("%x", sha256.Sum256(u2[:]))))
}
