package handlers

import (
	"fmt"
	"net/http"

	"github.com/Bevs-n-Devs/lilyshiddenparadise/logs"
	"github.com/Bevs-n-Devs/lilyshiddenparadise/middleware"
	"github.com/Bevs-n-Devs/lilyshiddenparadise/utils"
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

	// set cookie to logout landlord
	sessionToken, err := utils.CheckSessionToken(r)
	if err != nil {
		http.Redirect(w, r, "/login/landlord?authenticationError=UNAUTHORIZED+401:+Error+authenticating+landlord.+Failed+to+get+session+token", http.StatusSeeOther)
		return
	}

	csrfToken, err := utils.CheckCSRFToken(r)
	if err != nil {
		http.Redirect(w, r, "/login/landlord?authenticationError=UNAUTHORIZED+401:+Error+authenticating+landlord.+Failed+to+get+CSRF+token", http.StatusSeeOther)
		return
	}

	// set cookies to landlord dashboard tenants page
	createLandlordDashboardTenantSessionCookie := middleware.LandlordDashboardTenantsSessionCookie(w, sessionToken)
	if !createLandlordDashboardTenantSessionCookie {
		logs.Logs(logErr, "Failed to create session for the landlord tenants dashboard page. Redirecting back to login page")
		http.Redirect(w, r, "/login/landlord?authenticationError=UNAUTHORIZED+401:+Error+authenticating+landlord.+Failed+to+create+session+cookie", http.StatusSeeOther)
		return
	}

	createLandlordDashboardTenantCSRFTokenCookie := middleware.LandlordDashboardTenantsCSRFTokenCookie(w, csrfToken)
	if !createLandlordDashboardTenantCSRFTokenCookie {
		logs.Logs(logErr, "Failed to create CSRF token for the landlord tenants dashboard page. Redirecting back to login page")
		http.Redirect(w, r, "/login/landlord?authenticationError=UNAUTHORIZED+401:+Error+authenticating+landlord.+Failed+to+create+CSRF+token+cookie", http.StatusSeeOther)
	}

	// set cookies to logout
	createSessionCookie := middleware.LogoutLandlordSessionCookie(w, sessionToken)
	if !createSessionCookie {
		logs.Logs(logErr, "Failed to create session cookie for landlord. Redirecting to landlord login page")
		http.Redirect(w, r, "/login/landlord?authenticationError=UNAUTHORIZED+401:+Error+authenticating+landlord.+Failed+to+create+session+cookie", http.StatusSeeOther)
		return
	}
	createCSRFTokenCookie := middleware.LogoutLandlordCSRFTokenCookie(w, csrfToken)
	if !createCSRFTokenCookie {
		logs.Logs(logErr, "Failed to create CSRF token cookie for landlord. Redirecting to landlord login page")
		http.Redirect(w, r, "/login/landlord?authenticationError=UNAUTHORIZED+401:+Error+authenticating+landlord.+Failed+to+create+CSRF+token+cookie", http.StatusSeeOther)
		return
	}

	// direct user to protected dashboard
	err = Templates.ExecuteTemplate(w, "landlordDashboard.html", nil)
	if err != nil {
		logs.Logs(logErr, fmt.Sprintf("Unable to load landlord dashboard: %s", err.Error()))
		http.Error(w, fmt.Sprintf("Unable to load landlord dashboard: %s", err.Error()), http.StatusInternalServerError)
	}
}
