package web

import (
	"fmt"
	"log"
	"strconv"
	bt "studio/basic_types"
	"time"

	"github.com/gin-gonic/gin"
)

func (web *Web) allCookiesExists(c *gin.Context) bool {
	_, loginErr := c.Cookie("login")
	_, sessionErr := c.Cookie("session_id")
	return loginErr == nil && sessionErr == nil
}

func (web *Web) addSession(login, sessionID string) {
	entity, _ := web.st.Login(login)
	err := web.rdb.HSet(web.ctx, "session:"+sessionID, map[string]interface{}{
		"id":       fmt.Sprint(entity.GetId()),
		"login":    entity.GetLogin(),
		"fullname": entity.FullName(),
		"acclevel": uint(entity.AccessLevel()),
	}).Err()
	if err != nil {
		log.Println("add from cookies error:", err)
	}
	err = web.rdb.Expire(web.ctx, "session:"+sessionID, 24*time.Hour).Err()
	if err != nil {
		log.Println("set expire error:", err)
	}
}

func (web *Web) entityFromSession(c *gin.Context) (entity bt.Entity) {
	if web.allCookiesExists(c) {
		loginCookie, _ := c.Cookie("login")
		sessionCookie, _ := c.Cookie("session_id")
		sessionKey := "session:" + sessionCookie

		result, err := web.rdb.HGetAll(web.ctx, sessionKey).Result()
		if err != nil {
			log.Println("reading redis error:", err)
			return nil
		}

		if len(result) == 0 {
			web.addSession(loginCookie, sessionCookie)

			result, err = web.rdb.HGetAll(web.ctx, sessionKey).Result()
			if err != nil || len(result) == 0 {
				log.Println("session not found after add:", err)
				return nil
			}
		}

		id, err := strconv.Atoi(result["id"])
		if err != nil {
			log.Println("invalid id in session:", err)
			return nil
		}

		accessLevel, err := strconv.Atoi(result["acclevel"])
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
