package handlers

import (
	"fmt"
	"net/http"

	"github.com/Bevs-n-Devs/lilyshiddenparadise/db"
	"github.com/Bevs-n-Devs/lilyshiddenparadise/logs"
	"github.com/Bevs-n-Devs/lilyshiddenparadise/middleware"
	"github.com/Bevs-n-Devs/lilyshiddenparadise/utils"
)

func SubmitLoginTenant(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		logs.Logs(logErr, fmt.Sprintf("Invalid request method: %s. Redirecting back to login tenant page.", r.Method))
		http.Redirect(w, r, "/login/tenant?badRequest=BAD+REQUEST+400:+Invalid+request+method", http.StatusBadRequest)
	}

	// parse form data
	err := r.ParseForm()
	if err != nil {
		logs.Logs(logErr, fmt.Sprintf("Error parsing form data: %s", err.Error()))
		http.Error(w, fmt.Sprintf("Error parsing form data: %s", err.Error()), http.StatusInternalServerError)
		return
	}

	// extract data from form
	tenantEmail := r.FormValue("tenantEmail")
	tenantPassword := r.FormValue("tenantPassword")

	// hash username & password
	hashUsername := utils.HashData(tenantEmail)
	hashPassword := utils.HashData(tenantPassword)

	authenticate, err := db.AuthenticateTenant(hashUsername, hashPassword)
	if err != nil {
		logs.Logs(logErr, fmt.Sprintf("Error authenticating tenant: %s. Redirecting back to tenant login page", err.Error()))
		http.Redirect(w, r, "/login/tenant?authenticationError=UNAUTHORIZED+401:+Error+authenticating+tenant", http.StatusSeeOther)
		return
	}

	if !authenticate {
		logs.Logs(logErr, "Tenant does not exist. Please try again. Redirecting back to tenant login page...")
		http.Redirect(w, r, "/login/tenant?notFound=NOT+FOUND+404:+Tenant+does_not-exist.+Please+try+again.", http.StatusNotFound)
		return
	}

	// add tokens & expiry time to database
	sessionToken, csrfToken, expiryTime, err := db.UpdateTenantSessionTokens(hashUsername)
	if err != nil {
		logs.Logs(logErr, fmt.Sprintf("Error updating tenant session tokens: %s", err.Error()))
		http.Error(w, fmt.Sprintf("Error updating tenant session tokens: %s", err.Error()), http.StatusInternalServerError)
		return
	}

	// set session cookie for tenant dashboard
	createTenantDashboardSessionCookie := middleware.TenantDashboardSessionCookie(w, sessionToken, expiryTime)
	if !createTenantDashboardSessionCookie {
		logs.Logs(logErr, "Failed to create session cookie. Redirecting back to tenant login page...")
		http.Redirect(w, r, "/login/tenant?internalServerError=INTERNAL+SERVER+ERROR+500:+Failed+to+create+session+cookie", http.StatusInternalServerError)
		return
	}

	// set csrf cookie for tenant dashboard
	createTenantDashboardCsrfCookie := middleware.TenantDashboardCSRFTokenCookie(w, csrfToken, expiryTime)
	if !createTenantDashboardCsrfCookie {
		logs.Logs(logErr, "Failed to create CSRF cookie. Redirecting back to tenant login page...")
		http.Redirect(w, r, "/login/tenant?internalServerError=INTERNAL+SERVER+ERROR+500:+Failed+to+create+CSRF+cookie", http.StatusInternalServerError)
		return
	}

	// redirect to tenant dashboard if authentication is successful
	http.Redirect(w, r, "/tenant/dashboard", http.StatusSeeOther)
}
