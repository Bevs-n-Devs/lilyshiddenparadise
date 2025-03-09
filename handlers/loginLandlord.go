package handlers

import (
	"fmt"
	"net/http"

	"github.com/Bevs-n-Devs/lilyshiddenparadise/logs"
)

func LoginLandlord(w http.ResponseWriter, r *http.Request) {
	// get any error messages
	badRequestError := r.URL.Query().Get("badRequest")
	notFoundError := r.URL.Query().Get("notFound")
	authenticationError := r.URL.Query().Get("authenticationError")

	data := ErrorMessages{
		BadRequestError:     badRequestError,
		NotFoundError:       notFoundError,
		AuthenticationError: authenticationError,
	}

	err := Templates.ExecuteTemplate(w, "loginLandlord.html", data)
	if err != nil {
		logs.Logs(logErr, fmt.Sprintf("Unable to load login landlord page: %s", err.Error()))
		http.Error(w, fmt.Sprintf("Unable to load page to login landlord: %s", err.Error()), http.StatusInternalServerError)
	}
}
