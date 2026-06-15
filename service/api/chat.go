package api

import (
	"encoding/json"
	"net/http"

	"git.sapienzaapps.it/fantasticcoffee/fantastic-coffee-decaffeinated/service/api/reqcontext"
	"github.com/julienschmidt/httprouter"
)

func (rt *_router) createChat(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
	var body struct {
		TargetUserId string `json:"targetUserId"`
	}

	err := json.NewDecoder(r.Body).Decode(&body)
	if err != nil {
		http.Error(w, "Formato JSON non valido", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	newChat, err := rt.db.CreateChat(ctx.UserId, body.TargetUserId)
	if err != nil {
		// 500
		http.Error(w, "Errore durante la creazione della chat", http.StatusInternalServerError)
		return
	}


	w.Header().Set("Content-Type", "application/json")
	// 201
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(newChat)
}