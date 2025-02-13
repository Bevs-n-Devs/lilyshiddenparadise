package handlers

import (
	"fmt"
	"net/http"

	"github.com/Bevs-n-Devs/lilyshiddenparadise/logs"
)

func SubmitForm(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		logs.Logs(3, fmt.Sprintf("Invalid request method: %s. Redirecting back to home page.", r.Method)) // this can be redirected back to the form page later on..
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	// Parse form data - vaidate the data can be taken from the form
	err := r.ParseForm()
	if err != nil {
		logs.Logs(3, fmt.Sprintf("Error parsing form data: %s", err.Error()))
		http.Error(w, "Error parsing form data: "+err.Error(), http.StatusBadRequest)
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
	instabelIncome := r.FormValue("instabelIncome")
	incomeReason := r.FormValue("incomeReason")

	if ifEvicted == "yes" {
		result := checkIfEvicted(ifEvicted, evictedReason)
		// redirect to form page with error message
		if !result {
			logs.Logs(3, fmt.Sprintf("Invalid form data: %s", err.Error()))
			http.Redirect(w, r, "/form", http.StatusSeeOther)
			return
		}
	}

	if ifConvicted == "yes" {
		result := checkIfConvicted(ifConvicted, convictedReason)
		if !result {
			logs.Logs(3, fmt.Sprintf("Invalid form data: %s", err.Error()))
			http.Redirect(w, r, "/form", http.StatusSeeOther)
			return
		}
	}

	if ifVehicle == "yes" {
		result := checkIfVehiclke(ifVehicle, vehicleReg)
		if !result {
			logs.Logs(3, fmt.Sprintf("Invalid form data: %s", err.Error()))
			http.Redirect(w, r, "/form", http.StatusSeeOther)
			return
		}
	}

	if haveChildren == "yes" {
		result := checkIfHaveChildren(haveChildren, children)
		if !result {
			logs.Logs(3, fmt.Sprintf("Invalid form data: %s", err.Error()))
			http.Redirect(w, r, "/form", http.StatusSeeOther)
			return
		}
	}

	if refusedRent == "yes" {
		result := checkIfRefusedRent(refusedRent, refusedRentReason)
		if !result {
			logs.Logs(3, fmt.Sprintf("Invalid form data: %s", err.Error()))
			http.Redirect(w, r, "/form", http.StatusSeeOther)
			return
		}
	}

	if instabelIncome == "yes" {
		result := checkIfStableIncome(instabelIncome, incomeReason)
		if !result {
			logs.Logs(3, fmt.Sprintf("Invalid form data: %s", err.Error()))
			http.Redirect(w, r, "/form", http.StatusSeeOther)
			return
		}
	}

	// log form data (save to databse later on)
	logMessage := fmt.Sprintf("Form data - Full Name: %s, Date of Birth: %s, Passport Number: %s, Phone Number: %s, Email: %s, Occupation: %s, Employer: %s, Employee Number: %s, Emergency Contact Name: %s, Emergency Contact Number: %s, Emergency Contact Address: %s, If Evicted: %s, Evicted Reason: %s, If Convicted: %s, Convicted Reason: %s, Smoke: %s, Pets: %s, If Vehicle: %s, Vehicle Reg: %s, Have Children: %s, Children: %s, Refused Rent: %s, Refused Rent Reason: %s, Instabel Income: %s, Income Reason: %s", fullName, dateOfBirth, passportNumber, phoneNumber, email, occupation, employer, employeeNumber, emergencyContactName, emergencyContactNumber, emergencyContactAddress, ifEvicted, evictedReason, ifConvicted, convictedReason, smoke, pets, ifVehicle, vehicleReg, haveChildren, children, refusedRent, refusedRentReason, instabelIncome, incomeReason)
	logs.Logs(1, logMessage)

	// redirect to home page
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

/*
Checks if ifEvicted is yes and evictedReason is empty;
returns false if invalid.
*/
// Returns true if ifEvicted is yes and evictedReason is not empty
func checkIfEvicted(ifEvicted, evictedReason string) bool {
	if ifEvicted == "yes" && evictedReason == "" {
		return false
	}
	return true
}

/*
Checks if ifConvicted is yes and convictedReason is empty;
returns false if invalid.
*/
// Returns true if ifConvicted is yes and convictedReason is not empty
func checkIfConvicted(ifConvicted, convictedReason string) bool {
	if ifConvicted == "yes" && convictedReason == "" {
		return false
	}
	return true
}

/*
Checks if ifVehicle is yes and vehicleReg is empty;
returns false if invalid.
*/
// Returns true if ifVehicle is yes and vehicleReg is not empty
func checkIfVehiclke(ifVehicle, vehicleReg string) bool {
	if ifVehicle == "yes" && vehicleReg == "" {
		return false
	}
	return true
}

/*
Checks if haveChildren is yes and children is empty;
returns false if invalid.
*/
// Returns true if haveChildren is yes and children is not empty
func checkIfHaveChildren(haveChildren, children string) bool {
	if haveChildren == "yes" && children == "" {
		return false
	}
	return true
}

/*
Checks if refusedRent is yes and refusedRentReason is empty;
returns false if invalid.
*/
// Returns true if refusedRent is yes and refusedRentReason is not empty
func checkIfRefusedRent(refusedRent, refusedRentReason string) bool {
	if refusedRent == "yes" && refusedRentReason == "" {
		return false
	}
	return true
}

/*
Checks if instabelIncome is yes and incomeReason is empty;
returns false if invalid.
*/
// Returns true if stableIncome is no and incomeReason is not empty
func checkIfStableIncome(instabelIncome, incomeReason string) bool {
	if instabelIncome == "yes" && incomeReason == "" {
		return false
	}
	return true
}
