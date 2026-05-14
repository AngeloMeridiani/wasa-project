package api

import (
	"encoding/json"
	"fmt"
	"net/http"
)

// Mappa globale simulata per salvare gli utenti temporaneamente (in assenza di DB reale per questo endpoint)
var userDatabase = make(map[string]string)

func doLogin(w http.ResponseWriter, r *http.Request) {
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
	if lunghezzaNome < 3 || lunghezzaNome > 16 {
		http.Error(w, "Il nome deve essere tra 3 e 16 caratteri", http.StatusBadRequest)
		return
	}

	var idUtente string

	if idEsistente, esiste := userDatabase[richiesta.Name]; esiste {
		fmt.Printf("Bentornato %s! Recupero ID esistente.\n", richiesta.Name)
		idUtente = idEsistente
	} else {
		fmt.Printf("Nuovo utente %s! Creazione in corso...\n", richiesta.Name)
		idUtente = "abcd012345"
		userDatabase[richiesta.Name] = idUtente
	}

	
	risposta := LoginResponse{
		Identifier: idUtente,
		Token:      "GciOiJIUzI1NiIsInR5cCI6IkpXVCJ", 
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated) 
	json.NewEncoder(w).Encode(risposta)
}
