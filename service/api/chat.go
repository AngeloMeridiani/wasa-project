package api

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"
	"time"

	"git.sapienzaapps.it/fantasticcoffee/fantastic-coffee-decaffeinated/service/api/reqcontext"
	"git.sapienzaapps.it/fantasticcoffee/fantastic-coffee-decaffeinated/service/database"
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
	_ = json.NewEncoder(w).Encode(newChat)
}

func (rt *_router) getMyConversations(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
	chats, err := rt.db.GetMyChats(ctx.UserId)
	if err != nil {
		http.Error(w, "Errore durante il recupero delle chat", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(map[string]interface{}{"chats": chats})
}

func (rt *_router) getConversation(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
	chatId := ps.ByName("chatId")
	limit := 10
	if limitStr := r.URL.Query().Get("limit"); limitStr != "" {
		if val, err := strconv.Atoi(limitStr); err == nil && val > 0 {
			limit = val
		}
	}

	chat, messages, err := rt.db.GetCurrentChat(chatId, ctx.UserId, limit)
	if err != nil {
		if strings.Contains(err.Error(), "chat non trovata") || strings.Contains(err.Error(), "non autorizzato") {
			http.Error(w, "Chat non trovata o accesso negato", http.StatusNotFound)
			return
		}
		http.Error(w, "Errore interno del server", http.StatusInternalServerError)
		return
	}

	if messages == nil {
		messages = []database.Message{}
	}

	response := struct {
		ChatId     string             `json:"chatId"`
		GroupId    *string            `json:"groupId,omitempty"`
		Title      string             `json:"title,omitempty"`
		CreatedAt  time.Time          `json:"createdAt"`
		ModifiedAt time.Time          `json:"modifiedAt"`
		Messages   []database.Message `json:"messages"`
	}{
		ChatId:     chat.ChatId,
		GroupId:    chat.GroupId,
		Title:      chat.Title,
		CreatedAt:  chat.CreatedAt,
		ModifiedAt: chat.ModifiedAt,
		Messages:   messages,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(response)
}

func (rt *_router) deleteChat(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
	chatId := ps.ByName("chatId")

	err := rt.db.DeleteChat(chatId, ctx.UserId)
	if err != nil {
		if strings.Contains(err.Error(), "non trovata") {
			http.Error(w, "Chat non trovata", http.StatusNotFound)
			return
		}
		http.Error(w, "Errore interno durante l'eliminazione della chat", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
