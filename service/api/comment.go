package api

import (
	"encoding/json"
	"net/http"

	"git.sapienzaapps.it/fantasticcoffee/fantastic-coffee-decaffeinated/service/api/reqcontext"
	"github.com/julienschmidt/httprouter"
)

func (rt *_router) createComment(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
	messageId := ps.ByName("messageId")

	var body struct {
		Emoji string `json:"emoji"`
	}

	err := json.NewDecoder(r.Body).Decode(&body)
	if err != nil {
		// 400
		http.Error(w, "Formato JSON non valido", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	newComment, err := rt.db.CreateComment(ctx.UserId, messageId, body.Emoji)
	if err != nil {
		// 500
		http.Error(w, "Errore durante la creazione del commento", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	// 201
	w.WriteHeader(http.StatusCreated)
	_ = json.NewEncoder(w).Encode(newComment)
}

func (rt *_router) deleteComment(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
	commentId := ps.ByName("commentId")

	err := rt.db.DeleteComment(commentId, ctx.UserId)
	if err != nil {
		// 404
		http.Error(w, "Commento non trovato o non autorizzato", http.StatusNotFound)
		return
	}
	// 204
	w.WriteHeader(http.StatusNoContent)
}
