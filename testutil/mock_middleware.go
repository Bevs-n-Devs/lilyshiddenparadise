package testutil

import (
	"net/http"
	"time"
)

// MockMiddleware provides mock implementations of middleware functions for testing
type MockMiddleware struct{}

// LandlordDashboardSessionCookie is a mock implementation for testing
func (m *MockMiddleware) LandlordDashboardSessionCookie(w http.ResponseWriter, sessionToken string, expiryTime time.Time) bool {
	http.SetCookie(w, &http.Cookie{
		Name:     "session_token",
		Value:    sessionToken,
		Expires:  expiryTime,
		HttpOnly: true,
		Path:     "/landlord/dashboard",
	})
	return true
}

// LandlordDashboardCSRFTokenCookie is a mock implementation for testing
func (m *MockMiddleware) LandlordDashboardCSRFTokenCookie(w http.ResponseWriter, csrfToken string, expiryTime time.Time) bool {
	http.SetCookie(w, &http.Cookie{
		Name:     "csrf_token",
		Value:    csrfToken,
		Expires:  expiryTime,
		HttpOnly: false,
		Path:     "/landlord/dashboard",
	})
	return true
}
