package web

import (
	"fmt"
	"log"

	bt "github.com/gh0st17/studio/basic_types"

	"github.com/gin-gonic/gin"
)

func loadOrders(web *Web, entity bt.Entity) []bt.Order {
	var (
		orders []bt.Order
		key    string = "orders:0"
		err    error
	)

	if entity.AccessLevel() == bt.CUSTOMER {
		key = "orders:" + fmt.Sprint(entity.GetId())
	}

	if web.rdbPresent.Load() {
		orders, err = loadFromRedis[bt.Order](web, key)
		if err == nil {
			return orders
		}
	}

	orders, err = web.st.Orders(entity)
	if err != nil {
		log.Println("load orders error:", err)
		return nil
	}

	if len(orders) > 0 && web.rdbPresent.Load() {
		go saveToRedis(web, key, orders)
	}

	return orders
}

func (web *Web) loadOrderItems(orderId uint, c *gin.Context) (orderItems []bt.OrderItem) {
	entity := web.loadEntity(c)
	key := fmt.Sprintf("orderItems:%d:%d", entity.GetId(), orderId)

	var err error

	if web.rdbPresent.Load() && redisArrayExists(web, key) {
		orderItems, err = loadFromRedis[bt.OrderItem](web, key)
		if err == nil {
			return orderItems
		}
	}

	orders := loadOrders(web, entity)

	orderItems, err = web.st.OrderItems(entity, uint(orderId), orders)
	if err != nil {
		log.Println("load orders items error:", err)
		return nil
	}

	if web.rdbPresent.Load() {
		go saveToRedis(web, key, orderItems)
	}

	return orderItems
}
