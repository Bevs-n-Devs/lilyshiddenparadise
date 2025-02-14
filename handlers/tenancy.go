package handlers

import (
	"fmt"
	"net/http"

	"github.com/Bevs-n-Devs/lilyshiddenparadise/logs"
)

func Tenancy(w http.ResponseWriter, r *http.Request) {
	evictedError := r.URL.Query().Get("evictedError")

	data := struct {
		EvictedError string
	}{
		EvictedError: evictedError,
	}

	err := Templates.ExecuteTemplate(w, "tenancy.html", data)
	if err != nil {
		logs.Logs(3, fmt.Sprintf("Unable to load tenancy page: %s", err))
		http.Error(w, "Unable to load tenancy page: "+err.Error(), http.StatusInternalServerError)
		return
	}
}
