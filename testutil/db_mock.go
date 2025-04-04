package testutil

import (
	"time"
)

// Global variables to store mock implementations
var (
	MockAuthenticateLandlord = func(email, password string) (bool, error) {
		if email == "test@example.com" && password == "password123" {
			return true, nil
		}
		return false, nil
	}
	
	MockUpdateLandlordSessionTokens = func(email string) (string, string, time.Time, error) {
		return "test_session_token", "test_csrf_token", time.Now().Add(30 * time.Minute), nil
	}
)
