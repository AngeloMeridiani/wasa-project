package api

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"time"

	"git.sapienzaapps.it/fantasticcoffee/fantastic-coffee-decaffeinated/service/api/reqcontext"
	"git.sapienzaapps.it/fantasticcoffee/fantastic-coffee-decaffeinated/service/database"
	"github.com/julienschmidt/httprouter"
)

func (rt *_router) sendMessage(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
	contentType := r.Header.Get("Content-Type")

	var sentMessage *database.Message
	var err error

	if strings.HasPrefix(contentType, "multipart/form-data") {
		// Messaggio Immagine
		err = r.ParseMultipartForm(5 << 20) // 5 MB max
		if err != nil {
			http.Error(w, "Errore caricamento file", http.StatusBadRequest)
			return
		}

		chatId := r.FormValue("chatId")
		if chatId == "" {
			http.Error(w, "chatId mancante", http.StatusBadRequest)
			return
		}

		var replyToPtr *string
		replyTo := r.FormValue("replyTo")
		if replyTo != "" {
			replyToPtr = &replyTo
		}

		text := r.FormValue("text")

		file, header, err := r.FormFile("image")
		if err != nil {
			http.Error(w, "Immagine mancante", http.StatusBadRequest)
			return
		}
		defer file.Close()

		fileBytes, err := io.ReadAll(file)
		if err != nil {
			http.Error(w, "Errore lettura immagine", http.StatusBadRequest)
			return
		}

		// Determina l'estensione e il nome file
		ext := ".jpg"
		if strings.Contains(header.Header.Get("Content-Type"), "png") {
			ext = ".png"
		}

		// Usiamo il timestamp per rendere il nome file univoco, o usiamo un ID
		fileName := fmt.Sprintf("msg_%d%s", time.Now().UnixNano(), ext)
		filePath := fmt.Sprintf("./images/media/%s", fileName)
		photoUrl := fmt.Sprintf("/images/media/%s", fileName)

		err = os.WriteFile(filePath, fileBytes, 0644)
		if err != nil {
			ctx.Logger.WithError(err).Error("Errore durante il salvataggio dell'immagine del messaggio")
			http.Error(w, "Errore interno durante il salvataggio", http.StatusInternalServerError)
			return
		}

		sentMessage, err = rt.db.SendImageMessage(ctx.UserId, chatId, text, photoUrl, replyToPtr)

	} else {
		// Messaggio di Testo (JSON)
		var body struct {
			ChatId  string  `json:"chatId"`
			Text    string  `json:"text"`
			ReplyTo *string `json:"replyTo"`
		}

		err = json.NewDecoder(r.Body).Decode(&body)
		if err != nil {
			http.Error(w, "Formato JSON non valido", http.StatusBadRequest)
			return
		}
		defer r.Body.Close()

		sentMessage, err = rt.db.SendMessage(ctx.UserId, body.ChatId, body.Text, body.ReplyTo)
	}

	if err != nil {
		ctx.Logger.WithError(err).Error("Errore salvataggio messaggio nel DB")
		http.Error(w, "Errore interno del server", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
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
