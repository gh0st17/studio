package database

import (
	"fmt"
)

var (
	ErrOpenDB      = fmt.Errorf("ошибка открытия базы данных")
	ErrPingDB      = fmt.Errorf("ошибка подключения к базе данных")
	ErrCloseDB     = fmt.Errorf("ошибка закрытия базы данных")
	ErrExec        = fmt.Errorf("ошибка при вставке или обновлении в базу данных")
	ErrQuery       = fmt.Errorf("ошибка при запросе в базе данных")
	ErrFetchTable  = fmt.Errorf("dest должен быть указателем на срез")
	ErrReadDB      = fmt.Errorf("ошибка чтения базы данных")
	ErrDelete      = fmt.Errorf("ошибка при удалении из базы данных")
	ErrLogin       = fmt.Errorf("ошибка авторизации")
	ErrBegin       = fmt.Errorf("не удалось начать транзакцию")
	ErrNotPending  = fmt.Errorf("заказ не находится в состоянии ожидания")
	ErrStatusRange = fmt.Errorf("изменение статуса выходит за пределы допустимого")
)
