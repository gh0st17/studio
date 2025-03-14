package web

import (
	"fmt"
	bt "studio/basic_types"
)

type Web struct{}

func (w *Web) Run(ent bt.Entity) {
	panic("not implemented")

	//fmt.Println("Запуск web-интерфейса...")
}

func (w *Web) Login() string {
	panic("not implemented")
}

func (w *Web) Main() string {
	return ""
}

func (w *Web) Registration(login string) (customer bt.Customer) {
	panic("not implemented")
}

func (w *Web) DisplayOrders(orders []bt.Order) {
	fmt.Println("Список заказов (web-интерфейс)")
}

func (w *Web) SelectOrderId() (uint, error) {
	fmt.Println("Список заказов (web-интерфейс)")

	return 0, nil
}

func (w *Web) DisplayOrderItems(oI []bt.OrderItem) {
	fmt.Println("Просмотр заказа (web-интерфейс)")
}

func (w *Web) CancelOrder(id uint) {
	fmt.Println("Отмена заказа (web-интерфейс)")
}

func (w *Web) CreateOrder() {
	fmt.Println("Создание заказа через web-интерфейс")
}

func (w *Web) CompleteOrder(id uint) {
	fmt.Printf("Выполнение заказа %d через web-интерфейс\n", id)
}

func (w *Web) BackupDB() {
	fmt.Println("Резервное копирование через web-интерфейс")
}
