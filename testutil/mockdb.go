package testutil

import (
	"database/sql"
	"errors"
	"sync"
	"time"
)

// MockDB is a mock implementation of the database for testing
type MockDB struct {
	mu                      sync.RWMutex
	landlords               map[string]*Landlord
	tenants                 map[string]*Tenant
	tenantApplications      map[int]*TenantApplication
	messages                []Message
	nextTenantAppID         int
	nextTenantID            int
	nextMessageID           int
	failNextOperation       bool
	failNextOperationReason string
}

// Landlord represents a landlord in the mock database
type Landlord struct {
	ID           int
	Email        string
	Password     string
	SessionToken string
	CSRFToken    string
	TokenExpiry  time.Time
}

// Tenant represents a tenant in the mock database
type Tenant struct {
	ID                int
	LandlordID        int
	HashEmail         string
	HashPassword      string
	EncryptTenantName []byte
	EncryptEmail      []byte
	EncryptPassword   []byte
	EncryptRoomType   []byte
	EncryptMoveInDate []byte
	EncryptRentDue    []byte
	EncryptMonthlyRent []byte
	Currency          string
	SessionToken      string
	CSRFToken         string
	TokenExpiry       time.Time
	CreatedAt         time.Time
}

// TenantApplication represents a tenant application in the mock database
type TenantApplication struct {
	ID                      int
	LandlordID              int
	HashFullName            string
	HashDob                 string
	HashPassportNumber      string
	HashEmail               string
	Status                  string
	EncryptFullName         []byte
	EncryptDob              []byte
	EncryptPassportNumber   []byte
	EncryptPhoneNumber      []byte
	EncryptEmail            []byte
	EncryptOccupation       []byte
	EncryptEmployer         []byte
	EncryptEmployerNumber   []byte
	EncryptEmergencyContact []byte
	EncryptEmergencyNumber  []byte
	EncryptEmergencyAddress []byte
	EncryptIfEvicted        []byte
	EncryptEvictedReason    []byte
	EncryptIfConvicted      []byte
	EncryptConvictedReason  []byte
	EncryptSmoke            []byte
	EncryptPets             []byte
	EncryptIfVehicle        []byte
	EncryptVehicleReg       []byte
	EncryptHaveChildren     []byte
	EncryptChildren         []byte
	EncryptRefusedRent      []byte
	EncryptRefusedReason    []byte
	EncryptUnstableIncome   []byte
	EncryptIncomeReason     []byte
	CreatedAt               time.Time
}

// Message represents a message in the mock database
type Message struct {
	ID             int
	SenderID       int
	SenderType     string
	ReceiverID     int
	ReceiverType   string
	EncryptMessage []byte
	SentAt         time.Time
}

// NewMockDB creates a new mock database
func NewMockDB() *MockDB {
	return &MockDB{
		landlords:          make(map[string]*Landlord),
		tenants:            make(map[string]*Tenant),
		tenantApplications: make(map[int]*TenantApplication),
		messages:           []Message{},
		nextTenantAppID:    1,
		nextTenantID:       1,
		nextMessageID:      1,
	}
}

// SetFailNextOperation sets the mock database to fail the next operation
func (m *MockDB) SetFailNextOperation(reason string) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.failNextOperation = true
	m.failNextOperationReason = reason
}

// checkFailNextOperation checks if the next operation should fail
func (m *MockDB) checkFailNextOperation() error {
	m.mu.Lock()
	defer m.mu.Unlock()
	if m.failNextOperation {
		m.failNextOperation = false
		reason := m.failNextOperationReason
		m.failNextOperationReason = ""
		return errors.New(reason)
	}
	return nil
}

// CreateNewLandlord creates a new landlord in the mock database
func (m *MockDB) CreateNewLandlord(landlordEmail, landlordPassword string) error {
	if err := m.checkFailNextOperation(); err != nil {
		return err
	}

	m.mu.Lock()
	defer m.mu.Unlock()

	// Check if landlord already exists
	if _, exists := m.landlords[landlordEmail]; exists {
		return errors.New("landlord already exists")
	}

	m.landlords[landlordEmail] = &Landlord{
		ID:       len(m.landlords) + 1,
		Email:    landlordEmail,
		Password: landlordPassword, // In a real implementation, this would be hashed
	}

	return nil
}

