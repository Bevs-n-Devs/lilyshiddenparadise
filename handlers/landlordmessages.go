package handlers

import (
	"fmt"
	"net/http"

	"github.com/Bevs-n-Devs/lilyshiddenparadise/db"
	"github.com/Bevs-n-Devs/lilyshiddenparadise/logs"
	"github.com/Bevs-n-Devs/lilyshiddenparadise/middleware"
	"github.com/Bevs-n-Devs/lilyshiddenparadise/utils"
)

func LandlordMessages(w http.ResponseWriter, r *http.Request) {
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

	// get session cookie
	sessionToken, err := utils.CheckSessionToken(r)
	if err != nil {
		http.Redirect(w, r, "/login/landlord?authenticationError=UNAUTHORIZED+401:+Error+authenticating+landlord.+Failed+to+get+session+token", http.StatusSeeOther)
		return
	}

	// get landlord emial from session cookie
	landlordEmail, err := db.GetEmailFromLandlordSessionToken(sessionToken.Value)
	if err != nil {
		http.Redirect(w, r, "/login/landlord?authenticationError=UNAUTHORIZED+401:+Error+authenticating+landlord.+Failed+to+get+landlord+email+from+session+token", http.StatusSeeOther)
		return
	}

	// update the landlord's session token, CSRF token and expiry time in the database
	// this will be doen for each request
	newSessionToken, newCsrfToken, newExpiryTime, err := db.UpdateLandlordSessionTokens(landlordEmail)
	if err != nil {
		logs.Logs(logErr, fmt.Sprintf("Error updating landlord session tokens: %s. Redirecting to landlord login page", err.Error()))
		http.Redirect(w, r, "/login/landlord?authenticationError=UNAUTHORIZED+401:+Error+authenticating+landlord.+Failed+to+update+session+tokens", http.StatusSeeOther)
		return
	}

	// TODO! Set cookies for each available page

	// set new cookies for landlord dashboard
	createLandlordDashboardSessionCookie := middleware.LandlordDashboardSessionCookie(w, newSessionToken, newExpiryTime)
	if !createLandlordDashboardSessionCookie {
		logs.Logs(logErr, "Failed to get session cookie for landlord dashboard. Redirecting to landlord login page")
		http.Redirect(w, r, "/login/landlord?authenticationError=UNAUTHORIZED+401:+Error+authenticating+landlord.+Failed+to+get+session+cookie", http.StatusSeeOther)
		return
	}
	createLandordDashboardCSRFTokenCookie := middleware.LandlordDashboardCSRFTokenCookie(w, newCsrfToken, newExpiryTime)
	if !createLandordDashboardCSRFTokenCookie {
		logs.Logs(logErr, "Failed to get CSRF token cookie for landlord dashboard. Redirecting to landlord login page")
		http.Redirect(w, r, "/login/landlord?authenticationError=UNAUTHORIZED+401:+Error+authenticating+landlord.+Failed+to+get+CSRF+token+cookie", http.StatusSeeOther)
		return
	}

	// set cookies to landlord dashboard tenants page
	createLandlordDashboardTenantSessionCookie := middleware.LandlordDashboardTenantsSessionCookie(w, newSessionToken, newExpiryTime)
	if !createLandlordDashboardTenantSessionCookie {
		logs.Logs(logErr, "Failed to create session for the landlord tenants dashboard page. Redirecting back to login page")
		http.Redirect(w, r, "/login/landlord?authenticationError=UNAUTHORIZED+401:+Error+authenticating+landlord.+Failed+to+create+session+cookie", http.StatusSeeOther)
		return
	}
	createLandlordDashboardTenantCSRFTokenCookie := middleware.LandlordDashboardTenantsCSRFTokenCookie(w, newCsrfToken, newExpiryTime)
	if !createLandlordDashboardTenantCSRFTokenCookie {
		logs.Logs(logErr, "Failed to create CSRF token for the landlord tenants dashboard page. Redirecting back to login page")
		http.Redirect(w, r, "/login/landlord?authenticationError=UNAUTHORIZED+401:+Error+authenticating+landlord.+Failed+to+create+CSRF+token+cookie", http.StatusSeeOther)
	}

	// set cookies to landlord messages page
	createLandlordMessagesDashboardSessionCookie := middleware.LandlordMessagesDashboardSessionCookie(w, newSessionToken, newExpiryTime)
	if !createLandlordMessagesDashboardSessionCookie {
		logs.Logs(logErr, "Failed to create session for the landlord messages dashboard page. Redirecting back to login page")
		http.Redirect(w, r, "/login/landlord?authenticationError=UNAUTHORIZED+401:+Error+authenticating+landlord.+Failed+to+create+session+cookie", http.StatusSeeOther)
		return
	}
	createLandlordMessagesDashboardCSRFTokenCookie := middleware.LandlordMessagesDashboardCSRFTokenCookie(w, newCsrfToken, newExpiryTime)
	if !createLandlordMessagesDashboardCSRFTokenCookie {
		logs.Logs(logErr, "Failed to create CSRF token for the landlord messages dashboard page. Redirecting back to login page")
		http.Redirect(w, r, "/login/landlord?authenticationError=UNAUTHORIZED+401:+Error+authenticating+landlord.+Failed+to+create+CSRF+token+cookie", http.StatusSeeOther)
	}

	// set cookies to logout
	logoutSessionCookie := middleware.LogoutLandlordSessionCookie(w, newSessionToken)
	if !logoutSessionCookie {
		logs.Logs(logErr, "Failed to create session cookie for landlord. Redirecting to landlord login page")
		http.Redirect(w, r, "/login/landlord?authenticationError=UNAUTHORIZED+401:+Error+authenticating+landlord.+Failed+to+create+session+cookie", http.StatusSeeOther)
		return
	}
	logoutCSRFTokenCookie := middleware.LogoutLandlordCSRFTokenCookie(w, newCsrfToken)
	if !logoutCSRFTokenCookie {
		logs.Logs(logErr, "Failed to create CSRF token cookie for landlord. Redirecting to landlord login page")
		http.Redirect(w, r, "/login/landlord?authenticationError=UNAUTHORIZED+401:+Error+authenticating+landlord.+Failed+to+create+CSRF+token+cookie", http.StatusSeeOther)
		return
	}

	// get landlord tenant names
	encryptedEncryptedTenantNames, err := db.GetTenantsByLandlordEmail(landlordEmail)
	if err != nil {
		logs.Logs(logErr, fmt.Sprintf("Failed to get tenants: %s", err.Error()))
		http.Error(w, fmt.Sprintf("Failed to get tenants: %s", err.Error()), http.StatusInternalServerError)
		return
	}

	// decrypt tenant names
	var showTenantNames []ShowLandlordTenants

	for _, encryptedTenantName := range encryptedEncryptedTenantNames {
		var showTenantName ShowLandlordTenants
		tenantName, err := utils.Decrypt(encryptedTenantName.EncryptTenantName)
		if err != nil {
			logs.Logs(logErr, fmt.Sprintf("Failed to decrypt tenant name: %s", err.Error()))
			http.Error(w, fmt.Sprintf("Failed to decrypt tenant name: %s", err.Error()), http.StatusInternalServerError)
			return
		}
		showTenantName.DecryptTenantName = string(tenantName)
		showTenantNames = append(showTenantNames, showTenantName)
	}

	err = Templates.ExecuteTemplate(w, "landlordMessages.html", showTenantNames)
	if err != nil {
		logs.Logs(logErr, fmt.Sprintf("Unable to load landlord dashboard: %s", err.Error()))
		http.Error(w, fmt.Sprintf("Unable to load landlord dashboard: %s", err.Error()), http.StatusInternalServerError)
	}
}
