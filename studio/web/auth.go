package web

import (
	"log"
	"net/http"
	"time"

	"github.com/google/uuid"
)

func (web *Web) loginHandler(w http.ResponseWriter, r *http.Request) {
	loginCookie, _ := r.Cookie("login")
	var login string
	if loginCookie != nil {
		login = loginCookie.Value
	}

	if web.allCookiesExists(r) {
		sessionCookie, _ := r.Cookie("session_id")
		web.addFromCookies(login, sessionCookie.Value)
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	web.execTemplate("login.html", w, struct{ Login string }{login})
}

func (web *Web) doLoginHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		if err := r.ParseForm(); err != nil {
			http.Error(w, "ParseForm() error", http.StatusBadRequest)
			return
		}

		login := r.FormValue("login")
		entity, err := web.st.Login(login)

		if err != nil {
			web.execTemplate("alert.html", w, struct{ Msg string }{err.Error()})
			log.Println("login error:", err)
			return
		}

		sessionID := uuid.New().String()
		web.sessionMutex.Lock()
		web.sessionStore[sessionID] = entity
		web.sessionMutex.Unlock()

		http.SetCookie(w, &http.Cookie{
			Name:     "session_id",
			Value:    sessionID,
			Path:     "/",
			Expires:  time.Now().Add(time.Hour * 24),
			HttpOnly: true,
		})
		http.SetCookie(w, &http.Cookie{
			Name:     "login",
			Value:    login,
			Path:     "/",
			HttpOnly: true,
		})
	}

	http.Redirect(w, r, "/", http.StatusSeeOther)
}
