package handlers

import (
	"fmt"
	"net/http"
	"os"

	"github.com/Bevs-n-Devs/lilyshiddenparadise/db"
	"github.com/Bevs-n-Devs/lilyshiddenparadise/email"
	"github.com/Bevs-n-Devs/lilyshiddenparadise/env"
	"github.com/Bevs-n-Devs/lilyshiddenparadise/logs"
	"github.com/Bevs-n-Devs/lilyshiddenparadise/middleware"
	"github.com/Bevs-n-Devs/lilyshiddenparadise/utils"
)

func SendMessageToLandlord(w http.ResponseWriter, r *http.Request) {
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

	// update the tenant's session token, CSRF token and expiry time in the database
	// this will be done for each request
	newSessionToken, newCsrfToken, newExpiryTime, err := db.UpdateTenantSessionTokens(tenantEmail)
	if err != nil {
		logs.Logs(logErr, fmt.Sprintf("Error updating tenant session tokens: %s. Redirecting to tenant login page", err.Error()))
		http.Redirect(w, r, "/login/tenant?authenticationError=UNAUTHORIZED+401:+Error+authenticating+tenant.+Failed+to+update+session+tokens", http.StatusSeeOther)
		return
	}

	// TODO! Set cookies for each available page
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

	// TODO: get form data and send message to landlord via email => ID
	// parse form data
	err = r.ParseForm()
	if err != nil {
		logs.Logs(logErr, fmt.Sprintf("Error parsing form data: %s", err.Error()))
		http.Error(w, fmt.Sprintf("Error parsing form data: %s", err.Error()), http.StatusInternalServerError)
		return
	}

	// extract tenant message from form
	tenantMessage := r.FormValue("tenantMessage")

	// TODO: get landlord ID via landlord email
	// get landlord email via environment variable
	if os.Getenv("LANDLORD_EMAIL") == "" {
		logs.Logs(logWarn, "Could not get landlord email from hosting platform. Loading from .env file...")
		err := env.LoadEnv("env/.env")
		if err != nil {
			logs.Logs(logErr, fmt.Sprintf("Could not load environment variables from .env file: %s", err.Error()))
			http.Error(w, fmt.Sprintf("Could not load environment variables from .env file: %s", err.Error()), http.StatusInternalServerError)
			return
		}
	}

	landlordEmail := os.Getenv("LANDLORD_EMAIL")
	if landlordEmail == "" {
		logs.Logs(logErr, "Landlord email is empty!")
		http.Error(w, "Landlord email is empty!", http.StatusInternalServerError)
		return
	}

	// get landlord id
	landlordId, err := db.GetLandlordIdByEmail(landlordEmail)
	if err != nil {
		logs.Logs(logErr, fmt.Sprintf("Failed to get landlord ID: %s", err.Error()))
		http.Error(w, fmt.Sprintf("Failed to get landlord ID: %s", err.Error()), http.StatusInternalServerError)
		return
	}

	// TODO: get tenant ID via tenant email
	tenantId, err := db.GetTenantIdByEmail(tenantEmail)
	if err != nil {
		logs.Logs(logErr, fmt.Sprintf("Failed to get tenant ID: %s", err.Error()))
		http.Error(w, fmt.Sprintf("Failed to get tenant ID: %s", err.Error()), http.StatusInternalServerError)
		return
	}

	// TODO: save message to database [lhp_messages table]
	err = db.SendMessage(tenantId, TENANT, landlordId, LANDLORD, tenantMessage)
	if err != nil {
		logs.Logs(logErr, fmt.Sprintf("Failed to send message to landlord: %s", err.Error()))
		http.Error(w, fmt.Sprintf("Failed to send message to landlord: %s", err.Error()), http.StatusInternalServerError)
		return
	}

	// TODO: send email notification to landlord
	encryptedTeanntName, err := db.GetTenantNameByHashEmail(tenantEmail)
	if err != nil {
		logs.Logs(logErr, fmt.Sprintf("Failed to get encrypted tenant name: %s", err.Error()))
		http.Error(w, fmt.Sprintf("Failed to get encrypted tenant name: %s", err.Error()), http.StatusInternalServerError)
		return
	}

	tenantFullName, err := utils.Decrypt([]byte(encryptedTeanntName))
	if err != nil {
		logs.Logs(logErr, fmt.Sprintf("Failed to decrypt tenant name: %s", err.Error()))
		http.Error(w, fmt.Sprintf("Failed to decrypt tenant name: %s", err.Error()), http.StatusInternalServerError)
		return
	}

	// send email to landlord
	err = email.NotifyLandlordNewMessageFromTenant(string(tenantFullName), landlordEmail, tenantMessage)
	if err != nil {
		logs.Logs(logErr, fmt.Sprintf("Failed to send email to landlord: %s", err.Error()))
		http.Error(w, fmt.Sprintf("Failed to send email to landlord: %s", err.Error()), http.StatusInternalServerError)
		return
	}

	// TODO: redirect backl to messages page if message sent successfully
	http.Redirect(w, r, "/tenant/dashboard/messages", http.StatusSeeOther)
}
