package api

import (
	"encoding/json"
	"net/http"

	"git.sapienzaapps.it/fantasticcoffee/fantastic-coffee-decaffeinated/service/api/reqcontext"
	"github.com/julienschmidt/httprouter"
)

func (rt *_router) doLogin(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
	if r.Method != http.MethodPost {
		http.Error(w, "Metodo non consentito", http.StatusMethodNotAllowed)
		return
	}

	var loginRequest LoginRequest
	err := json.NewDecoder(r.Body).Decode(&loginRequest)
	if err != nil {
		http.Error(w, "JSON non valido", http.StatusBadRequest)
		return
	}

	lunghezzaNome := len(loginRequest.Name)
	if lunghezzaNome < 3 || lunghezzaNome > 25 {
		http.Error(w, "Il nome deve essere tra 3 e 25 caratteri", http.StatusBadRequest)
		return
	}

	userId, err := rt.db.DoLogin(loginRequest.Name)
	if err != nil {
		rt.baseLogger.WithError(err).Error("Errore durante il login/creazione utente")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	risposta := LoginResponse{
		Identifier: userId,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated) 
	json.NewEncoder(w).Encode(risposta)
}
