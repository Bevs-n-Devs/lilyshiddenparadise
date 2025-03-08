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

// struct for tenancy form error messages
type FormError struct {
	EvictedError        string
	ConvictedError      string
	VehicleError        string
	ChildrenError       string
	RefusedRentError    string
	UnstableIncomeError string
}

// TODO: Create a struct for the tenancy form data
