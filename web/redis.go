package web

import (
	"encoding/json"
	"fmt"
	"log"
	"time"
)

func saveToRedis[T any](web *Web, key string, data []T) error {
	bytes, err := json.Marshal(data)
	if err != nil {
		return err
	}

	err = web.rdb.Set(web.ctx, key, bytes, time.Hour).Err()
	if err != nil {
		return err
	}
	log.Printf("key '%s' saved to Redis", key)

	return nil
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

	log.Printf("key '%s' loaded from Redis", key)
	return out, nil
}

func redisArrayExists(web *Web, key string) bool {
	exists, _ := web.rdb.Exists(web.ctx, key).Result()
	return exists == 1
}

func invalidateOrdersCache(web *Web, customer_id uint) {
	web.rdb.Del(web.ctx, "orders:0")
	web.rdb.Del(web.ctx, fmt.Sprintf("orders:%d", customer_id))
}
