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
	ID                      int    `json:"id"`
	Status                  string `json:"status"`
	FullName                []byte `json:"full_name"`
	Dob                     []byte `json:"dob"`
	PassportNumber          []byte `json:"passport_number"`
	PhoneNumber             []byte `json:"phone_number"`
	Email                   []byte `json:"email"`
	Occupation              []byte `json:"occupation"`
	Employer                []byte `json:"employer"`
	EmployerNumber          []byte `json:"employer_number"`
	EmergencyContact        []byte `json:"emergency_contact"`
	EmergencyContactNumber  []byte `json:"emergency_contact_number"`
	EmergencyContactAddress []byte `json:"emergency_contact_address"`
	Evicted                 []byte `json:"if_evicted"`
	EvictedReason           []byte `json:"evicted_reason"`
	Convicted               []byte `json:"if_convicted"`
	ConvictedReason         []byte `json:"convicted_reason"`
	Smoke                   []byte `json:"smoke"`
	Pets                    []byte `json:"pets"`
	Vehicle                 []byte `json:"vehicle"`
	VehicleReg              []byte `json:"vehicle_reg"`
	HaveChildren            []byte `json:"have_children"`
	Children                []byte `json:"children"`
	RefusedRent             []byte `json:"refused_rent"`
	RefusedReason           []byte `json:"refused_reason"`
	UnstableIncome          []byte `json:"unstable_income"`
	UnstableReason          []byte `json:"unstable_reason"`
	CreatedAt               string `json:"created_at"`
}

type GetTenantInformation struct {
	Email       []byte `json:"encrypt_email"`
	RoomType    []byte `json:"encrypt_room_type"`
	MoveInDate  []byte `json:"encrypt_move_in_date"`
	RentDueDate []byte `json:"encrypt_rent_due"`
	MonthlyRent []byte `json:"encrypt_monthly_rent"`
	Currency    string `json:"currency"`
}

type LandlordTenants struct {
	ID                int    `json:"id"`
	EncryptTenantName []byte `json:"encrypt_tenant_name"`
}
