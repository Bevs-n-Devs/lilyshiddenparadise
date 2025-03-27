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

func GetTenantNameByEmail(email string) (string, error) {
	if db == nil {
		logs.Logs(logDbErr, "Database connection is empty!")
		return "", errors.New("database connection not established")
	}

	hashEmail := utils.HashData(email)

	var encryptedTenantName string
	query := `
	SELECT encrypt_full_name
	FROM lhp_tenant_application 
	WHERE hash_email = $1;
	`
	err := db.QueryRow(query, hashEmail).Scan(&encryptedTenantName)
	if err != nil {
		logs.Logs(logDbErr, fmt.Sprintf("Error getting tenant name: %s", err.Error()))
		return "", err
	}

	return encryptedTenantName, nil
}

/*
CreateNewTenant creates a new tenant record in the database.

This function checks if the database connection is initialized and retrieves the landlord email from environment variables.
If the landlord email is not found, it attempts to load it from a .env file. It then retrieves the landlord ID using the
landlord email. The tenant's email and password are hashed and encrypted along with other tenant details such as room type,
move-in date, rent due, and monthly rent. These details are then inserted into the lhp_tenants table.

Arguments:

- tenantEmail: The email address of the tenant.

- tenantPassword: The password for the tenant's account.

- roomType: The type of room assigned to the tenant.

- moveInDate: The date the tenant will move in.

- rentDue: The date rent payment is due each month.

- monthlyRent: The amount of rent due each month.

- currency: The currency in which the rent is paid.

Returns:

- error: An error object if the tenant cannot be created.
*/
func CreateNewTenant(tenantEmail, tenantPassword, roomType, moveInDate, rentDue, monthlyRent, currency string) error {
	if db == nil {
		logs.Logs(logDbErr, "Database connection is not initialized")
		return errors.New("database connection is not initialized")
	}

	// get landlord email via environment variable
	if os.Getenv("LANDLORD_EMAIL") == "" {
		logs.Logs(logWarning, "Could not get landlord email from hosting platform. Loading from .env file...")
		err := env.LoadEnv("env/.env")
		if err != nil {
			logs.Logs(logErr, fmt.Sprintf("Could not load environment variables from .env file: %s", err.Error()))
			return err
		}
	}

	landlordEmail := os.Getenv("LANDLORD_EMAIL")
	if landlordEmail == "" {
		logs.Logs(logDbErr, "Landlord email is empty!")
		return errors.New("landlord email is empty")
	}

	// get landlord id
	landlordId, err := GetLandlordIdByEmail(landlordEmail)
	if err != nil {
		logs.Logs(logDbErr, fmt.Sprintf("Failed to get landlord ID: %s", err.Error()))
		return err
	}

	// get encrypted tenant name via tenantEmail
	encryptedTenantName, err := GetTenantNameByEmail(tenantEmail)
	if err != nil {
		logs.Logs(logDbErr, fmt.Sprintf("Failed to get tenant name: %s", err.Error()))
		return err
	}

	// hash & encrypt identifiers
	hashEmail := utils.HashData(tenantEmail)
	hashPassword := utils.HashData(tenantPassword)

	encrypt_email, err := utils.Encrypt([]byte(tenantEmail))
	if err != nil {
		logs.Logs(logDbErr, fmt.Sprintf("Failed to encrypt email: %s", err.Error()))
		return err
	}

	encrypt_password, err := utils.Encrypt([]byte(tenantPassword))
	if err != nil {
		logs.Logs(logDbErr, fmt.Sprintf("Failed to encrypt password: %s", err.Error()))
		return err
	}

	encrypt_room_type, err := utils.Encrypt([]byte(roomType))
	if err != nil {
		logs.Logs(logDbErr, fmt.Sprintf("Failed to encrypt room type: %s", err.Error()))
		return err
	}

	encrypt_move_in_date, err := utils.Encrypt([]byte(moveInDate))
	if err != nil {
		logs.Logs(logDbErr, fmt.Sprintf("Failed to encrypt move in date: %s", err.Error()))
		return err
	}

	encrypt_rent_due, err := utils.Encrypt([]byte(rentDue))
	if err != nil {
		logs.Logs(logDbErr, fmt.Sprintf("Failed to encrypt rent due: %s", err.Error()))
		return err
	}

	encrypt_monthly_rent, err := utils.Encrypt([]byte(monthlyRent))
	if err != nil {
		logs.Logs(logDbErr, fmt.Sprintf("Failed to encrypt monthly rent: %s", err.Error()))
		return err
	}

	query := `
	INSERT INTO lhp_tenants (
		landlord_id,
		hash_email,
		hash_password,
		created_at,
		encrypt_email,
		encrypt_password,
		encrypt_room_type,
		encrypt_move_in_date,
		encrypt_rent_due,
		encrypt_monthly_rent,
		currency,
		encrypt_tenant_name,
	)
	VALUES ($1, $2, $3, NOW(), $4, $5, $6, $7, $8, $9, $10, $11);
	`
	_, err = db.Exec(query, landlordId, hashEmail, hashPassword, encrypt_email, encrypt_password, encrypt_room_type, encrypt_move_in_date, encrypt_rent_due, encrypt_monthly_rent, currency, encryptedTenantName)
	if err != nil {
		logs.Logs(logDbErr, fmt.Sprintf("Failed to create new tenant: %s", err.Error()))
		return err
	}
	return nil
}

