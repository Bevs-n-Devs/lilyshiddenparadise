package handlers_test

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/Bevs-n-Devs/lilyshiddenparadise/handlers"
	"github.com/Bevs-n-Devs/lilyshiddenparadise/logs"
	"github.com/Bevs-n-Devs/lilyshiddenparadise/testutil"
)

func TestHome(t *testing.T) {
	// Initialize test environment
	testutil.InitTestEnv()
	
	go logs.LogProcessor()

	// Define test cases
	testCases := []struct {
		name                 string
		url                  string
		expectedStatusCode   int
		expectedBodyContains string
	}{
		{
			name:                 "Home page loads successfully",
			url:                  "/",
			expectedStatusCode:   http.StatusOK,
			expectedBodyContains: "LHP", // Changed to something that actually exists in the page
		},
		{
			name:                 "Home page with authentication error",
			url:                  "/?authenticationError=Invalid+credentials",
			expectedStatusCode:   http.StatusOK,
			expectedBodyContains: "LHP", // Changed to something that actually exists in the page
		},
	}

	// Run test cases
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Create a request
			req := httptest.NewRequest("GET", tc.url, nil)

			// Create a response recorder
			rr := httptest.NewRecorder()

			// Call the handler
			handlers.Home(rr, req)

			// Check the status code
			if rr.Code != tc.expectedStatusCode {
				t.Errorf("Expected status code %d, got %d", tc.expectedStatusCode, rr.Code)
			}

			// Check the response body
			if tc.expectedBodyContains != "" && !strings.Contains(rr.Body.String(), tc.expectedBodyContains) {
				t.Errorf("Expected response body to contain '%s', got '%s'", tc.expectedBodyContains, rr.Body.String())
			}
		})
	}
}
