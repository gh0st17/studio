package studio

import (
	"log"
	bt "studio/basic_types"
	db "studio/database"
	"studio/errtype"
	"studio/ui"
)

type Studio struct {
	ui         ui.UI
	ent        bt.Entity
	sDB        db.StudioDB
	customers  []bt.Customer
	orders     []bt.Order
	orderItems map[uint][]bt.OrderItem
	materials  []bt.Material
	models     []bt.Model
}

func New(ui ui.UI) (*Studio, error) {
	s := Studio{
		ui: ui,
	}

	return &s, nil
}

func (s *Studio) updateOrders(accLevel bt.AccessLevel) (err error) {
	switch accLevel {
	case bt.CUSTOMER:
		if s.orders, err = s.sDB.FetchOrdersByCid(s.ent.GetId()); err != nil {
			return err
		}
	case bt.OPERATOR:
		if s.orders, err = s.sDB.FetchOrders(); err != nil {
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

func (s *Studio) Run(dbPath string, reg bool) (err error) {
	if err = s.sDB.LoadDB(dbPath); err != nil {
		return err
	}

	login := s.ui.Login()

	if reg {
		customer := s.ui.Registration(login)
		if err = s.sDB.Registration(customer); err != nil {
			return err
		}
	}

	if s.ent, err = s.sDB.Login(login); err != nil {
		return err
	}

	if err = s.initTables(s.ent.GetAccessLevel()); err != nil {
		return errtype.ErrRuntime(errtype.Join(ErrInitTables, err))
	}

	s.ui.Run(s.ent)

	for {
		choice := s.ui.Main()
		switch choice {
		case "Создать заказ":
			s.CreateOrder()
		case "Просмотреть заказы":
			s.DisplayOrders()
		case "Просмотреть содержимое заказa":
			id, _ := s.ui.ReadNumbers("Выберите id заказа")
			s.DisplayOrderItems(id[0])
		case "Отменить заказ":
			id, _ := s.ui.ReadNumbers("Выберите id заказа")
			err = s.CancelOrder(id[0])
		case "Выполнить заказ":
			id, _ := s.ui.ReadNumbers("Выберите id заказа")
			err = s.ProcessOrder(id[0])
		case "Выдача заказа":
			id, _ := s.ui.ReadNumbers("Выберите id заказа")
			err = s.ReleaseOrder(id[0])
		case "Выход":
			if err = s.sDB.CloseDB(); err != nil {
				return err
			}
			return nil
		}

		if err != nil {
			return err
		}
	}
}

func (s *Studio) CreateOrder() error {
	var cartModels []bt.Model

	s.ui.DisplayTable(s.models)
	ids, _ := s.ui.ReadNumbers("Выберите модели по номеру артикула")

	for _, id := range ids {
		cartModels = append(cartModels, s.models[id-1])
	}

	err := s.sDB.CreateOrder(s.ent.GetId(), cartModels)
	if err != nil {
		s.ui.Alert("Не удалось создать заказ")
		log.Fatalf("Не удалось создать заказ: %v\n", err)
		return nil
	}

	return s.updateOrders(s.ent.GetAccessLevel())
}

func (s *Studio) DisplayOrders() {
	s.ui.DisplayTable(s.orders)
}

func (s *Studio) DisplayOrderItems(id uint) {
	s.ui.DisplayTable(s.orderItems[id])
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
