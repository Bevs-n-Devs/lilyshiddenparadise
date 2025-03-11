package handlers

import "html/template"

const (
	localPort = ":9001"
	logInfo   = 1
	logWarn   = 2
	logErr    = 3
)

var (
	Templates       *template.Template // global Templates variable to hold all HTML templates
	LandlordSession string
	LandlordCSRF    string
)

// struct for error messages to display user via HTML template
type ErrorMessages struct {
	// HTTP server error messages
	BadRequestError     string
	NotFoundError       string
	AuthenticationError string
	InternalServerError string
	CookieError         string
	// Tenancy form error messages
	AgeError             string
	EvictedError         string
	ConvictedError       string
	VehicleError         string
	ChildrenError        string
	RefusedRentError     string
	UnstableIncomeError  string
	ConfirmPasswordError string
	DatabaseError        string
}

// TODO: Create a struct for the tenancy form data

type ShowLandlordApplications struct {
	Status                 string `json:"status"`
	FullName               string `json:"full_name"`
	Dob                    string `json:"dob"`
	PassportNumber         string `json:"passport_number"`
	PhoneNumber            string `json:"phone_number"`
	Email                  string `json:"email"`
	Occupation             string `json:"occupation"`
	Employer               string `json:"employer"`
	EmployerNumber         string `json:"employer_number"`
	EmergencyContact       string `json:"emergency_contact"`
	EmergencyContactNumber string `json:"emergency_contact_number"`
	Evicted                string `json:"if_evicted"`
	EvictedReason          string `json:"evicted_reason"`
	Convicted              string `json:"if_convicted"`
	ConvictedReason        string `json:"convicted_reason"`
	Smoke                  string `json:"smoke"`
	Pets                   string `json:"pets"`
	Vehicle                string `json:"vehicle"`
	VehicleReg             string `json:"vehicle_reg"`
	HaveChildren           string `json:"have_children"`
	Children               string `json:"children"`
	RefusedRent            string `json:"refused_rent"`
	RefusedReason          string `json:"refused_reason"`
	UnstableIncome         string `json:"unstable_income"`
	UnstableReason         string `json:"unstable_reason"`
	CreatedAt              string `json:"created_at"`
}
