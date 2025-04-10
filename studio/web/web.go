package web

import (
	"fmt"
	"net/http"
	bt "studio/basic_types"
	"studio/studio"
	"sync"
)

type Web struct {
	st       *studio.Studio
	accLevel bt.AccessLevel
	userName string
	opt      []string

	writer  http.ResponseWriter
	request *http.Request

	sessionStore map[string]bt.Entity
	sessionMutex sync.RWMutex
}

func New() *Web {
	w := &Web{}
	w.st = &studio.Studio{}
	w.sessionStore = make(map[string]bt.Entity)

	return w
}

func (w *Web) Run(dbPath string) error {
	if err := w.st.Run(dbPath); err != nil {
		return err
	}

	http.HandleFunc("/", w.homeHandler)
	http.HandleFunc("/login", w.loginHandler)
	http.HandleFunc("/do_login", w.doLoginHandler)
	http.HandleFunc("/register", w.registerHandler)
	http.HandleFunc("/main", w.mainHandler)
	http.HandleFunc("/orders", w.ordersHandler)
	http.HandleFunc("/order-items", w.orderItemsHandler)
	http.HandleFunc("/create-order", w.createOrderHandler)
	http.HandleFunc("/alert", w.alertHandler)

	fmt.Println("Запуск веб-интерфейса...")
	return http.ListenAndServe(":8080", nil)
}

func (w *Web) Login() string {
	http.Redirect(w.writer, w.request, "/login", http.StatusSeeOther)
	return ""
}

func (w *Web) Registration(login string) (customer bt.Customer) {
	// В вебе будет форма регистрации
	fmt.Println("Отображение формы регистрации на вебе")
	customer.Login = login
	customer.FirstName = "" // получить из формы
	customer.LastName = ""  // получить из формы
	return customer
}

func (w *Web) Main() string {
	// Веб-страница с основными действиями
	fmt.Println("Отображение главной страницы с действиями")
	return "" // выбрать действие из формы / url path / кнопки
}

func (w *Web) DisplayTable(table interface{}) {
	// Рендер HTML таблицы в веб-интерфейсе
	switch data := table.(type) {
	case []bt.Order:
		w.displayOrders(data)
	case []bt.OrderItem:
		w.displayOrderItems(data)
	case []bt.Model:
		w.displayModels(data)
	default:
		panic("Неизвестный тип таблицы")
	}
}

func (w *Web) displayOrders(orders []bt.Order) {
	fmt.Println("Рендеринг списка заказов в вебе")
}

func (w *Web) displayOrderItems(orderItems []bt.OrderItem) {
	fmt.Println("Рендеринг содержимого заказа в вебе")
}

func (w *Web) displayModels(models []bt.Model) {
	fmt.Println("Рендеринг списка моделей в вебе")
}

func (w *Web) ReadNumbers(prompt string) ([]uint, error) {
	fmt.Printf("Отображение формы для ввода чисел: %s\n", prompt)
	return nil, nil // здесь будет получение чисел из формы
}

func (w *Web) CreateOrder() {
	fmt.Println("Отображение формы создания заказа в вебе")
}

func (w *Web) Alert(msg string) {
	http.Redirect(w.writer, w.request,
		fmt.Sprintf("/alert?msg=%s&next=/main", msg),
		http.StatusSeeOther,
	)
}

func (w *Web) SetAccessLevel(accLevel bt.AccessLevel) {
	w.accLevel = accLevel
	switch w.accLevel {
	case bt.CUSTOMER:
		w.opt = customerOptions()
	case bt.OPERATOR:
		w.opt = operatorOptions()
	}
}

func (w *Web) SetUserName(userName string) {
	w.userName = userName
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
