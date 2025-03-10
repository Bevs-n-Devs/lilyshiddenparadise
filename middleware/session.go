package middleware

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/Bevs-n-Devs/lilyshiddenparadise/db"
	"github.com/Bevs-n-Devs/lilyshiddenparadise/logs"
	"github.com/Bevs-n-Devs/lilyshiddenparadise/utils"
)

const (
	logInfo = 1
	logErr  = 3
)

var errAuth = errors.New("user not authenticated")

/*
AuthenticateLandlordRequest authenticates a landlord request by validating the session token
and CSRF token in the request. If the tokens are invalid, an error is returned.

	If the session token is missing, errAuth is returned with a message indicating the session
	token is missing.

	If the session token or CSRF token do not exist in the database or the database connection
	is not initialized, an error is returned.

	If the session token or CSRF token is invalid, errAuth is returned with a message indicating
	the token is invalid.
*/
func AuthenticateLandlordRequest(r *http.Request) error {
	// get the session token from the cookie

	sessionToken, err := utils.CheckSessionToken(r)
	if err != nil {
		return fmt.Errorf("%s Session token is missing", err.Error())
	}

	email, err := db.GetEmailFromLandlordSessionToken(sessionToken.Value)
	if err != nil {
		logs.Logs(logErr, fmt.Sprintf("Failed to get email from session token: %s", err.Error()))
		return fmt.Errorf("%s! Failed to get email from session token", errAuth)
	}

	// check if email and session token exists in the database
	exists, err := db.ValidateLandlordSessionToken(email, sessionToken.Value)
	if err != nil {
		logs.Logs(logErr, fmt.Sprintf("Failed to validate landlord session token: %s", err.Error()))
		return fmt.Errorf("%s! Failed to validate landlord session token", errAuth)
	}
	if !exists {
		logs.Logs(logErr, "Invlaid landlord session token.")
		return fmt.Errorf("%s! Invalid landlord session token", errAuth)
	}
	logs.Logs(logInfo, fmt.Sprintf("Landlord session validation result: %t", exists))

	// get csrf token from the cookie
	csrfToken, err := utils.CheckCSRFToken(r)
	if err != nil {
		return fmt.Errorf("%s Failed to get CSRF token", err.Error())
	}

	// check if email and csrf token exists in the database
	exists, err = db.ValidateLandlordCSRFToken(email, csrfToken.Value)
	if err != nil {
		logs.Logs(logErr, fmt.Sprintf("Failed to validate landlord CSRF token: %s", err.Error()))
		return fmt.Errorf("%s! Failed to validate landlord CSRF token", errAuth)
	}
	if !exists {
		logs.Logs(logErr, "Invalid landlord CSRF token.")
		return fmt.Errorf("%s! Invalid landlord CSRF token", errAuth)
	}
	logs.Logs(logInfo, fmt.Sprintf("Landlord CSRF token validation result: %t", exists))

	return nil
}
