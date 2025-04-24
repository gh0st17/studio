package web

import (
	"bytes"
	"log"
	"time"

	"github.com/vmihailenco/msgpack/v5"

	"github.com/gin-gonic/gin"
)

func (web *Web) dataCookiesExists(c *gin.Context) bool {
	sessionDataCookie, err := c.Cookie("session_data")

	return err == nil && sessionDataCookie != ""
}

func (web *Web) addSession(sessionKey string) {
	err := web.rdb.Set(web.ctx, sessionKey, "", 24*time.Hour).Err()
	if err != nil {
		log.Println("session add error:", err)
	}
}

func (web *Web) deleteSession(sessionKey string, c *gin.Context) {
	c.SetCookie("session_id", "", -1, "/", "", false, true)
	c.SetCookie("session_data", "", -1, "/", "", false, true)

	if web.rdbPresent.Load() {
		web.rdb.Del(web.ctx, sessionKey)
	}
}

func (web *Web) sessionExists(c *gin.Context) bool {
	if web.rdbPresent.Load() {
		sessionCookie, _ := c.Cookie("session_id")
		if sessionCookie == "" {
			return false
		}

		sessionKey := "session:" + sessionCookie
		if web.rdb.Exists(web.ctx, sessionKey).Val() == 1 {
			return true
		}
	}
	return false
}

func (web *Web) userFromCookies(c *gin.Context) *User {
	sessionDataCookie, _ := c.Cookie("session_data")
	buf := bytes.NewBufferString(sessionDataCookie)

	user := &User{}
	err := msgpack.Unmarshal(buf.Bytes(), user)
	if err != nil {
		sidCookie, _ := c.Cookie("seesion_id")
		web.deleteSession(sidCookie, c)
	}

	return user
}
