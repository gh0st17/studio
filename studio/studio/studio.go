package studio

import (
	"log"
	bt "studio/basic_types"
	db "studio/database"
	"studio/errtype"
)

type Studio struct {
	sDB        db.StudioDB
	orderItems map[uint][]bt.OrderItem
	materials  []bt.Material
	models     []bt.Model
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

func New(dbPath string) (st *Studio, err error) {
	st = &Studio{}
	if err = st.sDB.LoadDB(dbPath); err != nil {
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
	var cartModels []bt.Model
	for _, id := range ids {
		cartModels = append(cartModels, s.models[id-1])
	}

	err := s.sDB.CreateOrder(ent.GetId(), cartModels)
	if err != nil {
		log.Fatalf("Не удалось создать заказ: %v\n", err)
		return nil
	}

	return nil
}

func (s *Studio) Models() []bt.Model {
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

func (s *Studio) OrderItems(id uint) []bt.OrderItem {
	return s.orderItems[id]
}

func (s *Studio) CancelOrder(id uint) error {
	if err := s.sDB.SetOrderStatus(id, bt.Canceled); err != nil {
		return err
	}

	return nil
}

func (s *Studio) ProcessOrder(id uint) error {
	if err := s.sDB.SetOrderStatus(id, bt.Processing); err != nil {
		return err
	}

	return nil
}

func (s *Studio) ReleaseOrder(id uint) error {
	if err := s.sDB.SetOrderStatus(id, bt.Released); err != nil {
		return err
	}

	return nil
}
