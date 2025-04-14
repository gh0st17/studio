package database

import (
	"database/sql"
	"fmt"
	"log"
	"strings"
	"studio/errtype"

	_ "github.com/lib/pq"
)

type StudioDB struct {
	sDB *sql.DB
}

// Представляет критрии для подстановки в условие
// SQL запроса
type whereClause struct {
	key          string
	op           string
	value        string
	postOperator string
}

type joinClause struct {
	jType string
	table string
	on    string
}

type selectParams struct {
	cols      string
	table     string
	sortcol   string
	joins     []joinClause
	criteries []whereClause
}

type insertParams struct {
	table     string
	cols      string
	values    []string
	returning string
}

type updateParams struct {
	table     string
	set       map[string]string
	criteries []whereClause
}

// Общая функция для запросов в базе данных
func (db *StudioDB) query(sp selectParams) (rows *sql.Rows, err error) {
	query := fmt.Sprintf("SELECT %s FROM %s ", sp.cols, sp.table)

	if len(sp.joins) > 0 {
		for _, j := range sp.joins {
			query += fmt.Sprintf("%s %s %s ", j.jType, j.table, j.on)
		}
	}

	var args []interface{}
	if len(sp.criteries) > 0 {
		query += "WHERE "
		for i, c := range sp.criteries {
			query += fmt.Sprintf("%s%s$%d %s ", c.key, c.op, i+1, c.postOperator)
			args = append(args, c.value)
		}
	}

	if sp.sortcol != "" {
		query += fmt.Sprintf("ORDER BY %s ASC", sp.sortcol)
	}

	log.Printf("[db.query]: %s | args: %v", query, args)
	if rows, err = db.sDB.Query(query, args...); err != nil {
		return nil, errtype.ErrDataBase(errtype.Join(ErrQuery, err))
	}

	return rows, nil
}

// Общая функция для вставки в базу данных
func (db *StudioDB) insert(ip insertParams, tx *sql.Tx) (row *sql.Row, err error) {
	query := fmt.Sprintf("INSERT INTO %s (%s) VALUES ", ip.table, ip.cols)

	if len(ip.values) == 0 {
		return nil, errtype.ErrDataBase(ErrInsert)
	}

	var (
		args []interface{}
		i    int = 1
	)
	for k, v := range ip.values {
		splittedVals := strings.Split(v, ",")

		query += "("
		for j, sVal := range splittedVals {
			query += fmt.Sprintf("$%d", i)
			args = append(args, sVal)
			i++

			if len(splittedVals) != j+1 {
				query += ","
			}
		}
		query += ")"

		if len(ip.values) != k+1 {
			query += ","
		}
	}

	if ip.returning != "" {
		query += " RETURNING " + ip.returning
		log.Printf("[db.insert]: %s | args: %v", query, args)
		row = db.queryInsert(tx, query, args)
		return row, row.Err()
	} else {
		log.Printf("[db.insert]: %s | args: %v", query, args)
		err = db.execInsert(tx, query, args)
		return nil, err
	}
}

func (db *StudioDB) queryInsert(tx *sql.Tx, query string, args []interface{}) *sql.Row {
	if tx == nil {
		return db.sDB.QueryRow(query, args...)
	} else {
		return tx.QueryRow(query, args...)
	}
}

func (db *StudioDB) execInsert(tx *sql.Tx, query string, args []interface{}) (err error) {
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

// Общая функция для обновления в базе данных
func (db *StudioDB) update(up updateParams) (err error) {
	query := fmt.Sprintf("UPDATE %s SET ", up.table)

	if len(up.set) == 0 {
		return errtype.ErrDataBase(ErrUpdate)
	}

	var (
		set  []string
		args []interface{}
		i    int = 1
	)

	for k, v := range up.set {
		set = append(set, fmt.Sprintf("%s=$%d", k, i))
		args = append(args, v)
		i++
	}

	for i, s := range set {
		query += s
		if len(set) != i+1 {
			query += ", "
		} else {
			query += " "
		}
	}

	if len(up.criteries) > 0 {
		query += "WHERE "
		for _, c := range up.criteries {
			query += fmt.Sprintf("%s%s$%d %s ", c.key, c.op, i, c.postOperator)
			args = append(args, c.value)
		}
	}

	log.Printf("[db.update]: %s | args: %v", query, args)
	if _, err = db.sDB.Exec(query, args...); err != nil {
		return errtype.ErrDataBase(errtype.Join(ErrUpdate, err))
	}

	return nil
}
