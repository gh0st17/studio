package database

import (
	"database/sql"
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
type whereClause struct {
	key          string
	op           string
	value        any
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

// Загружает локальную базу данных из файла
func (db *StudioDB) LoadDB(fileName string) error {
	var err error
	db.sDB, err = sql.Open("sqlite3", fileName)
	if err != nil {
		return errtype.ErrDataBase(errtype.Join(ErrOpenDB, err))
	}

	return nil
}

// Закрывает базу данных
func (db *StudioDB) CloseDB() error {
	if err := db.sDB.Close(); err != nil {
		return errtype.ErrDataBase(errtype.Join(ErrCloseDB, err))
	}

	return nil
}

func (db *StudioDB) Login(login string) (bt.Entity, error) {
	var sp selectParams

	sp = selectParams{
		"accLevel", "users", "",
		[]whereClause{{"login", "=", "'" + login + "'", ""}},
	}

	rows, err := db.query(sp)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var accLevel uint
	if !rows.Next() {
		return nil, errtype.ErrDataBase(
			errtype.Join(ErrLogin, err),
		)
	} else {
		if err = rows.Scan(&accLevel); err != nil {
			return nil, errtype.ErrDataBase(errtype.Join(ErrReadDB, err))
		}
	}

	if bt.AccessLevel(accLevel) == bt.CUSTOMER {
		return db.loginCustomer(login)
	} else {
		return db.loginEmployee(login)
	}
}

func (db *StudioDB) Registration(customer bt.Customer) error {
	ip := insertParams{
		"users",
		"login,accLevel",
		[]string{
			fmt.Sprintf("'%s',1", customer.Login),
		},
	}

	if err := db.insert(ip); err != nil {
		return err
	}

	ip = insertParams{
		"customers",
		"first_name,last_name,login",
		[]string{
			fmt.Sprintf(
				"'%s','%s','%s'", customer.FirstName,
				customer.LastName, customer.Login,
			),
		},
	}

	if err := db.insert(ip); err != nil {
		return err
	}

	return nil
}

func (db *StudioDB) FetchCustomers() (customers []bt.Customer, err error) {
	sp := selectParams{
		"*", "customers", "first_name, last_name", []whereClause{},
	}

	if err = db.fetchTable(sp, &customers); err != nil {
		return nil, err
	}

	return customers, nil
}

func (db *StudioDB) FetchOrdersByCid(cid uint) ([]bt.Order, error) {
	sp := selectParams{
		"*", "orders", "id",
		[]whereClause{{"c_id", "=", fmt.Sprintf("%d", cid), ""}},
	}

	return db.fetchOrders(sp)
}

func (db *StudioDB) FetchOrders() (orders []bt.Order, err error) {
	sp := selectParams{"*", "orders", "id", []whereClause{}}
	return db.fetchOrders(sp)
}

func (db *StudioDB) FetchOrderItems(orders []bt.Order, models []bt.Model) (map[uint][]bt.OrderItem, error) {
	orderItems := make(map[uint][]bt.OrderItem)

	var (
		rawOrderItems []bt.RawOrderItem
		orderItemsArr []bt.OrderItem
	)

	for _, order := range orders {
		sp := selectParams{
			"*", "order_items", "id",
			[]whereClause{{"o_id", "=", fmt.Sprintf("%d", order.Id), ""}},
		}
		if err := db.fetchTable(sp, &rawOrderItems); err != nil {
			return nil, err
		}
		orderItemsArr = db.unrawOrdersItems(rawOrderItems, models)
		orderItems[order.Id] = orderItemsArr
	}

	return orderItems, nil
}

func (db *StudioDB) FetchMaterials() (materials []bt.Material, err error) {
	sp := selectParams{
		"*", "materials", "id",
		[]whereClause{},
	}

	if err = db.fetchTable(sp, &materials); err != nil {
		return nil, err
	}

	return materials, nil
}

func (db *StudioDB) FetchModels() (models []bt.Model, err error) {
	return db.fetchModels()
}
