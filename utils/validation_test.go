package utils_test

import (
	"strings"
	"testing"

	"github.com/Bevs-n-Devs/lilyshiddenparadise/logs"
	"github.com/Bevs-n-Devs/lilyshiddenparadise/testutil"
	"github.com/Bevs-n-Devs/lilyshiddenparadise/utils"
)

func TestCheckIfEvicted(t *testing.T) {
	// Initialize test environment
	testutil.InitTestEnv()
	go logs.LogProcessor()

	// Test cases
	testCases := []struct {
		name           string
		ifEvicted      string
		evictedReason  string
		expectedResult bool
	}{
		{
			name:           "Not evicted",
			ifEvicted:      "no",
			evictedReason:  "",
			expectedResult: true,
		},
		{
			name:           "Evicted with reason",
			ifEvicted:      "yes",
			evictedReason:  "Late payment",
			expectedResult: true,
		},
		{
			name:           "Evicted without reason",
			ifEvicted:      "yes",
			evictedReason:  "",
			expectedResult: false,
		},
		{
			name:           "Empty eviction status with reason",
			ifEvicted:      "",
			evictedReason:  "Some reason",
			expectedResult: true,
		},
		{
			name:           "Empty eviction status without reason",
			ifEvicted:      "",
			evictedReason:  "",
			expectedResult: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := utils.CheckIfEvicted(tc.ifEvicted, tc.evictedReason)
			if result != tc.expectedResult {
				t.Errorf("Expected result %v, got %v", tc.expectedResult, result)
			}
		})
	}
}

func TestCheckIfConvicted(t *testing.T) {
	// Initialize test environment
	testutil.InitTestEnv()
	go logs.LogProcessor()

	// Test cases
	testCases := []struct {
		name           string
		ifConvicted    string
		convictedReason string
		expectedResult bool
	}{
		{
			name:           "Not convicted",
			ifConvicted:    "no",
			convictedReason: "",
			expectedResult: true,
		},
		{
			name:           "Convicted with reason",
			ifConvicted:    "yes",
			convictedReason: "Theft",
			expectedResult: true,
		},
		{
			name:           "Convicted without reason",
			ifConvicted:    "yes",
			convictedReason: "",
			expectedResult: false,
		},
		{
			name:           "Empty conviction status with reason",
			ifConvicted:    "",
			convictedReason: "Some reason",
			expectedResult: true,
		},
		{
			name:           "Empty conviction status without reason",
			ifConvicted:    "",
			convictedReason: "",
			expectedResult: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := utils.CheckIfConvicted(tc.ifConvicted, tc.convictedReason)
			if result != tc.expectedResult {
				t.Errorf("Expected result %v, got %v", tc.expectedResult, result)
			}
		})
	}
}

func TestCheckIfVehicle(t *testing.T) {
	// Initialize test environment
	testutil.InitTestEnv()
	go logs.LogProcessor()

	// Test cases
	testCases := []struct {
		name           string
		ifVehicle      string
		vehicleReg     string
		expectedResult bool
	}{
		{
			name:           "No vehicle",
			ifVehicle:      "no",
			vehicleReg:     "",
			expectedResult: true,
		},
		{
			name:           "Has vehicle with registration",
			ifVehicle:      "yes",
			vehicleReg:     "ABC123",
			expectedResult: true,
		},
		{
			name:           "Has vehicle without registration",
			ifVehicle:      "yes",
			vehicleReg:     "",
			expectedResult: false,
		},
		{
			name:           "Empty vehicle status with registration",
			ifVehicle:      "",
			vehicleReg:     "ABC123",
			expectedResult: true,
		},
		{
			name:           "Empty vehicle status without registration",
			ifVehicle:      "",
			vehicleReg:     "",
			expectedResult: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := utils.CheckIfVehicle(tc.ifVehicle, tc.vehicleReg)
			if result != tc.expectedResult {
				t.Errorf("Expected result %v, got %v", tc.expectedResult, result)
			}
		})
	}
}

