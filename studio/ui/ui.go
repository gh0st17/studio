package ui

import bt "studio/basic_types"

type UI interface {
	Run(bt.Entity)                             // Запуск интерфейса
	Login() string                             // Авторизация пользователя системы
	Registration(string) bt.Customer           // Регистрация клиента
	Main() string                              // Страница/список доступных действии
	DisplayTable(interface{})                  // Отображение любой таблицы пользователя
	ReadNumbers(prompt string) ([]uint, error) // Чтение числа
	CancelOrder(uint)                          // Отмена заказа
	CreateOrder()                              // Создание нового заказа
	CompleteOrder(uint)                        // Выполнение заказа
	Alert(string)                              // Вывод сообщения пользователю
}
