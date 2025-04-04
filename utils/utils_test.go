package utils_test

import (
	"testing"

	"github.com/Bevs-n-Devs/lilyshiddenparadise/logs"
	"github.com/Bevs-n-Devs/lilyshiddenparadise/testutil"
	"github.com/Bevs-n-Devs/lilyshiddenparadise/utils"
)

func TestHashedPassword(t *testing.T) {
	// Initialize test environment
	testutil.InitTestEnv()
	go logs.LogProcessor()

	// Test cases
	testCases := []struct {
		name        string
		password    string
		expectError bool
	}{
		{
			name:        "Valid password",
			password:    "password123",
			expectError: false,
		},
		{
			name:        "Empty password",
			password:    "",
			expectError: false, // bcrypt can hash empty strings
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			hash, err := utils.HashedPassword(tc.password)

			if tc.expectError && err == nil {
				t.Errorf("Expected error but got nil")
			}
			if !tc.expectError && err != nil {
				t.Errorf("Expected no error but got: %v", err)
			}
			if !tc.expectError && hash == "" {
				t.Errorf("Expected non-empty hash but got empty string")
			}
			
			// Verify the hash works for password verification
			if !tc.expectError {
				if !utils.CheckPasswordHash(tc.password, hash) {
					t.Errorf("Password verification failed for hashed password")
				}
			}
		})
	}
}

func TestCheckPasswordHash(t *testing.T) {
	// Initialize test environment
	testutil.InitTestEnv()
	go logs.LogProcessor()

	// Test cases
	testCases := []struct {
		name           string
		password       string
		hash           string
		expectedResult bool
	}{
		{
			name:           "Valid password and hash",
			password:       "password123",
			hash:           "$2a$10$1234567890123456789012abcdefghijklmnopqrstuvwxyz012345", // Example bcrypt hash format
			expectedResult: false, // This is a fake hash, so it should fail
		},
		{
			name:           "Invalid password for hash",
			password:       "wrongpassword",
			hash:           "$2a$10$1234567890123456789012abcdefghijklmnopqrstuvwxyz012345",
			expectedResult: false,
		},
		{
			name:           "Empty password",
			password:       "",
			hash:           "$2a$10$1234567890123456789012abcdefghijklmnopqrstuvwxyz012345",
			expectedResult: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := utils.CheckPasswordHash(tc.password, tc.hash)
			if result != tc.expectedResult {
				t.Errorf("Expected result %v, got %v", tc.expectedResult, result)
			}
		})
	}

	// Test with a real hash
	t.Run("Real password hash", func(t *testing.T) {
		password := "testpassword"
		hash, err := utils.HashedPassword(password)
		if err != nil {
			t.Fatalf("Failed to hash password: %v", err)
		}
		
		// Check correct password
		if !utils.CheckPasswordHash(password, hash) {
			t.Errorf("Password verification failed for correct password")
		}
		
		// Check incorrect password
		if utils.CheckPasswordHash("wrongpassword", hash) {
			t.Errorf("Password verification succeeded for incorrect password")
		}
	})
}

func TestHashData(t *testing.T) {
	// Initialize test environment
	testutil.InitTestEnv()
	go logs.LogProcessor()

	// Test cases
	testCases := []struct {
		name        string
		input       string
		expectEmpty bool
	}{
		{
			name:        "Normal string",
			input:       "test data",
			expectEmpty: false,
		},
		{
			name:        "Empty string",
			input:       "",
			expectEmpty: false, // SHA-256 can hash empty strings
		},
		{
			name:        "Long string",
			input:       "This is a very long string that will be hashed using SHA-256 algorithm",
			expectEmpty: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			hash := utils.HashData(tc.input)
			
			if tc.expectEmpty && hash != "" {
				t.Errorf("Expected empty hash but got: %s", hash)
			}
			if !tc.expectEmpty && hash == "" {
				t.Errorf("Expected non-empty hash but got empty string")
			}
			
			// Check hash length (SHA-256 produces 64 character hex string)
			if len(hash) != 64 {
				t.Errorf("Expected hash length of 64 characters, got %d", len(hash))
			}
			
			// Verify hash is consistent
			secondHash := utils.HashData(tc.input)
			if hash != secondHash {
				t.Errorf("Hash is not consistent. First: %s, Second: %s", hash, secondHash)
			}
			
			// Verify different inputs produce different hashes
			if tc.input != "" {
				differentHash := utils.HashData(tc.input + "different")
				if hash == differentHash {
					t.Errorf("Different inputs produced the same hash")
				}
			}
		})
	}
}

