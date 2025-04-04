package middleware_test

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/Bevs-n-Devs/lilyshiddenparadise/logs"
	"github.com/Bevs-n-Devs/lilyshiddenparadise/middleware"
	"github.com/Bevs-n-Devs/lilyshiddenparadise/testutil"
)

func TestLandlordDashboardCookies(t *testing.T) {
	// Initialize test environment
	testutil.InitTestEnv()
	go logs.LogProcessor()

	// Test cases for session and CSRF cookies
	testCases := []struct {
		name           string
		cookieFunction func(http.ResponseWriter, string, time.Time) bool
		cookieName     string
		cookiePath     string
		httpOnly       bool
	}{
		{
			name:           "LandlordDashboardSessionCookie",
			cookieFunction: middleware.LandlordDashboardSessionCookie,
			cookieName:     "session_token",
			cookiePath:     "/landlord/dashboard",
			httpOnly:       true,
		},
		{
			name:           "LandlordDashboardCSRFTokenCookie",
			cookieFunction: middleware.LandlordDashboardCSRFTokenCookie,
			cookieName:     "csrf_token",
			cookiePath:     "/landlord/dashboard",
			httpOnly:       false,
		},
		{
			name:           "LandlordDashboardTenantsSessionCookie",
			cookieFunction: middleware.LandlordDashboardTenantsSessionCookie,
			cookieName:     "session_token",
			cookiePath:     "/landlord/dashboard/tenants",
			httpOnly:       true,
		},
		{
			name:           "LandlordDashboardTenantsCSRFTokenCookie",
			cookieFunction: middleware.LandlordDashboardTenantsCSRFTokenCookie,
			cookieName:     "csrf_token",
			cookiePath:     "/landlord/dashboard/tenants",
			httpOnly:       false,
		},
		{
			name:           "TenantDashboardSessionCookie",
			cookieFunction: middleware.TenantDashboardSessionCookie,
			cookieName:     "session_token",
			cookiePath:     "/tenant/dashboard",
			httpOnly:       true,
		},
		{
			name:           "TenantDashboardCSRFTokenCookie",
			cookieFunction: middleware.TenantDashboardCSRFTokenCookie,
			cookieName:     "csrf_token",
			cookiePath:     "/tenant/dashboard",
			httpOnly:       false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Create a response recorder
			w := httptest.NewRecorder()

			// Set the cookie
			tokenValue := "test_token_value"
			expiryTime := time.Now().Add(30 * time.Minute)
			result := tc.cookieFunction(w, tokenValue, expiryTime)

			// Check the result
			if !result {
				t.Errorf("Expected cookie function to return true")
			}

			// Get the cookies from the response
			response := w.Result()
			cookies := response.Cookies()

			// Check if the cookie was set correctly
			var found bool
			for _, cookie := range cookies {
				if cookie.Name == tc.cookieName {
					found = true
					if cookie.Value != tokenValue {
						t.Errorf("Expected cookie value to be '%s', got '%s'", tokenValue, cookie.Value)
					}
					if cookie.Path != tc.cookiePath {
						t.Errorf("Expected cookie path to be '%s', got '%s'", tc.cookiePath, cookie.Path)
					}
					if cookie.HttpOnly != tc.httpOnly {
						t.Errorf("Expected cookie HttpOnly to be %v, got %v", tc.httpOnly, cookie.HttpOnly)
					}
					if cookie.SameSite != http.SameSiteStrictMode {
						t.Errorf("Expected cookie SameSite to be %v, got %v", http.SameSiteStrictMode, cookie.SameSite)
					}
					// Check expiry time with some tolerance
					timeDiff := cookie.Expires.Sub(expiryTime)
					if timeDiff < -time.Second || timeDiff > time.Second {
						t.Errorf("Expected cookie expiry to be close to %v, got %v", expiryTime, cookie.Expires)
					}
				}
			}

			if !found {
				t.Errorf("Cookie '%s' not found in response", tc.cookieName)
			}
		})
	}
}

func TestDeleteCookies(t *testing.T) {
	// Initialize test environment
	testutil.InitTestEnv()
	go logs.LogProcessor()

	// Test cases for delete cookie functions
	testCases := []struct {
		name           string
		cookieFunction func(http.ResponseWriter) bool
		cookieName     string
	}{
		{
			name:           "DeleteLandlordSessionCookie",
			cookieFunction: middleware.DeleteLandlordSessionCookie,
			cookieName:     "session_token",
		},
		{
			name:           "DeleteLandlordCSRFCookie",
			cookieFunction: middleware.DeleteLandlordCSRFCookie,
			cookieName:     "csrf_token",
		},
		{
			name:           "DeleteTenantSessionCookie",
			cookieFunction: middleware.DeleteTenantSessionCookie,
			cookieName:     "session_token",
		},
		{
			name:           "DeleteTenantCSRFCookie",
			cookieFunction: middleware.DeleteTenantCSRFCookie,
			cookieName:     "csrf_token",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Create a response recorder
			w := httptest.NewRecorder()

			// Delete the cookie
			result := tc.cookieFunction(w)

			// Check the result
			if !result {
				t.Errorf("Expected cookie function to return true")
			}

			// Get the cookies from the response
			response := w.Result()
			cookies := response.Cookies()

			// Check if the cookie was deleted correctly
			var found bool
			for _, cookie := range cookies {
				if cookie.Name == tc.cookieName {
					found = true
					if cookie.Value != "" {
						t.Errorf("Expected cookie value to be empty, got '%s'", cookie.Value)
					}
					if cookie.Path != "/" {
						t.Errorf("Expected cookie path to be '/', got '%s'", cookie.Path)
					}
					// Check if expiry is in the past
					if !cookie.Expires.Before(time.Now()) {
						t.Errorf("Expected cookie expiry to be in the past, got %v", cookie.Expires)
					}
				}
			}

			if !found {
				t.Errorf("Cookie '%s' not found in response", tc.cookieName)
			}
		})
	}
}

