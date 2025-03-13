package utils

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/Bevs-n-Devs/lilyshiddenparadise/logs"
	"golang.org/x/crypto/bcrypt"
)

/*
HashedPassword generates a hashed password with a cost of 2^10.

It takes a string representation of a password and returns a byte slice
representing the hashed password. The second return value is an error type
that is non-nil if an error occurs while hashing the password.

Arguments:

- password: A string representation of the password to hash.

Returns:

- string: A string representation of the hashed password.

- error: An error type that is non-nil if an error occurs while hashing the password.
*/
func HashedPassword(password string) (string, error) {
	// byte representation of the password string, password hashed 2^10 times
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 10)
	return string(bytes), err
}

/*
Checks a hashed password against a given password.
Returns true if the given password matches the hashed password, false if not.

Arguments:

- password: A string representation of the password to check.

- hash: A string representation of the hashed password to check against.

Returns:

- bool: True if the given password matches the hashed password, false if not.
*/
func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

/*
HashData takes a string and returns a SHA-256 hash of the string.

The resulting hash is a fixed-size 256-bit string, represented as
a 64-character hexadecimal string.
*/
func HashData(data string) string {
	hash := sha256.Sum256([]byte(data))
	return hex.EncodeToString(hash[:])
}

/*
VerifyHash takes an input string and a stored hash string and
returns true if the two hashes match, or false if they do not.

This is a simple equality check and does not provide any
additional security features. It is the responsibility of the
caller to ensure that the input string and stored hash are valid
and have been secured appropriately.
*/
func VerifyHash(input string, storedHash string) bool {
	return input == storedHash // Compare with the stored hash
}

/*
Encrypts any identifiable data the user enters.
Will need the MASTER_KEY from envrionment variable to work.

We need to convert the data into bytes to encrypt it.

Return a list of bytes or an error.
*/
func Encrypt(data []byte) ([]byte, error) {
	// create a new AES cipher block using the master key
	block, err := aes.NewCipher(masterKey)
	if err != nil {
		logs.Logs(logErr, fmt.Sprintf("Error creating AES cipher block: %s", err.Error()))
		return nil, err // Return error if key is invalid
	}

	// Create a GCM (Galois Counter Mode) cipher from the AES block
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		logs.Logs(logErr, fmt.Sprintf("Error creating GCM cipher: %s", err.Error()))
		return nil, err // Return error if GCM initialization fails
	}

	// Generate a nonce (unique number used only once) of required size
	nonce := make([]byte, gcm.NonceSize())   // GCM nonce should be unique per encryption
	_, err = io.ReadFull(rand.Reader, nonce) // Fill nonce with random bytes
	if err != nil {
		logs.Logs(logErr, fmt.Sprintf("Error generating nonce: %s", err.Error()))
		return nil, err // Return error if random generation fails
	}

	// Encrypt the data using AES-GCM
	// Seal appends encrypted data to nonce (authentication tag included)
	ciphertext := gcm.Seal(nil, nonce, data, nil)

	// Return the concatenated nonce + ciphertext
	logs.Logs(logInfo, "Data encrypted successfully")
	return append(nonce, ciphertext...), nil
}

/*
Decrypt decrypts the given encrypted data using AES-GCM with the master key.
It expects the data to contain the nonce followed by the ciphertext.

Parameters:

	data ([]byte): The encrypted data containing the nonce and ciphertext.

Returns:

	([]byte): The decrypted plaintext if successful.
	(error): An error if the decryption process fails, such as an invalid key or corrupted data.
*/
func Decrypt(data []byte) ([]byte, error) {
	// Create a new AES cipher block using the same master key
	block, err := aes.NewCipher(masterKey)
	if err != nil {
		logs.Logs(logErr, fmt.Sprintf("Error creating AES cipher block: %s", err.Error()))
		return nil, err // Return error if key is invalid
	}

	// Create a GCM cipher from the AES block
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		logs.Logs(logErr, fmt.Sprintf("Error creating GCM cipher: %s", err.Error()))
		return nil, err // Return error if GCM initialization fails
	}

	// Extract the nonce from the start of the encrypted data
	nonceSize := gcm.NonceSize()
	nonce, ciphertext := data[:nonceSize], data[nonceSize:]

	// Decrypt the ciphertext using AES-GCM
	plaintext, err := gcm.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		logs.Logs(logErr, fmt.Sprintf("Error decrypting data: %s", err.Error()))
		return nil, err // Return error if decryption fails
	}

	// Return the decrypted plaintext
	logs.Logs(logInfo, "Data decrypted successfully")
	return plaintext, nil
}

/*
GenerateToken generates a cryptographically secure random token of a given length.
It takes a single int argument, the length of the token to generate.
It returns a string representation of the generated token and an error. The
error is non-nil if an error occurs while generating the token.

Arguments:

- length: An int representing the length of the token to generate.

Returns:

- string: A string representation of the generated token.

- error: An error type that is non-nil if an error occurs while generating the token.
*/
func GenerateToken(length int) (string, error) {
	bytes := make([]byte, length)
	if _, err := rand.Read(bytes); err != nil {
		logs.Logs(logErr, fmt.Sprintf("Error generating token: %s", err.Error()))
		return "", err
	}
	return base64.URLEncoding.EncodeToString(bytes), nil
}

