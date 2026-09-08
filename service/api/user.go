package api

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"

	"git.sapienzaapps.it/fantasticcoffee/fantastic-coffee-decaffeinated/service/api/reqcontext"
	"git.sapienzaapps.it/fantasticcoffee/fantastic-coffee-decaffeinated/service/globaltime"
	"github.com/julienschmidt/httprouter"
)

func (rt *_router) SetMyUserName(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {

	userId := ps.ByName("userId")

	if ctx.UserId != userId {
		// 403
		http.Error(w, "Non sei autorizzato a modificare questo utente", http.StatusForbidden)
		return
	}

	var body struct {
		Name string `json:"name"`
	}

	err := json.NewDecoder(r.Body).Decode(&body)
	if err != nil {
		// 400
		http.Error(w, "Formato JSON non valido", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	err = rt.db.SetUserName(userId, body.Name)
	if err != nil {
		if err.Error() == "utente non trovato" {
			// 404
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}
		// 400
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	user, err := rt.db.GetUser(userId)
	if err != nil {
		// 500
		http.Error(w, "Errore interno del server", http.StatusInternalServerError)
		return
	}

	// 200
	_ = json.NewEncoder(w).Encode(user)
}

func (rt *_router) setMyPhoto(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
	// Estraggo l'id dell'utente dall'URL
	userId := ps.ByName("userId")

	if ctx.UserId != userId {
		// 403
		http.Error(w, "Non sei autorizzato a modificare la foto di questo utente", http.StatusForbidden)
		return
	}

	// Controllo che l'utente esista
	_, err := rt.db.GetUser(userId)
	if err != nil {
		// 404
		http.Error(w, "Utente non trovato", http.StatusNotFound)
		return
	}

	// Controllo il Content-Type
	contentType := r.Header.Get("Content-Type")
	if contentType != "image/png" && contentType != "image/jpeg" {
		http.Error(w, "Formato non supportato. Usa image/png o image/jpeg", http.StatusBadRequest)
		return
	}

	// Check limite dimensione massima dell'immagine
	r.Body = http.MaxBytesReader(w, r.Body, 5*1024*1024)

	// Leggo i byte dell'immagine
	fileBytes, err := io.ReadAll(r.Body)
	if err != nil {
		// 400
		http.Error(w, "Errore nella lettura del file o file troppo grande", http.StatusBadRequest)
		return
	}

	ctx.Logger.Infof("Ricevuta immagine per User %s! Dimensione: %d bytes. Tipo: %s\n", userId, len(fileBytes), contentType)

	// Determina l'estensione e il nome file
	ext := ".jpg"
	if contentType == "image/png" {
		ext = ".png"
	}
	fileName := fmt.Sprintf("img_user_%s%s", userId, ext)
	filePath := fmt.Sprintf("./images/pfp/%s", fileName)
	photoUrl := fmt.Sprintf("/images/pfp/%s?t=%d", fileName, globaltime.Now().UnixNano())

	if err := os.MkdirAll("./images/pfp", 0755); err != nil {
		ctx.Logger.WithError(err).Error("Errore durante la creazione della cartella immagini")
		http.Error(w, "Errore interno durante il salvataggio", http.StatusInternalServerError)
		return
	}

	// Scrivi i byte su disco
	err = os.WriteFile(filePath, fileBytes, 0644)
	if err != nil {
		// 500
		ctx.Logger.WithError(err).Error("Errore durante il salvataggio dell'immagine")
		http.Error(w, "Errore interno durante il salvataggio", http.StatusInternalServerError)
		return
	}

	// Aggiorna il DB
	err = rt.db.SetMyPhoto(userId, photoUrl)
	if err != nil {
		ctx.Logger.WithError(err).Error("Errore aggiornamento photoUrl nel DB")
		http.Error(w, "Errore interno durante l'aggiornamento del DB", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(map[string]string{"photoUrl": photoUrl, "message": "Foto profilo aggiornata correttamente"})
}

func (rt *_router) getUser(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
	// Estraggo l'id dell'utente dall'URL
	userId := ps.ByName("userId")

	if ctx.UserId != userId {
		// 403
		http.Error(w, "Non sei autorizzato ad ottenere le informazioni di questo utente", http.StatusForbidden)
		return
	}

	user, err := rt.db.GetUser(userId)
	if err != nil {
		// 404
		http.Error(w, "Utente non trovato", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(user)
}

func (rt *_router) getUsers(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
	name := r.URL.Query().Get("name")

	users, err := rt.db.GetUsersByName(name)
	if err != nil {
		http.Error(w, "Errore interno del server", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(users)
}
