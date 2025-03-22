package handlers

import (
	"fmt"
	"net/http"

	"github.com/Bevs-n-Devs/lilyshiddenparadise/db"
	"github.com/Bevs-n-Devs/lilyshiddenparadise/logs"
	"github.com/Bevs-n-Devs/lilyshiddenparadise/middleware"
	"github.com/Bevs-n-Devs/lilyshiddenparadise/utils"
)

func UpdateTenantPassword(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		logs.Logs(logErr, fmt.Sprintf("Invalid request method: %s. Redirecting back to tenant login page.", r.Method))
		http.Redirect(w, r, "/login/tenant?badRequest=BAD+REQUEST+400:+Invalid+request+method", http.StatusBadRequest)
		return
	}

	// deny the request if the authorization fails
	err := middleware.AuthenticateTenantRequest(r)
	if err != nil {
		logs.Logs(logErr, fmt.Sprintf("Error authenticating tenant: %s. Redirecting to tenant login page", err.Error()))
		http.Redirect(w, r, "/login/tenant?authenticationError=UNAUTHORIZED+401:+Error+authenticating+tenant", http.StatusSeeOther)
		return
	}

	// get session cookie
	sessionToken, err := utils.CheckSessionToken(r)
	if err != nil {
		logs.Logs(logErr, fmt.Sprintf("Error getting session token: %s. Redirecting to tenant login page", err.Error()))
		http.Redirect(w, r, "/login/tenant?authenticationError=UNAUTHORIZED+401:+Error+authenticating+tenant.+Failed+to+get+session+token", http.StatusSeeOther)
		return
	}

	// get tenant email from session cookie
	tenantEmail, err := db.GetHashedEmailFromTenantSessionToken(sessionToken.Value)
	if err != nil {
		logs.Logs(logErr, fmt.Sprintf("Error getting tenant email from session token: %s. Redirecting to tenant login page", err.Error()))
		http.Redirect(w, r, "/login/tenant?authenticationError=UNAUTHORIZED+401:+Error+authenticating+tenant.+Failed+to+get+tenant+email+from+session+token", http.StatusSeeOther)
		return
	}

	newSessionToken, newCsrfToken, newExpiryTime, err := db.UpdateTenantSessionTokens(tenantEmail)
	if err != nil {
		logs.Logs(logErr, fmt.Sprintf("Error updating tenant session tokens: %s. Redirecting to tenant login page", err.Error()))
		http.Redirect(w, r, "/login/tenant?authenticationError=UNAUTHORIZED+401:+Error+authenticating+tenant.+Failed+to+update+session+tokens", http.StatusSeeOther)
		return
	}

	// set cookie to redirect the user back to tenant dashboard
	createTenantDashboardSessionCookie := middleware.TenantDashboardSessionCookie(w, newSessionToken, newExpiryTime)
	if !createTenantDashboardSessionCookie {
		logs.Logs(logErr, "Error creating tenant dashboard session cookie. Redirecting to tenant login page")
		http.Redirect(w, r, "/login/tenant?authenticationError=UNAUTHORIZED+401:+Error+authenticating+tenant.+Failed+to+create+tenant+dashboard+session+cookie", http.StatusSeeOther)
		return
	}
	createTenantDashboardCsrfCookie := middleware.TenantDashboardCSRFTokenCookie(w, newCsrfToken, newExpiryTime)
	if !createTenantDashboardCsrfCookie {
		logs.Logs(logErr, "Error creating tenant dashboard csrf cookie. Redirecting to tenant login page")
		http.Redirect(w, r, "/login/tenant?authenticationError=UNAUTHORIZED+401:+Error+authenticating+tenant.+Failed+to+create+tenant+dashboard+csrf+cookie", http.StatusSeeOther)
		return
	}

	// parse form data
	err = r.ParseForm()
	if err != nil {
		logs.Logs(logErr, fmt.Sprintf("Error parsing form data: %s", err.Error()))
		http.Error(w, fmt.Sprintf("Error parsing form data: %s", err.Error()), http.StatusInternalServerError)
		return
	}

	// extract data from form
	formTenantEmail := r.FormValue("tenantEmail")
	formOldPassword := r.FormValue("oldPassword")
	formNewPassword := r.FormValue("newPassword")
	formConfirmPassword := r.FormValue("confirmPassword")

	// TODO: make sure old password is correct
	hashTenantEmail := utils.HashData(formTenantEmail)
	hashTenantPassword := utils.HashData(formOldPassword)
	exists, err := db.AuthenticateTenant(hashTenantEmail, hashTenantPassword)
	if err != nil {
		logs.Logs(logErr, fmt.Sprintf("Error authenticating tenant: %s. Redirecting to tenant login page", err.Error()))
		http.Redirect(w, r, "/login/tenant?authenticationError=UNAUTHORIZED+401:+Error+authenticating+tenant", http.StatusSeeOther)
		return
	}

	if !exists {
		logs.Logs(logErr, "Tenant email or password is incorrect. Redirecting to tenant login page")
		http.Redirect(w, r, "/login/tenant?authenticationError=UNAUTHORIZED+401:+Tenant+email+or+password+is+incorrect", http.StatusSeeOther)
		return
	}

	// TODO: validate new password is the same as confirmed password
	match := utils.ValidateNewPassword(formNewPassword, formConfirmPassword)
	if !match {
		logs.Logs(logErr, "New password and confirmed password do not match. Redirecting to tenant back to tenant account password page")
		http.Redirect(w, r, "/tenant/dashboard/account?authenticationError=UNAUTHORIZED+401:+New+password+and+confirmed+password+do+not+match", http.StatusSeeOther)
		return
	}

	// TODO: update password in database (will also need to encrypt updated password)
	err = db.UpdateTenantPassword(hashTenantEmail, hashTenantPassword, formNewPassword)
	if err != nil {
		logs.Logs(logErr, fmt.Sprintf("Error updating tenant password: %s. Redirecting to tenant login page", err.Error()))
		http.Redirect(w, r, "/login/tenant?authenticationError=UNAUTHORIZED+401:+Error+updating+tenant+password", http.StatusSeeOther)
		return
	}

	// redirect to tenant dashboard
	http.Redirect(w, r, "/tenant/dashboard", http.StatusSeeOther)
}
