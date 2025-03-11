package params

import (
	"fmt"
)

var (
	ErrUnknownIType = fmt.Errorf("неизвестный тип интерфейса")
	ErrMissingIType = fmt.Errorf("тип интерфейса не указан")
	ErrDBPath       = fmt.Errorf("путь до файла базы данных не указан")
	ErrDBExists     = fmt.Errorf("путь до файла базы данных не существует")
)
