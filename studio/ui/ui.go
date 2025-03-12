package ui

import bt "studio/basic_types"

type UI interface {
	Run(ent bt.Entity)                   // Запуск интерфейса
	Login() string                       // Авторизация пользователя
	Main() string                        // Страница/список доступных действии
	DisplayOrders(orders []bt.Order)     // Отображение списка заказов
	SelectOrderId() (uint, error)        // Выбор id заказа
	DisplayOrderItems(oI []bt.OrderItem) // Отображение содержимого заказа
	CancelOrder(id uint)                 // Отмена заказа
	CreateOrder()                        // Создание нового заказа
	EditOrder(id uint)                   // Редактирование заказа
	ProcessOrder(id uint)                // Исполнение заказа
	ReleaseOrder(id uint)                // Выдача заказа
	BackupDB()                           // Резервное копирование БД
}
