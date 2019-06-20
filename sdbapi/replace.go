package sdbapi

import (
	"bytes"
	"fmt"
	"strings"

	"github.com/gocraft/dbr"
)

type ReplaceStmt struct {
	dialect        string
	runner         dbr.SessionRunner
	table          string
	primaryColumns []string
	primaryValues  []interface{}
	columns        []string
	values         []interface{}
}

func (p *Tx) ReplaceInto(table string) *ReplaceStmt {
	return &ReplaceStmt{
		dialect: p.Dialect,
		runner:  p,
		table:   table,
	}
}

func (p *Session) ReplaceInto(table string) *ReplaceStmt {
	return &ReplaceStmt{
		dialect: p.Dialect,
		runner:  p,
		table:   table,
	}
}

func (p *ReplaceStmt) PrimaryColumns(columns ...string) *ReplaceStmt {
	p.primaryColumns = columns
	return p
}

func (p *ReplaceStmt) PrimaryValues(values ...interface{}) *ReplaceStmt {
	p.primaryValues = values
	return p
}

func (p *ReplaceStmt) Columns(columns ...string) *ReplaceStmt {
	p.columns = columns
	return p
}

func (p *ReplaceStmt) Values(values ...interface{}) *ReplaceStmt {
	p.values = values
	return p
}

func (p *ReplaceStmt) Exec() error {
	var (
		tpl                       bytes.Buffer
		insertColumnsPlaceHolders bytes.Buffer
		updateColumnsPlaceHolders bytes.Buffer
		updateValues              []interface{}
		i                         int
		err                       error
	)

	insertColumnsPlaceHolders.WriteString("?")
	for i = 1; i < len(p.primaryColumns)+len(p.columns); i++ {
		insertColumnsPlaceHolders.WriteString(",?")
	}

	updateColumnsPlaceHolders.WriteString(fmt.Sprintf("%s = ?", p.columns[0]))
	for i = 1; i < len(p.columns); i++ {
		updateColumnsPlaceHolders.WriteString(fmt.Sprintf(", %s = ?", p.columns[i]))
	}

	tpl.WriteString(
		fmt.Sprintf("INSERT INTO %s (%s, %s) VALUES (%s) ", p.table,
			strings.Join(p.primaryColumns, ","), strings.Join(p.columns, ","),
			insertColumnsPlaceHolders.String()))

	if p.dialect == "sqlite3" {
		tpl.WriteString(fmt.Sprintf(" ON CONFLICT(%s) DO UPDATE SET %s;",
			strings.Join(p.primaryColumns, ","),
			updateColumnsPlaceHolders.String()))
		updateValues = append(updateValues, p.primaryValues...)
		updateValues = append(updateValues, p.values...)
		updateValues = append(updateValues, p.values...)
	} else {
		tpl.WriteString(fmt.Sprintf(" ON DUPLICATE KEY UPDATE %s;",
			updateColumnsPlaceHolders.String()))
		updateValues = append(updateValues, p.primaryValues...)
		updateValues = append(updateValues, p.values...)
		updateValues = append(updateValues, p.values...)
	}

	_, err = p.runner.UpdateBySql(tpl.String(), updateValues...).Exec()
	return err
}
