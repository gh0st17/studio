package web

import (
	"log"
	"net/http"
	"strconv"
	bt "studio/basic_types"
	"time"

	"github.com/google/uuid"
)

func (web *Web) loginHandler(w http.ResponseWriter, r *http.Request) {
	cookie, _ := r.Cookie("session_id")
	if cookie != nil && web.isSessionExists(cookie.Value) {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	web.execTemplate("login.html", w, nil)
}

func (web *Web) doLoginHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		if err := r.ParseForm(); err != nil {
			http.Error(w, "ParseForm() error", http.StatusBadRequest)
			return
		}

		login := r.FormValue("login")
		entity, err := web.st.Login(login)

		if err != nil {
			web.execTemplate("alert.html", w, struct{ Msg string }{err.Error()})
			log.Println("login error:", err)
			return
		}

		sessionID := uuid.New().String()
		web.sessionMutex.Lock()
		web.sessionStore[sessionID] = entity
		web.sessionMutex.Unlock()

		http.SetCookie(w, &http.Cookie{
			Name:  "session_id",
			Value: sessionID,
			Path:  "/",
		})
	}

	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func (web *Web) registerHandler(w http.ResponseWriter, r *http.Request) {
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
	if r.Method == http.MethodPost {
		action := r.FormValue("action")
		cookie, _ := r.Cookie("session_id")
		sessionId := cookie.Value

		switch action {
		case "Создать заказ":
			http.Redirect(w, r, "/create-order", http.StatusSeeOther)
		case "Просмотреть заказы":
			http.Redirect(w, r, "/orders", http.StatusSeeOther)
		case "Выход":
			web.sessionMutex.Lock()
			delete(web.sessionStore, sessionId)
			web.sessionMutex.Unlock()
			http.Redirect(w, r, "/login", http.StatusSeeOther)
		}

		return
	}

	cookie, err := r.Cookie("session_id")
	if err != nil || !web.isSessionExists(cookie.Value) {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	web.sessionMutex.RLock()
	entity := web.sessionStore[cookie.Value]
	web.sessionMutex.RUnlock()
	if entity == nil {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

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
	cookie, err := r.Cookie("session_id")
	if err != nil || !web.isSessionExists(cookie.Value) {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	sessionID := cookie.Value
	web.sessionMutex.RLock()
	ent := web.sessionStore[sessionID]
	web.sessionMutex.RUnlock()

	if r.Method == http.MethodPost {
		action := r.FormValue("action")
		orderId := r.FormValue("order_id")

		oId, _ := strconv.ParseUint(orderId, 10, 32)

		err := func() error {
			switch action {
			case "process":
				return web.st.ProcessOrder(ent, uint(oId))
			case "release":
				return web.st.ReleaseOrder(ent, uint(oId))
			case "cancel":
				return web.st.CancelOrder(ent, uint(oId))
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

	if rawOrders, err := web.st.Orders(ent); err != nil {
		http.Error(w, "Ошибка просмотра заказов", http.StatusInternalServerError)
		log.Println("orders error:", err)
	} else {
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

		web.sessionMutex.RLock()
		accLevel := web.sessionStore[sessionID].AccessLevel()
		web.sessionMutex.RUnlock()

		if accLevel == bt.CUSTOMER {
			web.execTemplate("orders.html", w, struct{ Orders []Order }{orders})
		} else {
			web.execTemplate("orders-operator.html", w, struct{ Orders []Order }{orders})
		}
	}
}

func (web *Web) orderItemsHandler(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("session_id")
	if err != nil || !web.isSessionExists(cookie.Value) {
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

		sessionID := cookie.Value
		web.sessionMutex.RLock()
		ent := web.sessionStore[sessionID]
		web.sessionMutex.RUnlock()

		if orderItems, err := web.st.OrderItems(ent, uint(o_id)); err != nil {
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
	cookie, err := r.Cookie("session_id")
	if err != nil || !web.isSessionExists(cookie.Value) {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	if r.Method == http.MethodPost {
		sessionID := cookie.Value

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

		web.sessionMutex.RLock()
		ent := web.sessionStore[sessionID]
		web.sessionMutex.RUnlock()

		err := web.st.CreateOrder(ent, modelIds)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			log.Println("create order error:", err)
		}

		http.Redirect(w, r, "/orders", http.StatusSeeOther)
	}
	web.execTemplate("create-order.html", w, web.st.Models())
}

func (web *Web) viewModelHandler(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("session_id")
	if err != nil || !web.isSessionExists(cookie.Value) {
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
