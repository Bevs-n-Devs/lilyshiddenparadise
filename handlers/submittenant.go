package handlers

import (
	"fmt"
	"net/http"

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
	email := r.FormValue("email")
	occupation := r.FormValue("occupation")
	employer := r.FormValue("employedBy")
	employeeNumber := r.FormValue("workNumber")
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

	// log form data (save to databse later on)
	logMessage := fmt.Sprintf("Form data - Full Name: %s, Date of Birth: %s, Passport Number: %s, Phone Number: %s, Email: %s, Occupation: %s, Employer: %s, Employee Number: %s, Emergency Contact Name: %s, Emergency Contact Number: %s, Emergency Contact Address: %s, If Evicted: %s, Evicted Reason: %s, If Convicted: %s, Convicted Reason: %s, Smoke: %s, Pets: %s, If Vehicle: %s, Vehicle Reg: %s, Have Children: %s, Children: %s, Refused Rent: %s, Refused Rent Reason: %s, Instabel Income: %s, Income Reason: %s", fullName, dateOfBirth, passportNumber, phoneNumber, email, occupation, employer, employeeNumber, emergencyContactName, emergencyContactNumber, emergencyContactAddress, ifEvicted, evictedReason, ifConvicted, convictedReason, smoke, pets, ifVehicle, vehicleReg, haveChildren, children, refusedRent, refusedRentReason, unstableIncome, incomeReason)
	logs.Logs(logInfo, logMessage)

	// redirect to home page
	logs.Logs(logInfo, "Form data saved successfully. Redirecting to home page.")
	http.Redirect(w, r, "/", http.StatusSeeOther)
}
