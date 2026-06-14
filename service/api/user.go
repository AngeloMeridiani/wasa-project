package api

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"git.sapienzaapps.it/fantasticcoffee/fantastic-coffee-decaffeinated/service/api/reqcontext"
	"github.com/julienschmidt/httprouter"
)

func (rt *_router) SetMyUserName(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {

	userId := ps.ByName("userId")

	if userId == "" {
		http.Error(w, "userId mancante", http.StatusBadRequest)
		return
	}

	var body struct {
		Name string `json:"name"`
	}
	err := json.NewDecoder(r.Body).Decode(&body)
	if err != nil {
		http.Error(w, "Formato JSON non valido", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	err = rt.db.SetUserName(userId, body.Name)
	if err != nil {
		if err.Error() == "utente non trovato" {
			// 404
			return
		}
		// 400
		return
	}

	user, err := rt.db.GetUser(userId)
	if err != nil {
		// 500
		return
	}

	// 200 OK
	json.NewEncoder(w).Encode(user)
}

func (rt *_router) setMyPhoto(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
	userId := ps.ByName("userId")

	// A. Controllo che l'utente esista. Se no errore
	_, err := rt.db.GetUser(userId)
	if err != nil {
		http.Error(w, "Utente non trovato", http.StatusNotFound)
		return
	}

	// B. Controllo Content-Type (image/png o image/jpeg)
	contentType := r.Header.Get("Content-Type")
	if contentType != "image/png" && contentType != "image/jpeg" {
		http.Error(w, "Formato non supportato. Usa image/png o image/jpeg", http.StatusBadRequest)
		return
	}

	// Check limite dimensione massima immagine
	r.Body = http.MaxBytesReader(w, r.Body, 5*1024*1024)

	// Leggo i byte dell'immagine
	fileBytes, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Errore nella lettura del file o file troppo grande", http.StatusBadRequest)
		return
	}

	fmt.Printf("Ricevuta immagine per User %s! Dimensione: %d bytes. Tipo: %s\n", userId, len(fileBytes), contentType)

	// Risposta 200 OK
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "Foto profilo aggiornata correttamente"})
}
