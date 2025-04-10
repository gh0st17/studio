package studio

import (
	"log"
	bt "studio/basic_types"
	db "studio/database"
	"studio/errtype"
)

type Studio struct {
	ent        bt.Entity
	sDB        db.StudioDB
	customers  []bt.Customer
	orders     []bt.Order
	orderItems map[uint][]bt.OrderItem
	materials  []bt.Material
	models     []bt.Model
}

func (s *Studio) updateOrders(accLevel bt.AccessLevel) (err error) {
	switch accLevel {
	case bt.CUSTOMER:
		if s.orders, err = s.sDB.FetchOrders(s.ent.GetId()); err != nil {
			return err
		}
	case bt.OPERATOR:
		if s.orders, err = s.sDB.FetchOrders(0); err != nil {
			return err
		}
	}

	if s.orderItems, err = s.sDB.FetchOrderItems(s.orders, s.models); err != nil {
		return err
	}

	return nil
}

func (s *Studio) initTables(accLevel bt.AccessLevel) (err error) {
	if accLevel == bt.OPERATOR {
		if s.customers, err = s.sDB.FetchCustomers(); err != nil {
			return err
		}
	}

	if s.materials, err = s.sDB.FetchMaterials(); err != nil {
		return err
	}
	if s.models, err = s.sDB.FetchModels(); err != nil {
		return err
	}

	return s.updateOrders(accLevel)
}

func (s *Studio) Run(dbPath string) (err error) {
	if err = s.sDB.LoadDB(dbPath); err != nil {
		return err
	}

	return nil
}

func (s *Studio) Shutdown() (err error) {
	return s.sDB.CloseDB()
}

func (s *Studio) Registration(customer bt.Customer) (err error) {
	if err = s.sDB.Registration(customer); err != nil {
		return err
	}

	s.ent = &customer

	if err = s.initTables(s.ent.GetAccessLevel()); err != nil {
		return errtype.ErrRuntime(errtype.Join(ErrInitTables, err))
	}

	return nil
}

func (s *Studio) Login(login string) (ent bt.Entity, err error) {
	if s.ent, err = s.sDB.Login(login); err != nil {
		return nil, err
	}

	if err = s.initTables(s.ent.GetAccessLevel()); err != nil {
		return nil, errtype.ErrRuntime(errtype.Join(ErrInitTables, err))
	}

	return s.ent, nil
}

func (s *Studio) CreateOrder(ids []uint) error {
	var cartModels []bt.Model
	for _, id := range ids {
		cartModels = append(cartModels, s.models[id-1])
	}

	err := s.sDB.CreateOrder(s.ent.GetId(), cartModels)
	if err != nil {
		log.Fatalf("Не удалось создать заказ: %v\n", err)
		return nil
	}

	return s.updateOrders(s.ent.GetAccessLevel())
}

func (s *Studio) Models() []bt.Model {
	return s.models
}

func (s *Studio) Orders() []bt.Order {
	return s.orders
}

func (s *Studio) OrderItems(id uint) []bt.OrderItem {
	return s.orderItems[id]
}

func (s *Studio) CancelOrder(id uint) error {
	if err := s.sDB.SetOrderStatus(id, bt.Canceled); err != nil {
		return err
	}
	return s.updateOrders(s.ent.GetAccessLevel())
}

func (s *Studio) ProcessOrder(id uint) error {
	if err := s.sDB.SetOrderStatus(id, bt.Processing); err != nil {
		return err
	}
	return s.updateOrders(s.ent.GetAccessLevel())
}

func (s *Studio) ReleaseOrder(id uint) error {
	if err := s.sDB.SetOrderStatus(id, bt.Released); err != nil {
		return err
	}
	return s.updateOrders(s.ent.GetAccessLevel())
}
