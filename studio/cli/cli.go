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

	"github.com/AlecAivazis/survey/v2"
)

type CLI struct {
	st  *studio.Studio
	ent bt.Entity
	opt []string
}

func New(dbPath string) (c *CLI, err error) {
	c = &CLI{}
	if c.st, err = studio.New(dbPath); err != nil {
		return nil, err
	}

	return c, nil
}

func (c *CLI) Run() (err error) {
	for {
		c.ent, err = c.st.Login(c.Login())

		if err != nil {
			fmt.Println(err)
			continue
		} else {
			break
		}
	}

	switch c.ent.AccessLevel() {
	case bt.CUSTOMER:
		c.opt = customerOptions()
	case bt.OPERATOR:
		c.opt = operatorOptions()
	}
	fmt.Println("Запуск консольного интерфейса...")
	fmt.Printf("С возвращением, %s!\n", c.ent.FirstLastName())
	pause()

	for {
		choice := c.main()
		switch choice {
		case "Создать заказ":
			err = c.createOrder()
		case "Просмотреть заказы":
			if orders, err := c.st.Orders(c.ent); err == nil {
				c.displayOrders(orders)
			}
		case "Просмотреть содержимое заказa":
			id, _ := userinput.PromptUint("Выберите id заказа")
			if orderItems, err := c.st.OrderItems(c.ent, id[0]); err == nil {
				c.displayOrderItems(orderItems)
			}
		case "Отменить заказ":
			id, _ := userinput.PromptUint("Выберите id заказа")
			err = c.st.CancelOrder(c.ent, id[0])
		case "Выполнить заказ":
			id, _ := userinput.PromptUint("Выберите id заказа")
			err = c.st.ProcessOrder(c.ent, id[0])
		case "Выдача заказа":
			id, _ := userinput.PromptUint("Выберите id заказа")
			err = c.st.ReleaseOrder(c.ent, id[0])
		case "Выход":
			if err = c.st.Shutdown(); err != nil {
				return err
			}
			return nil
		}

		switch err {
		case nil:
			continue
		case db.ErrNotPending, db.ErrStatusRange, studio.ErrPerm:
			c.alert(fmt.Sprint(err))
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

func (c *CLI) main() (choice string) {
	clearScreen()
	prompt := &survey.Select{
		Message:  "Выберите действие:",
		Options:  c.opt,
		PageSize: 5,
	}
	survey.AskOne(prompt, &choice)

	return choice
}

func (c *CLI) displayTable(table interface{}) {
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

func (c *CLI) displayOrders(orders []bt.Order) {
	if len(orders) == 0 {
		fmt.Println("Вы еще не совершали заказов")
	}

	if c.ent.AccessLevel() == bt.CUSTOMER {
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

func (c *CLI) createOrder() (err error) {
	c.displayTable(c.st.Models())
	ids, _ := userinput.PromptUint("Выберите модели по номеру артикула")

	return c.st.CreateOrder(c.ent, ids)
}

func (c *CLI) alert(msg string) {
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