// AuthenticateLandlord checks if the provided email and password match
func (m *MockDB) AuthenticateLandlord(email, password string) (bool, error) {
	if err := m.checkFailNextOperation(); err != nil {
		return false, err
	}

	m.mu.RLock()
	defer m.mu.RUnlock()

	landlord, exists := m.landlords[email]
	if !exists {
		return false, sql.ErrNoRows
	}

	// In a real implementation, this would use bcrypt to compare hashed passwords
	if landlord.Password != password {
		return false, errors.New("invalid password")
	}

	return true, nil
}

// UpdateLandlordSessionTokens updates the session and CSRF tokens for a landlord
func (m *MockDB) UpdateLandlordSessionTokens(email string) (string, string, time.Time, error) {
	if err := m.checkFailNextOperation(); err != nil {
		return "", "", time.Time{}, err
	}

	m.mu.Lock()
	defer m.mu.Unlock()

	landlord, exists := m.landlords[email]
	if !exists {
		return "", "", time.Time{}, sql.ErrNoRows
	}

	sessionToken := "session_token_" + email
	csrfToken := "csrf_token_" + email
	expiry := time.Now().Add(30 * time.Minute)

	landlord.SessionToken = sessionToken
	landlord.CSRFToken = csrfToken
	landlord.TokenExpiry = expiry

	return sessionToken, csrfToken, expiry, nil
}

// GetEmailFromLandlordSessionToken gets the email associated with a session token
func (m *MockDB) GetEmailFromLandlordSessionToken(sessionToken string) (string, error) {
	if err := m.checkFailNextOperation(); err != nil {
		return "", err
	}

	m.mu.RLock()
	defer m.mu.RUnlock()

	for email, landlord := range m.landlords {
		if landlord.SessionToken == sessionToken {
			return email, nil
		}
	}

	return "", errors.New("user not found")
}

// ValidateLandlordSessionToken validates a landlord's session token
func (m *MockDB) ValidateLandlordSessionToken(email, sessionToken string) (bool, error) {
	if err := m.checkFailNextOperation(); err != nil {
		return false, err
	}

	m.mu.RLock()
	defer m.mu.RUnlock()

	landlord, exists := m.landlords[email]
	if !exists {
		return false, errors.New("user not found")
	}

	return landlord.SessionToken == sessionToken, nil
}

// ValidateLandlordCSRFToken validates a landlord's CSRF token
func (m *MockDB) ValidateLandlordCSRFToken(email, csrfToken string) (bool, error) {
	if err := m.checkFailNextOperation(); err != nil {
		return false, err
	}

	m.mu.RLock()
	defer m.mu.RUnlock()

	landlord, exists := m.landlords[email]
	if !exists {
		return false, errors.New("user not found")
	}

	return landlord.CSRFToken == csrfToken, nil
}

// LogoutLandlord removes a landlord's session token, CSRF token, and expiry time
func (m *MockDB) LogoutLandlord(email string) error {
	if err := m.checkFailNextOperation(); err != nil {
		return err
	}

	m.mu.Lock()
	defer m.mu.Unlock()

	landlord, exists := m.landlords[email]
	if !exists {
		return errors.New("user not found")
	}

	landlord.SessionToken = ""
	landlord.CSRFToken = ""
	landlord.TokenExpiry = time.Time{}

	return nil
}

// GetLandlordIdByEmail gets a landlord's ID by email
func (m *MockDB) GetLandlordIdByEmail(email string) (int, error) {
	if err := m.checkFailNextOperation(); err != nil {
		return 0, err
	}

	m.mu.RLock()
	defer m.mu.RUnlock()

	landlord, exists := m.landlords[email]
	if !exists {
		return 0, sql.ErrNoRows
	}

	return landlord.ID, nil
}

