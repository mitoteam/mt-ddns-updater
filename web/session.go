package web

import (
	"time"

	"github.com/gorilla/sessions"
	"github.com/mitoteam/mbr"
	"github.com/mitoteam/mt-ddns-updater/app"
	"github.com/mitoteam/mttools"
)

var sessionStore *sessions.CookieStore
var sessionName = "mt-ddns-updater"
var sessionTimeoutField = "session_timeout"

func GetSession(ctx *mbr.MbrContext) *sessions.Session {
	if sessionStore == nil {
		sessionStore = sessions.NewCookieStore([]byte(app.App.AppSettings.(*app.AppSettingsType).WebserverCookieSecret))
	}

	session, _ := sessionStore.Get(ctx.Request(), sessionName)

	return session
}

// returns true if the user is authenticated (session is valid and not expired). ctx - context of the request.
func IsAuthenticated(ctx *mbr.MbrContext) bool {
	session := GetSession(ctx)
	sessionTimeout := time.Now()
	if sessionTimeoutValue, ok := session.Values[sessionTimeoutField]; ok {
		if v, err := time.Parse(time.RFC3339, mttools.AnyToString(sessionTimeoutValue)); err == nil {
			sessionTimeout = v
		}
	}

	return sessionTimeout.After(time.Now())
}

func Logout(ctx *mbr.MbrContext) error {
	ctx.Set(isAuthenticatedCtxField, false)

	session := GetSession(ctx)
	delete(session.Values, sessionTimeoutField)
	return session.Save(ctx.Request(), ctx.Writer())
}
