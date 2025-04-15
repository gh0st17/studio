package database

import (
	"database/sql"
	"fmt"
	"log"
	"reflect"
	bt "studio/basic_types"
	"studio/errtype"
)

func (db *StudioDB) login(query string, args []any, dest interface{}) (bt.Entity, error) {
	err := db.fetchTable(query, args, dest)
	if err != nil {
		return nil, err
	}

	slice := reflect.ValueOf(dest).Elem()
	if slice.Len() == 0 {
		return nil, fmt.Errorf("no entity found")
	}

	entity := slice.Index(0).Addr().Interface().(bt.Entity)
	return entity, nil
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

func (db *StudioDB) fetchModels() (models map[uint]bt.Model, err error) {
	var (
		rows, mmRows *sql.Rows
		m_id         uint
		title        string
		leng, price  float64
		model        bt.Model
	)
	models = make(map[uint]bt.Model)

	if rows, err = db.queryRows(fetchModelsQuery, []any{}); err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		if err := rows.Scan(&m_id, &title, &price); err != nil {
			return nil, errtype.ErrDataBase(errtype.Join(ErrReadDB, err))
		}

		model.Id = m_id
		model.Title = title
		model.Price = price
		model.Materials = make(map[uint]bt.Material)
		model.MatLeng = make(map[uint]float64)

		if mmRows, err = db.queryRows(fetchModelMatQuery, []any{m_id}); err != nil {
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
		models[model.Id] = model

	}

	return models, nil
}
