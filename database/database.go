package database

import (
	"database/sql"
	"strings"

	bt "github.com/gh0st17/studio/basic_types"
	"github.com/gh0st17/studio/errtype"

	_ "github.com/lib/pq"
)

// Выполняет подключение к базе данных
func (db *StudioDB) LoadDB(socket string) (err error) {
	var soc string

	socParts := strings.Split(socket, ":")
	if len(socParts) != 2 {
		return errtype.ErrDataBase(errtype.Join(ErrOpenDB, err))
	}

	soc += "host=" + socParts[0] + " port=" + socParts[1] + " "
	db.sDB, err = sql.Open("postgres", soc+connStr)
	if err != nil {
		return errtype.ErrDataBase(errtype.Join(ErrOpenDB, err))
	}

	if err = db.sDB.Ping(); err != nil {
		return errtype.ErrDataBase(errtype.Join(ErrPingDB, err))
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
	rows, err := db.queryRows(loginQuery, []any{login})
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
		return db.login(fetchCustLoginQuery, []any{login}, &[]bt.Customer{})
	} else {
		return db.login(fetchEmplLoginQuery, []any{login}, &[]bt.Employee{})
	}
}

func (db *StudioDB) Registration(c bt.Customer) (err error) {
	var tx *sql.Tx

	if tx, err = db.sDB.Begin(); err != nil {
		return errtype.ErrDataBase(errtype.Join(ErrBegin, err))
	}

	if err = db.exec(insertUserQuery, []any{c.Login, 1}, tx); err != nil {
		tx.Rollback()
		return err
	}

	args := []any{c.FirstName, c.LastName, c.Login}
	if err = db.exec(insertCustQuery, args, tx); err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit()
}

func (db *StudioDB) FetchOrders(cid uint) (orders []bt.Order, err error) {
	var (
		query string
		args  []any
	)
	if cid > 0 {
		query = fetchOrdersQueryCid
		args = append(args, cid)
	} else {
		query = fetchOrdersQueryAll
	}

	if err = db.fetchTable(query, args, &orders); err != nil {
		return nil, err
	}

	return orders, nil
}

func (db *StudioDB) FetchOrderItems(o_id uint, models map[uint]bt.Model) ([]bt.OrderItem, error) {
	type RawOrderItem struct {
		Id, O_id, Model uint
		UnitPrice       float64
	}

	var rawOrderItems []RawOrderItem
	if err := db.fetchTable(fetchOrderItemsQuery, []any{o_id}, &rawOrderItems); err != nil {
		return nil, err
	}

	orderItems := make([]bt.OrderItem, len(rawOrderItems))
	for i, rawOrderItem := range rawOrderItems {
		orderItems[i] = bt.OrderItem{
			Id:        rawOrderItem.Id,
			O_id:      rawOrderItem.O_id,
			Model:     models[rawOrderItem.Model],
			UnitPrice: rawOrderItem.UnitPrice,
		}
	}

	return orderItems, nil
}

func (db *StudioDB) FetchMaterials() (materials map[uint]bt.Material, err error) {
	var matSlice []bt.Material
	if err = db.fetchTable(fetchMatQuery, []any{}, &matSlice); err != nil {
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
		order_id uint
		tx       *sql.Tx
	)

	if tx, err = db.sDB.Begin(); err != nil {
		return errtype.ErrDataBase(errtype.Join(ErrBegin, err))
	}

	if row := db.queryRow(insertOrderQuery, []any{cid}, tx); row.Err() != nil {
		tx.Rollback()
		return row.Err()
	} else {
		row.Scan(&order_id)
	}

	for _, m := range models {
		err = db.exec(
			insertOrderItemsQuery,
			[]any{order_id, m.Id, m.Price}, tx,
		)
		if err != nil {
			tx.Rollback()
			return err
		}
	}

	return tx.Commit()
}

func (db *StudioDB) SetOrderStatus(id uint, newStatus bt.OrderStatus) (err error) {
	var (
		tx          *sql.Tx
		orderStatus bt.OrderStatus
	)

	if tx, err = db.sDB.Begin(); err != nil {
		return errtype.ErrDataBase(errtype.Join(ErrBegin, err))
	}

	if row := db.queryRow(fetchOrderStatus, []any{id}, nil); row.Err() != nil {
		return row.Err()
	} else {
		row.Scan(&orderStatus)
	}

	if newStatus == bt.Canceled && orderStatus > 1 {
		return ErrNotPending
	} else if newStatus != bt.Canceled && newStatus-orderStatus != 1 {
		return ErrStatusRange
	}

	if err := db.exec(updateStatusQuery, []any{newStatus, id}, tx); err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit()
}

func (db *StudioDB) SetOperator(eId, oId uint) (err error) {
	var tx *sql.Tx

	if tx, err = db.sDB.Begin(); err != nil {
		return errtype.ErrDataBase(errtype.Join(ErrBegin, err))
	}

	if err := db.exec(updateOrderEmplQuery, []any{eId, oId}, tx); err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit()
}
