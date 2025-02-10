package handlers

import (
	"fmt"
	"net/http"

	"github.com/Bevs-n-Devs/lilyshiddenparadise/logs"
)

func Home(w http.ResponseWriter, r *http.Request) {
	err := Templates.ExecuteTemplate(w, "index.html", nil)
	if err != nil {
		logs.Logs(3, fmt.Sprintf("Unable to load home page: %s", err))
		http.Error(w, "Unable to load home page: "+err.Error(), http.StatusInternalServerError)
	}
}
