package ui

type UI interface {
	Start()                 // Запуск интерфейса
	Stop()                  // Завершение работы интерфейса
	DisplayOrderStat()      // Отображение статуса заказа
	DisplayOrders()         // Отображение списка заказов
	CreateOrder()           // Создание нового заказа
	EditOrder(id uint)      // Редактирование заказа
	ProcessOrder(id uint)   // Исполнение заказа
	ReleaseOrder(id uint)   // Выдача заказа
	ShowMessage(msg string) // Отображение сообщений пользователю
	BackupDB()              // Резервное копирование БД
}
