package handlers_test

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

	"github.com/Bevs-n-Devs/lilyshiddenparadise/handlers"
	"github.com/Bevs-n-Devs/lilyshiddenparadise/logs"
	"github.com/Bevs-n-Devs/lilyshiddenparadise/testutil"
)

func TestLandlordMessages(t *testing.T) {
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
			name: "Valid session",
			setupRequest: func(req *http.Request) {
				// Set up session cookie
				cookie := &http.Cookie{
					Name:  "session_token",
					Value: "session_token_test@example.com",
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
			expectedLocation:   "/login/landlord?authenticationError=Please login to access this page",
		},
		{
			name: "Database error when getting tenants",
			setupRequest: func(req *http.Request) {
				// Set up session cookie
				cookie := &http.Cookie{
					Name:  "session_token",
					Value: "session_token_test@example.com",
				}
				req.AddCookie(cookie)
			},
			expectedStatusCode: http.StatusSeeOther,
			expectedLocation:   "/landlord/dashboard?databaseError=Failed to get tenants",
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
			req := httptest.NewRequest(http.MethodGet, "/landlord/messages", nil)
			
			// Set up the request
			tc.setupRequest(req)
			
			// Create a response recorder
			rr := httptest.NewRecorder()
			
			// Call the handler
			handlers.LandlordMessages(rr, req)
			
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

func TestSendMessageToTenant(t *testing.T) {
	// Use the global test environment
	testutil.TestMain(nil)
	
	go logs.LogProcessor()
	
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
					Value: "session_token_test@example.com",
				}
				req.AddCookie(cookie)
				
				// Set up CSRF token
				cookie = &http.Cookie{
					Name:  "csrf_token",
					Value: "csrf_token_test@example.com",
				}
				req.AddCookie(cookie)
			},
			formData: map[string]string{
				"tenant_id":  "1",
				"message":    "Hello tenant, this is a test message.",
				"csrf_token": "csrf_token_test@example.com",
			},
			expectedStatusCode: http.StatusSeeOther,
			expectedLocation:   "/landlord/tenant/1/messages",
		},
		{
			name: "Missing session cookie",
			setupRequest: func(req *http.Request) {
				// No session cookie
				
				// Set up CSRF token
				cookie := &http.Cookie{
					Name:  "csrf_token",
					Value: "csrf_token_test@example.com",
				}
				req.AddCookie(cookie)
			},
			formData: map[string]string{
				"tenant_id":  "1",
				"message":    "Hello tenant, this is a test message.",
				"csrf_token": "csrf_token_test@example.com",
			},
			expectedStatusCode: http.StatusSeeOther,
			expectedLocation:   "/login/landlord?authenticationError=Please login to access this page",
		},
		{
			name: "Invalid CSRF token",
			setupRequest: func(req *http.Request) {
				// Set up session cookie
				cookie := &http.Cookie{
					Name:  "session_token",
					Value: "session_token_test@example.com",
				}
				req.AddCookie(cookie)
				
				// Set up CSRF token
				cookie = &http.Cookie{
					Name:  "csrf_token",
					Value: "csrf_token_test@example.com",
				}
				req.AddCookie(cookie)
			},
			formData: map[string]string{
				"tenant_id":  "1",
				"message":    "Hello tenant, this is a test message.",
				"csrf_token": "invalid_csrf_token",
			},
			expectedStatusCode: http.StatusSeeOther,
			expectedLocation:   "/landlord/messages?validationError=Invalid CSRF token",
		},
		{
			name: "Empty message",
			setupRequest: func(req *http.Request) {
				// Set up session cookie
				cookie := &http.Cookie{
					Name:  "session_token",
					Value: "session_token_test@example.com",
				}
				req.AddCookie(cookie)
				
				// Set up CSRF token
				cookie = &http.Cookie{
					Name:  "csrf_token",
					Value: "csrf_token_test@example.com",
				}
				req.AddCookie(cookie)
			},
			formData: map[string]string{
				"tenant_id":  "1",
				"message":    "",
				"csrf_token": "csrf_token_test@example.com",
			},
			expectedStatusCode: http.StatusSeeOther,
			expectedLocation:   "/landlord/tenant/1/messages?validationError=Message cannot be empty",
		},
		{
			name: "Database error when sending message",
			setupRequest: func(req *http.Request) {
				// Set up session cookie
				cookie := &http.Cookie{
					Name:  "session_token",
					Value: "session_token_test@example.com",
				}
				req.AddCookie(cookie)
				
				// Set up CSRF token
				cookie = &http.Cookie{
					Name:  "csrf_token",
					Value: "csrf_token_test@example.com",
				}
				req.AddCookie(cookie)
			},
			formData: map[string]string{
				"tenant_id":  "1",
				"message":    "Hello tenant, this is a test message.",
				"csrf_token": "csrf_token_test@example.com",
			},
			expectedStatusCode: http.StatusSeeOther,
			expectedLocation:   "/landlord/tenant/1/messages?databaseError=Failed to send message",
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
			req := httptest.NewRequest(http.MethodPost, "/landlord/send-message-to-tenant", strings.NewReader(form.Encode()))
			req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
			
			// Set up the request
			tc.setupRequest(req)
			
			// Create a response recorder
			rr := httptest.NewRecorder()
			
			// Call the handler
			handlers.SendMessageToTenant(rr, req)
			
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
