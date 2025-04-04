package testutil

import (
	"time"
)

// DBInterface defines the interface for database operations
// This allows us to mock the database for testing
type DBInterface interface {
	// Landlord operations
	CreateNewLandlord(landlordEmail, landlordPassword string) error
	AuthenticateLandlord(email, password string) (bool, error)
	UpdateLandlordSessionTokens(email string) (string, string, time.Time, error)
	GetEmailFromLandlordSessionToken(sessionToken string) (string, error)
	ValidateLandlordSessionToken(email, sessionToken string) (bool, error)
	ValidateLandlordCSRFToken(email, csrfToken string) (bool, error)
	LogoutLandlord(email string) error
	GetLandlordIdByEmail(email string) (int, error)

	// Tenant operations
	CreateNewTenant(tenantEmail, tenantPassword, roomType, moveInDate, rentDue, monthlyRent, currency string) error
	AuthenticateTenant(username, password string) (bool, error)
	UpdateTenantSessionTokens(hashEmail string) (string, string, time.Time, error)
	GetHashedEmailFromTenantSessionToken(sessionToken string) (string, error)
	ValidateTenantSessionToken(hashEmail, sessionToken string) (bool, error)
	ValidateTenantCSRFToken(hashEmail, csrfToken string) (bool, error)
	LogoutTenant(hashEmail string) error
	GetTenantInformationByHashEmail(hashEmail string) (Tenant, error)

	// Tenant application operations
	SaveTenantApplicationForm(
		fullName,
		dateOfBirth,
		passportNumber,
		phoneNumber,
		email,
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
		incomeReason string,
	) error
	GetAllTenantApplications() ([]TenantApplication, error)
	UpdateTenantApplicationStatus(id string, status string) error

	// Message operations
	SendMessage(senderID int, senderType string, receiverID int, receiverType string, message string) error
	GetMessageBetweenLandlordsAndTenant(tenantID string) ([]Message, error)
}
