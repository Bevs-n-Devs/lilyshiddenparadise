package db

import (
	"database/sql"
	"errors"
	"fmt"
	"os"
	"time"

	_ "embed"

	_ "github.com/lib/pq"

	"github.com/Bevs-n-Devs/lilyshiddenparadise/env"
	"github.com/Bevs-n-Devs/lilyshiddenparadise/logs"
	"github.com/Bevs-n-Devs/lilyshiddenparadise/utils"
)

/*
ConnectDB connects to the PostgreSQL database via the DATABASE_URL environment
variable. If this variable is empty, it attempts to load the environment
variables from the .env file. The function logs the progress of the
connection attempt and returns an error if the connection cannot be
established.

Returns:

- error: An error object if the connection cannot be established.
*/
func ConnectDB() error {
	var err error

	// connect to database via environment variable
	if os.Getenv("DATABASE_URL") == "" {
		logs.Logs(logWarning, "Could not get database URL from hosting platform. Loading from .env file...")
		err := env.LoadEnv("env/.env")
		if err != nil {
			logs.Logs(logErr, fmt.Sprintf("Could not load environment variables from .env file: %s", err.Error()))
			return err
		}
	}

	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		logs.Logs(logDbErr, "Database URL is empty!")
		return fmt.Errorf("database URL is empty")
	}

	logs.Logs(logDb, "Connecting to database...")
	db, err = sql.Open("postgres", dbURL) // open db connection from global db variable
	if err != nil {
		logs.Logs(logDbErr, fmt.Sprintf("Could not connect to database: %s", err.Error()))
		return err
	}

	// verify connection
	logs.Logs(logDb, "Verifying database connection...")
	if db == nil {
		logs.Logs(logDbErr, "Database connection is empty!")
		return errors.New("database connection not established")
	}
	err = db.Ping()
	if err != nil {
		logs.Logs(logDbErr, fmt.Sprintf("Cannot ping database: %s", err.Error()))
		return err
	}
	logs.Logs(logDb, "Database connection established.")
	return nil
}

/*
CreateNewLandlord creates a new landlord in the database.

Arguments:

- landlordEmail: The landlord's email address to store in the database.

- landlordPassword: The landlord's password to store in the database.

Returns:

- error: An error object if the landlord cannot be created.
*/
func CreateNewLandlord(landlordEmail, landlordPassword string) error {
	if db == nil {
		logs.Logs(logDbErr, "Database connection is empty!")
		return errors.New("database connection not established")
	}

	hashPassword, err := utils.HashedPassword(landlordPassword)
	if err != nil {
		logs.Logs(logErr, fmt.Sprintf("Error hashing password: %s", err.Error()))
		return err
	}

	query := `
	INSERT INTO lhp_landlords (email, password, created_at)
	VALUES ($1, $2, NOW());
	`
	_, err = db.Exec(query, landlordEmail, hashPassword)
	if err != nil {
		logs.Logs(logDbErr, fmt.Sprintf("Error creating new landlord: %s", err.Error()))
		return err
	}
	return nil
}

/*
AuthenticateLandlord checks if the provided email and password match the stored credentials.

It returns true if the credentials are correct, otherwise false. An error is returned if the query fails.

Returns:

- bool: True if the credentials are correct, otherwise false.

- error: An error if the query fails.
*/
func AuthenticateLandlord(email, password string) (bool, error) {
	if db == nil {
		logs.Logs(logDbErr, "Database connection is not initialized")
		return false, errors.New("database connection is not initialized")
	}

	var hashedPassword string
	query := `
	SELECT password 
	FROM lhp_landlords 
	WHERE email=$1;
	`
	err := db.QueryRow(query, email).Scan(&hashedPassword)
	if err != nil {
		return false, err
	}

	ok := utils.CheckPasswordHash(password, hashedPassword)
	if !ok {
		return false, errors.New("invalid password")
	}

	return true, nil
}

