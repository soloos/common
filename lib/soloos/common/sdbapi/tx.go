package sdbapi

import "github.com/gocraft/dbr"

type Tx struct {
	Dialect string
	*dbr.Tx
}
