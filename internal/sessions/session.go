package sessions

import (
	"net/http"
	"time"

	"github.com/alexedwards/scs/v2"
)

func NewSessionManager() *scs.SessionManager {
	s := scs.New()
	s.Lifetime = 3 * time.Hour
	s.IdleTimeout = 20 * time.Minute
	s.Cookie.Name = "session_id"
	s.Cookie.Domain = "jolt.engehost.net"
	s.Cookie.HttpOnly = true
	s.Cookie.Path = "/"
	s.Cookie.Persist = true
	s.Cookie.SameSite = http.SameSiteStrictMode
	s.Cookie.Secure = true

	return s
}