func TestCheckIfHaveChildren(t *testing.T) {
	// Initialize test environment
	testutil.InitTestEnv()
	go logs.LogProcessor()

	// Test cases
	testCases := []struct {
		name           string
		haveChildren   string
		children       string
		expectedResult bool
	}{
		{
			name:           "No children",
			haveChildren:   "no",
			children:       "",
			expectedResult: true,
		},
		{
			name:           "Has children with details",
			haveChildren:   "yes",
			children:       "2 children",
			expectedResult: true,
		},
		{
			name:           "Has children without details",
			haveChildren:   "yes",
			children:       "",
			expectedResult: false,
		},
		{
			name:           "Empty children status with details",
			haveChildren:   "",
			children:       "2 children",
			expectedResult: true,
		},
		{
			name:           "Empty children status without details",
			haveChildren:   "",
			children:       "",
			expectedResult: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := utils.CheckIfHaveChildren(tc.haveChildren, tc.children)
			if result != tc.expectedResult {
				t.Errorf("Expected result %v, got %v", tc.expectedResult, result)
			}
		})
	}
}

func TestCheckIfRefusedRent(t *testing.T) {
	// Initialize test environment
	testutil.InitTestEnv()
	go logs.LogProcessor()

	// Test cases
	testCases := []struct {
		name             string
		refusedRent      string
		refusedRentReason string
		expectedResult   bool
	}{
		{
			name:             "Not refused rent",
			refusedRent:      "no",
			refusedRentReason: "",
			expectedResult:   true,
		},
		{
			name:             "Refused rent with reason",
			refusedRent:      "yes",
			refusedRentReason: "Financial issues",
			expectedResult:   true,
		},
		{
			name:             "Refused rent without reason",
			refusedRent:      "yes",
			refusedRentReason: "",
			expectedResult:   false,
		},
		{
			name:             "Empty refused rent status with reason",
			refusedRent:      "",
			refusedRentReason: "Some reason",
			expectedResult:   true,
		},
		{
			name:             "Empty refused rent status without reason",
			refusedRent:      "",
			refusedRentReason: "",
			expectedResult:   true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := utils.CheckIfRefusedRent(tc.refusedRent, tc.refusedRentReason)
			if result != tc.expectedResult {
				t.Errorf("Expected result %v, got %v", tc.expectedResult, result)
			}
		})
	}
}

func TestCheckIfStableIncome(t *testing.T) {
	// Initialize test environment
	testutil.InitTestEnv()
	go logs.LogProcessor()

	// Test cases
	testCases := []struct {
		name           string
		unstableIncome string
		incomeReason   string
		expectedResult bool
	}{
		{
			name:           "Stable income",
			unstableIncome: "no",
			incomeReason:   "",
			expectedResult: true,
		},
		{
			name:           "Unstable income with reason",
			unstableIncome: "yes",
			incomeReason:   "Freelance work",
			expectedResult: true,
		},
		{
			name:           "Unstable income without reason",
			unstableIncome: "yes",
			incomeReason:   "",
			expectedResult: false,
		},
		{
			name:           "Empty income status with reason",
			unstableIncome: "",
			incomeReason:   "Some reason",
			expectedResult: true,
		},
		{
			name:           "Empty income status without reason",
			unstableIncome: "",
			incomeReason:   "",
			expectedResult: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := utils.CheckIfStableIncome(tc.unstableIncome, tc.incomeReason)
			if result != tc.expectedResult {
				t.Errorf("Expected result %v, got %v", tc.expectedResult, result)
			}
		})
	}
}

