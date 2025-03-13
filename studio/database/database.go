package database

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
	"reflect"
	bt "studio/basic_types"
	"studio/errtype"
	"time"

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

func (db *StudioDB) loginCustomer(login string) (bt.Entity, error) {
	cols, table := "id, first_name, last_name", "customers"
	sp := selectParams{
		cols, table, "",
		[]whereClause{{"login", "=", "'" + login + "'", ""}},
	}
	rows, err := db.query(sp)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	ent := &bt.Customer{}
	rows.Next()
	if err = rows.Scan(&ent.Id, &ent.FirstName, &ent.LastName); err != nil {
		return nil, errtype.ErrDataBase(errtype.Join(ErrReadDB, err))
	}

	return ent, nil
}

func (db *StudioDB) loginEmployee(login string) (bt.Entity, error) {
	cols, table := "id, first_name, last_name, job_id", "employees"
	sp := selectParams{
		cols, table, "",
		[]whereClause{{"login", "=", "'" + login + "'", ""}},
	}
	rows, err := db.query(sp)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	ent := &bt.Employee{}
	rows.Next()
	if err = rows.Scan(&ent.Id, &ent.FirstName, &ent.LastName, &ent.JobId); err != nil {
		return nil, errtype.ErrDataBase(errtype.Join(ErrReadDB, err))
	}

	return ent, nil
}

func (db *StudioDB) FetchCustomers() (customers []bt.Customer, err error) {
	sp := selectParams{
		"id, first_name, last_name",
		"customers",
		"first_name, last_name",
		[]whereClause{},
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

func (*StudioDB) unrawOrdersItems(rawOI []bt.RawOrderItem, models []bt.Model) []bt.OrderItem {
	var orderItems []bt.OrderItem
	orderItemMap := make(map[uint]*bt.OrderItem) // Храним ссылки на уже добавленные заказы

	modelsMap := make(map[uint]bt.Model)
	for _, model := range models {
		modelsMap[model.Id] = model
	}

	for _, roi := range rawOI {
		orderItem, exists := orderItemMap[roi.Id]
		if !exists {
			orderItem = &bt.OrderItem{
				Id:        roi.Id,
				O_id:      roi.O_id,
				Model:     []bt.Model{}, // Инициализируем пустым срезом
				UnitPrice: roi.UnitPrice,
			}
			orderItemMap[roi.Id] = orderItem
		}

		if model, found := modelsMap[roi.Model]; found {
			orderItemMap[roi.Id].Model = append(orderItem.Model, model)
		}
	}

	for _, orderItem := range orderItemMap {
		orderItems = append(orderItems, *orderItem)
	}

	return orderItems
}

func (db *StudioDB) fetchTable(sp selectParams, dest interface{}) error {
	rows, err := db.query(sp)
	if err != nil {
		return err
	}
	defer rows.Close()

	// Получаем тип среза, куда будем записывать данные
	destSlice := reflect.ValueOf(dest)
	if destSlice.Kind() != reflect.Ptr || destSlice.Elem().Kind() != reflect.Slice {
		return errors.New("dest must be a pointer to a slice")
	}

	elemType := destSlice.Elem().Type().Elem()

	for rows.Next() {
		// Создаём новый экземпляр структуры
		elem := reflect.New(elemType).Elem()

		// Создаём массив для сканирования полей
		numFields := elem.NumField()
		scanDest := make([]interface{}, numFields)

		for i := 0; i < numFields; i++ {
			scanDest[i] = elem.Field(i).Addr().Interface()
		}

		if err := rows.Scan(scanDest...); err != nil {
			return errtype.ErrDataBase(errtype.Join(ErrReadDB, err))
		}

		// Добавляем элемент в срез
		destSlice.Elem().Set(reflect.Append(destSlice.Elem(), elem))

		log.Printf("[db.fetchTable]: Added element %v type of %v\n", elem, elemType)
	}

	return nil
}

func (db *StudioDB) fetchModels() (models []bt.Model, err error) {
	var (
		sp           selectParams
		rows, mmRows *sql.Rows
		m_id         uint
		title        string
		leng         float64
		price        float64
		model        bt.Model
	)

	sp = selectParams{"id, title", "models", "id", []whereClause{}}
	if rows, err = db.query(sp); err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		if err := rows.Scan(&m_id, &title); err != nil {
			return nil, errtype.ErrDataBase(errtype.Join(ErrReadDB, err))
		}

		model.Id = m_id
		model.Title = title
		model.Materials = make(map[uint]bt.Material)
		model.MatLeng = make(map[uint]float64)

		sp = selectParams{
			"m.id, m.title, mm.leng, m.price",
			"model_materials mm JOIN materials m ON mm.material_id = m.id", "",
			[]whereClause{{"mm.model_id", "=", m_id, ""}},
		}
		if mmRows, err = db.query(sp); err != nil {
			return nil, err
		}
		for mmRows.Next() {
			if err := mmRows.Scan(&m_id, &title, &leng, &price); err != nil {
				return nil, errtype.ErrDataBase(errtype.Join(ErrReadDB, err))
			}

			model.Materials[m_id] = bt.Material{
				Id:    m_id,
				Title: title,
				Price: price,
			}
			model.MatLeng[m_id] = leng
		}
		models = append(models, model)

	}

	return models, nil
}

func (db *StudioDB) fetchOrders(sp selectParams) (orders []bt.Order, _ error) {
	var rawOrders []bt.RawOrder
	if err := db.fetchTable(sp, &rawOrders); err != nil {
		return nil, err
	}

	var order bt.Order
	for _, rawOrder := range rawOrders {
		order = bt.Order{
			Id:          rawOrder.Id,
			C_id:        rawOrder.C_id,
			E_id:        rawOrder.E_id,
			Status:      rawOrder.Status,
			TotalPrice:  rawOrder.TotalPrice,
			CreateDate:  time.Unix(rawOrder.CreateDate, 0),
			ReleaseDate: time.Unix(rawOrder.ReleaseDate, 0),
		}

		orders = append(orders, order)
	}
	return orders, nil
}

// Общая функция для запросов в базе данных
func (db *StudioDB) query(sp selectParams) (*sql.Rows, error) {
	var (
		err   error
		query string
		rows  *sql.Rows
	)

	query = fmt.Sprintf("SELECT %s FROM %s ", sp.cols, sp.table)

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

func (db *StudioDB) insert(ip insertParams) error {
	var (
		err   error
		query string
	)

	query = fmt.Sprintf("INSERT INTO %s (%s) VALUES ", ip.table, ip.cols)

	if len(ip.values) == 0 {
		return errtype.ErrDataBase(ErrInsert)
	}

	for i, c := range ip.values {
		query += fmt.Sprintf("(%s)", c)

		if len(ip.values) != i+1 {
			query += ","
		}
	}

	log.Printf("[db.query]: %s", query)
	if _, err = db.sDB.Exec(query); err != nil {
		return errtype.ErrDataBase(errtype.Join(ErrInsert, err))
	}

	return nil
}
