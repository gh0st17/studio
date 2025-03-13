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

func (c *CLI) Registration(login string) (customer bt.Customer) {
	customer.FirstName, _ = userinput.PromptString("Введите свое имя")
	customer.LastName, _ = userinput.PromptString("Введите свою фамилию")
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

func (c *CLI) DisplayOrders(o []bt.Order) {
	fmt.Print(Orders(o))
	pause()
}

func (c *CLI) SelectOrderId() (uint, error) {
	return userinput.PromptUint("Выберите id заказа")
}

func (c *CLI) DisplayOrderItems(oI []bt.OrderItem) {
	fmt.Print(OrderItems(oI))
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
	var (
		ctime, rtime string
	)
	s = fmt.Sprintf(
		"  # Статус заказа %9s %19s %19s\n",
		"Сумма", "Создан", "Выдан",
	)

	for _, o := range orders {
		ctime = o.CreateDate.Format(dateFormat)

		if o.ReleaseDate != time.Unix(0, 0) {
			rtime = o.ReleaseDate.Format(dateFormat)
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
	Model      []bt.Model
)

func (ois OrderItems) String() (s string) {
	var sum float64 = 0.0

	for _, oi := range ois {
		s += Model(oi.Model).String()
		s += fmt.Sprintf("Cтоимость изготовления %2.2f\n\n", oi.UnitPrice)
		sum += oi.UnitPrice
	}
	s += fmt.Sprintln("Общая стоимость заказа: ", sum)

	return s
}

func (mod Model) String() (s string) {
	for _, m := range mod {
		s += fmt.Sprintf("%s:\n", m.Title)

		for _, mat := range m.Materials {
			s += fmt.Sprintf("\t%s стоимостью %2.2f за погонный метр длиной %2.2f метра\n",
				mat.Title, mat.Price, m.MatLeng[m.Id],
			)
		}
	}

	return s
}
