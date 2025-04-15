package database

import (
	"database/sql"
	"log"
	"studio/errtype"

	_ "github.com/lib/pq"
)

type StudioDB struct {
	sDB *sql.DB
}

func (db *StudioDB) queryRows(query string, args []any) (rows *sql.Rows, err error) {
	log.Printf("[db.queryRows]: %s | args: %v", query, args)
	if rows, err = db.sDB.Query(query, args...); err != nil {
		return nil, errtype.ErrDataBase(errtype.Join(ErrQuery, err))
	}

	return rows, nil
}

func (db *StudioDB) exec(query string, args []any, tx *sql.Tx) (err error) {
	log.Printf("[db.execTx]: %s | args: %v", query, args)
	if tx == nil {
		if _, err = db.sDB.Exec(query, args...); err != nil {
			return errtype.ErrDataBase(errtype.Join(ErrInsert, err))
		}
	} else {
		if _, err = tx.Exec(query, args...); err != nil {
			return errtype.ErrDataBase(errtype.Join(ErrInsert, err))
		}
	}

	return nil
}

func (db *StudioDB) queryRow(query string, args []any, tx *sql.Tx) (row *sql.Row) {
	log.Printf("[db.queryTx]: %s | args: %v", query, args)
	if tx == nil {
		return db.sDB.QueryRow(query, args...)
	} else {
		return tx.QueryRow(query, args...)
	}
}
