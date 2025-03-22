package handlers

import (
	"fmt"
	"net/http"

	"github.com/Bevs-n-Devs/lilyshiddenparadise/db"
	"github.com/Bevs-n-Devs/lilyshiddenparadise/email"
	"github.com/Bevs-n-Devs/lilyshiddenparadise/logs"
	"github.com/Bevs-n-Devs/lilyshiddenparadise/middleware"
	"github.com/Bevs-n-Devs/lilyshiddenparadise/utils"
)

func LandlordManageApplications(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		logs.Logs(logErr, "Invalid request method. Redirecting back to landlord login page.")
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

	// TODO! Set new cookies for each available page via tenant applications page

	// set new cookies for landlord dashboard - after successful form submission
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

	// set new cookies for landlord manage applications - if values from submission form are missing
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

	// parse form data
	err = r.ParseForm()
	if err != nil {
		logs.Logs(logErr, fmt.Sprintf("Error parsing form data: %s", err.Error()))
		http.Error(w, fmt.Sprintf("Error parsing form data: %s", err.Error()), http.StatusInternalServerError)
		return
	}

	// extract data from form
	applicationId := r.FormValue("applicationId")
	applicationResult := r.FormValue("applicationResult")
	roomType := r.FormValue("roomType")
	moveInDate := r.FormValue("moveInDate")
	rentDue := r.FormValue("rentDue")
	monthlyRent := r.FormValue("monthlyRent")
	currency := r.FormValue("currency")

	// if the application has been denied
	if applicationResult == "denied" {
		err = db.UpdateTenantApplicationStatus(applicationId, applicationResult)
		if err != nil {
			logs.Logs(logErr, fmt.Sprintf("Error updating tenant application status: %s", err.Error()))
			http.Error(w, fmt.Sprintf("Error updating tenant application status: %s", err.Error()), http.StatusInternalServerError)
			return
		}
		http.Redirect(w, r, "/landlord/dashboard", http.StatusSeeOther)
	}

	// validate the form data
	err = utils.ValidateManageTenantApplication(applicationResult, roomType, moveInDate, rentDue, monthlyRent, currency)
	if err != nil {
		logs.Logs(logErr, fmt.Sprintf("Error validating tenant application form data: %s. Redirecting back to landlord tenant applications page", err.Error()))
		http.Redirect(w, r, "/landlord/dashboard/tenant-applications?validationError=BAD+REQUEST+400:+Error+validating+tenant+application+form+data.+Missing+parameters.", http.StatusSeeOther)
		return
	}

	// update the tenant application status
	err = db.UpdateTenantApplicationStatus(applicationId, applicationResult)
	if err != nil {
		logs.Logs(logErr, fmt.Sprintf("Error updating tenant application status: %s", err.Error()))
		http.Error(w, fmt.Sprintf("Error updating tenant application status: %s", err.Error()), http.StatusInternalServerError)
		return
	}

	// TODO: get email & passport number via applicationID from database
	encryptEmail, encryptPassportNumber, err := db.GetTenantEmailAndPassportNumberViaApplicationID(applicationId)
	if err != nil {
		logs.Logs(logErr, fmt.Sprintf("Error getting tenant email & passport number from database: %s", err.Error()))
		http.Redirect(w, r, "/landlord/dashboard?internalServerError=INTERNAL+SERVER+ERROR+500:+Error+getting+tenant+email+&+passport+number+from+database", http.StatusSeeOther)
		return
	}

	tenantEmail, err := (utils.Decrypt([]byte(encryptEmail)))
	if err != nil {
		logs.Logs(logErr, fmt.Sprintf("Error decrypting email: %s", err.Error()))
		http.Error(w, fmt.Sprintf("Error decrypting email: %s", err.Error()), http.StatusInternalServerError)
		return
	}
	strTenantEmail := string(tenantEmail)

	tenantPassportNumber, err := (utils.Decrypt([]byte(encryptPassportNumber)))
	if err != nil {
		logs.Logs(logErr, fmt.Sprintf("Error decrypting passport number: %s", err.Error()))
		http.Error(w, fmt.Sprintf("Error decrypting passport number: %s", err.Error()), http.StatusInternalServerError)
		return
	}
	strTenantPassportNumber := string(tenantPassportNumber)

	// TODO: generate new tenant username & password
	tenantUsername, tenantPassword, err := utils.GenerateTenantUsernamePassportNumberAndPassword(strTenantEmail, strTenantPassportNumber)
	if err != nil {
		logs.Logs(logErr, fmt.Sprintf("Error generating tenant username & password: %s", err.Error()))
		http.Error(w, fmt.Sprintf("Error generating tenant username & password: %s", err.Error()), http.StatusInternalServerError)
		return
	}

	logs.Logs(logInfo, fmt.Sprintf("New tenant username: %s, New tenant password: %s", tenantUsername, tenantPassword))

	// TODO: send email to tenant with new username & password
	err = email.NotifyTenantNewAccount(tenantUsername, tenantPassword, roomType, moveInDate, rentDue, monthlyRent, currency)
	if err != nil {
		logs.Logs(logErr, fmt.Sprintf("Failed to send email notification to tenant: %s", err.Error()))
		http.Error(w, fmt.Sprintf("Failed to send email notification to tenant: %s", err.Error()), http.StatusInternalServerError)
		return
	}

	// TODO: send email to landlord with tenant username & password
	err = email.NotifyLandlordNewAccount(tenantUsername, tenantPassword, roomType, moveInDate, rentDue, monthlyRent, currency)
	if err != nil {
		logs.Logs(logErr, fmt.Sprintf("Failed to send email notification to landlord: %s", err.Error()))
		http.Error(w, fmt.Sprintf("Failed to send email notification to landlord: %s", err.Error()), http.StatusInternalServerError)
		return
	}

	// TODO: Save tenant to database
	err = db.CreateNewTenant(tenantUsername, tenantPassword, roomType, moveInDate, rentDue, monthlyRent, currency)
	if err != nil {
		logs.Logs(logErr, fmt.Sprintf("Failed to save tenant to database: %s", err.Error()))
		http.Error(w, fmt.Sprintf("Failed to save tenant to database: %s", err.Error()), http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/landlord/dashboard", http.StatusSeeOther)
}
