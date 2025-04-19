package database

import (
	"database/sql"
	"fmt"
	"reflect"

	bt "github.com/gh0st17/studio/basic_types"
	"github.com/gh0st17/studio/errtype"
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
