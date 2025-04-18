package web

import (
	"encoding/json"
	"fmt"
	"time"
)

func saveToRedis[T any](web *Web, key string, data []T) error {
	bytes, err := json.Marshal(data)
	if err != nil {
		return err
	}
	return web.rdb.Set(web.ctx, key, bytes, time.Hour).Err()
}

func loadFromRedis[T any](web *Web, key string) ([]T, error) {
	val, err := web.rdb.Get(web.ctx, key).Result()
	if err != nil {
		return nil, err
	}
	var out []T
	err = json.Unmarshal([]byte(val), &out)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func redisArrayExists(web *Web, key string) (bool, error) {
	exists, err := web.rdb.Exists(web.ctx, key).Result()
	if err != nil {
		return false, err
	}
	return exists == 1, nil
}

func invalidateOrdersCache(web *Web, customer_id uint) {
	web.rdb.Del(web.ctx, "orders:0")
	web.rdb.Del(web.ctx, fmt.Sprintf("orders:%d", customer_id))
}
