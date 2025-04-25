package web

import (
	"log"
	"net/http"
	"strconv"

	bt "github.com/gh0st17/studio/basic_types"

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
			web.alert(c, http.StatusInternalServerError, registrationError)
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
		switch action {
		case "Создать заказ":
			c.Redirect(http.StatusSeeOther, "/create-order")
		case "Просмотреть заказы":
			c.Redirect(http.StatusSeeOther, "/orders")
		case "Выход":
			sidCookie, _ := c.Cookie("session_id")
			web.deleteSession(sidCookie, c)
			c.Redirect(http.StatusSeeOther, "/login")
		}
		return
	}

	var opt []string

	user := web.userFromCookies(c)
	if user.AccLevel == bt.CUSTOMER {
		opt = customerOptions()
	} else {
		opt = operatorOptions()
	}

	data := struct {
		UserName  string
		MenuItems []string
	}{
		UserName:  user.UserFullName,
		MenuItems: opt,
	}

	c.HTML(http.StatusOK, "main.html", data)
}

func (web *Web) ordersHandler(c *gin.Context) {
	if !web.sessionExists(c) {
		web.alert(c, http.StatusServiceUnavailable, "Сервис временно недоступен")
		return
	}

	user := web.userFromCookies(c)
	if c.Request.Method == http.MethodPost {
		actionForm := c.PostForm("action")
		orderIdForm := c.PostForm("order_id")
		customerIdForm := c.PostForm("c_id")

		var (
			customerId uint64
			orderId    uint64
			err        error
		)

		if customerIdForm == "" {
			customerId = uint64(user.Id)
		} else {
			customerId, err = strconv.ParseUint(customerIdForm, 10, 32)
			if err != nil {
				web.alert(c, http.StatusBadRequest, wrongClientID)
				return
			}
		}

		orderId, err = strconv.ParseUint(orderIdForm, 10, 32)
		if err != nil {
			web.alert(c, http.StatusBadRequest, wrongOrderID)
			return
		}

		orders := loadOrders(web, user)
		err = func() error {
			switch actionForm {
			case "process":
				return web.st.ProcessOrder(user, uint(orderId))
			case "release":
				return web.st.ReleaseOrder(user, uint(orderId))
			case "cancel":
				return web.st.CancelOrder(user, uint(orderId), orders)
			default:
				return nil
			}
		}()

		if err != nil {
			web.alert(c, http.StatusInternalServerError, err.Error())
			log.Println("change status error:", err)
			return
		}

		invalidateOrdersCache(web, uint(customerId))
	}

	orders := loadOrders(web, user)
	if len(orders) == 0 {
		web.alert(c, http.StatusOK, emptyOrders)
		return
	}

	if user.AccessLevel() == bt.CUSTOMER {
		c.HTML(http.StatusOK, "orders.html", gin.H{"Orders": orders})
	} else {
		c.HTML(http.StatusOK, "orders-operator.html", gin.H{"Orders": orders})
	}
}

func (web *Web) orderItemsHandler(c *gin.Context) {
	if !web.sessionExists(c) {
		web.alert(c, http.StatusServiceUnavailable, "Сервис временно недоступен")
		return
	}

	orderIdGet := c.Query("id")

	if orderIdGet == "" {
		web.alert(c, http.StatusBadRequest, missingOrderID)
		return
	}

	orderId, err := strconv.ParseUint(orderIdGet, 10, 32)
	if err != nil {
		web.alert(c, http.StatusBadRequest, wrongOrderID)
		return
	}

	orderItems := web.loadOrderItems(uint(orderId), c)
	if orderItems == nil {
		web.alert(c, http.StatusInternalServerError, errLoadOrderItems)
		return
	}

	var totalPrice float64
	for _, item := range orderItems {
		totalPrice += item.UnitPrice
	}

	c.HTML(http.StatusOK, "order-items.html", gin.H{
		"OrderItems": orderItems,
		"TotalPrice": totalPrice,
	})
}

func (web *Web) createOrderHandler(c *gin.Context) {
	user := web.userFromCookies(c)

	if user.AccLevel == bt.OPERATOR {
		c.Redirect(http.StatusSeeOther, "/")
	}

	if c.Request.Method == http.MethodPost {
		if !web.sessionExists(c) {
			web.alert(c, http.StatusServiceUnavailable, "Сервис временно недоступен")
			return
		}

		if err := c.Request.ParseForm(); err != nil {
			web.alert(c, http.StatusBadRequest, err.Error())
			log.Println("create order parsing error:", err)
			return
		}

		modelIdsStr := c.PostFormArray("model_ids")
		modelIds := make([]uint, len(modelIdsStr))
		for i, mid := range modelIdsStr {
			if id, err := strconv.ParseUint(mid, 10, 32); err == nil {
				modelIds[i] = uint(id)
			}
		}

		err := web.st.CreateOrder(user, modelIds)
		if err != nil {
			web.alert(c, http.StatusBadRequest, err.Error())
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
			web.alert(c, http.StatusBadRequest, wrongModelID)
			log.Println("view model error:", err)
			return
		}

		if model, ok := web.st.Models()[uint(mid)]; ok {
			c.HTML(http.StatusOK, "model.html", model)
		} else {
			web.alert(c, http.StatusBadRequest, modelNotFound)
		}

		return
	}

	c.Redirect(http.StatusSeeOther, "/create-order")
}
