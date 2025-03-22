package handlers

import (
	"fmt"
	"net/http"

	"github.com/Bevs-n-Devs/lilyshiddenparadise/db"
	"github.com/Bevs-n-Devs/lilyshiddenparadise/logs"
	"github.com/Bevs-n-Devs/lilyshiddenparadise/middleware"
)

func LogoutTenant(w http.ResponseWriter, r *http.Request) {
	sessionToken, err := r.Cookie("session_token")
	if err != nil || sessionToken.Value == "" {
		logs.Logs(logErr, "Failed to get session token. Redirecting to home page")
		http.Redirect(w, r, "/?authenticationError=UNAUTHORIZED+401:+Error+authenticating+user", http.StatusSeeOther)
		return
	}

	email, err := db.GetHashedEmailFromTenantSessionToken(sessionToken.Value)
	if err != nil {
		logs.Logs(logErr, fmt.Sprintf("Failed to get email from session token: %s", err.Error()))
		http.Error(w, "Failed to get email from session token", http.StatusInternalServerError)
	}

	// delete the session token, CSRF token and expiry time from the database
	err = db.LogoutTenant(email)
	if err != nil {
		logs.Logs(logErr, fmt.Sprintf("Failed to logout user: %s", err.Error()))
		http.Redirect(w, r, "/tenant/dashboard", http.StatusSeeOther)
	}

	// delete the session token, CSRF token and expiry time from the cookie
	deleteSessionCookie := middleware.DeleteTenantSessionCookie(w)
	if !deleteSessionCookie {
		logs.Logs(logWarn, "Failed to delete session token cookie. Redirecting to home page")
		http.Redirect(w, r, "/?cookieError=COOKIE+ERROR+500:+Failed+to+delete+session+token+cookie", http.StatusSeeOther)
		return
	}
	deleteCSRFCookie := middleware.DeleteTenantCSRFCookie(w)
	if !deleteCSRFCookie {
		logs.Logs(logWarn, "Failed to delete CSRF token cookie. Redirecting to home page")
		http.Redirect(w, r, "/?cookieError=COOKIE+ERROR+500:+Failed+to+delete+CSRF+token+cookie", http.StatusSeeOther)
		return
	}

	logs.Logs(logInfo, "Tenant logged out successfully")
	http.Redirect(w, r, "/", http.StatusSeeOther)
}
