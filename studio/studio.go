package studio

import (
	"log"
	bt "studio/basic_types"
	db "studio/database"
	"studio/errtype"
)

type Studio struct {
	sDB       db.StudioDB
	materials map[uint]bt.Material
	models    map[uint]bt.Model
}

func New() (st *Studio, err error) {
	st = &Studio{}
	if err = st.sDB.LoadDB(); err != nil {
		return nil, err
	}

	if err = st.initTables(); err != nil {
		return nil, errtype.ErrRuntime(errtype.Join(ErrInitTables, err))
	}

	return st, nil
}

func (s *Studio) Shutdown() (err error) {
	return s.sDB.CloseDB()
}

func (s *Studio) Registration(customer bt.Customer) (err error) {
	if err = s.sDB.Registration(customer); err != nil {
		return err
	}
	return nil
}

func (s *Studio) Login(login string) (ent bt.Entity, err error) {
	if ent, err = s.sDB.Login(login); err != nil {
		return nil, err
	}

	return ent, nil
}

func (s *Studio) CreateOrder(ent bt.Entity, ids []uint) error {
	if len(ids) == 0 {
		return errtype.ErrRuntime(ErrEmptyCart)
	}

	var cartModels []bt.Model
	for _, id := range ids {
		cartModels = append(cartModels, s.models[id])
	}

	err := s.sDB.CreateOrder(ent.GetId(), cartModels)
	if err != nil {
		log.Printf("Не удалось создать заказ: %v\n", err)
		return nil
	}

	return nil
}

func (s *Studio) Models() map[uint]bt.Model {
	return s.models
}

func (s *Studio) Orders(ent bt.Entity) ([]bt.Order, error) {
	switch ent.AccessLevel() {
	case bt.OPERATOR:
		return s.sDB.FetchOrders(0)
	default:
		return s.sDB.FetchOrders(ent.GetId())
	}
}

func (s *Studio) OrderItems(ent bt.Entity, id uint) ([]bt.OrderItem, error) {
	if ok, err := s.checkOrder(ent, id); err != nil {
		return nil, err
	} else if ok {
		return s.sDB.FetchOrderItems(id, s.models)
	} else {
		return nil, nil
	}
}

func (s *Studio) CancelOrder(ent bt.Entity, id uint) error {
	if ok, err := s.checkOrder(ent, id); err != nil {
		return err
	} else if !ok {
		return ErrPerm
	}

	if err := s.sDB.SetOrderStatus(id, bt.Canceled); err != nil {
		return err
	}

	return nil
}

func (s *Studio) ProcessOrder(ent bt.Entity, id uint) error {
	if ent.AccessLevel() != bt.OPERATOR {
		return ErrPerm
	}

	if err := s.sDB.SetOrderStatus(id, bt.Processing); err != nil {
		return err
	}

	if err := s.sDB.SetOperator(ent.GetId(), id); err != nil {
		return err
	}

	return nil
}

func (s *Studio) ReleaseOrder(ent bt.Entity, id uint) error {
	if ent.AccessLevel() != bt.OPERATOR {
		return ErrPerm
	}

	if err := s.sDB.SetOrderStatus(id, bt.Released); err != nil {
		return err
	}

	return nil
}

func (s *Studio) initTables() (err error) {
	if s.materials, err = s.sDB.FetchMaterials(); err != nil {
		return err
	}
	if s.models, err = s.sDB.FetchModels(); err != nil {
		return err
	}

	return nil
}

func (s *Studio) checkOrder(ent bt.Entity, id uint) (bool, error) {
	orders, err := s.Orders(ent)
	if err != nil {
		return false, err
	}

	for _, o := range orders {
		if o.Id == id {
			return true, nil
		}
	}

	return false, ErrOrder404
}
