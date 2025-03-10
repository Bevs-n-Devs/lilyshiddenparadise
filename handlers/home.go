package handlers

import (
	"fmt"
	"net/http"

	"github.com/Bevs-n-Devs/lilyshiddenparadise/logs"
)

func Home(w http.ResponseWriter, r *http.Request) {
	// get any error messages
	authenticationError := r.URL.Query().Get("authenticationError")

	data := ErrorMessages{
		AuthenticationError: authenticationError,
	}

	err := Templates.ExecuteTemplate(w, "home.html", data)
	if err != nil {
		logs.Logs(logErr, fmt.Sprintf("Unable to load home page: %s", err))
		http.Error(w, fmt.Sprintf("Unable to load home page: %s", err.Error()), http.StatusInternalServerError)
	}
}
