package handlers_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/Bevs-n-Devs/lilyshiddenparadise/handlers"
	"github.com/Bevs-n-Devs/lilyshiddenparadise/logs"
	"github.com/Bevs-n-Devs/lilyshiddenparadise/testutil"
)

func TestLogoutLandlord(t *testing.T) {
	// Use the global test environment
	testutil.TestMain(nil)
	
	go logs.LogProcessor()

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
			name: "Valid logout",
			setupRequest: func(req *http.Request) {
				// Set up session cookie
				cookie := &http.Cookie{
					Name:  "session_token",
					Value: "session_token_test@example.com",
				}
				req.AddCookie(cookie)
			},
			expectedStatusCode: http.StatusSeeOther,
			expectedLocation:   "/",
		},
		{
			name: "Missing session cookie",
			setupRequest: func(req *http.Request) {
				// No session cookie
			},
			expectedStatusCode: http.StatusSeeOther,
			expectedLocation:   "/login/landlord?authenticationError=Please login to access this page",
		},
		{
			name: "Database error during logout",
			setupRequest: func(req *http.Request) {
				// Set up session cookie
				cookie := &http.Cookie{
					Name:  "session_token",
					Value: "session_token_test@example.com",
				}
				req.AddCookie(cookie)
			},
			expectedStatusCode: http.StatusSeeOther,
			expectedLocation:   "/login/landlord?authenticationError=Internal server error",
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
			req := httptest.NewRequest(http.MethodGet, "/logout/landlord", nil)
			
			// Set up the request
			tc.setupRequest(req)
			
			// Create a response recorder
			rr := httptest.NewRecorder()
			
			// Call the handler
			handlers.LogoutLandlord(rr, req)
			
			// Check the status code
			if rr.Code != tc.expectedStatusCode {
				t.Errorf("Expected status code %d, got %d", tc.expectedStatusCode, rr.Code)
			}
			
			// Check the redirect location
			location := rr.Header().Get("Location")
			if location != tc.expectedLocation {
				t.Errorf("Expected redirect to '%s', got '%s'", tc.expectedLocation, location)
			}
			
			// Check that the session cookie is deleted if logout was successful
			if tc.expectedLocation == "/" {
				cookies := rr.Result().Cookies()
				found := false
				for _, cookie := range cookies {
					if cookie.Name == "session_token" && cookie.MaxAge < 0 {
						found = true
						break
					}
				}
				if !found {
					t.Errorf("Expected session_token cookie to be deleted")
				}
			}
		})
	}
}

func TestLogoutTenant(t *testing.T) {
	// Use the global test environment
	testutil.TestMain(nil)
	
	go logs.LogProcessor()

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
			name: "Valid logout",
			setupRequest: func(req *http.Request) {
				// Set up session cookie
				cookie := &http.Cookie{
					Name:  "session_token",
					Value: "session_token_hash_tenant@example.com",
				}
				req.AddCookie(cookie)
			},
			expectedStatusCode: http.StatusSeeOther,
			expectedLocation:   "/",
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
			name: "Database error during logout",
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
			req := httptest.NewRequest(http.MethodGet, "/logout/tenant", nil)
			
			// Set up the request
			tc.setupRequest(req)
			
			// Create a response recorder
			rr := httptest.NewRecorder()
			
			// Call the handler
			handlers.LogoutTenant(rr, req)
			
			// Check the status code
			if rr.Code != tc.expectedStatusCode {
				t.Errorf("Expected status code %d, got %d", tc.expectedStatusCode, rr.Code)
			}
			
			// Check the redirect location
			location := rr.Header().Get("Location")
			if location != tc.expectedLocation {
				t.Errorf("Expected redirect to '%s', got '%s'", tc.expectedLocation, location)
			}
			
			// Check that the session cookie is deleted if logout was successful
			if tc.expectedLocation == "/" {
				cookies := rr.Result().Cookies()
				found := false
				for _, cookie := range cookies {
					if cookie.Name == "session_token" && cookie.MaxAge < 0 {
						found = true
						break
					}
				}
				if !found {
					t.Errorf("Expected session_token cookie to be deleted")
				}
			}
		})
	}
}
