package email

import "os"

const (
	logInfo  = 1
	logWarn  = 2
	logErr   = 3
	smptHost = "smtp.gmail.com"
	smptPort = "587"
)

var (
	smptUser     = os.Getenv("LHP_EMAIL")
	smptPassword = os.Getenv("LHP_EMAIL_PASSWORD")
	recipient    = os.Getenv("NOTIFY_LANDLORD_EMAIL") // 1st destination email
	ccEmail      = os.Getenv("LHP_EMAIL")             // 2nd destination email
)
