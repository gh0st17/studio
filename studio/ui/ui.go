package ui

import bt "studio/basic_types"

type UI interface {
	Run(bt.Entity)                    // Запуск интерфейса
	Login() string                    // Авторизация пользователя системы
	Registration(string) bt.Customer  // Регистрация клиента
	Main() string                     // Страница/список доступных действии
	DisplayOrders([]bt.Order)         // Отображение списка заказов
	SelectOrderId() (uint, error)     // Выбор id заказа
	DisplayOrderItems([]bt.OrderItem) // Отображение содержимого заказа
	CancelOrder(uint)                 // Отмена заказа
	CreateOrder()                     // Создание нового заказа
	CompleteOrder(uint)               // Выполнение заказа
}
