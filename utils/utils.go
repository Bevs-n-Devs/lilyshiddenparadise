package utils

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
