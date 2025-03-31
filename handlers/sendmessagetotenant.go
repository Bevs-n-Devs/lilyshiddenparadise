package handlers

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/Bevs-n-Devs/lilyshiddenparadise/db"
	"github.com/Bevs-n-Devs/lilyshiddenparadise/email"
	"github.com/Bevs-n-Devs/lilyshiddenparadise/logs"
	"github.com/Bevs-n-Devs/lilyshiddenparadise/middleware"
	"github.com/Bevs-n-Devs/lilyshiddenparadise/utils"
)

func SendMessageToTenant(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
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

	// get landlord id from email
	landlordId, err := db.GetLandlordIdByEmail(landlordEmail)
	if err != nil {
		logs.Logs(logErr, fmt.Sprintf("Error getting landlord ID: %s. Redirecting to landlord login page", err.Error()))
		http.Redirect(w, r, "/login/landlord?authenticationError=UNAUTHORIZED+401:+Error+authenticating+landlord.+Failed+to+get+landlord+ID", http.StatusSeeOther)
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
	tenantID := strings.TrimPrefix(r.URL.Path, "/landlord/send-message/")
	// convert tenant id to int
	tenantIdInt, err := strconv.Atoi(tenantID)
	if err != nil {
		logs.Logs(logErr, fmt.Sprintf("Error converting tenant ID to int: %s. Redirecting to landlord login page", err.Error()))
		http.Redirect(w, r, "/login/landlord?authenticationError=UNAUTHORIZED+401:+Error+authenticating+landlord.+Failed+to+convert+tenant+ID+to+int", http.StatusSeeOther)
		return
	}

	// TODO! Set cookies for each available page
	createLandlordTenantMessagesSessionCookie := middleware.LandlordTenantMessagesSessionCookie(w, tenantID, newSessionToken, newExpiryTime)
	if !createLandlordTenantMessagesSessionCookie {
		logs.Logs(logErr, "Error creating landlord tenant messages session cookie. Redirecting to landlord login page")
		http.Redirect(w, r, "/login/landlord?authenticationError=UNAUTHORIZED+401:+Error+authenticating+landlord.+Failed+to+create+landlord+tenant+messages+session+cookie", http.StatusSeeOther)
		return
	}
	createLandlordTenantMessagesCSRFTokenCookie := middleware.LandlordTenantMessagesCSRFTokenCookie(w, tenantID, newCsrfToken, newExpiryTime)
	if !createLandlordTenantMessagesCSRFTokenCookie {
		logs.Logs(logErr, "Error creating landlord tenant messages CSRF token cookie. Redirecting to landlord login page")
		http.Redirect(w, r, "/login/landlord?authenticationError=UNAUTHORIZED+401:+Error+authenticating+landlord.+Failed+to+create+landlord+tenant+messages+CSRF+token+cookie", http.StatusSeeOther)
		return
	}

	// TODO: extract data from form
	err = r.ParseForm()
	if err != nil {
		logs.Logs(logErr, fmt.Sprintf("Error parsing form data: %s", err.Error()))
		http.Error(w, fmt.Sprintf("Error parsing form data: %s", err.Error()), http.StatusInternalServerError)
		return
	}

	// TODO: save message to database
	landlordMessage := r.FormValue("landlordMessage")

	err = db.SendMessage(landlordId, LANDLORD, tenantIdInt, TENANT, landlordMessage)
	if err != nil {
		logs.Logs(logErr, fmt.Sprintf("Error sending message to tenant: %s", err.Error()))
		http.Error(w, fmt.Sprintf("Error sending message to tenant: %s", err.Error()), http.StatusInternalServerError)
		return
	}

	// TODO: send email notification to tenant
	// get tenant email
	encryptTenantEmail, err := db.GetTenantEncryptedEmailById(tenantIdInt)
	if err != nil {
		logs.Logs(logErr, fmt.Sprintf("Error getting tenant email: %s", err.Error()))
		http.Error(w, fmt.Sprintf("Error getting tenant email: %s", err.Error()), http.StatusInternalServerError)
		return
	}

	tenantEmail, err := utils.Decrypt([]byte(encryptTenantEmail))
	if err != nil {
		logs.Logs(logErr, fmt.Sprintf("Error decrypting tenant email: %s", err.Error()))
		http.Error(w, fmt.Sprintf("Error decrypting tenant email: %s", err.Error()), http.StatusInternalServerError)
		return
	}

	err = email.NotifyTenantNewMessageFromLandlord(string(tenantEmail), landlordMessage)
	if err != nil {
		logs.Logs(logErr, fmt.Sprintf("Error sending email notification to tenant: %s", err.Error()))
		http.Error(w, fmt.Sprintf("Error sending email notification to tenant: %s", err.Error()), http.StatusInternalServerError)
		return
	}

	// redirect back to selected tenant messages page
	logs.Logs(logInfo, "Message successfully sent to tenant. Redirecting back to tenant messages page.")
	http.Redirect(w, r, "/landlord/dashboard/messages/tenant/"+tenantID, http.StatusSeeOther)
}
