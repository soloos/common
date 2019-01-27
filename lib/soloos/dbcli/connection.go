package dbcli

import (
	"fmt"

	"github.com/gocraft/dbr"
)

type Connection struct {
	*dbr.Connection
	DBDriver string
	Dsn      string
}

func (p *Connection) Init(dbDriver, dsn string) error {
	var err error
	p.DBDriver = dbDriver
	p.Dsn = dsn
	p.Connection, err = dbr.Open(p.DBDriver, p.Dsn, nil)
	if err != nil {
		return err
	}

	return nil
}

func (p *Connection) ReplaceInto(
	tx *dbr.Tx,
	table string,
	primaryColKey string,
	updateColKey string,
	primaryColValue interface{},
	updateColValue interface{},
) error {
	var (
		tpl string
		err error
	)

	if p.DBDriver == "sqlite3" {
		tpl = fmt.Sprintf(`
		INSERT INTO %s (%s, %s)
		VALUES (?,?)
		ON CONFLICT(%s) DO UPDATE SET %s = ?;
		`, table, primaryColKey, updateColKey,
			primaryColKey, updateColKey)
		_, err = tx.UpdateBySql(tpl, primaryColValue, updateColValue, updateColValue).Exec()
	} else {
		tpl = fmt.Sprintf(`
		INSERT INTO %s (%s,%s) VALUES (?,?)
		ON DUPLICATE KEY UPDATE %s=?;
		`, table, primaryColKey, updateColKey,
			updateColKey)
		_, err = tx.UpdateBySql(tpl, primaryColValue, updateColValue, updateColValue).Exec()
	}

	if err != nil {
		goto QUERY_DONE
	}

QUERY_DONE:
	return err
}
