package handlers

import "html/template"

const (
	localPort = ":9001"
	logInfo   = 1
	logWarn   = 2
	logErr    = 3
)

var (
	Templates *template.Template // global Templates variable to hold all HTML templates
)

// struct for error messages to display user via HTML template
type ErrorMessages struct {
	EvictedError         string
	ConvictedError       string
	VehicleError         string
	ChildrenError        string
	RefusedRentError     string
	UnstableIncomeError  string
	ConfirmPasswordError string
	BadRequestError      string
	NotFoundError        string
}

// TODO: Create a struct for the tenancy form data
