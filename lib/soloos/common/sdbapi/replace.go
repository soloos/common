package sdbapi

import (
	"bytes"
	"fmt"
	"strings"
)

type ReplaceStmt struct {
	dialect       string
	tx            *Tx
	table         string
	primaryColumn string
	primaryValue  interface{}
	column        []string
	value         []interface{}
}

func (p *Tx) ReplaceInto(table string) *ReplaceStmt {
	return &ReplaceStmt{
		dialect: p.Dialect,
		tx:      p,
		table:   table,
	}
}

func (p *ReplaceStmt) PrimaryColumn(column string) *ReplaceStmt {
	p.primaryColumn = column
	return p
}

func (p *ReplaceStmt) PrimaryValue(value interface{}) *ReplaceStmt {
	p.primaryValue = value
	return p
}

func (p *ReplaceStmt) Columns(column ...string) *ReplaceStmt {
	p.column = column
	return p
}

func (p *ReplaceStmt) Values(value ...interface{}) *ReplaceStmt {
	p.value = value
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
	for i = 1; i < 1+len(p.column); i++ {
		insertColumnsPlaceHolders.WriteString(",?")
	}

	updateColumnsPlaceHolders.WriteString(fmt.Sprintf("%s = ?", p.column[0]))
	for i = 1; i < len(p.column); i++ {
		updateColumnsPlaceHolders.WriteString(fmt.Sprintf(", %s = ?", p.column[i]))
	}

	tpl.WriteString(
		fmt.Sprintf("INSERT INTO %s (%s, %s) VALUES (%s) ", p.table,
			p.primaryColumn, strings.Join(p.column, ","),
			insertColumnsPlaceHolders.String()))

	if p.dialect == "sqlite3" {
		tpl.WriteString(fmt.Sprintf(" ON CONFLICT(%s) DO UPDATE SET %s;",
			p.primaryColumn,
			updateColumnsPlaceHolders.String()))
		updateValues = append(updateValues, p.primaryValue)
		updateValues = append(updateValues, p.value...)
		updateValues = append(updateValues, p.value...)
	} else {
		tpl.WriteString(fmt.Sprintf(" ON DUPLICATE KEY UPDATE %s;",
			updateColumnsPlaceHolders.String()))
		updateValues = append(updateValues, p.primaryValue)
		updateValues = append(updateValues, p.value...)
		updateValues = append(updateValues, p.value...)
	}

	_, err = p.tx.UpdateBySql(tpl.String(), updateValues...).Exec()
	return err
}
