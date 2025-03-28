package handlers

import (
	"fmt"
	"net/http"

	"github.com/Bevs-n-Devs/lilyshiddenparadise/db"
	"github.com/Bevs-n-Devs/lilyshiddenparadise/logs"
	"github.com/Bevs-n-Devs/lilyshiddenparadise/middleware"
	"github.com/Bevs-n-Devs/lilyshiddenparadise/utils"
)

func LandlordNewTenant(w http.ResponseWriter, r *http.Request) {
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

	// TODO! Set cookies for each available page via landlord tenant dashboard page

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

	// set new cookies for landlord tenant applications
	createLandlordTenantApplicationsSessionCookie := middleware.LandlordDashboardTenantApplicationsSessionCookie(w, newSessionToken, newExpiryTime)
	if !createLandlordTenantApplicationsSessionCookie {
		logs.Logs(logErr, "Failed to get session cookie for landlord tenant applications. Redirecting to landlord login page")
		http.Redirect(w, r, "/login/landlord?authenticationError=UNAUTHORIZED+401:+Error+authenticating+landlord.+Failed+to+get+session+cookie", http.StatusSeeOther)
		return
	}
	createLandordTenantApplictionsCSRFTokenCookie := middleware.LandlordDashboardTenantApplicationsCSRFTokenCookie(w, newCsrfToken, newExpiryTime)
	if !createLandordTenantApplictionsCSRFTokenCookie {
		logs.Logs(logErr, "Failed to get CSRF token cookie for landlord tenant applications. Redirecting to landlord login page")
		http.Redirect(w, r, "/login/landlord?authenticationError=UNAUTHORIZED+401:+Error+authenticating+landlord.+Failed+to+get+CSRF+token+cookie", http.StatusSeeOther)
		return
	}

	// set cookies to create new tenant page
	createNewTenantSessionCookie := middleware.LandlordNewTenantSessionCookie(w, newSessionToken, newExpiryTime)
	if !createNewTenantSessionCookie {
		logs.Logs(logErr, "Failed to create session cookie for landlord new tenant page. Redirecting back to login page")
		http.Redirect(w, r, "/login/landlord?authenticationError=UNAUTHORIZED+401:+Error+authenticating+landlord.+Failed+to+create+session+cookie", http.StatusSeeOther)
		return
	}
	createNewTenantCSRFTokenCookie := middleware.LandlordNewTenantCSRFTokenCookie(w, newCsrfToken, newExpiryTime)
	if !createNewTenantCSRFTokenCookie {
		logs.Logs(logErr, "Failed to create CSRF token cookie for landlord new tenant page. Redirecting back to login page")
		http.Redirect(w, r, "/login/landlord?authenticationError=UNAUTHORIZED+401:+Error+authenticating+landlord.+Failed+to+create+CSRF+token+cookie", http.StatusSeeOther)
		return
	}

	// set cookies to submit new tenant handler
	createSubmitNewTenantSessionCookie := middleware.LandlordSubmitNewTenantSessionCookie(w, newSessionToken, newExpiryTime)
	if !createSubmitNewTenantSessionCookie {
		logs.Logs(logErr, "Failed to create session cookie for landlord submit new tenant page. Redirecting back to login page")
		http.Redirect(w, r, "/login/landlord?authenticationError=UNAUTHORIZED+401:+Error+authenticating+landlord.+Failed+to+create+session+cookie", http.StatusSeeOther)
		return
	}
	createSubmitNewTenantCSRFTokenCookie := middleware.LandlordSubmitNewTenantCSRFTokenCookie(w, newCsrfToken, newExpiryTime)
	if !createSubmitNewTenantCSRFTokenCookie {
		logs.Logs(logErr, "Failed to create CSRF token cookie for landlord submit new tenant page. Redirecting back to login page")
		http.Redirect(w, r, "/login/landlord?authenticationError=UNAUTHORIZED+401:+Error+authenticating+landlord.+Failed+to+create+CSRF+token+cookie", http.StatusSeeOther)
		return
	}

	// set cookie to logout landlord
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

	// direct user to protected landlord tenants page
	err = Templates.ExecuteTemplate(w, "newTenants.html", nil)
	if err != nil {
		logs.Logs(logErr, fmt.Sprintf("Unable to load create new tenant page: %s", err.Error()))
		http.Error(w, fmt.Sprintf("Unable to load create new tenant page: %s", err.Error()), http.StatusInternalServerError)
	}
}
