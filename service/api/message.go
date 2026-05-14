package api

import "net/http"

func (rt *_router) sendMessage(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
}

func (rt *_router) updateMessage(w http.ResponseWriter, r http.Request) {

}
