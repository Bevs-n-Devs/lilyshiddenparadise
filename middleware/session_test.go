package middleware_test

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/Bevs-n-Devs/lilyshiddenparadise/logs"
	"github.com/Bevs-n-Devs/lilyshiddenparadise/middleware"
	"github.com/Bevs-n-Devs/lilyshiddenparadise/testutil"
)

func TestAuthenticateLandlordRequest(t *testing.T) {
	// Initialize test environment
	testutil.InitTestEnv()
	go logs.LogProcessor()

	// Test cases
	testCases := []struct {
		name          string
		setupRequest  func(*http.Request)
		expectedError bool
		errorContains string
	}{
		{
			name: "Missing session token",
			setupRequest: func(req *http.Request) {
				// Only add CSRF token
				req.AddCookie(&http.Cookie{
					Name:  "csrf_token",
					Value: "csrf_token_test@example.com",
				})
			},
			expectedError: true,
			errorContains: "Session token is missing",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Create request
			req := httptest.NewRequest(http.MethodGet, "/", nil)
			tc.setupRequest(req)

			// Call the function
			err := middleware.AuthenticateLandlordRequest(req)

			// Check results
			if tc.expectedError && err == nil {
				t.Errorf("Expected error but got nil")
			}
			if !tc.expectedError && err != nil {
				t.Errorf("Expected no error but got: %v", err)
			}
			if tc.errorContains != "" && err != nil {
				if !strings.Contains(err.Error(), tc.errorContains) {
					t.Errorf("Expected error to contain '%s' but got: %v", tc.errorContains, err)
				}
			}
		})
	}
}

func TestAuthenticateTenantRequest(t *testing.T) {
	// Initialize test environment
	testutil.InitTestEnv()
	go logs.LogProcessor()

	// Test cases
	testCases := []struct {
		name          string
		setupRequest  func(*http.Request)
		expectedError bool
		errorContains string
	}{
		{
			name: "Missing session token",
			setupRequest: func(req *http.Request) {
				// Only add CSRF token
				req.AddCookie(&http.Cookie{
					Name:  "csrf_token",
					Value: "csrf_token_hash_tenant@example.com",
				})
			},
			expectedError: true,
			errorContains: "Session token is missing",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Create request
			req := httptest.NewRequest(http.MethodGet, "/", nil)
			tc.setupRequest(req)

			// Call the function
			err := middleware.AuthenticateTenantRequest(req)

			// Check results
			if tc.expectedError && err == nil {
				t.Errorf("Expected error but got nil")
			}
			if !tc.expectedError && err != nil {
				t.Errorf("Expected no error but got: %v", err)
			}
			if tc.errorContains != "" && err != nil {
				if !strings.Contains(err.Error(), tc.errorContains) {
					t.Errorf("Expected error to contain '%s' but got: %v", tc.errorContains, err)
				}
			}
		})
	}
}
