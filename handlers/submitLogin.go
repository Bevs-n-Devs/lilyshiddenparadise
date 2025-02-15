package handlers

import (
	"fmt"
	"net/http"

	"github.com/Bevs-n-Devs/lilyshiddenparadise/logs"
)

func SubmitLogin(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		logs.Logs(3, "Invalid request method")
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	// parse from data
	err := r.ParseForm()
	if err != nil {
		logs.Logs(3, fmt.Sprintf("Could not extract data from form: %s", err.Error()))
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// extract login data
	loginEmail := r.FormValue("loginEmail")
	password := r.FormValue("password")

	logs.Logs(1, fmt.Sprintf("Login attempt - username: %s", loginEmail))
	logs.Logs(1, fmt.Sprintf("Login attempt - password: %s", password))

	// redirect to home page
	logs.Logs(1, "Redirecting to home page")
	http.Redirect(w, r, "/", http.StatusSeeOther)
}
