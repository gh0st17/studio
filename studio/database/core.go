package database

import (
	"database/sql"
	"fmt"
	"log"
	"studio/errtype"
	"sync"

	_ "github.com/mattn/go-sqlite3"
)

type StudioDB struct {
	sDB *sql.DB
	mu  sync.Mutex
}

// Представляет критрии для подстановки в условие
// SQL запроса
type whereClause struct {
	key          string
	op           string
	value        string
	postOperator string
}

type selectParams struct {
	cols      string
	table     string
	sortcol   string
	criteries []whereClause
}

type insertParams struct {
	table  string
	cols   string
	values []string
}

type updateParams struct {
	table     string
	set       map[string]string
	criteries []whereClause
}

// Общая функция для запросов в базе данных
func (db *StudioDB) query(sp selectParams) (rows *sql.Rows, err error) {
	query := fmt.Sprintf("SELECT %s FROM %s ", sp.cols, sp.table)

	if len(sp.criteries) > 0 {
		query += "WHERE "
		for _, c := range sp.criteries {
			query += fmt.Sprintf("%s%s%v %s ", c.key, c.op, c.value, c.postOperator)
		}
	}

	if sp.sortcol != "" {
		query += fmt.Sprintf("ORDER BY %s ASC", sp.sortcol)
	}

	log.Printf("[db.query]: %s", query)
	if rows, err = db.sDB.Query(query); err != nil {
		return nil, errtype.ErrDataBase(errtype.Join(ErrQuery, err))
	}

	return rows, nil
}

// Общая функция для вставки в базу данных
func (db *StudioDB) insert(ip insertParams) (err error) {
	db.mu.Lock()
	defer db.mu.Unlock()

	query := fmt.Sprintf("INSERT INTO %s (%s) VALUES ", ip.table, ip.cols)

	if len(ip.values) == 0 {
		return errtype.ErrDataBase(ErrInsert)
	}

	for i, v := range ip.values {
		query += fmt.Sprintf("(%s)", v)

		if len(ip.values) != i+1 {
			query += ","
		}
	}

	log.Printf("[db.insert]: %s", query)
	if _, err = db.sDB.Exec(query); err != nil {
		return errtype.ErrDataBase(errtype.Join(ErrInsert, err))
	}

	return nil
}

// Общая функция для обновления в базе данных
func (db *StudioDB) update(up updateParams) (err error) {
	db.mu.Lock()
	defer db.mu.Unlock()

	query := fmt.Sprintf("UPDATE %s SET ", up.table)

	if len(up.set) == 0 {
		return errtype.ErrDataBase(ErrUpdate)
	}

	var set []string
	for k, v := range up.set {
		set = append(set, fmt.Sprintf("%s=%s", k, v))
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
			query += fmt.Sprintf("%s%s%v %s ", c.key, c.op, c.value, c.postOperator)
		}
	}

	log.Printf("[db.update]: %s", query)
	if _, err = db.sDB.Exec(query); err != nil {
		return errtype.ErrDataBase(errtype.Join(ErrUpdate, err))
	}

	return nil
}
