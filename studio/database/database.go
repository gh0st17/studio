package database

import (
	"database/sql"
	"fmt"
	bt "studio/basic_types"
	"studio/errtype"

	_ "github.com/mattn/go-sqlite3"
)

type StudioDB struct {
	sDB *sql.DB
}

// Представляет критрии для подстановки в условие
// SQL запроса
type Criteria struct {
	Key          string
	Value        any
	PostOperator string
}

func (s StudioDB) FetchCustomers() ([]bt.Customer, error) {
	rows, err := s.query(
		"*", "customers", "first_name, last_name", []Criteria{},
	)
	if err != nil {
		return nil, err
	}

	var (
		id         uint
		first_name string
		last_name  string
		customer   bt.Customer
		customers  []bt.Customer
	)

	for rows.Next() {
		err := rows.Scan(&id, &first_name, &last_name)
		if err != nil {
			return nil, errtype.ErrDataBase(errtype.Join(ErrReadDB, err))
		}

		customer = bt.Customer{
			Id:         id,
			First_name: first_name,
			Last_name:  last_name,
		}

		customers = append(customers, customer)
	}

	return customers, nil
}

func (StudioDB) FetchOrdersByCustId(cid uint) ([]bt.Order, error) {
	return nil, nil
}

func (StudioDB) FetchOrders() ([]bt.Order, error) {
	return nil, nil
}

func (StudioDB) FetchMaterials() ([]bt.Material, error) {
	return nil, nil
}

func (StudioDB) FetchModels() ([]bt.Model, error) {
	return nil, nil
}

// Общая функция для запросов в базе данных
func (db *StudioDB) query(cols string, table string, sortcol string, criteries []Criteria) (*sql.Rows, error) {
	var (
		err   error
		query string
		rows  *sql.Rows
	)

	query = fmt.Sprintf("SELECT %s FROM %s WHERE ", cols, table)
	for _, c := range criteries {
		query += fmt.Sprintf("%s=%v %s ", c.Key, c.Value, c.PostOperator)
	}

	if sortcol != "" {
		query += fmt.Sprintf("ORDER BY %s ASC", sortcol)
	}

	if rows, err = db.sDB.Query(query); err != nil {
		return nil, errtype.ErrDataBase(errtype.Join(ErrQuery, err))
	}

	return rows, nil
}
