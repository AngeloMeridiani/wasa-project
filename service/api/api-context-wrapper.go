package api

import (
	"net/http"
	"strings"

	"git.sapienzaapps.it/fantasticcoffee/fantastic-coffee-decaffeinated/service/api/reqcontext"
	"github.com/gofrs/uuid"
	"github.com/julienschmidt/httprouter"
	"github.com/sirupsen/logrus"
)

// httpRouterHandler is the signature for functions that accepts a reqcontext.RequestContext in addition to those
// required by the httprouter package.
type httpRouterHandler func(http.ResponseWriter, *http.Request, httprouter.Params, reqcontext.RequestContext)

// wrap parses the request and adds a reqcontext.RequestContext instance related to the request.
func (rt *_router) wrap(fn httpRouterHandler) func(http.ResponseWriter, *http.Request, httprouter.Params) {
	return func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
		reqUUID, err := uuid.NewV4()
		if err != nil {
			rt.baseLogger.WithError(err).Error("can't generate a request UUID")
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		var ctx = reqcontext.RequestContext{
			ReqUUID: reqUUID,
		}

		// Create a request-specific logger
		ctx.Logger = rt.baseLogger.WithFields(logrus.Fields{
			"reqid":     ctx.ReqUUID.String(),
			"remote-ip": r.RemoteAddr,
		})

		// Call the next handler in chain (usually, the handler function for the path)
		fn(w, r, ps, ctx)
	}
}
func (rt *_router) authWrap(fn httpRouterHandler) func(http.ResponseWriter, *http.Request, httprouter.Params) {
	// Restituiamo rt.wrap in modo da ereditare automaticamente il logger e l'UUID
	return rt.wrap(func(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
		// Estraiamo l'header "Authorization"
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			// 401
			http.Error(w, "Header Authorization mancante", http.StatusUnauthorized)
			return
		}
		const prefix = "Bearer "
		if !strings.HasPrefix(authHeader, prefix) {
			// 401
			http.Error(w, "Formato Authorization non valido (atteso Bearer)", http.StatusUnauthorized)
			return
		}

		token := strings.TrimPrefix(authHeader, prefix)

		// Verifichiamo che il token esista (il token coincide con l'identifier dell'utente)
		_, err := rt.db.GetUser(token)
		if err != nil {
			http.Error(w, "Token non valido o utente non trovato", http.StatusUnauthorized)
			return
		}

		// Aggiungiamo l'ID utente al contesto, così ogni handler saprà chi sta chiamando l'API
		ctx.UserId = token

		// Call the next handler in chain (usually, the handler function for the path)
		fn(w, r, ps, ctx)
	})
}
