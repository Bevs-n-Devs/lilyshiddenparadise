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
	// Tenancy form error messages
	EvictedError         string
	ConvictedError       string
	VehicleError         string
	ChildrenError        string
	RefusedRentError     string
	UnstableIncomeError  string
	ConfirmPasswordError string
}

// TODO: Create a struct for the tenancy form data
