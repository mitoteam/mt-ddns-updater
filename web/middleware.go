package web

import (
	"net/http"

	"github.com/mitoteam/mbr"
	"github.com/mitoteam/mttools"
)

const isAuthenticatedCtxField = "IsAuthenticated"

func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		//log.Println("DBG AuthMiddleware " + r.RequestURI)

		ctx := mbr.Context(r)
		mttools.AssertNotNil(ctx, "not in MbrContext")

		if IsAuthenticated(ctx) {
			ctx.Set(isAuthenticatedCtxField, true) // set flag in context

			// Call the next handler
			next.ServeHTTP(w, r)
		} else {
			//redirect to login page
			url := mbr.Url(RootCtl.Login, "destination", r.RequestURI)
			http.Redirect(w, r, url, http.StatusSeeOther)
			// and do not call other handlers
		}
	})
}
