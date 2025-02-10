package handlers

import (
	"fmt"
	"net/http"

	"github.com/Bevs-n-Devs/lilyshiddenparadise/logs"
)

func Tenancy(w http.ResponseWriter, r *http.Request) {
	err := Templates.ExecuteTemplate(w, "contact.html", nil)
	if err != nil {
		logs.Logs(3, fmt.Sprintf("Unable to load tenancy page: %s", err))
		http.Error(w, "Unable to load tenancy page: "+err.Error(), http.StatusInternalServerError)
	}
}