func TestSpecializedCookies(t *testing.T) {
	// Initialize test environment
	testutil.InitTestEnv()
	go logs.LogProcessor()

	// Test LandlordTenantMessagesSessionCookie and LandlordTenantMessagesCSRFTokenCookie
	t.Run("LandlordTenantMessagesCookies", func(t *testing.T) {
		// Create a response recorder
		w := httptest.NewRecorder()

		// Set the cookies
		tenantID := "123"
		sessionToken := "test_session_token"
		csrfToken := "test_csrf_token"
		expiryTime := time.Now().Add(30 * time.Minute)

		// Set session cookie
		result1 := middleware.LandlordTenantMessagesSessionCookie(w, tenantID, sessionToken, expiryTime)
		if !result1 {
			t.Errorf("Expected LandlordTenantMessagesSessionCookie to return true")
		}

		// Set CSRF cookie
		result2 := middleware.LandlordTenantMessagesCSRFTokenCookie(w, tenantID, csrfToken, expiryTime)
		if !result2 {
			t.Errorf("Expected LandlordTenantMessagesCSRFTokenCookie to return true")
		}

		// Get the cookies from the response
		response := w.Result()
		cookies := response.Cookies()

		// Expected paths
		expectedSessionPath := "/landlord/dashboard/messages/tenant/" + tenantID
		expectedCSRFPath := "/landlord/dashboard/messages/tenant/" + tenantID

		// Check if the cookies were set correctly
		var sessionFound, csrfFound bool
		for _, cookie := range cookies {
			if cookie.Name == "session_token" {
				sessionFound = true
				if cookie.Value != sessionToken {
					t.Errorf("Expected session cookie value to be '%s', got '%s'", sessionToken, cookie.Value)
				}
				if cookie.Path != expectedSessionPath {
					t.Errorf("Expected session cookie path to be '%s', got '%s'", expectedSessionPath, cookie.Path)
				}
				if !cookie.HttpOnly {
					t.Errorf("Expected session cookie HttpOnly to be true")
				}
			}
			if cookie.Name == "csrf_token" {
				csrfFound = true
				if cookie.Value != csrfToken {
					t.Errorf("Expected CSRF cookie value to be '%s', got '%s'", csrfToken, cookie.Value)
				}
				if cookie.Path != expectedCSRFPath {
					t.Errorf("Expected CSRF cookie path to be '%s', got '%s'", expectedCSRFPath, cookie.Path)
				}
				if cookie.HttpOnly {
					t.Errorf("Expected CSRF cookie HttpOnly to be false")
				}
			}
		}

		if !sessionFound {
			t.Errorf("Session cookie not found in response")
		}
		if !csrfFound {
			t.Errorf("CSRF cookie not found in response")
		}
	})

	// Test LogoutLandlordSessionCookie and LogoutLandlordCSRFTokenCookie
	t.Run("LogoutLandlordCookies", func(t *testing.T) {
		// Create a response recorder
		w := httptest.NewRecorder()

		// Set the cookies
		sessionToken := "test_session_token"
		csrfToken := "test_csrf_token"

		// Set logout cookies
		result1 := middleware.LogoutLandlordSessionCookie(w, sessionToken)
		if !result1 {
			t.Errorf("Expected LogoutLandlordSessionCookie to return true")
		}

		result2 := middleware.LogoutLandlordCSRFTokenCookie(w, csrfToken)
		if !result2 {
			t.Errorf("Expected LogoutLandlordCSRFTokenCookie to return true")
		}

		// Get the cookies from the response
		response := w.Result()
		cookies := response.Cookies()

		// Expected paths
		expectedPath := "/logout-landlord"

		// Check if the cookies were set correctly
		var sessionFound, csrfFound bool
		for _, cookie := range cookies {
			if cookie.Name == "session_token" {
				sessionFound = true
				if cookie.Value != sessionToken {
					t.Errorf("Expected session cookie value to be '%s', got '%s'", sessionToken, cookie.Value)
				}
				if cookie.Path != expectedPath {
					t.Errorf("Expected session cookie path to be '%s', got '%s'", expectedPath, cookie.Path)
				}
				if !cookie.HttpOnly {
					t.Errorf("Expected session cookie HttpOnly to be true")
				}
			}
			if cookie.Name == "csrf_token" {
				csrfFound = true
				if cookie.Value != csrfToken {
					t.Errorf("Expected CSRF cookie value to be '%s', got '%s'", csrfToken, cookie.Value)
				}
				if cookie.Path != expectedPath {
					t.Errorf("Expected CSRF cookie path to be '%s', got '%s'", expectedPath, cookie.Path)
				}
				if cookie.HttpOnly {
					t.Errorf("Expected CSRF cookie HttpOnly to be false")
				}
			}
		}

		if !sessionFound {
			t.Errorf("Session cookie not found in response")
		}
		if !csrfFound {
			t.Errorf("CSRF cookie not found in response")
		}
	})
}
