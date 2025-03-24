package studio

import (
	"bufio"
	"fmt"
	"os"
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

func New(ui ui.UI) (_ *Studio, err error) {
	s := Studio{
		ui: ui,
	}

	return &s, nil
}

func (s *Studio) initTables() (err error) {
	accLevel := s.ent.GetAccessLevel()

	switch accLevel {
	case bt.CUSTOMER:
		if s.orders, err = s.sDB.FetchOrdersByCid(s.ent.GetId()); err != nil {
			return err
		}
	case bt.OPERATOR:
		if s.customers, err = s.sDB.FetchCustomers(); err != nil {
			return err
		}
		if s.orders, err = s.sDB.FetchOrders(); err != nil {
			return err
		}
	}

	switch accLevel {
	case bt.CUSTOMER, bt.OPERATOR:
		if s.materials, err = s.sDB.FetchMaterials(); err != nil {
			return err
		}
		if s.models, err = s.sDB.FetchModels(); err != nil {
			return err
		}
		if s.orderItems, err = s.sDB.FetchOrderItems(s.orders, s.models); err != nil {
			return err
		}
	}

	return nil
}

func (s *Studio) Run(dbPath string, reg bool) (err error) {
	if err = s.sDB.LoadDB(dbPath); err != nil {
		return err
	}

	login := "user" // s.ui.Login()

	if reg {
		customer := s.ui.Registration(login)
		if err = s.sDB.Registration(customer); err != nil {
			return err
		}
	}

	if s.ent, err = s.sDB.Login(login); err != nil {
		return err
	}

	if s.ent.GetAccessLevel() != 3 {
		if err = s.initTables(); err != nil {
			return errtype.ErrRuntime(errtype.Join(ErrInitTables, err))
		}
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
			id, _ := s.ui.ReadNumber("Выберите id заказа")
			s.DisplayOrderItems(id)
		case "Отменить заказ":
			id, _ := s.ui.ReadNumber("Выберите id заказа")
			s.CancelOrder(id)
		case "Выполнение заказа":
			id, _ := s.ui.ReadNumber("Выберите id заказа")
			s.CompleteOrder(id)
		case "Выход":
			if err = s.sDB.CloseDB(); err != nil {
				return err
			}
			return nil
		}
	}
}

func (s *Studio) DisplayOrders() {
	s.ui.DisplayTable(s.orders)
}

func (s *Studio) DisplayOrderItems(id uint) {
	s.ui.DisplayTable(s.orderItems[id])
}

func (s *Studio) CancelOrder(id uint) {
	s.ui.CancelOrder(id)
}

func (s *Studio) CreateOrder() {
	var models []bt.Model

	s.ui.DisplayTable(s.models)
	id, _ := s.ui.ReadNumber("Выберите модель по номеру артикула")

	models = append(models, s.models[id-1])

	err := s.sDB.CreateOrder(s.ent.GetId(), models)
	if err != nil {
		fmt.Printf("Не удалось создать заказ: %v\n", err)
		bufio.NewReader(os.Stdin).ReadBytes('\n')
	}
}

func (s *Studio) CompleteOrder(id uint) {
	s.ui.CompleteOrder(id)
}