// SaveTenantApplicationForm saves a tenant application form
func (m *MockDB) SaveTenantApplicationForm(
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
) error {
	if err := m.checkFailNextOperation(); err != nil {
		return err
	}

	m.mu.Lock()
	defer m.mu.Unlock()

	// In a real implementation, this would hash and encrypt the data
	app := &TenantApplication{
		ID:                    m.nextTenantAppID,
		LandlordID:            1, // Default to landlord ID 1 for testing
		HashFullName:          "hash_" + fullName,
		HashDob:               "hash_" + dateOfBirth,
		HashPassportNumber:    "hash_" + passportNumber,
		HashEmail:             "hash_" + email,
		Status:                "pending",
		EncryptFullName:       []byte(fullName),
		EncryptDob:            []byte(dateOfBirth),
		EncryptPassportNumber: []byte(passportNumber),
		EncryptPhoneNumber:    []byte(phoneNumber),
		EncryptEmail:          []byte(email),
		EncryptOccupation:     []byte(occupation),
		EncryptEmployer:       []byte(employer),
		EncryptEmployerNumber: []byte(employerNumber),
		CreatedAt:             time.Now(),
		// Add other fields as needed
	}

	m.tenantApplications[m.nextTenantAppID] = app
	m.nextTenantAppID++

	return nil
}

// GetAllTenantApplications gets all tenant applications for a landlord
func (m *MockDB) GetAllTenantApplications() ([]TenantApplication, error) {
	if err := m.checkFailNextOperation(); err != nil {
		return nil, err
	}

	m.mu.RLock()
	defer m.mu.RUnlock()

	var applications []TenantApplication
	for _, app := range m.tenantApplications {
		applications = append(applications, *app)
	}

	return applications, nil
}

// UpdateTenantApplicationStatus updates the status of a tenant application
func (m *MockDB) UpdateTenantApplicationStatus(id string, status string) error {
	if err := m.checkFailNextOperation(); err != nil {
		return err
	}

	m.mu.Lock()
	defer m.mu.Unlock()

	appID := 0
	// Convert string ID to int
	for i := range m.tenantApplications {
		if i == appID {
			m.tenantApplications[i].Status = status
			return nil
		}
	}

	return errors.New("tenant application not found")
}

// CreateNewTenant creates a new tenant
func (m *MockDB) CreateNewTenant(tenantEmail, tenantPassword, roomType, moveInDate, rentDue, monthlyRent, currency string) error {
	if err := m.checkFailNextOperation(); err != nil {
		return err
	}

	m.mu.Lock()
	defer m.mu.Unlock()

	hashEmail := "hash_" + tenantEmail
	hashPassword := "hash_" + tenantPassword

	// Check if tenant already exists
	if _, exists := m.tenants[hashEmail]; exists {
		return errors.New("tenant already exists")
	}

	m.tenants[hashEmail] = &Tenant{
		ID:                m.nextTenantID,
		LandlordID:        1, // Default to landlord ID 1 for testing
		HashEmail:         hashEmail,
		HashPassword:      hashPassword,
		EncryptTenantName: []byte("Tenant " + tenantEmail),
		EncryptEmail:      []byte(tenantEmail),
		EncryptPassword:   []byte(tenantPassword),
		EncryptRoomType:   []byte(roomType),
		EncryptMoveInDate: []byte(moveInDate),
		EncryptRentDue:    []byte(rentDue),
		EncryptMonthlyRent: []byte(monthlyRent),
		Currency:          currency,
		CreatedAt:         time.Now(),
	}

	m.nextTenantID++

	return nil
}

// AuthenticateTenant checks if the provided username and password match
func (m *MockDB) AuthenticateTenant(username, password string) (bool, error) {
	if err := m.checkFailNextOperation(); err != nil {
		return false, err
	}

	m.mu.RLock()
	defer m.mu.RUnlock()

	hashUsername := "hash_" + username
	hashPassword := "hash_" + password

	tenant, exists := m.tenants[hashUsername]
	if !exists {
		return false, sql.ErrNoRows
	}

	if tenant.HashPassword != hashPassword {
		return false, errors.New("invalid password")
	}

	return true, nil
}

// UpdateTenantSessionTokens updates the session and CSRF tokens for a tenant
func (m *MockDB) UpdateTenantSessionTokens(hashEmail string) (string, string, time.Time, error) {
	if err := m.checkFailNextOperation(); err != nil {
		return "", "", time.Time{}, err
	}

	m.mu.Lock()
	defer m.mu.Unlock()

	tenant, exists := m.tenants[hashEmail]
	if !exists {
		return "", "", time.Time{}, sql.ErrNoRows
	}

	sessionToken := "session_token_" + hashEmail
	csrfToken := "csrf_token_" + hashEmail
	expiry := time.Now().Add(30 * time.Minute)

	tenant.SessionToken = sessionToken
	tenant.CSRFToken = csrfToken
	tenant.TokenExpiry = expiry

	return sessionToken, csrfToken, expiry, nil
}

