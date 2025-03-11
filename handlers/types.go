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
