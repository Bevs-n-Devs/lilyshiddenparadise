package email

import (
	"fmt"
	"net/smtp"
	"os"

	"github.com/Bevs-n-Devs/lilyshiddenparadise/env"
	"github.com/Bevs-n-Devs/lilyshiddenparadise/logs"
)

/*
NotifyLandlordNewApplication sends an email notification to the landlord when a new tenant application is submitted.

The function first checks if the email credentials are available in the environment variables.
If not, it loads them from the .env file.
Then, it checks if the email credentials are empty and returns an error if they are.
If the primary and secondary email addresses are the same, it sets the secondary email address to empty.
Finally, it creates an email message and sends it to the landlord using the smtp.SendMail function.
If the email cannot be sent, it logs an error and returns the error.

Returns:

- error: An error if the email cannot be sent.
*/
func NotifyLandlordNewApplication() error {
	if os.Getenv("LHP_EMAIL") == "" || os.Getenv("LHP_EMAIL_PASSWORD") == "" || os.Getenv("NOTIFY_LANDLORD_EMAIL") == "" {
		logs.Logs(logWarn, "Could not get email credentials from hosting platform. Loading from .env file...")
		err := env.LoadEnv("env/.env")
		if err != nil {
			logs.Logs(logErr, fmt.Sprintf("Unable to load environment variables: %s", err.Error()))
		}
	}

	if smptUser == "" || smptPassword == "" || recipient == "" || ccEmail == "" {
		logs.Logs(logErr, "Email credentials are empty!")
		return fmt.Errorf("email credentials are empty")
	}

	if recipient == ccEmail {
		logs.Logs(logWarn, "Primary and secondary email addresses are the same, skipping CC")
		ccEmail = ""
	}

	// create email message
	subject := "New Tenant Application"
	body := "A new tenant application has been submitted. Please login to the landlord dashboard to view the application."
	auth := smtp.PlainAuth("", smptUser, smptPassword, smptHost)
	err := smtp.SendMail(smptHost+":"+smptPort, auth, smptUser, []string{recipient, ccEmail}, []byte("Subject: "+subject+"\n\n"+body))
	if err != nil {
		logs.Logs(logErr, fmt.Sprintf("Failed to send email: %s", err.Error()))
		return err
	}

	logs.Logs(logInfo, "Email sent successfully. Landlord notified of new tenant application.")
	return nil
}