// GetHashedEmailFromTenantSessionToken gets the hashed email associated with a session token
func (m *MockDB) GetHashedEmailFromTenantSessionToken(sessionToken string) (string, error) {
	if err := m.checkFailNextOperation(); err != nil {
		return "", err
	}

	m.mu.RLock()
	defer m.mu.RUnlock()

	for hashEmail, tenant := range m.tenants {
		if tenant.SessionToken == sessionToken {
			return hashEmail, nil
		}
	}

	return "", errors.New("user not found")
}

// ValidateTenantSessionToken validates a tenant's session token
func (m *MockDB) ValidateTenantSessionToken(hashEmail, sessionToken string) (bool, error) {
	if err := m.checkFailNextOperation(); err != nil {
		return false, err
	}

	m.mu.RLock()
	defer m.mu.RUnlock()

	tenant, exists := m.tenants[hashEmail]
	if !exists {
		return false, errors.New("user not found")
	}

	return tenant.SessionToken == sessionToken, nil
}

// ValidateTenantCSRFToken validates a tenant's CSRF token
func (m *MockDB) ValidateTenantCSRFToken(hashEmail, csrfToken string) (bool, error) {
	if err := m.checkFailNextOperation(); err != nil {
		return false, err
	}

	m.mu.RLock()
	defer m.mu.RUnlock()

	tenant, exists := m.tenants[hashEmail]
	if !exists {
		return false, errors.New("user not found")
	}

	return tenant.CSRFToken == csrfToken, nil
}

// LogoutTenant removes a tenant's session token, CSRF token, and expiry time
func (m *MockDB) LogoutTenant(hashEmail string) error {
	if err := m.checkFailNextOperation(); err != nil {
		return err
	}

	m.mu.Lock()
	defer m.mu.Unlock()

	tenant, exists := m.tenants[hashEmail]
	if !exists {
		return errors.New("user not found")
	}

	tenant.SessionToken = ""
	tenant.CSRFToken = ""
	tenant.TokenExpiry = time.Time{}

	return nil
}

// GetTenantInformationByHashEmail gets tenant information by hashed email
func (m *MockDB) GetTenantInformationByHashEmail(hashEmail string) (Tenant, error) {
	if err := m.checkFailNextOperation(); err != nil {
		return Tenant{}, err
	}

	m.mu.RLock()
	defer m.mu.RUnlock()

	tenant, exists := m.tenants[hashEmail]
	if !exists {
		return Tenant{}, sql.ErrNoRows
	}

	return *tenant, nil
}

// SendMessage sends a message between a landlord and tenant
func (m *MockDB) SendMessage(senderID int, senderType string, receiverID int, receiverType string, message string) error {
	if err := m.checkFailNextOperation(); err != nil {
		return err
	}

	m.mu.Lock()
	defer m.mu.Unlock()

	m.messages = append(m.messages, Message{
		ID:             m.nextMessageID,
		SenderID:       senderID,
		SenderType:     senderType,
		ReceiverID:     receiverID,
		ReceiverType:   receiverType,
		EncryptMessage: []byte(message),
		SentAt:         time.Now(),
	})

	m.nextMessageID++

	return nil
}

// GetMessageBetweenLandlordsAndTenant gets messages between a landlord and tenant
func (m *MockDB) GetMessageBetweenLandlordsAndTenant(tenantID string) ([]Message, error) {
	if err := m.checkFailNextOperation(); err != nil {
		return nil, err
	}

	m.mu.RLock()
	defer m.mu.RUnlock()

	var messages []Message
	for _, msg := range m.messages {
		if (msg.SenderType == "landlord" && msg.ReceiverType == "tenant") ||
			(msg.SenderType == "tenant" && msg.ReceiverType == "landlord") {
			messages = append(messages, msg)
		}
	}

	return messages, nil
}
