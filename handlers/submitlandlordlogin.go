package handlers

import (
	"fmt"
	"net/http"

	"github.com/Bevs-n-Devs/lilyshiddenparadise/db"
	"github.com/Bevs-n-Devs/lilyshiddenparadise/logs"
)

func SubmitLoginLandlord(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		logs.Logs(logErr, fmt.Sprintf("Invalid request method: %s. Redirecting back to login landlord page.", r.Method))
		http.Redirect(w, r, "/login/landlord?badRequest=BAD+REQUEST+400:+Invalid+request+method", http.StatusBadRequest)
		return
	}

	// parse form data
	err := r.ParseForm()
	if err != nil {
		logs.Logs(logErr, fmt.Sprintf("Error parsing form data: %s", err.Error()))
		http.Error(w, fmt.Sprintf("Error parsing form data: %s", err.Error()), http.StatusInternalServerError)
		return
	}

	// extract data from form
	landlordEmail := r.FormValue("landlordEmail")
	landlordPassword := r.FormValue("landlordPassword")

	// check if landlord exists in database
	exists, err := db.AuthenticateLandlord(landlordEmail, landlordPassword)
	if err != nil {
		logs.Logs(logErr, fmt.Sprintf("Error authenticating landlord: %s", err.Error()))
		http.Error(w, fmt.Sprintf("Error authenticating landlord: %s", err.Error()), http.StatusInternalServerError)
		return
	}

	if !exists {
		logs.Logs(logErr, "Landlord does not exist. Please try again. Redirecting back to landlord login page...")
		http.Redirect(w, r, "/login/landlord?notFound=NOT+FOUND+404:+Landlord+does_not-exist.+Please+try+again.", http.StatusNotFound)
		return
	}

	// add tokens & expiry time to database
	sessionToken, csrfToken, expiryTime, err := db.UpdateLandlordSessionTokens(landlordEmail)
	if err != nil {
		logs.Logs(logErr, fmt.Sprintf("Error updating landlord session tokens: %s", err.Error()))
		http.Error(w, fmt.Sprintf("Error updating landlord session tokens: %s", err.Error()), http.StatusInternalServerError)
		return
	}

	// set session cookie
	http.SetCookie(w, &http.Cookie{
		Name:     "landlord_session_token",
		Value:    sessionToken,
		Expires:  expiryTime,
		HttpOnly: true,
	})

	// set csrf cookie
	http.SetCookie(w, &http.Cookie{
		Name:     "landlord_csrf_token",
		Value:    csrfToken,
		Expires:  expiryTime,
		HttpOnly: false,
	})

	// redirect to landlord dashboard if authentication is successful
	http.Redirect(w, r, "/landlord/dahboard", http.StatusSeeOther)
}
