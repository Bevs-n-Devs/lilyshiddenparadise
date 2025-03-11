package db

import "database/sql"

const (
	logWarning = 2
	logErr     = 3
	logDb      = 4
	logDbErr   = 5
)

var (
	db *sql.DB // global DB variable to hold DB connection
)

type GetLandlordApplications struct {
	Status                 string `json:"status"`
	FullName               []byte `json:"full_name"`
	Dob                    []byte `json:"dob"`
	PassportNumber         []byte `json:"passport_number"`
	PhoneNumber            []byte `json:"phone_number"`
	Email                  []byte `json:"email"`
	Occupation             []byte `json:"occupation"`
	Employer               []byte `json:"employer"`
	EmployerNumber         []byte `json:"employer_number"`
	EmergencyContact       []byte `json:"emergency_contact"`
	EmergencyContactNumber []byte `json:"emergency_contact_number"`
	Evicted                []byte `json:"if_evicted"`
	EvictedReason          []byte `json:"evicted_reason"`
	Convicted              []byte `json:"if_convicted"`
	ConvictedReason        []byte `json:"convicted_reason"`
	Smoke                  []byte `json:"smoke"`
	Pets                   []byte `json:"pets"`
	Vehicle                []byte `json:"vehicle"`
	VehicleReg             []byte `json:"vehicle_reg"`
	HaveChildren           []byte `json:"have_children"`
	Children               []byte `json:"children"`
	RefusedRent            []byte `json:"refused_rent"`
	RefusedReason          []byte `json:"refused_reason"`
	UnstableIncome         []byte `json:"unstable_income"`
	UnstableReason         []byte `json:"unstable_reason"`
	CreatedAt              string `json:"created_at"`
}
