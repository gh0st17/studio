// Пакет errtype предоставляет логику обработки
// и порождения ошибок в программе
package errtype

import (
	"fmt"
	"os"
)

type Error struct {
	text string
	code int // Код завершения после вывода ошибки
}

// Перевод и форматирование встроенных ошибок
func (e Error) Error() string {
	return e.text
}

// Возвращает общую ошибку времени выполнения
func ErrRuntime(err error) error {
	return &Error{
		text: err.Error(),
		code: 1,
	}
}

// Возвращает ошибки обработки входных параметром программы
func ErrArgument(err error) error {
	return &Error{
		text: err.Error(),
		code: 2,
	}
}

// Возвращает ошибки при использовании БД
func ErrDataBase(err error) error {
	return &Error{
		text: err.Error(),
		code: 3,
	}
}

// Объединяет описание ошибок в цепочку
//
// Копирует логику [errors.Join], но делает
// это в одну строку
func Join(errs ...error) error {
	n := 0
	for _, err := range errs {
		if err != nil {
			n++
		}
	}
	if n == 0 {
		return nil
	}

	var e error
	for _, err := range errs {
		if err == nil {
			continue
		}

		if e == nil {
			e = err
		} else {
			e = fmt.Errorf("%v: %v", e, err)
		}
	}
	return e
}

// Обработчик ошибок
func ErrorHandler(err error) {
	fmt.Println(err)
	if e, ok := err.(*Error); ok {
		os.Exit(e.code)
	} else {
		os.Exit(-1)
	}
}
