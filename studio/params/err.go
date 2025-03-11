package params

import (
	"fmt"
)

var (
	ErrUnknownIType  = fmt.Errorf("неизвестный тип интерфейса")
	ErrMissingIType  = fmt.Errorf("тип интерфейса не указан")
	ErrMissingDBPath = fmt.Errorf("путь до файла базы данных не указан")
	ErrDBNotExists   = fmt.Errorf("путь до файла базы данных не существует")
)
