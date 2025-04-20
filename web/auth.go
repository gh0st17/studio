package web

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
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
			c.String(http.StatusBadRequest, "ParseForm() error")
			return
		}

		login := c.PostForm("login")
		entity, err := web.st.Login(login)
		if entity == nil || err != nil {
			web.alert(c, http.StatusInternalServerError, err.Error())
			log.Println("login error:", err)
			return
		}

		sessionID := uuid.New().String()

		if web.rdbPresent.Load() {
			web.addSession(entity, "session:"+sessionID)
		}

		c.SetCookie("session_id", sessionID, 3600*24, "/", "", false, true)
		c.SetCookie("login", login, 0, "/", "", false, true)
		c.Redirect(http.StatusSeeOther, "/")
		return
	}

	c.Redirect(http.StatusSeeOther, "/login")
}
