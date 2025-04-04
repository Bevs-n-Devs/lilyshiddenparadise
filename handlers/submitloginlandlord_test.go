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

func TestSubmitLoginLandlord_Simple(t *testing.T) {
	// Initialize test environment
	testutil.InitTestEnv()

	go logs.LogProcessor()

	// Define test cases
	testCases := []struct {
		name               string
		method             string
		formValues         map[string]string
		expectedStatusCode int
		expectedLocation   string
	}{
		{
			name:   "Invalid method",
			method: http.MethodGet,
			formValues: map[string]string{
				"landlordEmail":    "test@example.com",
				"landlordPassword": "password123",
			},
			expectedStatusCode: http.StatusBadRequest,
			expectedLocation:   "/login/landlord?badRequest=BAD+REQUEST+400:+Invalid+request+method",
		},
	}

	// Run test cases
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Create form data
			form := url.Values{}
			for key, value := range tc.formValues {
				form.Add(key, value)
			}

			// Create a request
			req := httptest.NewRequest(tc.method, "/login/landlord/submit", strings.NewReader(form.Encode()))
			req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

			// Create a response recorder
			rr := httptest.NewRecorder()

			// Call the handler
			handlers.SubmitLoginLandlord(rr, req)

			// Check the status code
			if rr.Code != tc.expectedStatusCode {
				t.Errorf("Expected status code %d, got %d", tc.expectedStatusCode, rr.Code)
			}

			// Check for nil pointer before accessing the header
			if rr.Result() != nil {
				// Check the redirect location
				location := rr.Header().Get("Location")
				if location != tc.expectedLocation {
					t.Errorf("Expected redirect to '%s', got '%s'", tc.expectedLocation, location)
				}
			} else {
				t.Errorf("Expected non-nil result for response recorder")
			}
		})
	}
}
