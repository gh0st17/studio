package web

import (
	"fmt"
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
		sessionCookie, _ := c.Cookie("session_id")
		switch action {
		case "Создать заказ":
			c.Redirect(http.StatusSeeOther, "/create-order")
		case "Просмотреть заказы":
			c.Redirect(http.StatusSeeOther, "/orders")
		case "Выход":
			if web.rdbPresent.Load() {
				web.rdb.Del(web.ctx, "session:"+sessionCookie)
			}
			c.SetCookie("session_id", "", -1, "/", "", false, true)
			c.Redirect(http.StatusSeeOther, "/login")
		}
		return
	}

	entity := web.loadEntity(c)

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

	entity := web.loadEntity(c)
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
			customerId = uint64(entity.GetId())
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

		err = func() error {
			switch actionForm {
			case "process":
				return web.st.ProcessOrder(entity, uint(orderId))
			case "release":
				return web.st.ReleaseOrder(entity, uint(orderId))
			case "cancel":
				return web.st.CancelOrder(entity, uint(orderId))
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

	var orders []Order
	if orders = loadOrders(web, entity, web.rdbPresent.Load(), c); orders == nil {
		return
	}

	if entity.AccessLevel() == bt.CUSTOMER {
		c.HTML(http.StatusOK, "orders.html", gin.H{"Orders": orders})
	} else {
		c.HTML(http.StatusOK, "orders-operator.html", gin.H{"Orders": orders})
	}
}

func (web *Web) orderItemsHandler(c *gin.Context) {
	if c.Request.Method == http.MethodGet {
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

		var orderItems []bt.OrderItem
		entity := web.loadEntity(c)
		key := fmt.Sprintf("orderItems:%d:%d", entity.GetId(), orderId)

		loadFromDB := func() bool {
			if orderItems, err = web.st.OrderItems(entity, uint(orderId)); err != nil {
				web.alert(c, http.StatusBadRequest, err.Error())
				log.Println("load orders items error:", err)
				return false
			}

			if web.rdbPresent.Load() {
				go saveToRedis(web, key, orderItems)
			}

			return true
		}

		if web.rdbPresent.Load() && redisArrayExists(web, key) {
			if orderItems, err = loadFromRedis[bt.OrderItem](web, key); err != nil {
				if !loadFromDB() {
					return
				}
			}
		} else {
			if !loadFromDB() {
				return
			}
		}

		var totalPrice float64
		for _, item := range orderItems {
			totalPrice += item.UnitPrice
		}

		c.HTML(http.StatusOK, "order-items.html", gin.H{
			"OrderItems": orderItems,
			"TotalPrice": totalPrice,
		})
	} else {
		c.Redirect(http.StatusSeeOther, "/")
	}
}

func (web *Web) createOrderHandler(c *gin.Context) {
	entity := web.loadEntity(c)
	if entity.AccessLevel() == bt.OPERATOR {
		c.Redirect(http.StatusSeeOther, "/")
	}

	if c.Request.Method == http.MethodPost {
		if err := c.Request.ParseForm(); err != nil {
			web.alert(c, http.StatusBadRequest, err.Error())
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
