package cli

import (
	"fmt"
	"studio/basic_types"
	bt "studio/basic_types"

	"github.com/AlecAivazis/survey/v2"
)

type CLI struct {
	ent         bt.Entity
	entAccLevel bt.AccessLevel
}

func (c *CLI) Run(ent basic_types.Entity) {
	c.ent = ent
	fmt.Println("Запуск консольного интерфейса...")
}

func (c *CLI) Login() uint {

	return 0
}

func (c *CLI) Main() (choice string) {
	var opt []string
	switch c.entAccLevel {
	case bt.CUSTOMER:
		opt = customerOptions()
	case bt.OPERATOR:
		opt = operatorOptions()
	case bt.SYSADMIN:
		opt = sysAdminOptions()
	}

	prompt := &survey.Select{
		Message:  "Выберите действие:",
		Options:  opt,
		PageSize: 4,
	}
	survey.AskOne(prompt, &choice)

	fmt.Println("Вы выбрали:", choice)
	return choice
}

func (c *CLI) DisplayOrderStat() {
	fmt.Println("Список stat заказов (консольный интерфейс)")
}

func (c *CLI) DisplayOrders() {
	fmt.Println("Список заказов (консольный интерфейс)")
}

func (c *CLI) CreateOrder() {
	fmt.Println("Создание заказа через терминал")
}

func (c *CLI) EditOrder(id uint) {
	fmt.Printf("Редактирование заказа %d через терминал\n", id)
}

func (c *CLI) ProcessOrder(id uint) {
	fmt.Printf("Исполнение заказа %d через терминал\n", id)
}

func (c *CLI) ReleaseOrder(id uint) {
	fmt.Printf("Выдача заказа %d через терминал\n", id)
}

func (c *CLI) ShowMessage(msg string) {
	fmt.Println("Сообщение:", msg)
}

func (c *CLI) BackupDB() {
	fmt.Println("Резервное копирование через терминал")
}

func customerOptions() []string {
	return []string{
		"Создать заказ",
		"Отменить заказ",
		"Просмотреть статус заказов",
		"Выход",
	}
}

func operatorOptions() []string {
	return []string{
		"Просмотреть заказы",
		"Редактировать заказ",
		"Исполнение заказа",
		"Выдача заказа",
		"Выход",
	}
}

func sysAdminOptions() []string {
	return []string{
		"Создать заказ",
		"Просмотреть заказы",
		"Просмотреть статус заказов",
		"Редактировать заказ",
		"Исполнение заказа",
		"Выдача заказа",
		"Копирование БД",
		"Выход",
	}
}
