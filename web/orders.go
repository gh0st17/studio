package web

import (
	"fmt"
	"log"
	"time"

	bt "github.com/gh0st17/studio/basic_types"

	"github.com/gin-gonic/gin"
)

type Order struct {
	bt.Order
	CreateDate  string
	ReleaseDate string
}

func loadOrders(web *Web, entity bt.Entity) []Order {
	var (
		orders []Order
		key    string = "orders:0"
	)

	if entity.AccessLevel() == bt.CUSTOMER {
		key = "orders:" + fmt.Sprint(entity.GetId())
	}

	if web.rdbPresent.Load() {
		var err error
		orders, err = loadFromRedis[Order](web, key)
		if err == nil {
			return orders
		}
	}

	rawOrders, err := web.st.Orders(entity)
	if err != nil {
		log.Println("load orders error:", err)
		return nil
	}
	orders = transformOrders(rawOrders)

	if len(orders) > 0 && web.rdbPresent.Load() {
		go saveToRedis(web, key, orders)
	}

	return orders
}

func transformOrders(rawOrders []bt.Order) []Order {
	var orders []Order
	for _, rawO := range rawOrders {
		releaseDate := func() string {
			if rawO.ReleaseDate.Equal(time.Date(1970, 1, 1, 0, 0, 0, 0, time.UTC)) {
				return "---"
			} else {
				return rawO.ReleaseDate.Format(bt.DateFormat)
			}
		}()

		o := Order{
			Order:       rawO,
			CreateDate:  rawO.CreateDate.Format(bt.DateFormat),
			ReleaseDate: releaseDate,
		}

		orders = append(orders, o)
	}

	return orders
}

func (web *Web) loadOrderItems(orderId uint, c *gin.Context) (orderItems []bt.OrderItem) {
	entity := web.loadEntity(c)
	key := fmt.Sprintf("orderItems:%d:%d", entity.GetId(), orderId)

	var err error

	if web.rdbPresent.Load() && redisArrayExists(web, key) {
		orderItems, err = loadFromRedis[bt.OrderItem](web, key)
		if err != nil {
			return orderItems
		}
	}

	orderItems, err = web.st.OrderItems(entity, uint(orderId))
	if err != nil {
		log.Println("load orders items error:", err)
		return nil
	}

	if web.rdbPresent.Load() {
		go saveToRedis(web, key, orderItems)
	}

	return orderItems
}
