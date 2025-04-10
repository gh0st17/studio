package cli

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"runtime"
	bt "studio/basic_types"
	"studio/cli/userinput"
	db "studio/database"
	"studio/studio"
	"time"

	"github.com/AlecAivazis/survey/v2"
)

const dateFormat string = "02.01.2006 15:04:05"

type CLI struct {
	st       *studio.Studio
	accLevel bt.AccessLevel
	userName string
	opt      []string
}

func New() (c *CLI) {
	c = &CLI{}
	c.st = &studio.Studio{}

	return c
}

func (c *CLI) Run(dbPath string) error {
	if err := c.st.Run(dbPath); err != nil {
		return err
	}

	for {
		ent, err := c.st.Login(c.Login())

		if err == nil {
			c.accLevel = ent.GetAccessLevel()
			c.userName = ent.GetFirstLastName()
			break
		} else {
			fmt.Println(err)
			continue
		}
	}

	switch c.accLevel {
	case bt.CUSTOMER:
		c.opt = customerOptions()
	case bt.OPERATOR:
		c.opt = operatorOptions()
	}
	fmt.Println("Запуск консольного интерфейса...")
	fmt.Printf("С возвращением, %s!\n", c.userName)
	pause()

	var err error

	for {
		choice := c.Main()
		switch choice {
		case "Создать заказ":
			err = c.CreateOrder()
		case "Просмотреть заказы":
			c.displayOrders(c.st.Orders())
		case "Просмотреть содержимое заказa":
			id, _ := c.ReadNumbers("Выберите id заказа")
			c.displayOrderItems(c.st.OrderItems(id[0]))
		case "Отменить заказ":
			id, _ := c.ReadNumbers("Выберите id заказа")
			err = c.st.CancelOrder(id[0])
		case "Выполнить заказ":
			id, _ := c.ReadNumbers("Выберите id заказа")
			err = c.st.ProcessOrder(id[0])
		case "Выдача заказа":
			id, _ := c.ReadNumbers("Выберите id заказа")
			err = c.st.ReleaseOrder(id[0])
		case "Выход":
			if err = c.st.Shutdown(); err != nil {
				return err
			}
			return nil
		}

		switch err {
		case nil:
			continue
		case db.ErrNotPending, db.ErrStatusRange:
			c.Alert(fmt.Sprint(err))
			err = nil
		default:
			return err
		}
	}
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
		PageSize: 5,
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

func (c *CLI) SetAccessLevel(accLevel bt.AccessLevel) {
	c.accLevel = accLevel
}

func (c *CLI) SetUserName(userName string) {
	c.userName = userName
}

func (c *CLI) displayOrders(orders []bt.Order) {
	if len(orders) == 0 {
		fmt.Println("Вы еще не совершали заказов")
	}

	if c.accLevel == bt.CUSTOMER {
		fmt.Print(customerOrders(orders))
	} else {
		fmt.Print(employeeOrders(orders))
	}

	pause()
}

func (c *CLI) displayOrderItems(orderItems []bt.OrderItem) {
	if len(orderItems) > 0 {
		fmt.Print(OrderItems(orderItems))
	} else {
		fmt.Println("Заказа с таким номером не существует")
	}
	pause()
}

func (c *CLI) displayModels(models []bt.Model) {
	fmt.Print(Models(models))
}

func (c *CLI) ReadNumbers(prompt string) ([]uint, error) {
	return userinput.PromptUint(prompt)
}

func (c *CLI) CreateOrder() (err error) {
	c.DisplayTable(c.st.Models())
	ids, _ := c.ReadNumbers("Выберите модели по номеру артикула")

	return c.st.CreateOrder(ids)
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

type customerOrders []bt.Order

func (orders customerOrders) String() (s string) {
	var ctime, rtime string

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

type employeeOrders []bt.Order

func (orders employeeOrders) String() (s string) {
	var ctime, rtime string

	s = fmt.Sprintf(
		"  # Статус заказа %9s %19s %19s %s\n",
		"Сумма", "Создан", "Выдан", "# Клиента",
	)

	for _, o := range orders {
		ctime = o.LocalCreateDate().Format(dateFormat)

		if o.LocalReleaseDate() != time.Unix(0, 0) {
			rtime = o.LocalReleaseDate().Format(dateFormat)
		} else {
			rtime = "---"
		}

		s += fmt.Sprintf("%3d %13s %9.2f %19s %19s %9d\n",
			o.Id, o.Status, o.TotalPrice, ctime, rtime, o.C_id,
		)
	}

	return s
}

type (
	OrderItems []bt.OrderItem
	Model      bt.Model
	Models     []bt.Model
)

func (orderItems OrderItems) String() (s string) {
	var sum float64 = 0.0

	for i, oi := range orderItems {
		s += fmt.Sprintln("Позиция:", i+1)
		s += Model(oi.Model).String()
		sum += oi.UnitPrice
	}
	s += fmt.Sprintln("Общая стоимость заказа:", sum)

	return s
}

func (model Model) String() (s string) {
	s += fmt.Sprintf("%s (Артикул %d):\n", model.Title, model.Id)

	for _, mat := range model.Materials {
		s += fmt.Sprintf("\t%s стоимостью %2.2f за погонный метр длиной %2.2f метра\n",
			mat.Title, mat.Price, model.MatLeng[mat.Id],
		)
	}
	s += fmt.Sprintf("\tCтоимость изготовления %2.2f\n\n", model.Price)

	return s
}

func (models Models) String() (s string) {
	for _, m := range models {
		s += Model(m).String()
	}

	return s
}
