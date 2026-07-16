package api

import (
	"encoding/json"
	"net/http"
	"strings"
	"time"

	"git.sapienzaapps.it/fantasticcoffee/fantastic-coffee-decaffeinated/service/api/reqcontext"
	"git.sapienzaapps.it/fantasticcoffee/fantastic-coffee-decaffeinated/service/database"
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
	_ = json.NewEncoder(w).Encode(sentMessage)
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
	_ = json.NewEncoder(w).Encode(updatedMessage)
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
	_ = json.NewEncoder(w).Encode(forwardedMessage)
}

func (rt *_router) updateMessageStatus(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
	messageId := ps.ByName("messageId")

	var body struct {
		Status string `json:"status"`
	}

	err := json.NewDecoder(r.Body).Decode(&body)
	if err != nil {
		http.Error(w, "Formato JSON non valido", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	updatedAt, err := rt.db.UpdateMessageStatus(messageId, database.MessageState(body.Status))
	if err != nil {
		if strings.Contains(err.Error(), "messaggio non trovato") {
			http.Error(w, "Messaggio non trovato", http.StatusNotFound)
			return
		}
		if strings.Contains(err.Error(), "transizione") || strings.Contains(err.Error(), "non valido") {
			http.Error(w, "Transizione di stato non valida", http.StatusBadRequest)
			return
		}
		http.Error(w, "Errore interno del server", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(map[string]time.Time{"updatedAt": updatedAt})
}
