package handlers

import (
	"fmt"
	"net/http"

	"github.com/Bevs-n-Devs/lilyshiddenparadise/logs"
)

func LoginTenant(w http.ResponseWriter, r *http.Request) {
	err := Templates.ExecuteTemplate(w, "loginTenant.html", nil)
	if err != nil {
		logs.Logs(logErr, fmt.Sprintf("Unable to load contact page: %s", err))
		http.Error(w, "Unable to load contact page: "+err.Error(), http.StatusInternalServerError)
	}
}
