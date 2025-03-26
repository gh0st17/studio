package database

import (
	"fmt"
)

var (
	ErrOpenDB      = fmt.Errorf("ошибка открытия базы данных")
	ErrCloseDB     = fmt.Errorf("ошибка закрытия базы данных")
	ErrInsert      = fmt.Errorf("ошибка при вставке в базу данных")
	ErrUpdate      = fmt.Errorf("ошибка при обновлении в базе данных")
	ErrQuery       = fmt.Errorf("ошибка при запросе в базе данных")
	ErrReadDB      = fmt.Errorf("ошибка чтения базы данных")
	ErrDelete      = fmt.Errorf("ошибка при удалении из базы данных")
	ErrLogin       = fmt.Errorf("ошибка авторизации")
	ErrBegin       = fmt.Errorf("не удалось начать транзакцию")
	ErrNotPending  = fmt.Errorf("заказ не находится в состоянии ожидания")
	ErrStatusRange = fmt.Errorf("изменение статуса выходит за пределы допустимого")
)
