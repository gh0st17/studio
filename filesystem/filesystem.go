// Пакет filesystem предоставляет набор функции для работы
// с файловой системой и ее элементами
package filesystem

import "os"

func Exsists(path string) bool {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return false
	} else {
		return true
	}
}
