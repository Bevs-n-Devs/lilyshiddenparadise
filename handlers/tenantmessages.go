package handlers

import (
	"fmt"
	"net/http"

	"github.com/Bevs-n-Devs/lilyshiddenparadise/db"
	"github.com/Bevs-n-Devs/lilyshiddenparadise/logs"
	"github.com/Bevs-n-Devs/lilyshiddenparadise/middleware"
	"github.com/Bevs-n-Devs/lilyshiddenparadise/utils"
)

func TenantMessages(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
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

	// update the tenant's session token, CSRF token and expiry time in the database
	// this will be done for each request
	newSessionToken, newCsrfToken, newExpiryTime, err := db.UpdateTenantSessionTokens(tenantEmail)
	if err != nil {
		logs.Logs(logErr, fmt.Sprintf("Error updating tenant session tokens: %s. Redirecting to tenant login page", err.Error()))
		http.Redirect(w, r, "/login/tenant?authenticationError=UNAUTHORIZED+401:+Error+authenticating+tenant.+Failed+to+update+session+tokens", http.StatusSeeOther)
		return
	}

	// TODO! Set cookies for each available page via tenant dashboard page

	// set session cookies for tenant dashboard
	createTenantDashboardSessionCookie := middleware.TenantDashboardSessionCookie(w, newSessionToken, newExpiryTime)
	if !createTenantDashboardSessionCookie {
		logs.Logs(logErr, "Failed to create session cookie. Redirecting back to tenant login page...")
		http.Redirect(w, r, "/login/tenant?internalServerError=INTERNAL+SERVER+ERROR+500:+Failed+to+create+session+cookie", http.StatusInternalServerError)
		return
	}
	createTenantDashboardCsrfCookie := middleware.TenantDashboardCSRFTokenCookie(w, newCsrfToken, newExpiryTime)
	if !createTenantDashboardCsrfCookie {
		logs.Logs(logErr, "Failed to create CSRF cookie. Redirecting back to tenant login page...")
		http.Redirect(w, r, "/login/tenant?internalServerError=INTERNAL+SERVER+ERROR+500:+Failed+to+create+CSRF+cookie", http.StatusInternalServerError)
		return
	}

	// set cookie for tenant account page
	createTenantAccountSessionCookie := middleware.TenantDashboardAccountSessionCookie(w, newSessionToken, newExpiryTime)
	if !createTenantAccountSessionCookie {
		logs.Logs(logErr, "Failed to create session cookie. Redirecting back to tenant login page...")
		http.Redirect(w, r, "/login/tenant?internalServerError=INTERNAL+SERVER+ERROR+500:+Failed+to+create+session+cookie", http.StatusInternalServerError)
		return
	}
	createTenantAccountCsrfCookie := middleware.TenantDashboardAccountCSRFTokenCookie(w, newCsrfToken, newExpiryTime)
	if !createTenantAccountCsrfCookie {
		logs.Logs(logErr, "Failed to create CSRF cookie. Redirecting back to tenant login page...")
		http.Redirect(w, r, "/login/tenant?internalServerError=INTERNAL+SERVER+ERROR+500:+Failed+to+create+CSRF+cookie", http.StatusInternalServerError)
		return
	}

	// set session cookies to send message to landlord
	createSendLandlordMessageSessionCookie := middleware.SendMessageToLandlordSessionCookie(w, newSessionToken, newExpiryTime)
	if !createSendLandlordMessageSessionCookie {
		logs.Logs(logErr, "Failed to create session cookie. Redirecting back to tenant login page...")
		http.Redirect(w, r, "/login/tenant?internalServerError=INTERNAL+SERVER+ERROR+500:+Failed+to+create+session+cookie", http.StatusInternalServerError)
		return
	}
	createSendLandlordMessageCsrfCookie := middleware.SendMessageToLandlordCSRFTokenCookie(w, newCsrfToken, newExpiryTime)
	if !createSendLandlordMessageCsrfCookie {
		logs.Logs(logErr, "Failed to create CSRF cookie. Redirecting back to tenant login page...")
		http.Redirect(w, r, "/login/tenant?internalServerError=INTERNAL+SERVER+ERROR+500:+Failed+to+create+CSRF+cookie", http.StatusInternalServerError)
		return
	}

	// set session to submit message route
	createSubmitMessageToLandlordSessionCookie := middleware.SubmitMessageToLandlordSessionCookie(w, newSessionToken, newExpiryTime)
	if !createSubmitMessageToLandlordSessionCookie {
		logs.Logs(logErr, "Failed to create session cookie. Redirecting back to tenant login page...")
		http.Redirect(w, r, "/login/tenant?internalServerError=INTERNAL+SERVER+ERROR+500:+Failed+to+create+session+cookie", http.StatusInternalServerError)
		return
	}
	createSubmitMessageToLandlordCsrfTokenCookie := middleware.SubmitMessageToLandlordCSRFTokenCookie(w, newCsrfToken, newExpiryTime)
	if !createSubmitMessageToLandlordCsrfTokenCookie {
		logs.Logs(logErr, "Failed to create CSRF cookie. Redirecting back to tenant login page...")
		http.Redirect(w, r, "/login/tenant?internalServerError=INTERNAL+SERVER+ERROR+500:+Failed+to+create+CSRF+cookie", http.StatusInternalServerError)
		return
	}

	// set session cookies to logout tenant
	createTenantLogoutSessionCookie := middleware.LogoutTenantSessionCookie(w, newSessionToken)
	if !createTenantLogoutSessionCookie {
		logs.Logs(logErr, "Failed to create session cookie. Redirecting back to tenant login page...")
		http.Redirect(w, r, "/login/tenant?internalServerError=INTERNAL+SERVER+ERROR+500:+Failed+to+create+session+cookie", http.StatusInternalServerError)
		return
	}
	createTenantLogoutCsrfCookie := middleware.LogoutTenantCSRFTokenCookie(w, newCsrfToken)
	if !createTenantLogoutCsrfCookie {
		logs.Logs(logErr, "Failed to create CSRF cookie. Redirecting back to tenant login page...")
		http.Redirect(w, r, "/login/tenant?internalServerError=INTERNAL+SERVER+ERROR+500:+Failed+to+create+CSRF+cookie", http.StatusInternalServerError)
		return
	}

	// TODO: get all messages between landlord and tenant via their email => ID

	err = Templates.ExecuteTemplate(w, "messageLandlord.html", nil)
	if err != nil {
		logs.Logs(logErr, fmt.Sprintf("Unable to load message landlord: %s", err.Error()))
		http.Error(w, fmt.Sprintf("Unable to load message landlord: %s", err.Error()), http.StatusInternalServerError)
	}
}
