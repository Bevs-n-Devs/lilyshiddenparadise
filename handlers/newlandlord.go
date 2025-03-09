package handlers

import (
	"fmt"
	"net/http"

	"github.com/Bevs-n-Devs/lilyshiddenparadise/logs"
)

func NewLandlord(w http.ResponseWriter, r *http.Request) {
	// get error message (if any)
	confirmPasswordError := r.URL.Query().Get("confirmPasswordError")

	data := ErrorMessages{
		ConfirmPasswordError: confirmPasswordError,
	}

	// pass error message to HTML template
	err := Templates.ExecuteTemplate(w, "newLandlord.html", data)
	if err != nil {
		logs.Logs(logErr, fmt.Sprintf("Unable to load new landlord page: %s", err.Error()))
		http.Error(w, fmt.Sprintf("Unable to load page to create new landlord: %s", err.Error()), http.StatusInternalServerError)
	}
}
