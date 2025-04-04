package handlers_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/Bevs-n-Devs/lilyshiddenparadise/handlers"
	"github.com/Bevs-n-Devs/lilyshiddenparadise/testutil"
)

func TestTenantDashboard(t *testing.T) {
	// Use the global test environment
	testutil.TestMain(nil)

	// Define test cases
	testCases := []struct {
		name               string
		setupRequest       func(*http.Request)
		expectedStatusCode int
		expectedLocation   string
		failDB             bool
		failDBReason       string
	}{
		{
			name: "Valid session",
			setupRequest: func(req *http.Request) {
				// Set up session cookie
				cookie := &http.Cookie{
					Name:  "session_token",
					Value: "session_token_hash_tenant@example.com",
				}
				req.AddCookie(cookie)
			},
			expectedStatusCode: http.StatusOK,
		},
		{
			name: "Missing session cookie",
			setupRequest: func(req *http.Request) {
				// No session cookie
			},
			expectedStatusCode: http.StatusSeeOther,
			expectedLocation:   "/login/tenant?authenticationError=Please login to access this page",
		},
		{
			name: "Invalid session token",
			setupRequest: func(req *http.Request) {
				// Set up invalid session cookie
				cookie := &http.Cookie{
					Name:  "session_token",
					Value: "invalid_session_token",
				}
				req.AddCookie(cookie)
			},
			expectedStatusCode: http.StatusSeeOther,
			expectedLocation:   "/login/tenant?authenticationError=Please login to access this page",
		},
		{
			name: "Database error when validating session",
			setupRequest: func(req *http.Request) {
				// Set up session cookie
				cookie := &http.Cookie{
					Name:  "session_token",
					Value: "session_token_hash_tenant@example.com",
				}
				req.AddCookie(cookie)
			},
			expectedStatusCode: http.StatusSeeOther,
			expectedLocation:   "/login/tenant?authenticationError=Internal server error",
			failDB:             true,
			failDBReason:       "database error",
		},
	}

	// Run test cases
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Set up mock database to fail if needed
			if tc.failDB {
				testutil.TestEnvironment.DB.SetFailNextOperation(tc.failDBReason)
			}
			
			// Create a request
			req := httptest.NewRequest(http.MethodGet, "/tenant/dashboard", nil)
			
			// Set up the request
			tc.setupRequest(req)
			
			// Create a response recorder
			rr := httptest.NewRecorder()
			
			// Call the handler
			handlers.TenantDashboard(rr, req)
			
			// Check the status code
			if rr.Code != tc.expectedStatusCode {
				t.Errorf("Expected status code %d, got %d", tc.expectedStatusCode, rr.Code)
			}
			
			// Check the redirect location if applicable
			if tc.expectedLocation != "" {
				location := rr.Header().Get("Location")
				if location != tc.expectedLocation {
					t.Errorf("Expected redirect to '%s', got '%s'", tc.expectedLocation, location)
				}
			}
		})
	}
}
