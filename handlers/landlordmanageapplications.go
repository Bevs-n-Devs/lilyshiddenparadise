package handlers

import (
	"fmt"
	"net/http"

	"github.com/Bevs-n-Devs/lilyshiddenparadise/db"
	"github.com/Bevs-n-Devs/lilyshiddenparadise/logs"
)

func LandlordManageApplications(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		logs.Logs(logErr, "Invalid request method. Redirecting back to landlord login page.")
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

	applicationId := r.FormValue("applicationId")
	applicationResult := r.FormValue("applicationResult")

	// TODO: check if applicationResult is accepted or denied

	err = db.UpdateTenantApplicationStatus(applicationId, applicationResult)
	if err != nil {
		logs.Logs(logErr, fmt.Sprintf("Error updating tenant application status: %s", err.Error()))
		http.Error(w, fmt.Sprintf("Error updating tenant application status: %s", err.Error()), http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/landlord/dashboard", http.StatusSeeOther)
}
