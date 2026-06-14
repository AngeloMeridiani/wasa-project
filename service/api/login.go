package api

import (
	"database/sql"
	"encoding/json"
	"errors"
	"net/http"

	"git.sapienzaapps.it/fantasticcoffee/fantastic-coffee-decaffeinated/service/api/reqcontext"
	"github.com/julienschmidt/httprouter"
)

func (rt *_router) doLogin(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
	if r.Method != http.MethodPost {
		http.Error(w, "Metodo non consentito", http.StatusMethodNotAllowed)
		return
	}

	var richiesta LoginRequest
	err := json.NewDecoder(r.Body).Decode(&richiesta)
	if err != nil {
		http.Error(w, "JSON non valido", http.StatusBadRequest)
		return
	}

	lunghezzaNome := len(richiesta.Name)
	if lunghezzaNome < 3 || lunghezzaNome > 25 {
		http.Error(w, "Il nome deve essere tra 3 e 25 caratteri", http.StatusBadRequest)
		return
	}

	var userId string
	user, err := rt.db.GetUser(richiesta.Name)
	if err == nil {
		userId = user.UserId
	} else if errors.Is(err, sql.ErrNoRows) {
		newUser, errCreate := rt.db.CreateUser(richiesta.Name)
		if errCreate != nil {
			rt.baseLogger.WithError(errCreate).Error("Errore nella creazione dell'utente")
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		userId = newUser.UserId
	} else {
		rt.baseLogger.WithError(err).Error("Errore nella richiesta dell'utente")
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
