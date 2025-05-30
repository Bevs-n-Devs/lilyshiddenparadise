package handlers

import (
	"fmt"
	"net/http"

	"github.com/Bevs-n-Devs/lilyshiddenparadise/db"
	"github.com/Bevs-n-Devs/lilyshiddenparadise/logs"
	"github.com/Bevs-n-Devs/lilyshiddenparadise/middleware"
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
		logs.Logs(logErr, fmt.Sprintf("Error authenticating landlord: %s. Redirecting back to landlord login page", err.Error()))
		http.Redirect(w, r, "/login/landlord?authenticationError=UNAUTHORIZED+401:+Error+authenticating+landlord", http.StatusSeeOther)
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
	createLandlordDashboardSessionCookie := middleware.LandlordDashboardSessionCookie(w, sessionToken, expiryTime)
	if !createLandlordDashboardSessionCookie {
		logs.Logs(logErr, "Failed to create session cookie. Redirecting back to landlord login page...")
		http.Redirect(w, r, "/login/landlord?internalServerError=INTERNAL+SERVER+ERROR+500:+Failed+to+create+session+cookie", http.StatusInternalServerError)
		return
	}

	// set csrf cookie
	createLandlordDashboardCSRFCookie := middleware.LandlordDashboardCSRFTokenCookie(w, csrfToken, expiryTime)
	if !createLandlordDashboardCSRFCookie {
		logs.Logs(logErr, "Failed to create CSRF cookie. Redirecting back to landlord login page...")
		http.Redirect(w, r, "/login/landlord?internalServerError=INTERNAL+SERVER+ERROR+500:+Failed+to+create+CSRF+cookie", http.StatusInternalServerError)
		return
	}

	// redirect to landlord dashboard if authentication is successful
	http.Redirect(w, r, "/landlord/dashboard", http.StatusSeeOther)
}
