package studio

import "fmt"

var (
	ErrInitTables = fmt.Errorf("ошибка инициализации таблиц")
	ErrPerm       = fmt.Errorf("недостаточно прав для совершения этого действия")
)
