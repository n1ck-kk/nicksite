package lab

import (
	"net/http"

	"github.com/gorilla/sessions"
)

const sessionName = "lab-session"
const sessionKey = "authenticated"

var store *sessions.CookieStore

func InitStore(secret string, secure bool) {
	store = sessions.NewCookieStore([]byte(secret))
	store.Options = &sessions.Options{
		Path:     "/lab",
		MaxAge:   86400 * 7, // 7 days
		HttpOnly: true,
		Secure:   secure,
		SameSite: http.SameSiteLaxMode,
	}
}

func IsAuthenticated(r *http.Request) bool {
	session, err := store.Get(r, sessionName)
	if err != nil {
		return false
	}
	auth, ok := session.Values[sessionKey].(bool)
	return ok && auth
}

func SetAuthenticated(w http.ResponseWriter, r *http.Request) error {
	session, err := store.Get(r, sessionName)
	if err != nil {
		return err
	}
	session.Values[sessionKey] = true
	return session.Save(r, w)
}

func ClearSession(w http.ResponseWriter, r *http.Request) {
	session, err := store.Get(r, sessionName)
	if err != nil {
		return
	}
	session.Options.MaxAge = -1
	session.Save(r, w)
}
