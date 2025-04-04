package handlers_test

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

	"github.com/Bevs-n-Devs/lilyshiddenparadise/handlers"
	"github.com/Bevs-n-Devs/lilyshiddenparadise/testutil"
)

func TestTenantMessages(t *testing.T) {
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
			name: "Database error when getting messages",
			setupRequest: func(req *http.Request) {
				// Set up session cookie
				cookie := &http.Cookie{
					Name:  "session_token",
					Value: "session_token_hash_tenant@example.com",
				}
				req.AddCookie(cookie)
			},
			expectedStatusCode: http.StatusSeeOther,
			expectedLocation:   "/tenant/dashboard?databaseError=Failed to get messages",
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
			req := httptest.NewRequest(http.MethodGet, "/tenant/messages", nil)
			
			// Set up the request
			tc.setupRequest(req)
			
			// Create a response recorder
			rr := httptest.NewRecorder()
			
			// Call the handler
			handlers.TenantMessages(rr, req)
			
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

func TestSendMessageToLandlord(t *testing.T) {
	// Use the global test environment
	testutil.TestMain(nil)
	
	// Define test cases
	testCases := []struct {
		name               string
		setupRequest       func(*http.Request)
		formData           map[string]string
		expectedStatusCode int
		expectedLocation   string
		failDB             bool
		failDBReason       string
	}{
		{
			name: "Valid message send",
			setupRequest: func(req *http.Request) {
				// Set up session cookie
				cookie := &http.Cookie{
					Name:  "session_token",
					Value: "session_token_hash_tenant@example.com",
				}
				req.AddCookie(cookie)
				
				// Set up CSRF token
				cookie = &http.Cookie{
					Name:  "csrf_token",
					Value: "csrf_token_hash_tenant@example.com",
				}
				req.AddCookie(cookie)
			},
			formData: map[string]string{
				"message":    "Hello landlord, this is a test message.",
				"csrf_token": "csrf_token_hash_tenant@example.com",
			},
			expectedStatusCode: http.StatusSeeOther,
			expectedLocation:   "/tenant/messages",
		},
		{
			name: "Missing session cookie",
			setupRequest: func(req *http.Request) {
				// No session cookie
				
				// Set up CSRF token
				cookie := &http.Cookie{
					Name:  "csrf_token",
					Value: "csrf_token_hash_tenant@example.com",
				}
				req.AddCookie(cookie)
			},
			formData: map[string]string{
				"message":    "Hello landlord, this is a test message.",
				"csrf_token": "csrf_token_hash_tenant@example.com",
			},
			expectedStatusCode: http.StatusSeeOther,
			expectedLocation:   "/login/tenant?authenticationError=Please login to access this page",
		},
		{
			name: "Invalid CSRF token",
			setupRequest: func(req *http.Request) {
				// Set up session cookie
				cookie := &http.Cookie{
					Name:  "session_token",
					Value: "session_token_hash_tenant@example.com",
				}
				req.AddCookie(cookie)
				
				// Set up CSRF token
				cookie = &http.Cookie{
					Name:  "csrf_token",
					Value: "csrf_token_hash_tenant@example.com",
				}
				req.AddCookie(cookie)
			},
			formData: map[string]string{
				"message":    "Hello landlord, this is a test message.",
				"csrf_token": "invalid_csrf_token",
			},
			expectedStatusCode: http.StatusSeeOther,
			expectedLocation:   "/tenant/messages?validationError=Invalid CSRF token",
		},
		{
			name: "Empty message",
			setupRequest: func(req *http.Request) {
				// Set up session cookie
				cookie := &http.Cookie{
					Name:  "session_token",
					Value: "session_token_hash_tenant@example.com",
				}
				req.AddCookie(cookie)
				
				// Set up CSRF token
				cookie = &http.Cookie{
					Name:  "csrf_token",
					Value: "csrf_token_hash_tenant@example.com",
				}
				req.AddCookie(cookie)
			},
			formData: map[string]string{
				"message":    "",
				"csrf_token": "csrf_token_hash_tenant@example.com",
			},
			expectedStatusCode: http.StatusSeeOther,
			expectedLocation:   "/tenant/messages?validationError=Message cannot be empty",
		},
		{
			name: "Database error when sending message",
			setupRequest: func(req *http.Request) {
				// Set up session cookie
				cookie := &http.Cookie{
					Name:  "session_token",
					Value: "session_token_hash_tenant@example.com",
				}
				req.AddCookie(cookie)
				
				// Set up CSRF token
				cookie = &http.Cookie{
					Name:  "csrf_token",
					Value: "csrf_token_hash_tenant@example.com",
				}
				req.AddCookie(cookie)
			},
			formData: map[string]string{
				"message":    "Hello landlord, this is a test message.",
				"csrf_token": "csrf_token_hash_tenant@example.com",
			},
			expectedStatusCode: http.StatusSeeOther,
			expectedLocation:   "/tenant/messages?databaseError=Failed to send message",
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
			
			// Create form data
			form := url.Values{}
			for key, value := range tc.formData {
				form.Add(key, value)
			}
			
			// Create a request
			req := httptest.NewRequest(http.MethodPost, "/tenant/send-message-to-landlord", strings.NewReader(form.Encode()))
			req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
			
			// Set up the request
			tc.setupRequest(req)
			
			// Create a response recorder
			rr := httptest.NewRecorder()
			
			// Call the handler
			handlers.SendMessageToLandlord(rr, req)
			
			// Check the status code
			if rr.Code != tc.expectedStatusCode {
				t.Errorf("Expected status code %d, got %d", tc.expectedStatusCode, rr.Code)
			}
			
			// Check the redirect location
			location := rr.Header().Get("Location")
			if location != tc.expectedLocation {
				t.Errorf("Expected redirect to '%s', got '%s'", tc.expectedLocation, location)
			}
		})
	}
}
