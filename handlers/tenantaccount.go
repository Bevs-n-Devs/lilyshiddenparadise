package handlers

import (
	"fmt"
	"net/http"

	"github.com/Bevs-n-Devs/lilyshiddenparadise/db"
	"github.com/Bevs-n-Devs/lilyshiddenparadise/logs"
	"github.com/Bevs-n-Devs/lilyshiddenparadise/middleware"
	"github.com/Bevs-n-Devs/lilyshiddenparadise/utils"
)

func TenantAccount(w http.ResponseWriter, r *http.Request) {
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

	// set session cookies to update tenant password
	createUpdateTenantPasswordSessionCookie := middleware.UpdateTenantPasswordSessionCookie(w, newSessionToken, newExpiryTime)
	if !createUpdateTenantPasswordSessionCookie {
		logs.Logs(logErr, "Failed to create session cookie. Redirecting back to tenant login page...")
		http.Redirect(w, r, "/login/tenant?internalServerError=INTERNAL+SERVER+ERROR+500:+Failed+to+create+session+cookie", http.StatusInternalServerError)
		return
	}
	createUpdateTenantPasswordCsrfCookie := middleware.UpdateTenantPasswordCSRFTokenCookie(w, newCsrfToken, newExpiryTime)
	if !createUpdateTenantPasswordCsrfCookie {
		logs.Logs(logErr, "Failed to create CSRF cookie. Redirecting back to tenant login page...")
		http.Redirect(w, r, "/login/tenant?internalServerError=INTERNAL+SERVER+ERROR+500:+Failed+to+create+CSRF+cookie", http.StatusInternalServerError)
		return
	}

	// TODO: get the tenants application details / tenancy agreement
	tenantInfo, err := db.GetTenantInformationByHashEmail(tenantEmail)
	if err != nil {
		logs.Logs(logErr, fmt.Sprintf("Failed to get tenant information: %s", err.Error()))
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	// get any error messages
	authenticationError := r.URL.Query().Get("authenticationError")

	// decrypt encrypted tenant information
	var showData ShowTenantInformation

	showData.Error.AuthenticationError = authenticationError
	showData.Currency = tenantInfo.Currency

	getEmail, err := utils.Decrypt(tenantInfo.Email)
	if err != nil {
		logs.Logs(logErr, fmt.Sprintf("Failed to decrypt tenant email: %s", err.Error()))
		http.Error(w, fmt.Sprintf("Failed to decrypt tenant email: %s", err.Error()), http.StatusInternalServerError)
		return
	}
	showData.Email = string(getEmail)

	getRoomType, err := utils.Decrypt(tenantInfo.RoomType)
	if err != nil {
		logs.Logs(logErr, fmt.Sprintf("Failed to decrypt tenant room type: %s", err.Error()))
		http.Error(w, fmt.Sprintf("Failed to decrypt tenant room type: %s", err.Error()), http.StatusInternalServerError)
		return
	}
	showData.RoomType = string(getRoomType)

	getMoveInDate, err := utils.Decrypt(tenantInfo.MoveInDate)
	if err != nil {
		logs.Logs(logErr, fmt.Sprintf("Failed to decrypt tenant move in date: %s", err.Error()))
		http.Error(w, fmt.Sprintf("Failed to decrypt tenant move in date: %s", err.Error()), http.StatusInternalServerError)
		return
	}
	showData.MoveInDate = string(getMoveInDate)

	getRentDue, err := utils.Decrypt(tenantInfo.RentDueDate)
	if err != nil {
		logs.Logs(logErr, fmt.Sprintf("Failed to decrypt tenant rent due: %s", err.Error()))
		http.Error(w, fmt.Sprintf("Failed to decrypt tenant rent due: %s", err.Error()), http.StatusInternalServerError)
		return
	}
	showData.RentDueDate = string(getRentDue)

	getMonthlyRent, err := utils.Decrypt(tenantInfo.MonthlyRent)
	if err != nil {
		logs.Logs(logErr, fmt.Sprintf("Failed to decrypt tenant monthly rent: %s", err.Error()))
		http.Error(w, fmt.Sprintf("Failed to decrypt tenant monthly rent: %s", err.Error()), http.StatusInternalServerError)
		return
	}
	showData.MonthlyRent = string(getMonthlyRent)

	// direct user to protected tenant account page
	err = Templates.ExecuteTemplate(w, "tenantAccount.html", showData)
	if err != nil {
		logs.Logs(logErr, fmt.Sprintf("Unable to load tenant account page: %s", err.Error()))
		http.Error(w, fmt.Sprintf("Unable to load tenant account page: %s", err.Error()), http.StatusInternalServerError)
	}
}
