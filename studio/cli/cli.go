package cli

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"runtime"
	bt "studio/basic_types"
	"studio/cli/userinput"
	"time"

	"github.com/AlecAivazis/survey/v2"
)

const dateFormat string = "02.01.2006 15:04:05"

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
	}
	fmt.Println("Запуск консольного интерфейса...")
	fmt.Printf("С возвращением, %s!\n", c.ent.GetFirstLastName())
	pause()
}

func (c *CLI) Login() string {
	return userinput.PromptString("Введите Ваш логин")
}

func (c *CLI) Registration(login string) (customer bt.Customer) {
	customer.FirstName = userinput.PromptString("Введите свое имя")
	customer.LastName = userinput.PromptString("Введите свою фамилию")
	customer.Login = login

	return customer
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

func (c *CLI) DisplayTable(table interface{}) {
	if orders, ok := table.([]bt.Order); ok {
		c.displayOrders(orders)
	} else if orderItems, ok := table.([]bt.OrderItem); ok {
		c.displayOrderItems(orderItems)
	} else if models, ok := table.([]bt.Model); ok {
		c.displayModels(models)
	} else {
		panic("Неизвестный тип таблицы")
	}
}

func (c *CLI) displayOrders(o []bt.Order) {
	if len(o) > 0 {
		fmt.Print(Orders(o))
	} else {
		fmt.Println("Вы еще не совершали заказов")
	}
	pause()
}

func (c *CLI) displayOrderItems(oI []bt.OrderItem) {
	if len(oI) > 0 {
		fmt.Print(OrderItems(oI))
	} else {
		fmt.Println("Заказа с таким номером не существует")
	}
	pause()
}

func (c *CLI) displayModels(m []bt.Model) {
	fmt.Print(Models(m))
}

func (c *CLI) ReadNumbers(prompt string) ([]uint, error) {
	return userinput.PromptUint(prompt)
}

func (c *CLI) CreateOrder() {
	fmt.Println("Создание заказа через терминал")
	pause()
}

func (c *CLI) Alert(msg string) {
	fmt.Println(msg)
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
		"Выполнить заказ",
		"Выдача заказа",
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
	var (
		ctime, rtime string
	)
	s = fmt.Sprintf(
		"  # Статус заказа %9s %19s %19s\n",
		"Сумма", "Создан", "Выдан",
	)

	for _, o := range orders {
		ctime = o.LocalCreateDate().Format(dateFormat)

		if o.LocalReleaseDate() != time.Unix(0, 0) {
			rtime = o.LocalReleaseDate().Format(dateFormat)
		} else {
			rtime = "---"
		}

		s += fmt.Sprintf("%3d %13s %9.2f %19s %19s\n",
			o.Id, o.Status, o.TotalPrice, ctime, rtime,
		)
	}

	return s
}

type (
	OrderItems []bt.OrderItem
	Model      bt.Model
	Models     []bt.Model
)

func (ois OrderItems) String() (s string) {
	var sum float64 = 0.0

	for i, oi := range ois {
		s += fmt.Sprintln("Позиция:", i+1)
		s += Model(oi.Model).String()
		sum += oi.UnitPrice
	}
	s += fmt.Sprintln("Общая стоимость заказа:", sum)

	return s
}

func (m Model) String() (s string) {
	s += fmt.Sprintf("%s (Артикул %d):\n", m.Title, m.Id)

	for _, mat := range m.Materials {
		s += fmt.Sprintf("\t%s стоимостью %2.2f за погонный метр длиной %2.2f метра\n",
			mat.Title, mat.Price, m.MatLeng[m.Id],
		)
	}
	s += fmt.Sprintf("\tCтоимость изготовления %2.2f\n\n", m.Price)

	return s
}

func (mod Models) String() (s string) {
	for _, m := range mod {
		s += Model(m).String()
	}

	return s
}
