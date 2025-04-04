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

func TestLandlordManageApplications(t *testing.T) {
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
			name: "Database error when getting applications",
			setupRequest: func(req *http.Request) {
				// Set up session cookie
				cookie := &http.Cookie{
					Name:  "session_token",
					Value: "session_token_test@example.com",
				}
				req.AddCookie(cookie)
			},
			expectedStatusCode: http.StatusSeeOther,
			expectedLocation:   "/landlord/dashboard?databaseError=Failed to get tenant applications",
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
			req := httptest.NewRequest(http.MethodGet, "/landlord/manage-applications", nil)

			// Set up the request
			tc.setupRequest(req)

			// Create a response recorder
			rr := httptest.NewRecorder()

			// Call the handler
			handlers.LandlordManageApplications(rr, req)

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

func TestUpdateTenantApplicationStatus(t *testing.T) {
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
			name: "Valid status update - Approve",
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
				"application_id": "1",
				"status":         "approved",
				"csrf_token":     "csrf_token_test@example.com",
			},
			expectedStatusCode: http.StatusSeeOther,
			expectedLocation:   "/landlord/manage-applications",
		},
		{
			name: "Valid status update - Reject",
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
				"application_id": "1",
				"status":         "rejected",
				"csrf_token":     "csrf_token_test@example.com",
			},
			expectedStatusCode: http.StatusSeeOther,
			expectedLocation:   "/landlord/manage-applications",
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
				"application_id": "1",
				"status":         "approved",
				"csrf_token":     "csrf_token_test@example.com",
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
				"application_id": "1",
				"status":         "approved",
				"csrf_token":     "invalid_csrf_token",
			},
			expectedStatusCode: http.StatusSeeOther,
			expectedLocation:   "/landlord/manage-applications?validationError=Invalid CSRF token",
		},
		{
			name: "Database error when updating status",
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
				"application_id": "1",
				"status":         "approved",
				"csrf_token":     "csrf_token_test@example.com",
			},
			expectedStatusCode: http.StatusSeeOther,
			expectedLocation:   "/landlord/manage-applications?databaseError=Failed to update tenant application status",
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
			req := httptest.NewRequest(http.MethodPost, "/landlord/dashboard/manage-applications", strings.NewReader(form.Encode()))
			req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

			// Set up the request
			tc.setupRequest(req)

			// Create a response recorder
			rr := httptest.NewRecorder()

			// Call the handler
			handlers.LandlordManageApplications(rr, req)

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