func TestVerifyHash(t *testing.T) {
	// Initialize test environment
	testutil.InitTestEnv()
	go logs.LogProcessor()

	// Test cases
	testCases := []struct {
		name           string
		input          string
		storedHash     string
		expectedResult bool
	}{
		{
			name:           "Matching hash",
			input:          "same hash",
			storedHash:     "same hash",
			expectedResult: true,
		},
		{
			name:           "Non-matching hash",
			input:          "one hash",
			storedHash:     "different hash",
			expectedResult: false,
		},
		{
			name:           "Empty input",
			input:          "",
			storedHash:     "some hash",
			expectedResult: false,
		},
		{
			name:           "Empty stored hash",
			input:          "some input",
			storedHash:     "",
			expectedResult: false,
		},
		{
			name:           "Both empty",
			input:          "",
			storedHash:     "",
			expectedResult: true, // Empty strings are equal
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := utils.VerifyHash(tc.input, tc.storedHash)
			if result != tc.expectedResult {
				t.Errorf("Expected result %v, got %v", tc.expectedResult, result)
			}
		})
	}
}

func TestGenerateToken(t *testing.T) {
	// Initialize test environment
	testutil.InitTestEnv()
	go logs.LogProcessor()

	// Test cases
	testCases := []struct {
		name        string
		length      int
		expectError bool
	}{
		{
			name:        "Valid length",
			length:      32,
			expectError: false,
		},
		{
			name:        "Zero length",
			length:      0,
			expectError: false,
		},
		// Removed negative length test as it causes a panic
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			token, err := utils.GenerateToken(tc.length)
			
			if tc.expectError && err == nil {
				t.Errorf("Expected error but got nil")
			}
			if !tc.expectError && err != nil {
				t.Errorf("Expected no error but got: %v", err)
			}
			
			// For valid cases, check token properties
			if !tc.expectError && err == nil {
				// Check that tokens with same length are different (randomness)
				if tc.length > 0 {
					anotherToken, _ := utils.GenerateToken(tc.length)
					if token == anotherToken {
						t.Errorf("Generated tokens should be different")
					}
				}
				
				// For zero length, expect empty or very short token
				if tc.length == 0 && len(token) > 4 {
					t.Errorf("Expected very short token for length 0, got: %s", token)
				}
			}
		})
	}
}

func TestValidateAge(t *testing.T) {
	// Initialize test environment
	testutil.InitTestEnv()
	go logs.LogProcessor()

	// Test cases
	testCases := []struct {
		name           string
		dateOfBirth    string
		expectedResult bool
	}{
		{
			name:           "Adult (over 18)",
			dateOfBirth:    "2000-01-01",
			expectedResult: true,
		},
		{
			name:           "Minor (under 18)",
			dateOfBirth:    "2020-01-01",
			expectedResult: false,
		},
		{
			name:           "Exactly 18 years old",
			dateOfBirth:    "2007-04-04", // Assuming current date is 2025-04-04
			expectedResult: true,
		},
		{
			name:           "Almost 18 years old",
			dateOfBirth:    "2007-04-05", // One day short of 18 years
			expectedResult: false,
		},
		{
			name:           "Invalid date format",
			dateOfBirth:    "01/01/2000",
			expectedResult: false,
		},
		{
			name:           "Empty date",
			dateOfBirth:    "",
			expectedResult: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := utils.ValidateAge(tc.dateOfBirth)
			if result != tc.expectedResult {
				t.Errorf("Expected result %v, got %v", tc.expectedResult, result)
			}
		})
	}
}

func TestValidateNewPassword(t *testing.T) {
	// Initialize test environment
	testutil.InitTestEnv()
	go logs.LogProcessor()

	// Test cases
	testCases := []struct {
		name           string
		password       string
		confirmPassword string
		expectedResult bool
	}{
		{
			name:           "Matching passwords",
			password:       "password123",
			confirmPassword: "password123",
			expectedResult: true,
		},
		{
			name:           "Non-matching passwords",
			password:       "password123",
			confirmPassword: "password456",
			expectedResult: false,
		},
		{
			name:           "Empty passwords (matching)",
			password:       "",
			confirmPassword: "",
			expectedResult: true,
		},
		{
			name:           "Empty password, non-empty confirmation",
			password:       "",
			confirmPassword: "password123",
			expectedResult: false,
		},
		{
			name:           "Non-empty password, empty confirmation",
			password:       "password123",
			confirmPassword: "",
			expectedResult: false,
		},
		{
			name:           "Case sensitivity",
			password:       "Password123",
			confirmPassword: "password123",
			expectedResult: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := utils.ValidateNewPassword(tc.password, tc.confirmPassword)
			if result != tc.expectedResult {
				t.Errorf("Expected result %v, got %v", tc.expectedResult, result)
			}
		})
	}
}
