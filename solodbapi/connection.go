package solodbapi

import (
	"github.com/gocraft/dbr"
)

type Connection struct {
	*dbr.Connection
	DBDriver string
	Dialect  string
	Dsn      string
}

func (p *Connection) Init(dbDriver, dsn string) error {
	var err error
	p.DBDriver = dbDriver
	p.Dialect = dbDriver
	p.Dsn = dsn
	p.Connection, err = dbr.Open(p.DBDriver, p.Dsn, nil)
	if err != nil {
		return err
	}

	return nil
}
