package handlers

import (
	"fmt"
	"net/http"

	"github.com/Bevs-n-Devs/lilyshiddenparadise/db"
	"github.com/Bevs-n-Devs/lilyshiddenparadise/logs"
	"github.com/Bevs-n-Devs/lilyshiddenparadise/middleware"
	"github.com/Bevs-n-Devs/lilyshiddenparadise/utils"
)

func LandlordTenantApplications(w http.ResponseWriter, r *http.Request) {
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

	// get any error messages
	validationError := r.URL.Query().Get("validationError")
	data := ErrorMessages{
		ValidationError: validationError,
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

	// set new cookies for landlord manage applications
	createLandlordManageApplicationsSessionCookie := middleware.LandlordManageApplicationsSessionCookie(w, newSessionToken, newExpiryTime)
	if !createLandlordManageApplicationsSessionCookie {
		logs.Logs(logErr, "Failed to get session cookie for landlord manage applications. Redirecting to landlord login page")
		http.Redirect(w, r, "/login/landlord?authenticationError=UNAUTHORIZED+401:+Error+authenticating+landlord.+Failed+to+get+session+cookie", http.StatusSeeOther)
		return
	}
	createLandordManageApplicationsCSRFTokenCookie := middleware.LandlordManageApplicationsCSRFTokenCookie(w, newCsrfToken, newExpiryTime)
	if !createLandordManageApplicationsCSRFTokenCookie {
		logs.Logs(logErr, "Failed to get CSRF token cookie for landlord manage applications. Redirecting to landlord login page")
		http.Redirect(w, r, "/login/landlord?authenticationError=UNAUTHORIZED+401:+Error+authenticating+landlord.+Failed+to+get+CSRF+token+cookie", http.StatusSeeOther)
		return
	}

	// set new cookies for landlord tenants dashboard page
	createLandlordDashboardTenantsSessionCookie := middleware.LandlordDashboardTenantsSessionCookie(w, newSessionToken, newExpiryTime)
	if !createLandlordDashboardTenantsSessionCookie {
		logs.Logs(logErr, "Failed to get session cookie for landlord tenants dashboard. Redirecting to landlord login page")
		http.Redirect(w, r, "/login/landlord?authenticationError=UNAUTHORIZED+401:+Error+authenticating+landlord.+Failed+to+get+session+cookie", http.StatusSeeOther)
		return
	}
	createLandlordDashboardTenantsCSRFTokenCookie := middleware.LandlordDashboardTenantsCSRFTokenCookie(w, newCsrfToken, newExpiryTime)
	if !createLandlordDashboardTenantsCSRFTokenCookie {
		logs.Logs(logErr, "Failed to get CSRF token cookie for landlord tenants dashboard. Redirecting to landlord login page")
		http.Redirect(w, r, "/login/landlord?authenticationError=UNAUTHORIZED+401:+Error+authenticating+landlord.+Failed+to+get+CSRF+token+cookie", http.StatusSeeOther)
		return
	}

	// set new cookies to message dashboard page
	createLanlordMessageDashboardSessionCookies := middleware.LandlordMessagesDashboardSessionCookie(w, newSessionToken, newExpiryTime)
	if !createLanlordMessageDashboardSessionCookies {
		logs.Logs(logErr, "Failed to get session cookie for landlord message dashboard. Redirecting to landlord login page")
		http.Redirect(w, r, "/login/landlord?authenticationError=UNAUTHORIZED+401:+Error+authenticating+landlord.+Failed+to+get+session+cookie", http.StatusSeeOther)
		return
	}
	createLanlordMessageDashboardCSRFTokenCookie := middleware.LandlordMessagesDashboardCSRFTokenCookie(w, newCsrfToken, newExpiryTime)
	if !createLanlordMessageDashboardCSRFTokenCookie {
		logs.Logs(logErr, "Failed to get CSRF token cookie for landlord message dashboard. Redirecting to landlord login page")
		http.Redirect(w, r, "/login/landlord?authenticationError=UNAUTHORIZED+401:+Error+authenticating+landlord.+Failed+to+get+CSRF+token+cookie", http.StatusSeeOther)
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

	// get tenant applications from database
	getTenantApplications, err := db.GetAllTenantApplications()
	if err != nil {
		logs.Logs(logErr, fmt.Sprintf("Failed to get tenant applications: %s", err.Error()))
		http.Error(w, fmt.Sprintf("Failed to get tenant applications: %s", err.Error()), http.StatusInternalServerError)
		return
	}

	showTenantApplications := []ShowLandlordApplications{}

	// loop through tenant applications and decrypt data
	logs.Logs(logInfo, "Decrypting tenant applications...")
	for index := range getTenantApplications {
		var convertedData ShowLandlordApplications

		convertedData.ID = getTenantApplications[index].ID
		convertedData.Status = getTenantApplications[index].Status
		convertedData.CreatedAt = getTenantApplications[index].CreatedAt

		getTenantApplications[index].FullName, err = utils.Decrypt(getTenantApplications[index].FullName)
		if err != nil {
			logs.Logs(logErr, fmt.Sprintf("Failed to decrypt tenant application full name: %s", err.Error()))
			http.Error(w, fmt.Sprintf("Failed to decrypt tenant application full name: %s", err.Error()), http.StatusInternalServerError)
			return
		}
		convertedData.FullName = string(getTenantApplications[index].FullName)

		getTenantApplications[index].Dob, err = utils.Decrypt(getTenantApplications[index].Dob)
		if err != nil {
			logs.Logs(logErr, fmt.Sprintf("Failed to decrypt tenant application date of birth: %s", err.Error()))
			http.Error(w, fmt.Sprintf("Failed to decrypt tenant application date of birth: %s", err.Error()), http.StatusInternalServerError)
			return
		}
		convertedData.Dob = string(getTenantApplications[index].Dob)

		getTenantApplications[index].PassportNumber, err = utils.Decrypt(getTenantApplications[index].PassportNumber)
		if err != nil {
			logs.Logs(logErr, fmt.Sprintf("Failed to decrypt tenant application passport number: %s", err.Error()))
			http.Error(w, fmt.Sprintf("Failed to decrypt tenant application passport number: %s", err.Error()), http.StatusInternalServerError)
			return
		}
		convertedData.PassportNumber = string(getTenantApplications[index].PassportNumber)

		getTenantApplications[index].PhoneNumber, err = utils.Decrypt(getTenantApplications[index].PhoneNumber)
		if err != nil {
			logs.Logs(logErr, fmt.Sprintf("Failed to decrypt tenant application phone number: %s", err.Error()))
			http.Error(w, fmt.Sprintf("Failed to decrypt tenant application phone number: %s", err.Error()), http.StatusInternalServerError)
			return
		}
		convertedData.PhoneNumber = string(getTenantApplications[index].PhoneNumber)

		getTenantApplications[index].Email, err = utils.Decrypt(getTenantApplications[index].Email)
		if err != nil {
			logs.Logs(logErr, fmt.Sprintf("Failed to decrypt tenant application email: %s", err.Error()))
			http.Error(w, fmt.Sprintf("Failed to decrypt tenant application email: %s", err.Error()), http.StatusInternalServerError)
			return
		}
		convertedData.Email = string(getTenantApplications[index].Email)

		getTenantApplications[index].Occupation, err = utils.Decrypt(getTenantApplications[index].Occupation)
		if err != nil {
			logs.Logs(logErr, fmt.Sprintf("Failed to decrypt tenant application occupation: %s", err.Error()))
			http.Error(w, fmt.Sprintf("Failed to decrypt tenant application occupation: %s", err.Error()), http.StatusInternalServerError)
			return
		}
		convertedData.Occupation = string(getTenantApplications[index].Occupation)

		getTenantApplications[index].Employer, err = utils.Decrypt(getTenantApplications[index].Employer)
		if err != nil {
			logs.Logs(logErr, fmt.Sprintf("Failed to decrypt tenant application employer: %s", err.Error()))
			http.Error(w, fmt.Sprintf("Failed to decrypt tenant application employer: %s", err.Error()), http.StatusInternalServerError)
			return
		}
		convertedData.Employer = string(getTenantApplications[index].Employer)

		getTenantApplications[index].EmployerNumber, err = utils.Decrypt(getTenantApplications[index].EmployerNumber)
		if err != nil {
			logs.Logs(logErr, fmt.Sprintf("Failed to decrypt tenant application employer number: %s", err.Error()))
			http.Error(w, fmt.Sprintf("Failed to decrypt tenant application employer number: %s", err.Error()), http.StatusInternalServerError)
			return
		}
		convertedData.EmployerNumber = string(getTenantApplications[index].EmployerNumber)

		getTenantApplications[index].EmergencyContact, err = utils.Decrypt(getTenantApplications[index].EmergencyContact)
		if err != nil {
			logs.Logs(logErr, fmt.Sprintf("Failed to decrypt tenant application emergency contact: %s", err.Error()))
			http.Error(w, fmt.Sprintf("Failed to decrypt tenant application emergency contact: %s", err.Error()), http.StatusInternalServerError)
			return
		}
		convertedData.EmergencyContact = string(getTenantApplications[index].EmergencyContact)

		getTenantApplications[index].EmergencyContactNumber, err = utils.Decrypt(getTenantApplications[index].EmergencyContactNumber)
		if err != nil {
			logs.Logs(logErr, fmt.Sprintf("Failed to decrypt tenant application emergency contact number: %s", err.Error()))
			http.Error(w, fmt.Sprintf("Failed to decrypt tenant application emergency contact number: %s", err.Error()), http.StatusInternalServerError)
			return
		}
		convertedData.EmergencyContactNumber = string(getTenantApplications[index].EmergencyContactNumber)

		getTenantApplications[index].Evicted, err = utils.Decrypt(getTenantApplications[index].Evicted)
		if err != nil {
			logs.Logs(logErr, fmt.Sprintf("Failed to decrypt if tenant has been evicted: %s", err.Error()))
			http.Error(w, fmt.Sprintf("Failed to decrypt if tenant has been evicted: %s", err.Error()), http.StatusInternalServerError)
			return
		}
		convertedData.Evicted = string(getTenantApplications[index].Evicted)

		getTenantApplications[index].EvictedReason, err = utils.Decrypt(getTenantApplications[index].EvictedReason)
		if err != nil {
			logs.Logs(logErr, fmt.Sprintf("Failed to decrypt the reason for tenant eviction: %s", err.Error()))
			http.Error(w, fmt.Sprintf("Failed to decrypt the reason for tenant eviction: %s", err.Error()), http.StatusInternalServerError)
			return
		}
		convertedData.EvictedReason = string(getTenantApplications[index].EvictedReason)

		getTenantApplications[index].Convicted, err = utils.Decrypt(getTenantApplications[index].Convicted)
		if err != nil {
			logs.Logs(logErr, fmt.Sprintf("Failed to decrypt if tenant has been convicted: %s", err.Error()))
			http.Error(w, fmt.Sprintf("Failed to decrypt if tenant has been convicted: %s", err.Error()), http.StatusInternalServerError)
			return
		}
		convertedData.Convicted = string(getTenantApplications[index].Convicted)

		getTenantApplications[index].ConvictedReason, err = utils.Decrypt(getTenantApplications[index].ConvictedReason)
		if err != nil {
			logs.Logs(logErr, fmt.Sprintf("Failed to decrypt the reason for tenant conviction: %s", err.Error()))
			http.Error(w, fmt.Sprintf("Failed to decrypt the reason for tenant conviction: %s", err.Error()), http.StatusInternalServerError)
			return
		}
		convertedData.ConvictedReason = string(getTenantApplications[index].ConvictedReason)

		getTenantApplications[index].Smoke, err = utils.Decrypt(getTenantApplications[index].Smoke)
		if err != nil {
			logs.Logs(logErr, fmt.Sprintf("Failed to decrypt if tenant smokes: %s", err.Error()))
			http.Error(w, fmt.Sprintf("Failed to decrypt if tenant smokes: %s", err.Error()), http.StatusInternalServerError)
			return
		}
		convertedData.Smoke = string(getTenantApplications[index].Smoke)

		getTenantApplications[index].Pets, err = utils.Decrypt(getTenantApplications[index].Pets)
		if err != nil {
			logs.Logs(logErr, fmt.Sprintf("Failed to decrypt if tenant has pets: %s", err.Error()))
			http.Error(w, fmt.Sprintf("Failed to decrypt if tenant has pets: %s", err.Error()), http.StatusInternalServerError)
			return
		}
		convertedData.Pets = string(getTenantApplications[index].Pets)

		getTenantApplications[index].Vehicle, err = utils.Decrypt(getTenantApplications[index].Vehicle)
		if err != nil {
			logs.Logs(logErr, fmt.Sprintf("Failed to decrypt if tenant has a vehicle: %s", err.Error()))
			http.Error(w, fmt.Sprintf("Failed to decrypt if tenant has a vehicle: %s", err.Error()), http.StatusInternalServerError)
			return
		}
		convertedData.Vehicle = string(getTenantApplications[index].Vehicle)

		getTenantApplications[index].VehicleReg, err = utils.Decrypt(getTenantApplications[index].VehicleReg)
		if err != nil {
			logs.Logs(logErr, fmt.Sprintf("Failed to decrypt tenant vehicle registration: %s", err.Error()))
			http.Error(w, fmt.Sprintf("Failed to decrypt tenant vehicle registration: %s", err.Error()), http.StatusInternalServerError)
			return
		}
		convertedData.VehicleReg = string(getTenantApplications[index].VehicleReg)

		getTenantApplications[index].HaveChildren, err = utils.Decrypt(getTenantApplications[index].HaveChildren)
		if err != nil {
			logs.Logs(logErr, fmt.Sprintf("Failed to decrypt if tenant has children: %s", err.Error()))
			http.Error(w, fmt.Sprintf("Failed to decrypt if tenant has children: %s", err.Error()), http.StatusInternalServerError)
			return
		}
		convertedData.HaveChildren = string(getTenantApplications[index].HaveChildren)

		getTenantApplications[index].Children, err = utils.Decrypt(getTenantApplications[index].Children)
		if err != nil {
			logs.Logs(logErr, fmt.Sprintf("Failed to decrypt tenant children: %s", err.Error()))
			http.Error(w, fmt.Sprintf("Failed to decrypt tenant children: %s", err.Error()), http.StatusInternalServerError)
			return
		}
		convertedData.Children = string(getTenantApplications[index].Children)

		getTenantApplications[index].RefusedRent, err = utils.Decrypt(getTenantApplications[index].RefusedRent)
		if err != nil {
			logs.Logs(logErr, fmt.Sprintf("Failed to decrypt if tenant has refused to pay rent: %s", err.Error()))
			http.Error(w, fmt.Sprintf("Failed to decrypt if tenant has refused to pay rent: %s", err.Error()), http.StatusInternalServerError)
			return
		}
		convertedData.RefusedRent = string(getTenantApplications[index].RefusedRent)

		getTenantApplications[index].RefusedReason, err = utils.Decrypt(getTenantApplications[index].RefusedReason)
		if err != nil {
			logs.Logs(logErr, fmt.Sprintf("Failed to decrypt the reason for tenant refusing to pay rent: %s", err.Error()))
			http.Error(w, fmt.Sprintf("Failed to decrypt the reason for tenant refusing to pay rent: %s", err.Error()), http.StatusInternalServerError)
			return
		}
		convertedData.RefusedReason = string(getTenantApplications[index].RefusedReason)

		getTenantApplications[index].UnstableIncome, err = utils.Decrypt(getTenantApplications[index].UnstableIncome)
		if err != nil {
			logs.Logs(logErr, fmt.Sprintf("Failed to decrypt if tenant has unstable income: %s", err.Error()))
			http.Error(w, fmt.Sprintf("Failed to decrypt if tenant has unstable income: %s", err.Error()), http.StatusInternalServerError)
			return
		}
		convertedData.UnstableIncome = string(getTenantApplications[index].UnstableIncome)

		getTenantApplications[index].UnstableReason, err = utils.Decrypt(getTenantApplications[index].UnstableReason)
		if err != nil {
			logs.Logs(logErr, fmt.Sprintf("Failed to decrypt the reason for tenant having unstable income: %s", err.Error()))
			http.Error(w, fmt.Sprintf("Failed to decrypt the reason for tenant having unstable income: %s", err.Error()), http.StatusInternalServerError)
			return
		}
		convertedData.UnstableReason = string(getTenantApplications[index].UnstableReason)

		// append data to showTenanryApplications slice
		showTenantApplications = append(showTenantApplications, convertedData)
	}

	showData := struct {
		TenantApplications []ShowLandlordApplications
		ErrorMessage       string
	}{
		TenantApplications: showTenantApplications,
		ErrorMessage:       data.ValidationError,
	}

	// direct user to protected tenant applications
	err = Templates.ExecuteTemplate(w, "tenantApplications.html", showData)
	if err != nil {
		logs.Logs(logErr, fmt.Sprintf("Unable to load landlord tenant applications: %s", err.Error()))
		http.Error(w, fmt.Sprintf("Unable to load landlord tenant applications: %s", err.Error()), http.StatusInternalServerError)
	}
}
