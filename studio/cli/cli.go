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
		login, err = userinput.PromptString("Введите Ваш логин")
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

func (c *CLI) DisplayOrders(o []bt.Order) {
	fmt.Println(Orders(o))
	pause()
}

func (c *CLI) SelectOrderId() (uint, error) {
	return userinput.PromptUint("Выберите id заказа")
}

func (c *CLI) DisplayOrderItems(oI []bt.OrderItem) {
	fmt.Println("id oid mid unit_price")
	fmt.Println(OrderItems(oI))
	pause()
}

func (c *CLI) CancelOrder(id uint) {
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
		"Просмотреть заказы",
		"Просмотреть содержимое заказa",
		"Отменить заказ",
		"Выход",
	}
}

func operatorOptions() []string {
	return []string{
		"Просмотреть заказы",
		"Просмотреть содержимое заказa",
		"Редактировать заказ",
		"Выполнить заказ",
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

type Orders []bt.Order

func (orders Orders) String() (s string) {
	s = fmt.Sprintln("id cid eid status total_price created released")

	for _, o := range orders {
		s += fmt.Sprintln(
			o.Id, o.C_id, o.E_id, o.Status, o.TotalPrice, o.CreateDate, o.ReleaseDate,
		)
	}
	s += fmt.Sprintln()

	return s
}

type OrderItems []bt.OrderItem

func (ois OrderItems) String() (s string) {
	for _, oi := range ois {
		s += fmt.Sprintln(
			oi.Id, oi.O_id, oi.Model, oi.UnitPrice,
		)
	}
	s += fmt.Sprintln()

	return s
}
