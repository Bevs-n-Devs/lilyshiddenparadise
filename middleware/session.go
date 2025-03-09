package middleware

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/Bevs-n-Devs/lilyshiddenparadise/db"
	"github.com/Bevs-n-Devs/lilyshiddenparadise/logs"
)

const (
	logInfo = 1
	logErr  = 3
)

var ErrAuth = errors.New("user not authenticated")

func AuthenticateLandlordRequest(r *http.Request) error {
	// get the session token from the cookie
	sessionToken, err := r.Cookie("landlord_session_token")
	if err != nil {
		logs.Logs(logErr, fmt.Sprintf("Failed to get session token: %s", err.Error()))
		return fmt.Errorf("%s! Session token is missing", ErrAuth)
	}

	email, err := db.GetEmailFromLandlordSessionToken(sessionToken.Value)
	if err != nil {
		logs.Logs(logErr, fmt.Sprintf("Failed to get email from session token: %s", err.Error()))
		return fmt.Errorf("%s! Failed to get email from session token", ErrAuth)
	}

	// check if email and session token exists in the database
	exists, err := db.ValidateLandlordSessionToken(email, sessionToken.Value)
	if err != nil {
		logs.Logs(logErr, fmt.Sprintf("Failed to validate landlord session token: %s", err.Error()))
		return fmt.Errorf("%s! Failed to validate landlord session token", ErrAuth)
	}
	if !exists {
		logs.Logs(logErr, "Invlaid landlord session token.")
		return fmt.Errorf("%s! Invalid landlord session token", ErrAuth)
	}

	// get csrf token from the cookie
	csrfToken, err := r.Cookie("landlord_csrf_token")
	if err != nil {
		logs.Logs(logErr, fmt.Sprintf("Failed to get CSRF token: %s", err.Error()))
		return fmt.Errorf("%s! Failed to get CSRF token", ErrAuth)
	}

	// check if email and csrf token exists in the database
	exists, err = db.ValidateLandlordCSRFToken(email, csrfToken.Value)
	if err != nil {
		logs.Logs(logErr, fmt.Sprintf("Failed to validate landlord CSRF token: %s", err.Error()))
		return fmt.Errorf("%s! Failed to validate landlord CSRF token", ErrAuth)
	}
	if !exists {
		logs.Logs(logErr, "Invalid landlord CSRF token.")
		return fmt.Errorf("%s! Invalid landlord CSRF token", ErrAuth)
	}

	return nil

}