func TestValidateManageTenantApplication(t *testing.T) {
	// Initialize test environment
	testutil.InitTestEnv()
	go logs.LogProcessor()

	// Test cases
	testCases := []struct {
		name             string
		applicationResult string
		roomType         string
		moveInDate       string
		rentDue          string
		monthlyRent      string
		currency         string
		expectError      bool
	}{
		{
			name:             "Accepted application with all fields",
			applicationResult: "accepted",
			roomType:         "Single",
			moveInDate:       "2025-05-01",
			rentDue:          "1",
			monthlyRent:      "1000",
			currency:         "USD",
			expectError:      false,
		},
		{
			name:             "Rejected application",
			applicationResult: "rejected",
			roomType:         "",
			moveInDate:       "",
			rentDue:          "",
			monthlyRent:      "",
			currency:         "",
			expectError:      false,
		},
		{
			name:             "Accepted application missing room type",
			applicationResult: "accepted",
			roomType:         "",
			moveInDate:       "2025-05-01",
			rentDue:          "1",
			monthlyRent:      "1000",
			currency:         "USD",
			expectError:      true,
		},
		{
			name:             "Accepted application missing move in date",
			applicationResult: "accepted",
			roomType:         "Single",
			moveInDate:       "",
			rentDue:          "1",
			monthlyRent:      "1000",
			currency:         "USD",
			expectError:      true,
		},
		{
			name:             "Accepted application missing rent due",
			applicationResult: "accepted",
			roomType:         "Single",
			moveInDate:       "2025-05-01",
			rentDue:          "",
			monthlyRent:      "1000",
			currency:         "USD",
			expectError:      true,
		},
		{
			name:             "Accepted application missing monthly rent",
			applicationResult: "accepted",
			roomType:         "Single",
			moveInDate:       "2025-05-01",
			rentDue:          "1",
			monthlyRent:      "",
			currency:         "USD",
			expectError:      true,
		},
		{
			name:             "Accepted application missing currency",
			applicationResult: "accepted",
			roomType:         "Single",
			moveInDate:       "2025-05-01",
			rentDue:          "1",
			monthlyRent:      "1000",
			currency:         "",
			expectError:      true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			err := utils.ValidateManageTenantApplication(
				tc.applicationResult,
				tc.roomType,
				tc.moveInDate,
				tc.rentDue,
				tc.monthlyRent,
				tc.currency,
			)
			
			if tc.expectError && err == nil {
				t.Errorf("Expected error but got nil")
			}
			if !tc.expectError && err != nil {
				t.Errorf("Expected no error but got: %v", err)
			}
		})
	}
}

func TestGenerateTenantUsernamePassportNumberAndPassword(t *testing.T) {
	// Initialize test environment
	testutil.InitTestEnv()
	go logs.LogProcessor()

	// Test cases
	testCases := []struct {
		name           string
		email          string
		passportNumber string
		expectError    bool
	}{
		{
			name:           "Valid inputs",
			email:          "test@example.com",
			passportNumber: "AB123456",
			expectError:    false,
		},
		{
			name:           "Empty email",
			email:          "",
			passportNumber: "AB123456",
			expectError:    true,
		},
		{
			name:           "Empty passport number",
			email:          "test@example.com",
			passportNumber: "",
			expectError:    true,
		},
		{
			name:           "Both empty",
			email:          "",
			passportNumber: "",
			expectError:    true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			username, password, err := utils.GenerateTenantUsernamePassportNumberAndPassword(tc.email, tc.passportNumber)
			
			if tc.expectError && err == nil {
				t.Errorf("Expected error but got nil")
			}
			if !tc.expectError && err != nil {
				t.Errorf("Expected no error but got: %v", err)
			}
			
			// For valid cases, check username and password
			if !tc.expectError && err == nil {
				// Username should be the same as email
				if username != tc.email {
					t.Errorf("Expected username to be '%s', got '%s'", tc.email, username)
				}
				
				// Password should contain the passport number
				if !strings.Contains(password, tc.passportNumber) {
					t.Errorf("Expected password to contain passport number '%s', got '%s'", tc.passportNumber, password)
				}
				
				// Password should be longer than passport number (should have hash prefix)
				if len(password) <= len(tc.passportNumber) {
					t.Errorf("Expected password to be longer than passport number")
				}
			}
		})
	}
}
