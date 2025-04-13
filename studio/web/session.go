package web

import (
	"net/http"
	bt "studio/basic_types"
)

func (web *Web) allCookiesExists(r *http.Request) bool {
	loginCookie, _ := r.Cookie("login")
	sessionCookie, _ := r.Cookie("session_id")
	return loginCookie != nil && sessionCookie != nil
}

func (web *Web) addFromCookies(login, sessionID string) {
	entity, _ := web.st.Login(login)
	web.sessionMutex.Lock()
	web.sessionStore[sessionID] = entity
	web.sessionMutex.Unlock()
}

func (web *Web) entityFromSession(r *http.Request) (entity bt.Entity) {
	if web.allCookiesExists(r) {
		loginCookie, _ := r.Cookie("login")
		sessionCookie, _ := r.Cookie("session_id")
		if _, ok := web.sessionStore[sessionCookie.Value]; !ok {
			web.addFromCookies(loginCookie.Value, sessionCookie.Value)
			entity = web.sessionStore[sessionCookie.Value]
		} else {
			web.sessionMutex.RLock()
			entity = web.sessionStore[sessionCookie.Value]
			web.sessionMutex.RUnlock()
		}
	}

	return entity
}
