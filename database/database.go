package database

import (
	"database/sql"
	"fmt"
	bt "studio/basic_types"
	"studio/errtype"

	_ "github.com/lib/pq"
)

// Выполняет подключение к базе данных
func (db *StudioDB) LoadDB() error {
	var err error
	connStr := "host=localhost port=5432 user=studio "
	connStr += "password=studio dbname=studio "
	connStr += "sslmode=disable search_path=studio"
	db.sDB, err = sql.Open("postgres", connStr)
	if err != nil {
		return errtype.ErrDataBase(errtype.Join(ErrOpenDB, err))
	}

	return nil
}

// Закрывает подключение к базе данных
func (db *StudioDB) CloseDB() error {
	if err := db.sDB.Close(); err != nil {
		return errtype.ErrDataBase(errtype.Join(ErrCloseDB, err))
	}

	return nil
}

func (db *StudioDB) Login(login string) (bt.Entity, error) {
	sp := selectParams{
		"access_level", "users", "", []joinClause{},
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
		return db.login(login, "customers", &[]bt.Customer{})
	} else {
		return db.login(login, "employees", &[]bt.Employee{})
	}
}

func (db *StudioDB) Registration(customer bt.Customer) error {
	ip := insertParams{
		"users", "login,access_level",
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
		"*", "customers", "first_name, last_name", []joinClause{}, []whereClause{},
	}

	if err = db.fetchTable(sp, &customers); err != nil {
		return nil, err
	}

	return customers, nil
}

func (db *StudioDB) FetchOrders(cid uint) (orders []bt.Order, err error) {
	cols := "o.id, c.first_name || ' ' || c.last_name AS customer_name, "
	cols += "COALESCE(e.first_name || ' ' || e.last_name, '') AS employee_name, "
	cols += "o.status, (SELECT SUM(unit_price) FROM order_items WHERE o_id = o.id) AS total_price, "
	cols += "o.create_date, o.release_date"

	sp := selectParams{
		cols, "orders o", "o.id",
		[]joinClause{
			{"LEFT JOIN", "customers c", "ON o.c_id = c.id"},
			{"LEFT JOIN", "employees e", "ON o.e_id = e.id"},
		},
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

func (db *StudioDB) FetchOrderItems(o_id uint, models map[uint]bt.Model) ([]bt.OrderItem, error) {
	type RawOrderItem struct {
		Id, O_id, Model uint
		UnitPrice       float64
	}

	var orderItems []bt.OrderItem

	sp := selectParams{
		"*", "order_items", "id", []joinClause{},
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
				Model:     models[rawOrderItem.Model],
				UnitPrice: rawOrderItem.UnitPrice,
			},
		)
	}

	return orderItems, nil
}

func (db *StudioDB) FetchMaterials() (materials map[uint]bt.Material, err error) {
	sp := selectParams{
		"*", "materials", "id", []joinClause{}, []whereClause{},
	}

	var matSlice []bt.Material
	if err = db.fetchTable(sp, &matSlice); err != nil {
		return nil, err
	}

	materials = make(map[uint]bt.Material)
	for _, m := range matSlice {
		materials[m.Id] = m
	}

	return materials, nil
}

func (db *StudioDB) FetchModels() (models map[uint]bt.Model, err error) {
	return db.fetchModels()
}

func (db *StudioDB) CreateOrder(cid uint, models []bt.Model) (err error) {
	var (
		ip       insertParams
		order_id uint
		tx       *sql.Tx
	)

	ip = insertParams{
		"orders", "c_id",
		[]string{
			fmt.Sprintf(
				"%d", cid,
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
			"order_items", "o_id, model, unit_price",
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
		"status", "orders", "", []joinClause{},
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