func ManuallyCreateNewTenant(tenantFullName, tenantPassportID, tenantEmail, roomType, moveInDate, rentDue, monthlyRent, currency string) error {
	if db == nil {
		logs.Logs(logDbErr, "Database connection is not initialized")
		return errors.New("database connection is not initialized")
	}

	// get landlord email via environment variable
	if os.Getenv("LANDLORD_EMAIL") == "" {
		logs.Logs(logWarning, "Could not get landlord email from hosting platform. Loading from .env file...")
		err := env.LoadEnv("env/.env")
		if err != nil {
			logs.Logs(logErr, fmt.Sprintf("Could not load environment variables from .env file: %s", err.Error()))
			return err
		}
	}

	landlordEmail := os.Getenv("LANDLORD_EMAIL")
	if landlordEmail == "" {
		logs.Logs(logDbErr, "Landlord email is empty!")
		return errors.New("landlord email is empty")
	}

	// get landlord id
	landlordId, err := GetLandlordIdByEmail(landlordEmail)
	if err != nil {
		logs.Logs(logDbErr, fmt.Sprintf("Failed to get landlord ID: %s", err.Error()))
		return err
	}

	// hash & encrypt identifiers
	hashEmail, hashPassword, err := utils.GenerateTenantUsernamePassportNumberAndPassword(tenantEmail, tenantPassportID)
	if err != nil {
		logs.Logs(logDbErr, fmt.Sprintf("Failed to generate hash of tenant username & password: %s", err.Error()))
		return err
	}

	encryptName, err := utils.Encrypt([]byte(tenantFullName))
	if err != nil {
		logs.Logs(logDbErr, fmt.Sprintf("Failed to encrypt tenant name: %s", err.Error()))
		return err
	}

	encryptEmail, err := utils.Encrypt([]byte(tenantEmail))
	if err != nil {
		logs.Logs(logDbErr, fmt.Sprintf("Failed to encrypt email: %s", err.Error()))
		return err
	}

	encryptPassword, err := utils.Encrypt([]byte(tenantPassportID))
	if err != nil {
		logs.Logs(logDbErr, fmt.Sprintf("Failed to encrypt password: %s", err.Error()))
		return err
	}

	encryptRoomType, err := utils.Encrypt([]byte(roomType))
	if err != nil {
		logs.Logs(logDbErr, fmt.Sprintf("Failed to encrypt room type: %s", err.Error()))
		return err
	}

	encryptMoveInDate, err := utils.Encrypt([]byte(moveInDate))
	if err != nil {
		logs.Logs(logDbErr, fmt.Sprintf("Failed to encrypt move in date: %s", err.Error()))
		return err
	}

	encryptRentDue, err := utils.Encrypt([]byte(rentDue))
	if err != nil {
		logs.Logs(logDbErr, fmt.Sprintf("Failed to encrypt rent due: %s", err.Error()))
		return err
	}

	encryptMonthlyRent, err := utils.Encrypt([]byte(monthlyRent))
	if err != nil {
		logs.Logs(logDbErr, fmt.Sprintf("Failed to encrypt monthly rent: %s", err.Error()))
		return err
	}

	query := `
	INSERT INTO lhp_tenants (
		landlord_id,
		hash_email,
		hash_password,
		encrypt_tenant_name,
		encrypt_email,
		encrypt_password,
		encrypt_room_type,
		encrypt_move_in_date,
		encrypt_rent_due,
		encrypt_monthly_rent,
		currency,
		created_at,
	)
	VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, NOW());
	`
	_, err = db.Exec(query, landlordId, hashEmail, hashPassword, encryptName, encryptEmail, encryptPassword, encryptRoomType, encryptMoveInDate, encryptRentDue, encryptMonthlyRent, currency)
	if err != nil {
		logs.Logs(logDbErr, fmt.Sprintf("Failed to create new tenant: %s", err.Error()))
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
AuthenticateTenant checks if the provided email and password match the stored credentials.

It returns true if the credentials are correct, otherwise false. An error is returned if the query fails.

Returns:

- bool: True if the credentials are correct, otherwise false.

- error: An error if the query fails.
*/
func AuthenticateTenant(username, password string) (bool, error) {
	if db == nil {
		logs.Logs(logDbErr, "Database connection is not initialized")
		return false, errors.New("database connection is not initialized")
	}

	var (
		hashEmail    string
		hashPassword string
	)
	query := `
	SELECT hash_email, hash_password 
	FROM lhp_tenants 
	WHERE hash_email=$1
		AND hash_password=$2;
	`
	err := db.QueryRow(query, username, password).Scan(&hashEmail, &hashPassword)
	if err != nil {
		return false, err
	}

	// verify username and password
	verifyUsername := utils.VerifyHash(username, hashEmail)
	if !verifyUsername {
		return false, errors.New("invalid username")
	}
	verifyPassword := utils.VerifyHash(password, hashPassword)
	if !verifyPassword {
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
	expiry := time.Now().Add(30 * time.Second) // 30 seconds validity

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

func UpdateTenantSessionTokens(hash_email string) (string, string, time.Time, error) {
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
	expiry := time.Now().Add(30 * time.Second) // 30 seconds validity

	query := `
	UPDATE lhp_tenants 
	SET session_token=$1, csrf_token=$2, token_expiry=$3 
	WHERE hash_email=$4;
	`
	_, err = db.Exec(query, sessionToken, csrfToken, expiry, hash_email)
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

func GetEmailFromTenantSessionToken(sessionToken string) (string, error) {
	if db == nil {
		logs.Logs(logDbErr, "Database connection is not initialized")
		return "", errors.New("database connection is not initialized")
	}

	var email string
	query := `
	SELECT email 
	FROM lhp_tenants 
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

func ValidateTenantSessionToken(hashEmail, sessionToken string) (bool, error) {
	if db == nil {
		logs.Logs(logDbErr, "Database connection is not initialized")
		return false, errors.New("database connection is not initialized")
	}

	// query DB to get the stored session token
	var dbSessionToken string
	query := `
	SELECT session_token
	FROM lhp_tenants
	WHERE hash_email = $1;
	`
	err := db.QueryRow(query, hashEmail).Scan(&dbSessionToken)

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

func ValidateTenantCSRFToken(hashEmail, csrfToken string) (bool, error) {
	if db == nil {
		logs.Logs(logDbErr, "Database connection is not initialized")
		return false, errors.New("database connection is not initialized")
	}

	// query DB to get the stored CSRF token
	var dbCSRFToken string
	query := `
	SELECT csrf_token 
	FROM lhp_tenants 
	WHERE hash_email=$1;
	`
	err := db.QueryRow(query, hashEmail).Scan(&dbCSRFToken)
	if err != nil {
		return false, err
	}

	// compare the input CSRF token with DB CSRF token
	if csrfToken != dbCSRFToken {
		return false, nil
	}
	return true, nil
}

/*
LogoutLandlord removes a landlord's session token, CSRF token and expiry time from the database.

Arguments:

- email: The landlord's email address to remove the session token, CSRF token and expiry time for.

Returns:

- error: An error object if the logout operation fails (e.g. if the database connection is not initialized).
*/
func LogoutLandlord(email string) error {
	if db == nil {
		logs.Logs(logDbErr, "Database connection is not initialized")
		return errors.New("database connection is not initialized")
	}

	query := `
	UPDATE lhp_landlords 
	SET session_token=NULL, csrf_token=NULL, token_expiry=NULL 
	WHERE email=$1;
	`
	_, err := db.Exec(query, email)
	if err != nil {
		return err
	}
	return nil
}

func LogoutTenant(hashEmail string) error {
	if db == nil {
		logs.Logs(logDbErr, "Database connection is not initialized")
		return errors.New("database connection is not initialized")
	}

	query := `
	UPDATE lhp_tenants 
	SET session_token=NULL, csrf_token=NULL, token_expiry=NULL 
	WHERE hash_email=$1;
	`
	_, err := db.Exec(query, hashEmail)
	if err != nil {
		return err
	}
	return nil
}

/*
GetLandlordIdByEmail takes an email address as an argument and returns the landlord_id associated with that email address in the database.

Arguments:

- email: The email address to get the landlord_id for.

Returns:

- int: The landlord_id associated with the email address.

- error: An error object if the user is not found in the database or an error occurs while querying the database.
*/
func GetLandlordIdByEmail(email string) (int, error) {
	if db == nil {
		logs.Logs(logDbErr, "Database connection is not initialized")
		return 0, errors.New("database connection is not initialized")
	}

	var landlordId int
	query := `
	SELECT id 
	FROM lhp_landlords 
	WHERE email=$1;
	`
	err := db.QueryRow(query, email).Scan(&landlordId)
	if err != nil {
		return 0, err
	}
	return landlordId, nil
}

func GetTenantIdByEmail(email string) (int, error) {
	if db == nil {
		logs.Logs(logDbErr, "Database connection is not initialized")
		return 0, errors.New("database connection is not initialized")
	}

	var tenantId int
	query := `
	SELECT id 
	FROM lhp_tenants 
	WHERE email=$1;
	`
	err := db.QueryRow(query, email).Scan(&tenantId)
	if err != nil {
		return 0, err
	}
	return tenantId, nil
}

func SendMessage(senderId int, senderType string, receiverID int, receiverType string, message string) error {
	if db == nil {
		logs.Logs(logDbErr, "Database connection is not initialized")
		return errors.New("database connection is not initialized")
	}

	encryptMessage, err := utils.Encrypt([]byte(message))
	if err != nil {
		logs.Logs(logDbErr, fmt.Sprintf("Failed to encrypt message: %s", err.Error()))
		return err
	}

	query := `
	INSERT INTO lhp_messages 
		(
			sender_id, 
			sender_type, 
			receiver_id, 
			receiver_type, 
			encrypt_message, 
			sent_at
		)
	VALUES ($1, $2, $3, $4, $5, NOW());
	`
	_, err = db.Exec(
		query,
		senderId,
		senderType,
		receiverID,
		receiverType,
		encryptMessage,
	)
	if err != nil {
		logs.Logs(logDbErr, fmt.Sprintf("Failed to send message: %s", err.Error()))
		return err
	}

	logs.Logs(logDb, fmt.Sprintf("Message from %s successfully sent to %s", senderType, receiverType))
	return nil
}

/*
SaveTenantApplicationForm saves a tenant application form to the database.

The function takes in all the values from the tenant application form and stores them in the database.
The function first checks if the database connection is initialized and if the landlord email is not empty.
Then, it gets the landlord ID from the database and hashes the identifiers.
Next, the function encrypts the data using the aes encryption algorithm.
Finally, the function executes a SQL query to store the tenant application form to the database.

Arguments:

- fullName: The tenant's full name.

- dateOfBirth: The tenant's date of birth.

- passportNumber: The tenant's passport number.

- phoneNumber: The tenant's phone number.

- email: The tenant's email address.

- occupation: The tenant's occupation.

- employer: The tenant's employer.

- employerNumber: The tenant's employer's phone number.

- emergencyContactName: The tenant's emergency contact name.

- emergencyContactNumber: The tenant's emergency contact phone number.

- emergencyContactAddress: The tenant's emergency contact address.

- ifEvicted: Whether the tenant has been evicted before.

- evictedReason: The reason for eviction if applicable.

- ifConvicted: Whether the tenant has been convicted of a crime.

- convictedReason: The reason for conviction if applicable.

- smoke: Whether the tenant smokes.

- pets: Whether the tenant has pets.

- ifVehicle: Whether the tenant has a vehicle.

- vehicleReg: The vehicle registration number if applicable.

- haveChildren: Whether the tenant has children.

- children: The number of children the tenant has if applicable.

- refusedRent: Whether the tenant has refused rent before.

- refusedRentReason: The reason for refusing rent if applicable.

- unstableIncome: Whether the tenant has an unstable income.

- incomeReason: The reason for unstable income if applicable.

Returns:

- error: An error object if the database connection is not initialized or if the landlord email is empty.
*/
func SaveTenantApplicationForm(
	fullName,
	dateOfBirth,
	passportNumber,
	phoneNumber,
	email,
	occupation,
	employer,
	employerNumber,
	emergencyContactName,
	emergencyContactNumber,
	emergencyContactAddress,
	ifEvicted,
	evictedReason,
	ifConvicted,
	convictedReason,
	smoke,
	pets,
	ifVehicle,
	vehicleReg,
	haveChildren,
	children,
	refusedRent,
	refusedRentReason,
	unstableIncome,
	incomeReason string,
) error {
	if db == nil {
		logs.Logs(logDbErr, "Database connection is not initialized")
		return errors.New("database connection is not initialized")
	}

	// get landlord email via environment variable
	if os.Getenv("LANDLORD_EMAIL") == "" {
		logs.Logs(logWarning, "Could not get landlord email from hosting platform. Loading from .env file...")
		err := env.LoadEnv("env/.env")
		if err != nil {
			logs.Logs(logErr, fmt.Sprintf("Could not load environment variables from .env file: %s", err.Error()))
			return err
		}
	}

	landlordEmail := os.Getenv("LANDLORD_EMAIL")
	if landlordEmail == "" {
		logs.Logs(logDbErr, "Landlord email is empty!")
		return errors.New("landlord email is empty")
	}

	// get landlord id
	landlordId, err := GetLandlordIdByEmail(landlordEmail)
	if err != nil {
		logs.Logs(logDbErr, fmt.Sprintf("Failed to get landlord ID: %s", err.Error()))
		return err
	}

	// hash identifiers
	hashFullName := utils.HashData(fullName)
	hashDob := utils.HashData(dateOfBirth)
	hashPassportNumber := utils.HashData(passportNumber)
	hashEmail := utils.HashData(email)

	// set pending status
	const status = "pending"

	// encrypt data
	encryptFullName, err := utils.Encrypt([]byte(fullName))
	if err != nil {
		logs.Logs(logDbErr, fmt.Sprintf("Failed to encrypt full name: %s", err.Error()))
		return err
	}

	encryptDob, err := utils.Encrypt([]byte(dateOfBirth))
	if err != nil {
		logs.Logs(logDbErr, fmt.Sprintf("Failed to encrypt date of birth: %s", err.Error()))
		return err
	}

	encryptPassportNumber, err := utils.Encrypt([]byte(passportNumber))
	if err != nil {
		logs.Logs(logDbErr, fmt.Sprintf("Failed to encrypt passport number: %s", err.Error()))
		return err
	}

	encryptPhoneNumber, err := utils.Encrypt([]byte(phoneNumber))
	if err != nil {
		logs.Logs(logDbErr, fmt.Sprintf("Failed to encrypt phone number: %s", err.Error()))
		return err
	}

	encryptEmail, err := utils.Encrypt([]byte(email))
	if err != nil {
		logs.Logs(logDbErr, fmt.Sprintf("Failed to encrypt email: %s", err.Error()))
		return err
	}

	encryptOccupation, err := utils.Encrypt([]byte(occupation))
	if err != nil {
		logs.Logs(logDbErr, fmt.Sprintf("Failed to encrypt occupation: %s", err.Error()))
		return err
	}

	encryptEmployer, err := utils.Encrypt([]byte(employer))
	if err != nil {
		logs.Logs(logDbErr, fmt.Sprintf("Failed to encrypt employer: %s", err.Error()))
		return err
	}

	encryptEmployerNumber, err := utils.Encrypt([]byte(employerNumber))
	if err != nil {
		logs.Logs(logDbErr, fmt.Sprintf("Failed to encrypt employer number: %s", err.Error()))
		return err
	}

	encryptEmergencyContact, err := utils.Encrypt([]byte(emergencyContactName))
	if err != nil {
		logs.Logs(logDbErr, fmt.Sprintf("Failed to encrypt emergency contact name: %s", err.Error()))
		return err
	}

	encryptEmergencyNumber, err := utils.Encrypt([]byte(emergencyContactNumber))
	if err != nil {
		logs.Logs(logDbErr, fmt.Sprintf("Failed to encrypt emergency contact number: %s", err.Error()))
		return err
	}

	encryptEmergencyAddress, err := utils.Encrypt([]byte(emergencyContactAddress))
	if err != nil {
		logs.Logs(logDbErr, fmt.Sprintf("Failed to encrypt emergency contact address: %s", err.Error()))
		return err
	}

	encryptIfEvicted, err := utils.Encrypt([]byte(ifEvicted))
	if err != nil {
		logs.Logs(logDbErr, fmt.Sprintf("Failed to encrypt if evicted: %s", err.Error()))
		return err
	}

	encryptEvictedReason, err := utils.Encrypt([]byte(evictedReason))
	if err != nil {
		logs.Logs(logDbErr, fmt.Sprintf("Failed to encrypt evicted reason: %s", err.Error()))
		return err
	}

	encryptIfConvicted, err := utils.Encrypt([]byte(ifConvicted))
	if err != nil {
		logs.Logs(logDbErr, fmt.Sprintf("Failed to encrypt if convicted: %s", err.Error()))
		return err
	}

	encryptConvictedReason, err := utils.Encrypt([]byte(convictedReason))
	if err != nil {
		logs.Logs(logDbErr, fmt.Sprintf("Failed to encrypt convicted reason: %s", err.Error()))
		return err
	}

	encryptSmoke, err := utils.Encrypt([]byte(smoke))
	if err != nil {
		logs.Logs(logDbErr, fmt.Sprintf("Failed to encrypt smoke: %s", err.Error()))
		return err
	}

	encryptPets, err := utils.Encrypt([]byte(pets))
	if err != nil {
		logs.Logs(logDbErr, fmt.Sprintf("Failed to encrypt pets: %s", err.Error()))
		return err
	}

	encryptIfVechicle, err := utils.Encrypt([]byte(ifVehicle))
	if err != nil {
		logs.Logs(logDbErr, fmt.Sprintf("Failed to encrypt if vehicle: %s", err.Error()))
		return err
	}

	encryptVehicleReg, err := utils.Encrypt([]byte(vehicleReg))
	if err != nil {
		logs.Logs(logDbErr, fmt.Sprintf("Failed to encrypt vehicle registration: %s", err.Error()))
		return err
	}

	encryptHaveChildren, err := utils.Encrypt([]byte(haveChildren))
	if err != nil {
		logs.Logs(logDbErr, fmt.Sprintf("Failed to encrypt have children: %s", err.Error()))
		return err
	}

	encryptChildren, err := utils.Encrypt([]byte(children))
	if err != nil {
		logs.Logs(logDbErr, fmt.Sprintf("Failed to encrypt children: %s", err.Error()))
		return err
	}

	encryptRefusedRent, err := utils.Encrypt([]byte(refusedRent))
	if err != nil {
		logs.Logs(logDbErr, fmt.Sprintf("Failed to encrypt refused rent: %s", err.Error()))
		return err
	}

	encryptRefusedRentReason, err := utils.Encrypt([]byte(refusedRentReason))
	if err != nil {
		logs.Logs(logDbErr, fmt.Sprintf("Failed to encrypt refused rent reason: %s", err.Error()))
		return err
	}

	encryptUnstableIncome, err := utils.Encrypt([]byte(unstableIncome))
	if err != nil {
		logs.Logs(logDbErr, fmt.Sprintf("Failed to encrypt unstable income: %s", err.Error()))
		return err
	}

	encryptIncomeReason, err := utils.Encrypt([]byte(incomeReason))
	if err != nil {
		logs.Logs(logDbErr, fmt.Sprintf("Failed to encrypt income reason: %s", err.Error()))
		return err
	}

	query := `
	INSERT INTO lhp_tenant_application (
		landlord_id,
		hash_full_name,
		hash_dob,
		hash_passport_number,
		hash_email,
		status,
		encrypt_full_name,
		encrypt_dob,
		encrypt_passport_number,
		encrypt_phone_number,
		encrypt_email,
		encrypt_occupation,
		encrypt_employer,
		encrypt_employer_number,
		encrypt_emergency_contact,
		encrypt_emergency_number,
		encrypt_emergency_address,
		encrypt_if_evicted,
		encrypt_evicted_reason,
		encrypt_if_convicted,
		encrypt_convicted_reason,
		encrypt_smoke,
		encrypt_pets,
		encrypt_if_vehicle,
		encrypt_vehicle_reg,
		encrypt_have_children,
		encrypt_children,
		encrypt_refused_rent,
		encrypt_refused_rent_reason,
		encrypt_unstable_income,
		encrypt_income_reason,
		created_at
	)
	VALUES (
		$1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17, $18, $19, $20, $21, $22, $23, $24, $25, $26, $27, $28, $29, $30, $31, NOW() 
	);
	`

	// execute query
	_, err = db.Exec(
		query,
		landlordId,
		hashFullName,
		hashDob,
		hashPassportNumber,
		hashEmail,
		status,
		encryptFullName,
		encryptDob,
		encryptPassportNumber,
		encryptPhoneNumber,
		encryptEmail,
		encryptOccupation,
		encryptEmployer,
		encryptEmployerNumber,
		encryptEmergencyContact,
		encryptEmergencyNumber,
		encryptEmergencyAddress,
		encryptIfEvicted,
		encryptEvictedReason,
		encryptIfConvicted,
		encryptConvictedReason,
		encryptSmoke,
		encryptPets,
		encryptIfVechicle,
		encryptVehicleReg,
		encryptHaveChildren,
		encryptChildren,
		encryptRefusedRent,
		encryptRefusedRentReason,
		encryptUnstableIncome,
		encryptIncomeReason,
	)
	if err != nil {
		logs.Logs(logDbErr, fmt.Sprintf("Failed to store tenant application to database: %s", err.Error()))
		return err
	}

	logs.Logs(logDb, "Tenant application stored to database successfully")
	return nil
}

/*
GetAllTenantApplications retrieves all tenant applications associated with the current landlord.

This function checks if the database connection is initialized and retrieves the landlord's email from environment variables.
If the email is not found, it attempts to load it from a .env file. It then retrieves the landlord ID using the landlord email.
The function queries the database for all tenant applications related to this landlord ID, ordered by creation date in descending order.

Returns:

- []GetLandlordApplications: A slice containing the tenant applications.

- error: An error object if the tenant applications cannot be retrieved.
*/
func GetAllTenantApplications() ([]GetLandlordApplications, error) {
	if db == nil {
		logs.Logs(logDbErr, "Database connection is not initialized")
		return nil, errors.New("database connection is not initialized")
	}

	// get landlord email via environment variable
	if os.Getenv("LANDLORD_EMAIL") == "" {
		logs.Logs(logWarning, "Could not get landlord email from hosting platform. Loading from .env file...")
		err := env.LoadEnv("env/.env")
		if err != nil {
			logs.Logs(logErr, fmt.Sprintf("Could not load environment variables from .env file: %s", err.Error()))
			return nil, err
		}
	}

	landlordEmail := os.Getenv("LANDLORD_EMAIL")
	if landlordEmail == "" {
		logs.Logs(logDbErr, "Landlord email is empty!")
		return nil, errors.New("landlord email is empty")
	}

	// get landlord id
	landlordId, err := GetLandlordIdByEmail(landlordEmail)
	if err != nil {
		logs.Logs(logDbErr, fmt.Sprintf("Failed to get landlord ID: %s", err.Error()))
		return nil, err
	}

	query := `
	SELECT 
		id,
		status,
		encrypt_full_name,
		encrypt_dob,
		encrypt_passport_number,
		encrypt_phone_number,
		encrypt_email,
		encrypt_occupation,
		encrypt_employer,
		encrypt_employer_number,
		encrypt_emergency_contact,
		encrypt_emergency_number,
		encrypt_emergency_address,
		encrypt_if_evicted,
		encrypt_evicted_reason,
		encrypt_if_convicted,
		encrypt_convicted_reason,
		encrypt_smoke,
		encrypt_pets,
		encrypt_if_vehicle,
		encrypt_vehicle_reg,
		encrypt_have_children,
		encrypt_children,
		encrypt_refused_rent,
		encrypt_refused_rent_reason,
		encrypt_unstable_income,
		encrypt_income_reason,
		created_at
	FROM lhp_tenant_application
	WHERE landlord_id = $1
	ORDER BY created_at DESC;
	`
	rows, err := db.Query(query, landlordId)
	if err != nil {
		logs.Logs(logDbErr, fmt.Sprintf("Failed to get tenant applications: %s", err.Error()))
		return nil, err
	}
	defer rows.Close()

	var applicationsList []GetLandlordApplications

	for rows.Next() {
		var tenant GetLandlordApplications
		err := rows.Scan(
			&tenant.ID,
			&tenant.Status,
			&tenant.FullName,
			&tenant.Dob,
			&tenant.PassportNumber,
			&tenant.PhoneNumber,
			&tenant.Email,
			&tenant.Occupation,
			&tenant.Employer,
			&tenant.EmployerNumber,
			&tenant.EmergencyContact,
			&tenant.EmergencyContactNumber,
			&tenant.EmergencyContactAddress,
			&tenant.Evicted,
			&tenant.EvictedReason,
			&tenant.Convicted,
			&tenant.ConvictedReason,
			&tenant.Smoke,
			&tenant.Pets,
			&tenant.Vehicle,
			&tenant.VehicleReg,
			&tenant.HaveChildren,
			&tenant.Children,
			&tenant.RefusedRent,
			&tenant.RefusedReason,
			&tenant.UnstableIncome,
			&tenant.UnstableReason,
			&tenant.CreatedAt,
		)
		if err != nil {
			logs.Logs(logDbErr, fmt.Sprintf("Failed to scan tenant application: %s", err.Error()))
			return nil, err
		}
		applicationsList = append(applicationsList, tenant)
	}

	// check for errors
	err = rows.Err()
	if err != nil {
		logs.Logs(logDbErr, fmt.Sprintf("Failed to get tenant applications: %s", err.Error()))
		return nil, err
	}

	return applicationsList, nil
}

/*
UpdateTenantApplicationStatus updates the status of a tenant application in the database.

Arguments:

- id: The unique identifier of the tenant application to be updated.

- status: The new status to set for the tenant application.

Returns:

- error: An error if the status cannot be updated in the database.
*/
func UpdateTenantApplicationStatus(id string, status string) error {
	if db == nil {
		logs.Logs(logDbErr, "Database connection is not initialized")
		return errors.New("database connection is not initialized")
	}

	logs.Logs(logDb, "Updating tenant application status...")
	query := `
	UPDATE lhp_tenant_application
	SET status = $1
	WHERE id = $2;
	`
	_, err := db.Exec(query, status, id)
	if err != nil {
		logs.Logs(logDbErr, fmt.Sprintf("Failed to update tenant application status: %s", err.Error()))
		return err
	}
	logs.Logs(logDb, "Tenant application status updated successfully")
	return nil
}

/*
GetTenantEmailAndPassportNumberViaApplicationID retrieves the encrypted email and passport number for a tenant application.

Arguments:

- id: The unique identifier of the tenant application.

Returns:

- string: The encrypted email of the tenant.

- string: The encrypted passport number of the tenant.

- error: An error object if there is an issue retrieving the data from the database.
*/
func GetTenantEmailAndPassportNumberViaApplicationID(id string) (string, string, error) {
	if db == nil {
		logs.Logs(logDbErr, "Database connection is not initialized")
		return "", "", errors.New("database connection is not initialized")
	}

	var email, passportNumber string
	query := `
	SELECT encrypt_email, encrypt_passport_number
	FROM lhp_tenant_application
	WHERE id = $1;
	`
	err := db.QueryRow(query, id).Scan(&email, &passportNumber)
	if err != nil {
		return "", "", err
	}
	return email, passportNumber, nil
}

func GetHashedEmailFromTenantSessionToken(sessionToken string) (string, error) {
	if db == nil {
		logs.Logs(logDbErr, "Database connection is not initialized")
		return "", errors.New("database connection is not initialized")
	}

	var hashEmail string
	query := `
	SELECT hash_email 
	FROM lhp_tenants 
	WHERE session_token=$1;
	`
	err := db.QueryRow(query, sessionToken).Scan(&hashEmail)

	if err == sql.ErrNoRows {
		logs.Logs(logDbErr, "User not found")
		return "", errors.New("user not found")
	}

	if err != nil {
		logs.Logs(logDbErr, fmt.Sprintf("Failed to get session token: %s", err.Error()))
		return "", err
	}

	return hashEmail, nil
}

func UpdateTenantPassword(hashEmail, newPasswordHash, newPassword string) error {
	if db == nil {
		logs.Logs(logDbErr, "Database connection is not initialized")
		return errors.New("database connection is not initialized")
	}

	ecnryptNewPassword, err := utils.Encrypt([]byte(newPassword))
	if err != nil {
		logs.Logs(logDbErr, fmt.Sprintf("Failed to encrypt password: %s", err.Error()))
		return err
	}

	query := `
	UPDATE lhp_tenants
	SET hash_password = $1, encrypt_password = $2
	WHERE hash_email = $3;
	`
	_, err = db.Exec(query, newPasswordHash, ecnryptNewPassword, hashEmail)
	if err != nil {
		logs.Logs(logDbErr, fmt.Sprintf("Failed to update password: %s", err.Error()))
		return err
	}

	logs.Logs(logDb, "Password updated successfully")
	return nil
}

func GetTenantInformationByHashEmail(hashEmail string) (GetTenantInformation, error) {
	if db == nil {
		logs.Logs(logDbErr, "Database connection is not initialized")
		return GetTenantInformation{}, errors.New("database connection is not initialized")
	}

	query := `
	SELECT
		encrypt_email,
		encrypt_room_type,
		encrypt_move_in_date,
		encrypt_rent_due,
		encrypt_monthly_rent,
		currency
	FROM lhp_tenants
	WHERE hash_email = $1;
	`
	row, err := db.Query(query, hashEmail)
	if err != nil {
		logs.Logs(logDbErr, fmt.Sprintf("Failed to get tenant information: %s", err.Error()))
		return GetTenantInformation{}, err
	}
	defer row.Close()

	if !row.Next() {
		return GetTenantInformation{}, sql.ErrNoRows
	}

	var tenantInformation GetTenantInformation
	err = row.Scan(
		&tenantInformation.Email,
		&tenantInformation.RoomType,
		&tenantInformation.MoveInDate,
		&tenantInformation.RentDueDate,
		&tenantInformation.MonthlyRent,
		&tenantInformation.Currency,
	)
	if err != nil {
		logs.Logs(logDbErr, fmt.Sprintf("Failed to scan tenant information: %s", err.Error()))
		return GetTenantInformation{}, err
	}

	err = row.Err()
	if err != nil {
		logs.Logs(logDbErr, fmt.Sprintf("Failed to get tenant information: %s", err.Error()))
		return GetTenantInformation{}, err
	}

	return tenantInformation, nil
}
