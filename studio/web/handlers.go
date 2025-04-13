package web

import (
	"log"
	"net/http"
	"strconv"
	bt "studio/basic_types"
	"time"
)

func (web *Web) registerHandler(w http.ResponseWriter, r *http.Request) {
	if web.allCookiesExists(r) {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	if r.Method == http.MethodPost {
		customer := bt.Customer{}
		customer.FirstName = r.FormValue("first_name")
		customer.LastName = r.FormValue("last_name")
		customer.Login = r.FormValue("login")

		if err := web.st.Registration(customer); err != nil {
			http.Error(w, "Ошибка регистрации", http.StatusInternalServerError)
			log.Println("registration error:", err)
		}

		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	web.execTemplate("register.html", w, nil)
}

func (web *Web) mainHandler(w http.ResponseWriter, r *http.Request) {
	if !web.allCookiesExists(r) {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	if r.Method == http.MethodPost {
		action := r.FormValue("action")
		sessionCookie, _ := r.Cookie("session_id")
		switch action {
		case "Создать заказ":
			http.Redirect(w, r, "/create-order", http.StatusSeeOther)
		case "Просмотреть заказы":
			http.Redirect(w, r, "/orders", http.StatusSeeOther)
		case "Выход":
			web.sessionMutex.Lock()
			delete(web.sessionStore, sessionCookie.Value)
			web.sessionMutex.Unlock()
			c := &http.Cookie{
				Name:     "session_id",
				Value:    "",
				Path:     "/",
				Expires:  time.Unix(0, 0),
				HttpOnly: true,
			}
			http.SetCookie(w, c)
			http.Redirect(w, r, "/login", http.StatusSeeOther)
		}

		return
	}

	entity := web.entityFromSession(r)

	opt := func() []string {
		if entity.AccessLevel() == bt.CUSTOMER {
			return customerOptions()
		} else {
			return operatorOptions()
		}
	}()

	data := struct {
		UserName  string
		MenuItems []string
	}{
		UserName:  entity.FirstLastName(),
		MenuItems: opt,
	}

	web.execTemplate("main.html", w, data)
}

func (web *Web) ordersHandler(w http.ResponseWriter, r *http.Request) {
	if !web.allCookiesExists(r) {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	entity := web.entityFromSession(r)

	if r.Method == http.MethodPost {
		action := r.FormValue("action")
		orderId := r.FormValue("order_id")

		oId, _ := strconv.ParseUint(orderId, 10, 32)

		err := func() error {
			switch action {
			case "process":
				return web.st.ProcessOrder(entity, uint(oId))
			case "release":
				return web.st.ReleaseOrder(entity, uint(oId))
			case "cancel":
				return web.st.CancelOrder(entity, uint(oId))
			default:
				return nil
			}
		}()

		if err != nil {
			web.execTemplate("alert.html", w, struct{ Msg string }{err.Error()})
			log.Println("change status error:", err)
			return
		}
	}

	rawOrders, err := web.st.Orders(entity)
	if err != nil {
		http.Error(w, "Ошибка просмотра заказов", http.StatusInternalServerError)
		log.Println("orders error:", err)
		return
	}

	if len(rawOrders) == 0 {
		web.execTemplate("alert.html", w, struct{ Msg string }{"Вы еще не сделали ни одного заказа"})
		return
	}

	type Order struct {
		bt.Order
		CustomerName string
		EmployeeName string
		CreateDate   string
		ReleaseDate  string
		IsPending    bool
		Released     bool
		Processed    bool
		IsCanceled   bool
	}

	var orders []Order
	for _, rawO := range rawOrders {
		releaseDate := func() string {
			if rawO.ReleaseDate != 0 {
				return time.Unix(rawO.ReleaseDate, 0).Format(dateFormat)
			} else {
				return "---"
			}
		}()

		o := Order{
			Order:        rawO,
			CustomerName: web.st.FullName(rawO.C_id, bt.CUSTOMER),
			EmployeeName: web.st.FullName(rawO.E_id, bt.OPERATOR),
			CreateDate:   time.Unix(rawO.CreateDate, 0).Format(dateFormat),
			ReleaseDate:  releaseDate,
			IsPending:    rawO.Status == bt.Pending,
			Released:     rawO.Status == bt.Released,
			Processed:    rawO.Status == bt.Processing,
			IsCanceled:   rawO.Status == bt.Canceled,
		}

		orders = append(orders, o)
	}

	if entity.AccessLevel() == bt.CUSTOMER {
		web.execTemplate("orders.html", w, struct{ Orders []Order }{orders})
	} else {
		web.execTemplate("orders-operator.html", w, struct{ Orders []Order }{orders})
	}
}

func (web *Web) orderItemsHandler(w http.ResponseWriter, r *http.Request) {
	if !web.allCookiesExists(r) {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	if r.Method == http.MethodGet {
		orderID := r.URL.Query().Get("id")

		if orderID == "" {
			http.Error(w, "Не указан ID заказа", http.StatusBadRequest)
			return
		}

		o_id, err := strconv.ParseUint(orderID, 10, 32)
		if err != nil {
			http.Error(w, "ID заказа указан неверно", http.StatusBadRequest)
			return
		}

		entity := web.entityFromSession(r)
		if orderItems, err := web.st.OrderItems(entity, uint(o_id)); err != nil {
			http.Error(w, "Ошибка просмотра заказа", http.StatusInternalServerError)
			log.Println("orders items error:", err)
		} else {
			var totalPrice float64
			for _, item := range orderItems {
				totalPrice += item.UnitPrice
			}

			web.execTemplate("order-items.html", w,
				struct {
					OrderItems []bt.OrderItem
					TotalPrice float64
				}{
					OrderItems: orderItems,
					TotalPrice: totalPrice,
				},
			)
		}
	} else {
		http.Redirect(w, r, "/", http.StatusSeeOther)
	}
}

func (web *Web) createOrderHandler(w http.ResponseWriter, r *http.Request) {
	if !web.allCookiesExists(r) {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	if r.Method == http.MethodPost {
		if err := r.ParseForm(); err != nil {
			http.Error(w, "Ошибка парсинга формы", http.StatusBadRequest)
			log.Println("create order parsing error:", err)
			return
		}

		modelsIds := r.Form["model_ids"]
		var modelIds []uint
		for _, mid := range modelsIds {
			if id, err := strconv.ParseUint(mid, 10, 32); err == nil {
				modelIds = append(modelIds, uint(id))
			}
		}

		entity := web.entityFromSession(r)
		err := web.st.CreateOrder(entity, modelIds)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			log.Println("create order error:", err)
		}

		http.Redirect(w, r, "/orders", http.StatusSeeOther)
	}
	web.execTemplate("create-order.html", w, web.st.Models())
}

func (web *Web) viewModelHandler(w http.ResponseWriter, r *http.Request) {
	if !web.allCookiesExists(r) {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	if r.Method == http.MethodGet {
		modelID := r.URL.Query().Get("id")
		mid, err := strconv.ParseUint(modelID, 10, 32)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			log.Println("view model error:", err)
		}

		model := web.st.Models()[uint(mid)]

		web.execTemplate("model.html", w, model)
		return
	}

	http.Redirect(w, r, "/create-order.html", http.StatusSeeOther)
}
