package web

import (
	"fmt"
	"studio/basic_types"
)

type Web struct{}

func (w *Web) Run(ent basic_types.Entity) {
	panic("not implemented")

	//fmt.Println("Запуск web-интерфейса...")
}

func (w *Web) Login() string {
	return ""
}

func (w *Web) Main() string {
	return ""
}

func (w *Web) DisplayOrderStat() {
	fmt.Println("Статус заказов (web-интерфейс)")
}

func (w *Web) DisplayOrders() {
	fmt.Println("Список заказов (web-интерфейс)")
}

func (w *Web) CancelOrder() {
	fmt.Println("Отмена заказа (web-интерфейс)")
}

func (w *Web) CreateOrder() {
	fmt.Println("Создание заказа через web-интерфейс")
}

func (w *Web) EditOrder(id uint) {
	fmt.Printf("Редактирование заказа %d через web-интерфейс\n", id)
}

func (w *Web) ProcessOrder(id uint) {
	fmt.Printf("Исполнение заказа %d через web-интерфейс\n", id)
}

func (w *Web) ReleaseOrder(id uint) {
	fmt.Printf("Выдача заказа %d через web-интерфейс\n", id)
}

func (w *Web) BackupDB() {
	fmt.Println("Резервное копирование через web-интерфейс")
}