/*
Checks if a given password matches a confirmation password.
Returns true if the passwords match, false if not.

Arguments:

- password: A string representation of the password to check.

- confirmPassword: A string representation of the confirmation password to check against.

Returns:

- bool: True if the passwords match, false if not.
*/
func ValidateNewLandlordPassword(password, confirmPassword string) bool {
	return password == confirmPassword
}

/*
Checks if ifEvicted is yes and evictedReason is empty;
returns false if invalid.

Returns true if ifEvicted is yes and evictedReason is not empty
*/
func CheckIfEvicted(ifEvicted, evictedReason string) bool {
	if ifEvicted == "yes" && evictedReason == "" {
		return false
	}
	return true
}

/*
Checks if ifConvicted is yes and convictedReason is empty;
returns false if invalid.

Returns true if ifConvicted is yes and convictedReason is not empty
*/
func CheckIfConvicted(ifConvicted, convictedReason string) bool {
	if ifConvicted == "yes" && convictedReason == "" {
		return false
	}
	return true
}

/*
Checks if ifVehicle is yes and vehicleReg is empty;
returns false if invalid.

Returns true if ifVehicle is yes and vehicleReg is not empty
*/
func CheckIfVehicle(ifVehicle, vehicleReg string) bool {
	if ifVehicle == "yes" && vehicleReg == "" {
		return false
	}
	return true
}

/*
Checks if haveChildren is yes and children is empty;
returns false if invalid.

Returns true if haveChildren is yes and children is not empty
*/
func CheckIfHaveChildren(haveChildren, children string) bool {
	if haveChildren == "yes" && children == "" {
		return false
	}
	return true
}

/*
Checks if refusedRent is yes and refusedRentReason is empty;
returns false if invalid.

Returns true if refusedRent is yes and refusedRentReason is not empty
*/
func CheckIfRefusedRent(refusedRent, refusedRentReason string) bool {
	if refusedRent == "yes" && refusedRentReason == "" {
		return false
	}
	return true
}

/*
Checks if unstableIncome is yes and incomeReason is empty;
returns false if invalid.

Returns true if stableIncome is no and incomeReason is not empty
*/
func CheckIfStableIncome(unstableIncome, incomeReason string) bool {
	if unstableIncome == "yes" && incomeReason == "" {
		return false
	}
	return true
}

/*
Checks if a session token exists in the request and returns the session token as a Cookie
if it does. If the session token does not exist, an error is returned.

Returns:

- *http.Cookie: The session token as a Cookie if it exists.

- error: An error if the session token does not exist.
*/
func CheckSessionToken(r *http.Request) (*http.Cookie, error) {
	sessionToken, err := r.Cookie("session_token")
	if err != nil || sessionToken.Value == "" {
		return nil, fmt.Errorf("user not authenticated! failed to get session token: %s", err.Error())
	}
	return sessionToken, nil
}

/*
Checks if a CSRF token exists in the request and returns the CSRF token as a Cookie
if it does. If the CSRF token does not exist, an error is returned.

Returns:

- *http.Cookie: The CSRF token as a Cookie if it exists.

- error: An error if the CSRF token does not exist.
*/
func CheckCSRFToken(r *http.Request) (*http.Cookie, error) {
	csrfToken, err := r.Cookie("csrf_token")
	if err != nil || csrfToken.Value == "" {
		return nil, fmt.Errorf("user not authenticated! failed to get csrf token: %s", err.Error())
	}
	return csrfToken, nil
}

/*
Checks if the age of the user is 18 years or older.

Arguments:

- dateOfBirth: A string representation of the date of birth of the user.

Returns:

- bool: True if the user is 18 years or older, false if not.
*/
func ValidateAge(dateOfBirth string) bool {
	layout := "2006-01-02"
	dob, err := time.Parse(layout, dateOfBirth)
	if err != nil {
		return false
	}
	return dob.Before(time.Now().Local().AddDate(-18, 0, 0))
}

/*
ValidateManageTenantApplication validates the manage tenant application form data.

Arguments:

- applicationResult: The result of the tenant application. Accepted or Rejected.

- roomType: The type of room to be allocated to the tenant.

- moveInDate: The expected move in date of the tenant.

- rentDue: The rent due date.

- monthlyRent: The monthly rent for the room.

- currency: The currency of the monthly rent.

Returns:

- error: An error if any of the required fields are empty when the application result is "accepted".
*/
func ValidateManageTenantApplication(applicationResult, roomType, moveInDate, rentDue, monthlyRent, currency string) error {
	if applicationResult == "accepted" && roomType == "" {
		return fmt.Errorf("room type is required")
	}
	if applicationResult == "accepted" && moveInDate == "" {
		return fmt.Errorf("move in date is required")
	}
	if applicationResult == "accepted" && rentDue == "" {
		return fmt.Errorf("rent due is required")
	}
	if applicationResult == "accepted" && monthlyRent == "" {
		return fmt.Errorf("monthly rent is required")
	}
	if applicationResult == "accepted" && currency == "" {
		return fmt.Errorf("currency is required")
	}
	return nil
}
