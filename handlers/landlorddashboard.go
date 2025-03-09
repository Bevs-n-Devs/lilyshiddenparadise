package handlers

import (
	"fmt"
	"net/http"

	"github.com/Bevs-n-Devs/lilyshiddenparadise/logs"
	"github.com/Bevs-n-Devs/lilyshiddenparadise/middleware"
)

func LandlordDashboard(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		logs.Logs(logErr, fmt.Sprintf("Invalid request method: %s. Redirecting back to landlord login page.", r.Method))
		http.Redirect(w, r, "/login/landlord?badRequest=BAD+REQUEST+400:+Invalid+request+method", http.StatusBadRequest)
		return
	}

	// deny the request if the authorization fails
	err := middleware.AuthenticateLandlordRequest(r)
	if err != nil {
		logs.Logs(logErr, fmt.Sprintf("Error authenticating landlord: %s. Redirecting to landlord login page", err.Error()))
		http.Redirect(w, r, "/login/landlord?authenticationError=UNAUTHORIZED+401:+Error+authenticating+landlord", http.StatusSeeOther)
		return
	}

	// direct user to protected dashboard
	err = Templates.ExecuteTemplate(w, "landlordDashboard.html", nil)
	if err != nil {
		logs.Logs(logErr, fmt.Sprintf("Unable to load landlord dashboard: %s", err.Error()))
		http.Error(w, fmt.Sprintf("Unable to load landlord dashboard: %s", err.Error()), http.StatusInternalServerError)
	}
}
