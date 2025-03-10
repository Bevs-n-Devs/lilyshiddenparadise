package middleware

import (
	"net/http"
	"time"
)

/*
LandlordDashboardSessionCookie sets a cookie on the response with the session token, expiry time, and path set to /landlord/dashboard.
This is used to authenticate the landlord on the dashboard page.

Parameters:

- w: The http.ResponseWriter to set the cookie on.

- sessionToken: The session token to set in the cookie.

- expiryTime: The expiry time of the cookie.

Returns:

- bool: True if the cookie is set successfully, false otherwise.
*/
func LandlordDashboardSessionCookie(w http.ResponseWriter, sessionToken string, expiryTime time.Time) bool {
	http.SetCookie(w, &http.Cookie{
		Name:     "session_token",
		Value:    sessionToken,
		Expires:  expiryTime,
		HttpOnly: true,
		Path:     "/landlord/dashboard",
		SameSite: http.SameSiteStrictMode,
	})
	return true
}

func LandlordDashboardCSRFTokenCookie(w http.ResponseWriter, csrfToken string, expiryTime time.Time) bool {
	http.SetCookie(w, &http.Cookie{
		Name:     "csrf_token",
		Value:    csrfToken,
		Expires:  expiryTime,
		HttpOnly: true,
		Path:     "/landlord/dashboard",
		SameSite: http.SameSiteStrictMode,
	})
	return true
}

/*
LogoutLandlordSessionCookie sets a cookie on the response to log out the landlord by setting the session token
with the specified value and path to /logout-landlord.

Parameters:

- w: The http.ResponseWriter to set the cookie on.

- sessionToken: The session token retrieved from the request cookie.

Returns:

- bool: True if the cookie is set successfully.
*/
func LogoutLandlordSessionCookie(w http.ResponseWriter, sessionToken *http.Cookie) bool {
	http.SetCookie(w, &http.Cookie{
		Name:     "session_token",
		Value:    sessionToken.Value,
		HttpOnly: true,
		Path:     "/logout-landlord",
		SameSite: http.SameSiteStrictMode,
	})
	return true
}

/*
LogoutLandlordCSRFTokenCookie sets a cookie on the response to log out the landlord by setting the CSRF token
with the specified value and path to /logout-landlord.

Parameters:

- w: The http.ResponseWriter to set the cookie on.

- csrfToken: The CSRF token retrieved from the request cookie.

Returns:

- bool: True if the cookie is set successfully.
*/
func LogoutLandlordCSRFTokenCookie(w http.ResponseWriter, csrfToken *http.Cookie) bool {
	http.SetCookie(w, &http.Cookie{
		Name:     "csrf_token",
		Value:    csrfToken.Value,
		HttpOnly: true,
		Path:     "/logout-landlord",
		SameSite: http.SameSiteStrictMode,
	})
	return true
}

/*
DeleteLandlordSessionCookie deletes the session token cookie for the landlord by setting its value to empty
and its expiry time to the past. This is used to log out the landlord.

Parameters:

- w: The http.ResponseWriter to delete the cookie on.

Returns:

- bool: True if the cookie is deleted successfully.
*/
func DeleteLandlordSessionCookie(w http.ResponseWriter) bool {
	http.SetCookie(w, &http.Cookie{
		Name:     "session_token",
		Value:    "",
		Expires:  time.Now().Add(-time.Hour),
		HttpOnly: true,
		Path:     "/",
		SameSite: http.SameSiteStrictMode,
	})
	return true
}

/*
DeleteLandlordCSRFCookie deletes the CSRF token cookie for the landlord by setting its value to empty
and its expiry time to the past. This is used to log out the landlord.

Parameters:

- w: The http.ResponseWriter to delete the cookie on.

Returns:

- bool: True if the cookie is deleted successfully.
*/
func DeleteLandlordCSRFCookie(w http.ResponseWriter) bool {
	http.SetCookie(w, &http.Cookie{
		Name:     "csrf_token",
		Value:    "",
		Expires:  time.Now().Add(-time.Hour),
		HttpOnly: true,
		Path:     "/",
		SameSite: http.SameSiteStrictMode,
	})
	return true
}
