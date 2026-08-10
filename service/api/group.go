package api

import (
	"encoding/json"
	"net/http"
	"strings"

	"git.sapienzaapps.it/fantasticcoffee/fantastic-coffee-decaffeinated/service/api/reqcontext"
	"github.com/julienschmidt/httprouter"
)

func (rt *_router) createGroup(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
	var body struct {
		GroupName        string   `json:"groupName"`
		GroupDescription string   `json:"groupDescription"`
		UserId           string   `json:"userId"`
		PhotoUrl         string   `json:"photoUrl,omitempty"`
		Members          []string `json:"members,omitempty"`
	}

	err := json.NewDecoder(r.Body).Decode(&body)
	if err != nil {
		http.Error(w, "Formato JSON non valido", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	creatorId := body.UserId
	if creatorId == "" {
		creatorId = ctx.UserId
	}

	newGroup, err := rt.db.CreateGroup(body.GroupName, body.GroupDescription, body.PhotoUrl, creatorId, body.Members)
	if err != nil {
		if strings.Contains(err.Error(), "utente creatore non trovato") {
			http.Error(w, "Utente creatore non trovato", http.StatusNotFound)
			return
		}
		if strings.Contains(err.Error(), "caratteri") || strings.Contains(err.Error(), "superare") {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		http.Error(w, "Errore interno durante la creazione del gruppo", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	_ = json.NewEncoder(w).Encode(newGroup)
}

func (rt *_router) setGroupName(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
	groupId := ps.ByName("groupId")

	var body struct {
		Name   string `json:"name"`
		UserId string `json:"userId"`
	}

	err := json.NewDecoder(r.Body).Decode(&body)
	if err != nil {
		http.Error(w, "Formato JSON non valido", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	userId := body.UserId
	if userId == "" {
		userId = ctx.UserId
	}

	err = rt.db.SetGroupName(groupId, userId, body.Name)
	if err != nil {
		if strings.Contains(err.Error(), "non trovato") {
			http.Error(w, "Gruppo non trovato", http.StatusNotFound)
			return
		}
		if strings.Contains(err.Error(), "non autorizzato") {
			http.Error(w, "Accesso Negato - Non hai i permessi per modificare questo gruppo", http.StatusForbidden)
			return
		}
		if strings.Contains(err.Error(), "caratteri") {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		http.Error(w, "Errore interno del server", http.StatusInternalServerError)
		return
	}

	group, err := rt.db.GetGroup(groupId)
	if err != nil {
		http.Error(w, "Errore nel recupero del gruppo aggiornato", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(group)
}

func (rt *_router) setGroupPhoto(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
	groupId := ps.ByName("groupId")

	var body struct {
		PhotoUrl string `json:"photoUrl"`
	}

	err := json.NewDecoder(r.Body).Decode(&body)
	if err != nil {
		http.Error(w, "Formato JSON non valido", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	err = rt.db.SetGroupPhoto(groupId, ctx.UserId, body.PhotoUrl)
	if err != nil {
		if strings.Contains(err.Error(), "non trovato") {
			http.Error(w, "Gruppo non trovato", http.StatusNotFound)
			return
		}
		if strings.Contains(err.Error(), "non autorizzato") {
			http.Error(w, "Accesso Negato - Non hai i permessi per modificare questo gruppo", http.StatusForbidden)
			return
		}
		if strings.Contains(err.Error(), "caratteri") {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		http.Error(w, "Errore interno del server", http.StatusInternalServerError)
		return
	}

	group, err := rt.db.GetGroup(groupId)
	if err != nil {
		http.Error(w, "Errore nel recupero del gruppo aggiornato", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(group)
}

func (rt *_router) addToGroup(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
	groupId := ps.ByName("groupId")

	var body struct {
		UserId string `json:"userId"`
	}

	err := json.NewDecoder(r.Body).Decode(&body)
	if err != nil {
		http.Error(w, "Formato JSON non valido", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	err = rt.db.AddToGroup(groupId, body.UserId)
	if err != nil {
		if strings.Contains(err.Error(), "già membro") {
			http.Error(w, "L'utente è già un membro del gruppo", http.StatusConflict)
			return
		}
		if strings.Contains(err.Error(), "non trovato") {
			http.Error(w, "Utente o gruppo non trovato", http.StatusNotFound)
			return
		}
		http.Error(w, "Errore interno del server", http.StatusInternalServerError)
		return
	}

	group, err := rt.db.GetGroup(groupId)
	if err != nil {
		http.Error(w, "Errore nel recupero del gruppo aggiornato", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	_ = json.NewEncoder(w).Encode(group)
}

func (rt *_router) leaveGroup(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
	groupId := ps.ByName("groupId")
	memberId := ps.ByName("memberId")

	if ctx.UserId != memberId {
		http.Error(w, "Non puoi rimuovere un utente diverso da te stesso", http.StatusForbidden)
		return
	}

	err := rt.db.LeaveGroup(groupId, memberId)
	if err != nil {
		if strings.Contains(err.Error(), "non trovato") || strings.Contains(err.Error(), "inesistente") {
			http.Error(w, "Gruppo o utente non trovato", http.StatusNotFound)
			return
		}
		http.Error(w, "Errore interno del server", http.StatusInternalServerError)
		return
	}

	group, err := rt.db.GetGroup(groupId)
	if err != nil {
		w.WriteHeader(http.StatusOK)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(group)
}
