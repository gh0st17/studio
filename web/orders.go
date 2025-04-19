package web

import (
	"fmt"
	"log"
	"net/http"
	"time"

	bt "github.com/gh0st17/studio/basic_types"

	"github.com/gin-gonic/gin"
)

type Order struct {
	bt.Order
	CreateDate  string
	ReleaseDate string
}

func loadOrders(web *Web, entity bt.Entity, present bool, c *gin.Context) []Order {
	var (
		orders []Order
		key    string = "orders:0"
	)

	if entity.AccessLevel() == bt.CUSTOMER {
		key = "orders:" + fmt.Sprint(entity.GetId())
	}

	if present {
		var err error
		orders, err = loadFromRedis[Order](web, key)
		if err == nil {
			return orders
		}
	}

	rawOrders, err := web.st.Orders(entity)
	if err != nil {
		web.alert(c, http.StatusInternalServerError, err.Error())
		log.Println("orders error:", err)
		return nil
	}
	orders = transformOrders(rawOrders)

	if present {
		go saveToRedis(web, key, orders)
	}

	if len(orders) == 0 {
		web.alert(c, http.StatusOK, emptyOrders)
		return nil
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
