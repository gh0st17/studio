package cli

import (
	"fmt"

	"github.com/AlecAivazis/survey/v2"
)

type CLI struct{}

func (c *CLI) Run() {
	fmt.Println("Запуск консольного интерфейса...")
	printMainMenu()
}

func (c *CLI) DisplayOrders() {
	fmt.Println("Список заказов (консольный интерфейс)")
}

func (c *CLI) CreateOrder() {
	fmt.Println("Создание заказа через консоль")
}

func (c *CLI) EditOrder(id int) {
	fmt.Printf("Редактирование заказа %d через консоль\n", id)
}

func (c *CLI) DeleteOrder(id int) {
	fmt.Printf("Удаление заказа %d через консоль\n", id)
}

func (c *CLI) ShowMessage(msg string) {
	fmt.Println("Сообщение:", msg)
}

func printMainMenu() {
	var choice string
	prompt := &survey.Select{
		Message: "Выберите действие:",
		Options: []string{"Создать заказ", "Просмотреть заказы", "Выход"},
	}
	survey.AskOne(prompt, &choice)

	fmt.Println("Вы выбрали:", choice)
}
