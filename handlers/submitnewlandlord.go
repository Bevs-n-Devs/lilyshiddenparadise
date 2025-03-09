package handlers

import (
	"fmt"
	"net/http"

	"github.com/Bevs-n-Devs/lilyshiddenparadise/db"
	"github.com/Bevs-n-Devs/lilyshiddenparadise/logs"
	"github.com/Bevs-n-Devs/lilyshiddenparadise/utils"
)

func SubmitNewLandlord(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		logs.Logs(logWarn, fmt.Sprintf("Invalid request method: %s. Redirecting back to create new landlord page.", r.Method))
		http.Redirect(w, r, "/new/landlord", http.StatusSeeOther)
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
	confirmPassword := r.FormValue("confirmPassword")

	// validate passwords
	if !utils.ValidateNewLandlordPassword(landlordPassword, confirmPassword) {
		logs.Logs(logErr, "Passwords do not match. Please try again.")
		http.Redirect(w, r, "/new/landlord?confirmPasswordError=Passwords+do+not+match.+Please+try+again.", http.StatusSeeOther)
		return
	}

	// TODO: add logic to create new landlord in db
	err = db.CreateNewLandlord(landlordEmail, landlordPassword)
	if err != nil {
		logs.Logs(logErr, fmt.Sprintf("Error creating new landlord: %s", err.Error()))
		http.Error(w, fmt.Sprintf("Error creating new landlord: %s", err.Error()), http.StatusInternalServerError)
		return
	}

	logs.Logs(logInfo, "New landlord created successfully")
	http.Redirect(w, r, "/login/landlord", http.StatusSeeOther)

}
