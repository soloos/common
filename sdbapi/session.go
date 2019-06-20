package sdbapi

import (
	"github.com/gocraft/dbr"
)

type Session struct {
	Dialect string
	dbr.Session
}

func (p *Connection) InitSession(ret *Session) error {
	ret.Session.Connection = p.Connection
	ret.Session.EventReceiver = p.Connection.EventReceiver
	ret.Dialect = p.Dialect
	return nil
}

func (p *Connection) InitSessionWithTx(sess *Session, tx *Tx) error {
	var err error
	err = p.InitSession(sess)
	if err != nil {
		return err
	}

	err = sess.Begin(tx)
	if err != nil {
		return err
	}

	return nil
}

func (p *Session) Begin(ret *Tx) error {
	var err error
	ret.Tx, err = p.Session.Begin()
	ret.Dialect = p.Dialect
	return err
}
