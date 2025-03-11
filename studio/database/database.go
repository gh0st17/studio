package database

import (
	"database/sql"
	"errors"
	"fmt"
	bt "studio/basic_types"
	"studio/errtype"

	_ "github.com/mattn/go-sqlite3"
)

type StudioDB struct {
	entity bt.Entity
	sDB    *sql.DB
}

// Представляет критрии для подстановки в условие
// SQL запроса
type Criteria struct {
	Key          string
	Value        any
	PostOperator string
}

// Загружает локальную базу данных из файла
func (s *StudioDB) LoadDB(fileName string) error {
	var err error
	s.sDB, err = sql.Open("sqlite3", fileName)
	if err != nil {
		return errtype.ErrDataBase(errtype.Join(ErrOpenDB, err))
	}

	return nil
}

// Закрывает базу данных
func (s *StudioDB) CloseDB() error {
	if err := s.sDB.Close(); err != nil {
		return errtype.ErrDataBase(errtype.Join(ErrCloseDB, err))
	}

	return nil
}

func (s StudioDB) Login(login string) (bt.Entity, error) {
	rows, err := s.query(
		"accLevel", "users", "login", []Criteria{{"login", "'" + login + "'", ""}},
	)
	if err != nil {
		return nil, err
	}

	var accLevel uint
	if !rows.Next() {
		return nil, errtype.ErrDataBase(
			errtype.Join(errors.New("неправильный логин"), err),
		)
	} else {
		if err = rows.Scan(&accLevel); err != nil {
			return nil, errtype.ErrDataBase(errtype.Join(ErrReadDB, err))
		}
	}

	table := func() string {
		switch bt.AccessLevel(accLevel) {
		case bt.CUSTOMER:
			return "customers"
		case bt.OPERATOR:
			return "operators"
		default:
			return ""
		}
	}()
	rows.Close()

	if table == "" {
		return &bt.SysAdmin{}, nil
	}

	rows, err = s.query(
		"id, first_name, last_name", table, "id", []Criteria{{"login", "'" + login + "'", ""}},
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var (
		id         uint
		first_name string
		last_name  string
		entity     bt.Entity
	)

	rows.Next()
	if err = rows.Scan(&id, &first_name, &last_name); err != nil {
		return nil, errtype.ErrDataBase(errtype.Join(ErrReadDB, err))
	}

	entity = func() bt.Entity {
		switch bt.AccessLevel(accLevel) {
		case bt.CUSTOMER:
			return &bt.Customer{
				Id:         id,
				First_name: first_name,
				Last_name:  last_name,
			}
		case bt.OPERATOR:
			return &bt.Operator{
				Id:         id,
				First_name: first_name,
				Last_name:  last_name,
			}
		default:
			return nil
		}
	}()

	return entity, nil
}

func (s StudioDB) FetchCustomers() ([]bt.Customer, error) {
	rows, err := s.query(
		"id, first_name, last_name", "customers",
		"first_name, last_name", []Criteria{{"1", "1", ""}},
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

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
