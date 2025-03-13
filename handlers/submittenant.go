package handlers

import (
	"fmt"
	"net/http"

	"github.com/Bevs-n-Devs/lilyshiddenparadise/db"
	"github.com/Bevs-n-Devs/lilyshiddenparadise/email"
	"github.com/Bevs-n-Devs/lilyshiddenparadise/logs"
	"github.com/Bevs-n-Devs/lilyshiddenparadise/utils"
)

func SubmitTenantForm(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		logs.Logs(logErr, fmt.Sprintf("Invalid request method: %s. Redirecting back to home page.", r.Method)) // this can be redirected back to the form page later on..
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	// Parse form data - vaidate the data can be taken from the form
	err := r.ParseForm()
	if err != nil {
		logs.Logs(logErr, fmt.Sprintf("Error parsing form data: %s", err.Error()))
		http.Error(w, fmt.Sprintf("Error parsing form data: %s", err.Error()), http.StatusBadRequest)
		return
	}

	// extract data from form
	fullName := r.FormValue("fullName")
	dateOfBirth := r.FormValue("dob")
	passportNumber := r.FormValue("passportNumber")
	phoneNumber := r.FormValue("phoneNumber")
	tenantEmail := r.FormValue("email")
	occupation := r.FormValue("occupation")
	employer := r.FormValue("employedBy")
	employerNumber := r.FormValue("workNumber")
	emergencyContactName := r.FormValue("emergencyName")
	emergencyContactNumber := r.FormValue("emergencyNumber")
	emergencyContactAddress := r.FormValue("emergencyAddress")
	ifEvicted := r.FormValue("evicted")
	evictedReason := r.FormValue("evictedReason")
	ifConvicted := r.FormValue("convicted")
	convictedReason := r.FormValue("convictedReason")
	smoke := r.FormValue("smoke")
	pets := r.FormValue("pets")
	ifVehicle := r.FormValue("vehicle")
	vehicleReg := r.FormValue("vehicleReg")
	haveChildren := r.FormValue("children")
	children := r.FormValue("children")
	refusedRent := r.FormValue("refusedRent")
	refusedRentReason := r.FormValue("refusedRentReason")
	unstableIncome := r.FormValue("unstableIncome")
	incomeReason := r.FormValue("incomeReason")

	// validate the user age - check if over 18
	if !utils.ValidateAge(dateOfBirth) {
		logs.Logs(logErr, "Invalid form data: User is not 18 years or older.")
		http.Redirect(w, r, "/tenancy-form?ageError=User+is+not+18+years+or+older", http.StatusSeeOther)
		return
	}

	if ifEvicted == "yes" {
		result := utils.CheckIfEvicted(ifEvicted, evictedReason)
		// redirect to form page with error message
		if !result {
			logs.Logs(logErr, "Invalid form data: Evicted reason not given.")
			http.Redirect(w, r, "/tenancy-form?evictedError=Evicted+reason+not+given", http.StatusSeeOther)
			return
		}
	}

	if ifConvicted == "yes" {
		result := utils.CheckIfConvicted(ifConvicted, convictedReason)
		if !result {
			logs.Logs(logErr, "Invalid form data: Convicted reason not given.")
			http.Redirect(w, r, "/tenancy-form?convictedError=Conviction+information+not+given", http.StatusSeeOther)
			return
		}
	}

	if ifVehicle == "yes" {
		result := utils.CheckIfVehicle(ifVehicle, vehicleReg)
		if !result {
			logs.Logs(logErr, "Invalid form data: Vehicle registration not given.")
			http.Redirect(w, r, "/tenancy-form?vehicleError=Vehicle+registration+not+given", http.StatusSeeOther)
			return
		}
	}

	if haveChildren == "yes" {
		result := utils.CheckIfHaveChildren(haveChildren, children)
		if !result {
			logs.Logs(logErr, "Invalid form data: Children information not given.")
			http.Redirect(w, r, "/tenancy-form?childrenError=Children+information+not+given", http.StatusSeeOther)
			return
		}
	}

	if refusedRent == "yes" {
		result := utils.CheckIfRefusedRent(refusedRent, refusedRentReason)
		if !result {
			logs.Logs(logErr, "Invalid form data: Refused rent reason not given.")
			http.Redirect(w, r, "/tenancy-form?refusedRentError=Reason+for+refusing+rent+not+given", http.StatusSeeOther)
			return
		}
	}

	if unstableIncome == "yes" {
		result := utils.CheckIfStableIncome(unstableIncome, incomeReason)
		if !result {
			logs.Logs(logErr, "Invalid form data: Income reason not given.")
			http.Redirect(w, r, "/tenancy-form?unstableIncomeError=Reasons+for+unstable+income+not+given", http.StatusSeeOther)
			return
		}
	}

	// save form data to database
	err = db.SaveTenantApplicationForm(
		fullName,
		dateOfBirth,
		passportNumber,
		phoneNumber,
		tenantEmail,
		occupation,
		employer,
		employerNumber,
		emergencyContactName,
		emergencyContactNumber,
		emergencyContactAddress,
		ifEvicted,
		evictedReason,
		ifConvicted,
		convictedReason,
		smoke,
		pets,
		ifVehicle,
		vehicleReg,
		haveChildren,
		children,
		refusedRent,
		refusedRentReason,
		unstableIncome,
		incomeReason,
	)
	if err != nil {
		logs.Logs(logErr, fmt.Sprintf("Error saving form data to database: %s. Redirecting back to tenancy form page.", err.Error()))
		http.Redirect(w, r, "/tenancy-form?dbError=Error+saving+form+data+to+database", http.StatusSeeOther)
		return
	}

	err = email.NotifyLandlordNewApplication()
	if err != nil {
		logs.Logs(logErr, fmt.Sprintf("Failed to send email notification to landlord: %s", err.Error()))
		http.Redirect(w, r, "/tenancy-form?emailError=Failed+to+send+email+notification+to+landlord", http.StatusSeeOther)
		return
	}

	err = email.NotifyTenantApplicationProcessing(tenantEmail)
	if err != nil {
		logs.Logs(logErr, fmt.Sprintf("Failed to send email notification to tenant: %s", err.Error()))
		http.Redirect(w, r, "/tenancy-form?emailError=Failed+to+send+email+notification+to+tenant", http.StatusSeeOther)
		return
	}

	// redirect to home page
	logs.Logs(logInfo, "Form data saved successfully. Redirecting to home page.")
	http.Redirect(w, r, "/", http.StatusSeeOther)
}
