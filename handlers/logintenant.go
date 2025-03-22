package handlers

import (
	"fmt"
	"net/http"

	"github.com/Bevs-n-Devs/lilyshiddenparadise/logs"
)

func LoginTenant(w http.ResponseWriter, r *http.Request) {
	// get any error messages
	badRequestError := r.URL.Query().Get("badRequest")
	notFoundError := r.URL.Query().Get("notFound")
	authenticationError := r.URL.Query().Get("authenticationError")
	internalServerError := r.URL.Query().Get("internalServerError")

	data := ErrorMessages{
		BadRequestError:     badRequestError,
		NotFoundError:       notFoundError,
		AuthenticationError: authenticationError,
		InternalServerError: internalServerError,
	}

	err := Templates.ExecuteTemplate(w, "loginTenant.html", data)
	if err != nil {
		logs.Logs(logErr, fmt.Sprintf("Unable to load contact page: %s", err))
		http.Error(w, "Unable to load contact page: "+err.Error(), http.StatusInternalServerError)
	}
}
