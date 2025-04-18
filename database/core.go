package database

import (
	"database/sql"
	"log"
	"reflect"
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
	if tx == nil {
		log.Printf("[db.exec]: %s | args: %v", query, args)
		if _, err = db.sDB.Exec(query, args...); err != nil {
			return errtype.ErrDataBase(errtype.Join(ErrExec, err))
		}
	} else {
		log.Printf("[db.execTx]: %s | args: %v", query, args)
		if _, err = tx.Exec(query, args...); err != nil {
			return errtype.ErrDataBase(errtype.Join(ErrExec, err))
		}
	}

	return nil
}

func (db *StudioDB) queryRow(query string, args []any, tx *sql.Tx) (row *sql.Row) {
	if tx == nil {
		log.Printf("[db.query]: %s | args: %v", query, args)
		return db.sDB.QueryRow(query, args...)
	} else {
		log.Printf("[db.queryTx]: %s | args: %v", query, args)
		return tx.QueryRow(query, args...)
	}
}

func (db *StudioDB) fetchTable(query string, args []any, dest interface{}) error {
	rows, err := db.queryRows(query, args)
	if err != nil {
		return err
	}
	defer rows.Close()

	destSlice := reflect.ValueOf(dest)
	if destSlice.Kind() != reflect.Ptr || destSlice.Elem().Kind() != reflect.Slice {
		return errtype.ErrDataBase(ErrFetchTable)
	}

	elemType := destSlice.Elem().Type().Elem()

	for rows.Next() {
		elem := reflect.New(elemType).Elem()

		numFields := elem.NumField()
		scanDest := make([]interface{}, numFields)

		for i := 0; i < numFields; i++ {
			field := elem.Field(i)
			if field.Kind() == reflect.Uint {
				scanDest[i] = new(sql.NullInt64)
			} else {
				scanDest[i] = field.Addr().Interface()
			}
		}

		if err = rows.Scan(scanDest...); err != nil {
			return errtype.ErrDataBase(errtype.Join(ErrReadDB, err))
		}

		for i := 0; i < numFields; i++ {
			field := elem.Field(i)
			if field.Kind() == reflect.Uint {
				nullInt := scanDest[i].(*sql.NullInt64)
				if nullInt.Valid {
					field.SetUint(uint64(nullInt.Int64))
				} else {
					field.SetUint(0)
				}
			}
		}

		destSlice.Elem().Set(reflect.Append(destSlice.Elem(), elem))
		log.Printf("[db.fetchTable]: Added element %v type of %v\n", elem, elemType)
	}

	return nil
}
