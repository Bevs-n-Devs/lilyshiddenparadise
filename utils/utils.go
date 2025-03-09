package utils

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"

	"github.com/Bevs-n-Devs/lilyshiddenparadise/logs"
	"golang.org/x/crypto/bcrypt"
)

const (
	logErr = 3
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
