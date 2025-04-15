package web

import (
	"log"
	"net/http"
	"strconv"
	bt "studio/basic_types"

	"github.com/gin-gonic/gin"
)

func (web *Web) registerHandler(c *gin.Context) {
	if c.Request.Method == http.MethodPost {
		customer := bt.Customer{
			FirstName: c.PostForm("first_name"),
			LastName:  c.PostForm("last_name"),
			Login:     c.PostForm("login"),
		}

		if err := web.st.Registration(customer); err != nil {
			c.String(http.StatusInternalServerError, "Ошибка регистрации")
			log.Println("registration error:", err)
			return
		}

		c.Redirect(http.StatusSeeOther, "/login")
		return
	}

	c.HTML(http.StatusOK, "register.html", nil)
}

func (web *Web) mainHandler(c *gin.Context) {
	if c.Request.Method == http.MethodPost {
		action := c.PostForm("action")
		sessionCookie, _ := c.Cookie("session_id")
		switch action {
		case "Создать заказ":
			c.Redirect(http.StatusSeeOther, "/create-order")
		case "Просмотреть заказы":
			c.Redirect(http.StatusSeeOther, "/orders")
		case "Выход":
			web.rdb.Del(web.ctx, "session:"+sessionCookie)
			c.SetCookie("session_id", "", -1, "/", "", false, true)
			c.Redirect(http.StatusSeeOther, "/login")
		}
		return
	}

	entity := web.entityFromSession(c)

	var opt []string
	if entity.AccessLevel() == bt.CUSTOMER {
		opt = customerOptions()
	} else {
		opt = operatorOptions()
	}

	data := struct {
		UserName  string
		MenuItems []string
	}{
		UserName:  entity.FullName(),
		MenuItems: opt,
	}

	c.HTML(http.StatusOK, "main.html", data)
}

func (web *Web) ordersHandler(c *gin.Context) {
	entity := web.entityFromSession(c)

	if c.Request.Method == http.MethodPost {
		action := c.PostForm("action")
		orderId := c.PostForm("order_id")

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
			c.HTML(
				http.StatusInternalServerError,
				"alert.html",
				gin.H{"Msg": err.Error()},
			)
			log.Println("change status error:", err)
			return
		}
	}

	rawOrders, err := web.st.Orders(entity)
	if err != nil {
		c.HTML(
			http.StatusInternalServerError,
			"alert.html",
			gin.H{"Msg": err.Error()},
		)
		log.Println("orders error:", err)
		return
	}

	if len(rawOrders) == 0 {
		c.HTML(
			http.StatusOK,
			"alert.html",
			gin.H{
				"Msg": "Вы еще не сделали ни одного заказа",
			},
		)
		return
	}

	type Order struct {
		bt.Order
		CreateDate  string
		ReleaseDate string
		IsPending   bool
		Released    bool
		Processed   bool
		IsCanceled  bool
	}

	var orders []Order
	for _, rawO := range rawOrders {
		releaseDate := func() string {
			if rawO.ReleaseDate != 0 {
				return rawO.LocalReleaseDate().Format(bt.DateFormat)
			} else {
				return "---"
			}
		}()

		o := Order{
			Order:       rawO,
			CreateDate:  rawO.LocalCreateDate().Format(bt.DateFormat),
			ReleaseDate: releaseDate,
			IsPending:   rawO.Status == bt.Pending,
			Released:    rawO.Status == bt.Released,
			Processed:   rawO.Status == bt.Processing,
			IsCanceled:  rawO.Status == bt.Canceled,
		}

		orders = append(orders, o)
	}

	if entity.AccessLevel() == bt.CUSTOMER {
		c.HTML(http.StatusOK, "orders.html", gin.H{"Orders": orders})
	} else {
		c.HTML(http.StatusOK, "orders-operator.html", gin.H{"Orders": orders})
	}
}

func (web *Web) orderItemsHandler(c *gin.Context) {
	if c.Request.Method == http.MethodGet {
		orderID := c.Query("id")

		if orderID == "" {
			c.String(http.StatusBadRequest, "Не указан ID заказа")
			return
		}

		o_id, err := strconv.ParseUint(orderID, 10, 32)
		if err != nil {
			c.String(http.StatusBadRequest, "ID заказа указан неверно")
			return
		}

		entity := web.entityFromSession(c)
		if orderItems, err := web.st.OrderItems(entity, uint(o_id)); err != nil {
			c.HTML(http.StatusForbidden, "alert.html", gin.H{"Msg": err.Error()})
			log.Println("orders items error:", err)
		} else {
			var totalPrice float64
			for _, item := range orderItems {
				totalPrice += item.UnitPrice
			}

			c.HTML(http.StatusOK, "order-items.html", gin.H{
				"OrderItems": orderItems,
				"TotalPrice": totalPrice,
			})
		}
	} else {
		c.Redirect(http.StatusSeeOther, "/")
	}
}

func (web *Web) createOrderHandler(c *gin.Context) {
	entity := web.entityFromSession(c)

	if entity.AccessLevel() == bt.OPERATOR {
		c.Redirect(http.StatusSeeOther, "/")
	}

	if c.Request.Method == http.MethodPost {
		if err := c.Request.ParseForm(); err != nil {
			c.String(http.StatusBadRequest, "Ошибка парсинга формы")
			log.Println("create order parsing error:", err)
			return
		}

		modelsIds := c.PostFormArray("model_ids")
		var modelIds []uint
		for _, mid := range modelsIds {
			if id, err := strconv.ParseUint(mid, 10, 32); err == nil {
				modelIds = append(modelIds, uint(id))
			}
		}

		err := web.st.CreateOrder(entity, modelIds)
		if err != nil {
			c.String(http.StatusBadRequest, err.Error())
			log.Println("create order error:", err)
			return
		}

		c.Redirect(http.StatusSeeOther, "/orders")
		return
	}

	c.HTML(http.StatusOK, "create-order.html", web.st.Models())
}

func (web *Web) viewModelHandler(c *gin.Context) {
	if c.Request.Method == http.MethodGet {
		modelID := c.Query("id")
		mid, err := strconv.ParseUint(modelID, 10, 32)
		if err != nil {
			c.String(http.StatusBadRequest, err.Error())
			log.Println("view model error:", err)
			return
		}

		model := web.st.Models()[uint(mid)]

		c.HTML(http.StatusOK, "model.html", model)
		return
	}

	c.Redirect(http.StatusSeeOther, "/create-order")
}
