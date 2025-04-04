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

func TestUpdateTenantPassword(t *testing.T) {
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
			name: "Valid password update",
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
				"current_password": "password123",
				"new_password":     "newpassword123",
				"confirm_password": "newpassword123",
				"csrf_token":       "csrf_token_hash_tenant@example.com",
			},
			expectedStatusCode: http.StatusSeeOther,
			expectedLocation:   "/tenant/account?success=Password updated successfully",
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
				"current_password": "password123",
				"new_password":     "newpassword123",
				"confirm_password": "newpassword123",
				"csrf_token":       "csrf_token_hash_tenant@example.com",
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
				"current_password": "password123",
				"new_password":     "newpassword123",
				"confirm_password": "newpassword123",
				"csrf_token":       "invalid_csrf_token",
			},
			expectedStatusCode: http.StatusSeeOther,
			expectedLocation:   "/tenant/account?validationError=Invalid CSRF token",
		},
		{
			name: "Passwords don't match",
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
				"current_password": "password123",
				"new_password":     "newpassword123",
				"confirm_password": "differentpassword",
				"csrf_token":       "csrf_token_hash_tenant@example.com",
			},
			expectedStatusCode: http.StatusSeeOther,
			expectedLocation:   "/tenant/account?confirmPasswordError=Passwords do not match",
		},
		{
			name: "Empty new password",
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
				"current_password": "password123",
				"new_password":     "",
				"confirm_password": "",
				"csrf_token":       "csrf_token_hash_tenant@example.com",
			},
			expectedStatusCode: http.StatusSeeOther,
			expectedLocation:   "/tenant/account?validationError=New password cannot be empty",
		},
		{
			name: "Database error when updating password",
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
				"current_password": "password123",
				"new_password":     "newpassword123",
				"confirm_password": "newpassword123",
				"csrf_token":       "csrf_token_hash_tenant@example.com",
			},
			expectedStatusCode: http.StatusSeeOther,
			expectedLocation:   "/tenant/account?databaseError=Failed to update password",
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
			req := httptest.NewRequest(http.MethodPost, "/tenant/update-password", strings.NewReader(form.Encode()))
			req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
			
			// Set up the request
			tc.setupRequest(req)
			
			// Create a response recorder
			rr := httptest.NewRecorder()
			
			// Call the handler
			handlers.UpdateTenantPassword(rr, req)
			
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
