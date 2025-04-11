package database

import (
	"database/sql"
	"fmt"
	bt "studio/basic_types"
	"studio/errtype"

	_ "github.com/mattn/go-sqlite3"
)

// Загружает локальную базу данных из файла
func (db *StudioDB) LoadDB(fileName string) error {
	var err error
	db.sDB, err = sql.Open("sqlite3", fileName)
	if err != nil {
		return errtype.ErrDataBase(errtype.Join(ErrOpenDB, err))
	}

	db.sDB.Exec("PRAGMA journal_mode=WAL;")

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
	sp := selectParams{
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
		return nil, errtype.ErrDataBase(ErrLogin)
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

func (db *StudioDB) FetchOrders(cid uint) (orders []bt.Order, err error) {
	sp := selectParams{
		"*", "orders", "id",
		[]whereClause{},
	}

	if cid > 0 {
		sp.criteries = []whereClause{{"c_id", "=", fmt.Sprint(cid), ""}}
	}

	if err = db.fetchTable(sp, &orders); err != nil {
		return nil, err
	}

	return orders, nil
}

func (db *StudioDB) FetchOrderItems(o_id uint, models []bt.Model) ([]bt.OrderItem, error) {
	type RawOrderItem struct {
		Id, O_id, Model uint
		UnitPrice       float64
	}

	var orderItems []bt.OrderItem

	sp := selectParams{
		"*", "order_items", "id",
		[]whereClause{{"o_id", "=", fmt.Sprint(o_id), ""}},
	}
	var rawOrderItems []RawOrderItem
	if err := db.fetchTable(sp, &rawOrderItems); err != nil {
		return nil, err
	}

	for _, rawOrderItem := range rawOrderItems {
		orderItems = append(orderItems,
			bt.OrderItem{
				Id:        rawOrderItem.Id,
				O_id:      rawOrderItem.O_id,
				Model:     models[rawOrderItem.Model-1],
				UnitPrice: rawOrderItem.UnitPrice,
			},
		)
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

func (db *StudioDB) CreateOrder(cid uint, models []bt.Model) (err error) {
	var (
		ip         insertParams
		orderPrice float64
		order_id   uint
		tx         *sql.Tx
	)

	for _, m := range models {
		orderPrice += m.Price
	}

	ip = insertParams{
		"orders",
		"c_id, total_price",
		[]string{
			fmt.Sprintf(
				"%d,%f", cid, orderPrice,
			),
		},
	}

	if tx, err = db.sDB.Begin(); err != nil {
		return errtype.ErrDataBase(errtype.Join(ErrBegin, err))
	}

	if err = db.insert(ip); err != nil {
		tx.Rollback()
		return err
	}

	if order_id, err = db.getLastId(
		"orders",
		[]whereClause{{
			"c_id", "=",
			fmt.Sprint(cid), "",
		}},
	); err != nil {
		tx.Rollback()
		return err
	}

	for _, m := range models {
		ip = insertParams{
			"order_items",
			"o_id, model, unit_price",
			[]string{
				fmt.Sprintf(
					"%d,%d,%f", order_id, m.Id, m.Price,
				),
			},
		}

		if err = db.insert(ip); err != nil {
			tx.Rollback()
			return err
		}
	}

	tx.Commit()
	return nil
}

func (db *StudioDB) SetOrderStatus(id uint, newStatus bt.OrderStatus) error {
	sp := selectParams{
		"status", "orders", "",
		[]whereClause{{"id", "=", fmt.Sprint(id), ""}},
	}
	var orderStatus bt.OrderStatus
	if rows, err := db.query(sp); err != nil {
		return err
	} else {
		rows.Next()
		rows.Scan(&orderStatus)
		rows.Close()
	}

	if newStatus == bt.Canceled && orderStatus > 1 {
		return ErrNotPending
	} else if newStatus != bt.Canceled && newStatus-orderStatus != 1 {
		return ErrStatusRange
	}

	up := updateParams{
		"orders",
		map[string]string{"status": fmt.Sprint(int(newStatus))},
		[]whereClause{{"id", "=", fmt.Sprint(id), ""}},
	}

	if err := db.update(up); err != nil {
		return err
	}

	return nil
}

func (db *StudioDB) FetchFullName(id uint, accessLevel bt.AccessLevel) (name string) {
	table := func() string {
		if accessLevel == bt.OPERATOR {
			return "employees"
		} else {
			return "customers"
		}
	}()

	sp := selectParams{
		"first_name || ' ' || last_name AS full_name", table, "",
		[]whereClause{{"id", "=", fmt.Sprint(id), ""}},
	}

	rows, _ := db.query(sp)
	rows.Next()
	rows.Scan(&name)

	return name
}

func (db *StudioDB) SetOperator(eId, oId uint) error {
	up := updateParams{
		"orders",
		map[string]string{"e_id": fmt.Sprint(eId)},
		[]whereClause{{"id", "=", fmt.Sprint(oId), ""}},
	}

	if err := db.update(up); err != nil {
		return err
	}

	return nil
}
