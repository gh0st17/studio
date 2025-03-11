package database

import "fmt"

var (
	ErrOpenDB  = fmt.Errorf("ошибка открытия базы данных")
	ErrCloseDB = fmt.Errorf("ошибка закрытия базы данных")
	ErrInsert  = fmt.Errorf("ошибка при добавлении записи в базу данных")
	ErrQuery   = fmt.Errorf("ошибка при запросе в базе данных")
	ErrReadDB  = fmt.Errorf("ошибка чтения базы данных")
	ErrDelete  = fmt.Errorf("ошибка при запросе удаления групп в базе данных")
)
