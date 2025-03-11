package ui

import "studio/basic_types"

type UI interface {
	Run(ent basic_types.Entity) // Запуск интерфейса
	Login(login string) uint    // Авторизация пользователя
	Main() string               // Страница/список доступных действии
	DisplayOrderStat()          // Отображение статуса заказа
	DisplayOrders()             // Отображение списка заказов
	CreateOrder()               // Создание нового заказа
	EditOrder(id uint)          // Редактирование заказа
	ProcessOrder(id uint)       // Исполнение заказа
	ReleaseOrder(id uint)       // Выдача заказа
	ShowMessage(msg string)     // Отображение сообщений пользователю
	BackupDB()                  // Резервное копирование БД
}
