package database

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
	"reflect"
	bt "studio/basic_types"
	"studio/errtype"
)

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

	if !rows.Next() {
		return nil, errtype.ErrDataBase(ErrLogin)
	}

	ent := &bt.Customer{}
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

	if !rows.Next() {
		return nil, errtype.ErrDataBase(ErrLogin)
	}

	ent := &bt.Employee{}
	if err = rows.Scan(&ent.Id, &ent.FirstName, &ent.LastName, &ent.JobId); err != nil {
		return nil, errtype.ErrDataBase(errtype.Join(ErrReadDB, err))
	}

	return ent, nil
}

func (db *StudioDB) fetchTable(sp selectParams, dest interface{}) error {
	rows, err := db.query(sp)
	if err != nil {
		return err
	}
	defer rows.Close()

	destSlice := reflect.ValueOf(dest)
	if destSlice.Kind() != reflect.Ptr || destSlice.Elem().Kind() != reflect.Slice {
		return errors.New("dest must be a pointer to a slice")
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

func (db *StudioDB) fetchModels() (models []bt.Model, err error) {
	var (
		sp           selectParams
		rows, mmRows *sql.Rows
		m_id         uint
		title        string
		leng, price  float64
		model        bt.Model
	)

	sp = selectParams{"id, title, price", "models", "id", []whereClause{}}
	if rows, err = db.query(sp); err != nil {
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

		sp = selectParams{
			"m.id, m.title, mm.leng, m.price",
			"model_materials mm JOIN materials m ON mm.material_id = m.id", "",
			[]whereClause{{"mm.model_id", "=", fmt.Sprint(m_id), ""}},
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

func (db *StudioDB) getLastId(table string, w []whereClause) (uint, error) {
	type Id struct {
		Id uint
	}

	sp := selectParams{
		"id", table, "id", w,
	}

	var id []Id
	if err := db.fetchTable(sp, &id); err != nil {
		return 0, err
	}

	return id[len(id)-1].Id, nil
}
