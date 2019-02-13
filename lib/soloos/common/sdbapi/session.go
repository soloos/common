package sdbapi

import (
	"github.com/gocraft/dbr"
)

type Session struct {
	dbr.Session
}

func (p *Connection) InitSession(ret *Session) error {
	ret.Session.Connection = p.Connection
	ret.Session.EventReceiver = p.Connection.EventReceiver
	return nil
}

func (p *Session) Begin() (*Tx, error) {
	var (
		ret = new(Tx)
		err error
	)
	ret.Tx, err = p.Session.Begin()
	return ret, err
}