/*
UpdateLandlordSessionTokens generates new session and CSRF tokens for a given user and updates their expiry time in the database.

Arguments:

- email: A string representing the email for which the tokens should be updated.

Returns:

- string: The newly generated session token.

- string: The newly generated CSRF token.

- time.Time: The expiry time for the new tokens.

- error: An error object if the tokens cannot be generated or updated in the database.
*/
func UpdateLandlordSessionTokens(email string) (string, string, time.Time, error) {
	if db == nil {
		logs.Logs(logDbErr, "Database connection is not initialized")
		return "", "", time.Time{}, errors.New("database connection is not initialized")
	}

	sessionToken, err := utils.GenerateToken(32)
	if err != nil {
		logs.Logs(logDbErr, fmt.Sprintf("Failed to generate session token: %s", err.Error()))
		return "", "", time.Time{}, err
	}
	csrfToken, err := utils.GenerateToken(32)
	if err != nil {
		logs.Logs(logDbErr, fmt.Sprintf("Failed to generate CSRF token: %s", err.Error()))
		return "", "", time.Time{}, err
	}
	expiry := time.Now().Add(5 * time.Minute) // 5 minute validity

	query := `
	UPDATE lhp_landlords 
	SET session_token=$1, csrf_token=$2, token_expiry=$3 
	WHERE email=$4;
	`
	_, err = db.Exec(query, sessionToken, csrfToken, expiry, email)
	if err != nil {
		logs.Logs(logDbErr, fmt.Sprintf("Failed to update session tokens: %s", err.Error()))
		return "", "", time.Time{}, err
	}

	logs.Logs(logDb, "Session tokens updated successfully")
	return sessionToken, csrfToken, expiry, nil
}

/*
GetEmailFromLandlordSessionToken takes a session token as an argument and returns the email
address associated with that session token in the database.

Arguments:

- sessionToken: The session token to get the email address from the database.

Returns:

- string: The landlord's email address associated with the session token.

- error: An error object if the user is not found in the database or an error occurs while querying the database.
*/
func GetEmailFromLandlordSessionToken(sessionToken string) (string, error) {
	if db == nil {
		logs.Logs(logDbErr, "Database connection is not initialized")
		return "", errors.New("database connection is not initialized")
	}

	var email string
	query := `
	SELECT email 
	FROM lhp_landlords 
	WHERE session_token=$1;
	`
	err := db.QueryRow(query, sessionToken).Scan(&email)

	if err == sql.ErrNoRows {
		logs.Logs(logDbErr, "User not found")
		return "", errors.New("user not found")
	}

	if err != nil {
		logs.Logs(logDbErr, fmt.Sprintf("Failed to get session token: %s", err.Error()))
		return "", err
	}

	return email, nil
}

/*
ValidateLandlordSessionToken checks if a given session token matches the stored session token in the database
for a given landlord's email address.

Arguments:

- email: The landlord's email address to check.

- sessionToken: The session token to check against the stored session token.

Returns:

- bool: True if the session tokens match, false if not.

- error: An error object if the user is not found in the database or an error occurs while querying the database.
*/
func ValidateLandlordSessionToken(email, sessionToken string) (bool, error) {
	if db == nil {
		logs.Logs(logDbErr, "Database connection is not initialized")
		return false, errors.New("database connection is not initialized")
	}

	// query DB to get the stored session token
	var dbSessionToken string
	query := `
	SELECT session_token
	FROM lhp_landlords
	WHERE email = $1;
	`
	err := db.QueryRow(query, email).Scan(&dbSessionToken)

	if err == sql.ErrNoRows {
		logs.Logs(logDbErr, "User not found")
		return false, errors.New("user not found")
	}

	if err != nil {
		logs.Logs(logDbErr, fmt.Sprintf("Failed to get session token: %s", err.Error()))
		return false, err
	}

	// compare the input session token with DB session token
	if sessionToken != dbSessionToken {
		logs.Logs(logDbErr, "Invalid session token")
		return false, nil
	}

	return true, nil
}

/*
ValidateLandlordCSRFToken checks if a given CSRF token matches the stored CSRF token in the database
for a specified landlord's email address.

Arguments:

- email: The landlord's email address to check.

- csrfToken: The CSRF token to verify against the stored CSRF token.

Returns:

- bool: True if the CSRF tokens match, false otherwise.

- error: An error object if an error occurs while querying the database or if the database connection is not initialized.
*/
func ValidateLandlordCSRFToken(email, csrfToken string) (bool, error) {
	if db == nil {
		logs.Logs(logDbErr, "Database connection is not initialized")
		return false, errors.New("database connection is not initialized")
	}

	// query DB to get the stored CSRF token
	var dbCSRFToken string
	query := `
	SELECT csrf_token 
	FROM lhp_landlords 
	WHERE email=$1;
	`
	err := db.QueryRow(query, email).Scan(&dbCSRFToken)
	if err != nil {
		return false, err
	}

	// compare the input CSRF token with DB CSRF token
	if csrfToken != dbCSRFToken {
		return false, nil
	}
	return true, nil
}
