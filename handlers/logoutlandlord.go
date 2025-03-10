package handlers

import (
	"fmt"
	"net/http"
	"time"

	"github.com/Bevs-n-Devs/lilyshiddenparadise/db"
	"github.com/Bevs-n-Devs/lilyshiddenparadise/logs"
)

func LogoutLandlord(w http.ResponseWriter, r *http.Request) {
	sessionToken, err := r.Cookie("session_token")
	if err != nil || sessionToken.Value == "" {
		logs.Logs(logErr, fmt.Sprintf("Failed to get session token: %s. Redirecting to home page", err.Error()))
		http.Redirect(w, r, "/?authenticationError=UNAUTHORIZED+401:+Error+authenticating+user", http.StatusSeeOther)
		return
	}

	email, err := db.GetEmailFromLandlordSessionToken(sessionToken.Value)
	if err != nil {
		logs.Logs(logErr, fmt.Sprintf("Failed to get email from session token: %s", err.Error()))
		http.Error(w, "Failed to get email from session token", http.StatusInternalServerError)
	}

	// delete the session token, CSRF token and expiry time from the database
	err = db.LogoutLandlord(email)
	if err != nil {
		logs.Logs(logErr, fmt.Sprintf("Failed to logout user: %s", err.Error()))
		http.Redirect(w, r, "/landlord/dashboard", http.StatusSeeOther)
	}

	// delete the session token, CSRF token and expiry time from the cookie
	http.SetCookie(w, &http.Cookie{
		Name:     "session_token",
		Value:    "",
		Expires:  time.Now().Add(-time.Hour),
		HttpOnly: true,
		Path:     "/",
		SameSite: http.SameSiteStrictMode,
	})

	http.SetCookie(w, &http.Cookie{
		Name:     "csrf_token",
		Value:    "",
		Expires:  time.Now().Add(-time.Hour),
		HttpOnly: false,
		Path:     "/",
		SameSite: http.SameSiteStrictMode,
	})

	logs.Logs(logInfo, "Landlord logged out successfully")
	http.Redirect(w, r, "/", http.StatusSeeOther)
}
