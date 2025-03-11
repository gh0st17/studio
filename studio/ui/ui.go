package ui

import "studio/basic_types"

type UI interface {
	Run(ent basic_types.Entity) // Запуск интерфейса
	Login() string              // Авторизация пользователя
	Main() string               // Страница/список доступных действии
	DisplayOrderStat()          // Отображение статуса заказа
	DisplayOrders()             // Отображение списка заказов
	CancelOrder()               // Отмена заказа
	CreateOrder()               // Создание нового заказа
	EditOrder(id uint)          // Редактирование заказа
	ProcessOrder(id uint)       // Исполнение заказа
	ReleaseOrder(id uint)       // Выдача заказа
	BackupDB()                  // Резервное копирование БД
}
