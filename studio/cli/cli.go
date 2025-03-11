package cli

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"runtime"
	bt "studio/basic_types"
	"studio/cli/userinput"

	"github.com/AlecAivazis/survey/v2"
)

type CLI struct {
	ent bt.Entity
	opt []string
}

func (c *CLI) Run(ent bt.Entity) {
	c.ent = ent
	switch c.ent.GetAccessLevel() {
	case bt.CUSTOMER:
		c.opt = customerOptions()
	case bt.OPERATOR:
		c.opt = operatorOptions()
	case bt.SYSADMIN:
		c.opt = sysAdminOptions()
	}
	fmt.Println("Запуск консольного интерфейса...")
	fmt.Printf("С возвращением, %s!\n", c.ent.GetFirstLastName())
	pause()
}

func (c *CLI) Login() string {
	var (
		login string
		err   error
	)

	for {
		login, err = userinput.Prompt("Введите Ваш логин")
		if err == nil {
			break
		}
	}

	return login
}

func (c *CLI) Main() (choice string) {
	clearScreen()
	prompt := &survey.Select{
		Message:  "Выберите действие:",
		Options:  c.opt,
		PageSize: 4,
	}
	survey.AskOne(prompt, &choice)

	return choice
}

func (c *CLI) DisplayOrderStat() {
	fmt.Println("Статус заказов (консольный интерфейс)")
	pause()
}

func (c *CLI) DisplayOrders() {
	fmt.Println("Список заказов (консольный интерфейс)")
	pause()
}

func (c *CLI) CancelOrder() {
	fmt.Println("Отмена заказа (консольный интерфейс)")
	pause()
}

func (c *CLI) CreateOrder() {
	fmt.Println("Создание заказа через терминал")
	pause()
}

func (c *CLI) EditOrder(id uint) {
	fmt.Printf("Редактирование заказа %d через терминал\n", id)
	pause()
}

func (c *CLI) ProcessOrder(id uint) {
	fmt.Printf("Исполнение заказа %d через терминал\n", id)
	pause()
}

func (c *CLI) ReleaseOrder(id uint) {
	fmt.Printf("Выдача заказа %d через терминал\n", id)
	pause()
}

func (c *CLI) BackupDB() {
	fmt.Println("Резервное копирование через терминал")
	pause()
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
		"Копирование БД",
		"Выход",
	}
}

// Очищает консоль
func clearScreen() {
	switch runtime.GOOS {
	case "windows":
		cmd := exec.Command("cmd", "/c", "cls")
		cmd.Stdout = os.Stdout
		cmd.Run()
	default:
		cmd := exec.Command("clear")
		cmd.Stdout = os.Stdout
		cmd.Run()
	}
}

// Выводит сообщение и ждет нажатия клавиши
func pause() {
	fmt.Print("Нажмите Enter для продолжения...")

	if runtime.GOOS == "windows" {
		exec.Command("cmd", "/c", "pause >nul").Run()
	} else {
		bufio.NewReader(os.Stdin).ReadBytes('\n')
	}
}
