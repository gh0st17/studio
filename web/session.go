package web

import (
	"fmt"
	"log"
	"net/http"
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

func (web *Web) entityFromSession(c *gin.Context) (entity bt.Entity) {
	if web.allCookiesExists(c) {
		loginCookie, _ := c.Cookie("login")
		sessionCookie, _ := c.Cookie("session_id")
		sessionKey := "session:" + sessionCookie

		result, err := web.rdb.HGetAll(web.ctx, sessionKey).Result()
		if err != nil {
			log.Println("reading redis error:", err)
		}

		if len(result) == 0 {
			entity, err := web.st.Login(loginCookie)
			if err != nil {
				c.HTML(
					http.StatusInternalServerError,
					"alert.html",
					gin.H{"Msg": err.Error()},
				)
				log.Println(err)
				return nil
			}

			web.addSession(entity, sessionKey)
			result, err = web.rdb.HGetAll(web.ctx, sessionKey).Result()
			if err != nil || len(result) == 0 {
				log.Println("session not found after add:", err)
			}

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
