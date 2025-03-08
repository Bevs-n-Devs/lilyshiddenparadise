package handlers

import (
	"fmt"
	"net/http"

	"github.com/Bevs-n-Devs/lilyshiddenparadise/logs"
)

func Login(w http.ResponseWriter, r *http.Request) {
	err := Templates.ExecuteTemplate(w, "login.html", nil)
	if err != nil {
		logs.Logs(3, fmt.Sprintf("Unable to load contact page: %s", err))
		http.Error(w, "Unable to load contact page: "+err.Error(), http.StatusInternalServerError)
	}
}
