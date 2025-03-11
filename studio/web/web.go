package web

import "fmt"

type Web struct{}

func (w *Web) Run() {
	fmt.Println("Запуск консольного интерфейса...")
}
