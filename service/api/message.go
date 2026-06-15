package api

import (
	"encoding/json"
	"net/http"

	"git.sapienzaapps.it/fantasticcoffee/fantastic-coffee-decaffeinated/service/api/reqcontext"
	"github.com/julienschmidt/httprouter"
)

func (rt *_router) sendMessage(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
	var body struct {
		ChatId string `json:"chatId"`
		Text   string `json:"text"`
	}

	err := json.NewDecoder(r.Body).Decode(&body)
	if err != nil {
		// 400
		http.Error(w, "Formato JSON non valido", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	sentMessage, err := rt.db.SendMessage(ctx.UserId, body.ChatId, body.Text)
	if err != nil {
		// 500
		http.Error(w, "Errore interno del server", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	// 201
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(sentMessage)
}

func (rt *_router) updateMessage(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
	messageId := ps.ByName("messageId")

	var body struct {
		Text string `json:"text"`
	}
	
	err := json.NewDecoder(r.Body).Decode(&body)
	if err != nil {
		// 400
		http.Error(w, "Formato JSON non valido", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	updatedMessage, err := rt.db.UpdateMessage(messageId, body.Text, ctx.UserId)
	if err != nil {
		// 404
		http.Error(w, "Messaggio non trovato o non autorizzato", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	// 200
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(updatedMessage)
}

func (rt *_router) deleteMessage(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
	messageId := ps.ByName("messageId")

	err := rt.db.DeleteMessage(messageId, ctx.UserId)
	if err != nil {
		// 404
		http.Error(w, "Messaggio non trovato o non autorizzato", http.StatusNotFound)
		return
	}
	// 204
	w.WriteHeader(http.StatusNoContent)
}

func (rt *_router) forwardMessage(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
	messageId := ps.ByName("messageId")
	var body struct {
		TargetChatId string `json:"targetChatId"`
	}

	err := json.NewDecoder(r.Body).Decode(&body)
	if err != nil {
		// 400
		http.Error(w, "Formato JSON non valido", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	forwardedMessage, err := rt.db.ForwardMessage(messageId, body.TargetChatId, ctx.UserId)
	if err != nil {
		// 404
		http.Error(w, "Messaggio non trovato", http.StatusNotFound)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	// 200
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(forwardedMessage)
}