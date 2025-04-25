package web

import (
	"log"
	"net/http"

	bt "github.com/gh0st17/studio/basic_types"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/vmihailenco/msgpack/v5"
)

func (web *Web) loginHandler(c *gin.Context) {
	loginCookie, _ := c.Cookie("login")
	var login string
	if loginCookie != "" {
		login = loginCookie
	}

	c.HTML(http.StatusOK, "login.html", gin.H{"Login": login})
}

func (web *Web) doLoginHandler(c *gin.Context) {
	if c.Request.Method == http.MethodPost {
		if err := c.Request.ParseForm(); err != nil {
			web.alert(c, http.StatusBadRequest, err.Error())
			return
		}

		const expTime int = 3600 * 24
		sessionID := uuid.New().String()

		if web.rdbPresent.Load() {
			web.addSession("session:" + sessionID)
		} else {
			web.alert(c, http.StatusServiceUnavailable, "Сервис временно недоступен")
			c.Abort()
			return
		}

		login := c.PostForm("login")
		entity, err := web.st.Login(login)
		if entity == nil || err != nil {
			web.alert(c, http.StatusInternalServerError, err.Error())
			log.Println("login error:", err)
			return
		}

		id := func() uint {
			if entity.AccessLevel() == bt.CUSTOMER {
				return entity.GetId()
			} else {
				return 0
			}
		}

		user := &User{
			Id:           id(),
			Login:        login,
			AccLevel:     entity.AccessLevel(),
			UserFullName: entity.FullName(),
		}
		value, _ := msgpack.Marshal(user)

		c.SetCookie("session_id", sessionID, expTime, "/", "", false, true)
		c.SetCookie("session_data", string(value), expTime, "/", "", false, true)
		c.Redirect(http.StatusSeeOther, "/")
		return
	}

	c.Redirect(http.StatusSeeOther, "/login")
}
