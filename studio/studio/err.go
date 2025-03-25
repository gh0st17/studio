package studio

import "fmt"

var (
	ErrInitTables = fmt.Errorf("ошибка инициализации таблиц")
	ErrUpdOrders  = fmt.Errorf("ошибка обновления таблицы заказов")
)
