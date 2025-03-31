package handlers

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/Bevs-n-Devs/lilyshiddenparadise/db"
	"github.com/Bevs-n-Devs/lilyshiddenparadise/logs"
	"github.com/Bevs-n-Devs/lilyshiddenparadise/middleware"
	"github.com/Bevs-n-Devs/lilyshiddenparadise/utils"
)

func LandlordTenantMessages(w http.ResponseWriter, r *http.Request) {
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

	// get the URL prefix of the current page
	tenantID := strings.TrimPrefix(r.URL.Path, "/landlord/dashboard/messages/tenant/")

	// get all messages between landlord and tenant
	messages, err := db.GetMessageBetweenLandlordsAndTenant(tenantID)
	if err != nil {
		logs.Logs(logErr, fmt.Sprintf("Failed to get messages between landlords and tenants: %s", err.Error()))
		http.Error(w, fmt.Sprintf("Failed to get messages between landlords and tenants: %s", err.Error()), http.StatusInternalServerError)
		return
	}

	var showMessages []ShowMessages
	for _, message := range messages {
		var showMessage ShowMessages
		showMessage.LandlordID = message.LandlordID
		showMessage.TenantID = message.TenantID

		showMessage.SenderID = message.SenderID
		showMessage.SenderType = message.SenderType
		showMessage.ReceiverID = message.ReceiverID
		showMessage.ReceiverType = message.ReceiverType
		showMessage.SentAt = message.SentAt

		decryptMessage, err := utils.Decrypt(message.EncryptMessage)
		if err != nil {
			logs.Logs(logErr, fmt.Sprintf("Failed to decrypt message: %s", err.Error()))
			http.Error(w, fmt.Sprintf("Failed to decrypt message: %s", err.Error()), http.StatusInternalServerError)
			return
		}

		showMessage.Message = string(decryptMessage)

		showMessages = append(showMessages, showMessage)
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

	// set cookie to divert to selected tenant messages page
	createLandlordTenantMessagesSessionCookie := middleware.LandlordTenantMessagesSessionCookie(w, tenantID, newSessionToken, newExpiryTime)
	if !createLandlordTenantMessagesSessionCookie {
		logs.Logs(logErr, "Failed to get session cookie for landlord tenant messages. Redirecting to landlord login page")
		http.Redirect(w, r, "/login/landlord?authenticationError=UNAUTHORIZED+401:+Error+authenticating+landlord.+Failed+to+get+session+cookie", http.StatusSeeOther)
		return
	}
	createLandlordTenantMessagesCSRFTokenCookie := middleware.LandlordTenantMessagesCSRFTokenCookie(w, tenantID, newCsrfToken, newExpiryTime)
	if !createLandlordTenantMessagesCSRFTokenCookie {
		logs.Logs(logErr, "Failed to get CSRF token cookie for landlord tenant messages. Redirecting to landlord login page")
		http.Redirect(w, r, "/login/landlord?authenticationError=UNAUTHORIZED+401:+Error+authenticating+landlord.+Failed+to+get+CSRF+token+cookie", http.StatusSeeOther)
		return
	}

	// set cookie to process message for selected tenant
	createSubmitMessageFromLandlordSessionCookie := middleware.SubmitMessageFromLandlordSessionCookie(w, tenantID, newSessionToken, newExpiryTime)
	if !createSubmitMessageFromLandlordSessionCookie {
		logs.Logs(logErr, "Failed to get session cookie for landlord tenant messages. Redirecting to landlord login page")
		http.Redirect(w, r, "/login/landlord?authenticationError=UNAUTHORIZED+401:+Error+authenticating+landlord.+Failed+to+get+session+cookie", http.StatusSeeOther)
		return
	}
	createSubmitMessageFromLandlordCSRFTokenCookie := middleware.SubmitMessageFromLandlordCSRFTokenCookie(w, tenantID, newCsrfToken, newExpiryTime)
	if !createSubmitMessageFromLandlordCSRFTokenCookie {
		logs.Logs(logErr, "Failed to get CSRF token cookie for landlord tenant messages. Redirecting to landlord login page")
		http.Redirect(w, r, "/login/landlord?authenticationError=UNAUTHORIZED+401:+Error+authenticating+landlord.+Failed+to+get+CSRF+token+cookie", http.StatusSeeOther)
		return
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

	err = Templates.ExecuteTemplate(w, "messageTenant.html", showMessages)
	if err != nil {
		logs.Logs(logErr, fmt.Sprintf("Unable to load landlord dashboard: %s", err.Error()))
		http.Error(w, fmt.Sprintf("Unable to load landlord dashboard: %s", err.Error()), http.StatusInternalServerError)
	}
}
