package web

import (
	"fmt"
	"log"
	"strconv"
	"time"

	bt "github.com/gh0st17/studio/basic_types"

	"github.com/gin-gonic/gin"
)

func (web *Web) allCookiesExists(c *gin.Context) bool {
	l, _ := c.Cookie("login")
	s, _ := c.Cookie("session_id")
	return l != "" && s != ""
}

func (web *Web) addSession(entity bt.Entity, sessionKey string) {
	err := web.rdb.HSet(web.ctx, sessionKey, map[string]interface{}{
		"id":           fmt.Sprint(entity.GetId()),
		"login":        entity.GetLogin(),
		"fullname":     entity.FullName(),
		"access_level": uint(entity.AccessLevel()),
	}).Err()
	if err != nil {
		log.Println("add from cookies error:", err)
	}
	err = web.rdb.Expire(web.ctx, sessionKey, 24*time.Hour).Err()
	if err != nil {
		log.Println("set expire error:", err)
	}
}

func (web *Web) entityFromCookies(c *gin.Context) (entity bt.Entity) {
	if web.allCookiesExists(c) {
		loginCookie, _ := c.Cookie("login")
		entity, _ := web.st.Login(loginCookie)
		return entity
	}

	return entity
}

func (web *Web) entityFromRedis(c *gin.Context) (entity bt.Entity) {
	if web.allCookiesExists(c) {
		var (
			result map[string]string
			err    error
		)

		sessionCookie, _ := c.Cookie("session_id")
		sessionKey := "session:" + sessionCookie

		result, err = web.rdb.HGetAll(web.ctx, sessionKey).Result()
		if err != nil {
			log.Println("reading redis error:", err)
		}

		if len(result) == 0 {
			entity = web.entityFromCookies(c)
			go web.addSession(entity, sessionKey)

			return entity
		}

		id, err := strconv.Atoi(result["id"])
		if err != nil {
			log.Println("invalid id in session:", err)
			return nil
		}

		accessLevel, err := strconv.Atoi(result["access_level"])
		if err != nil {
			log.Println("invalid access level:", err)
			return nil
		}

		entity = &User{
			Id:           uint(id),
			Login:        result["login"],
			UserFullName: result["fullname"],
			AccLevel:     bt.AccessLevel(accessLevel),
		}
	}

	return entity
}

func (web *Web) loadEntity(c *gin.Context) bt.Entity {
	if web.rdbPresent.Load() {
		return web.entityFromRedis(c)
	} else {
		return web.entityFromCookies(c)
	}
}
