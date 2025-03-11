package studio

import (
	bt "studio/basic_types"
	db "studio/database"
	"studio/ui"
)

type Studio struct {
	ui        ui.UI
	ent       bt.Entity
	sDB       db.StudioDB
	customers []bt.Customer
	orders    []bt.Order
	materials []bt.Material
	models    []bt.Model
}

func New(ui ui.UI, ent bt.Entity) (_ *Studio, err error) {
	s := Studio{
		ui:  ui,
		ent: ent,
	}

	accLevel := ent.GetAccessLevel()

	switch accLevel {
	case bt.CUSTOMER:
		if s.orders, err = s.sDB.FetchOrdersByCustId(ent.GetId()); err != nil {
			return nil, err
		}
	case bt.OPERATOR:
		if s.customers, err = s.sDB.FetchCustomers(); err != nil {
			return nil, err
		}
		if s.orders, err = s.sDB.FetchOrders(); err != nil {
			return nil, err
		}
		if s.materials, err = s.sDB.FetchMaterials(); err != nil {
			return nil, err
		}
		if s.models, err = s.sDB.FetchModels(); err != nil {
			return nil, err
		}
	}

	return &s, nil
}

func (s *Studio) Run() error {
	s.ui.Run(s.ent)

	for {
		choice := s.ui.Main()
		switch choice {
		case "Создать заказ":
			s.CreateOrder()
		case "Просмотреть заказы":
			s.DisplayOrders()
		case "Просмотреть статус заказов":
			s.DisplayOrderStat()
		case "Редактировать заказ":
			s.EditOrder(1)
		case "Исполнение заказа":
			s.ProcessOrder(1)
		case "Выдача заказа":
			s.ReleaseOrder(1)
		case "Копирование БД":
			s.BackupDB()
		case "Выход":
			return nil
		}
	}
}

func (s *Studio) DisplayOrderStat() {
	s.ui.DisplayOrderStat()
}

func (s *Studio) DisplayOrders() {
	s.ui.DisplayOrders()
}

func (s *Studio) CreateOrder() {
	s.ui.CreateOrder()
}

func (s *Studio) EditOrder(id uint) {
	s.ui.EditOrder(id)
}

func (s *Studio) ProcessOrder(id uint) {
	s.ui.ProcessOrder(id)
}

func (s *Studio) ReleaseOrder(id uint) {
	s.ui.ReleaseOrder(id)
}

func (s *Studio) BackupDB() {
	s.ui.BackupDB()
}
