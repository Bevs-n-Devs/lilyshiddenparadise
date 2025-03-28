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

	// Update smptUser, smptPassword, recipient, and ccEmail variables
	smptUser = os.Getenv("LHP_EMAIL")
	smptPassword = os.Getenv("LHP_EMAIL_PASSWORD")
	recipient = os.Getenv("NOTIFY_LANDLORD_EMAIL")
	ccEmail = os.Getenv("LHP_EMAIL")

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

/*
NotifyTenantApplicationProcessing sends an email to the tenant to notify them that their application is being processed.

Arguments:

- tenantEmail: The tenant's email address to send the email to.

Returns:

- error: An error if the email cannot be sent.
*/
func NotifyTenantApplicationProcessing(tenantEmail string) error {
	if os.Getenv("LHP_EMAIL") == "" || os.Getenv("LHP_EMAIL_PASSWORD") == "" || os.Getenv("NOTIFY_LANDLORD_EMAIL") == "" {
		logs.Logs(logWarn, "Could not get email credentials from hosting platform. Loading from .env file...")
		err := env.LoadEnv("env/.env")
		if err != nil {
			logs.Logs(logErr, fmt.Sprintf("Unable to load environment variables: %s", err.Error()))
		}
	}

	// Update smptUser, smptPassword, recipient, and ccEmail variables
	smptUser = os.Getenv("LHP_EMAIL")
	smptPassword = os.Getenv("LHP_EMAIL_PASSWORD")
	ccEmail = os.Getenv("LHP_EMAIL")

	// set recipient to tenant email
	recipient := tenantEmail

	if smptUser == "" || smptPassword == "" || recipient == "" || ccEmail == "" {
		logs.Logs(logErr, "Email credentials are empty!")
		return fmt.Errorf("email credentials are empty")
	}

	if recipient == ccEmail {
		logs.Logs(logWarn, "Primary and secondary email addresses are the same, skipping CC")
		ccEmail = ""
	}

	// create email message
	subject := "Tenant Application Processing"
	body := "Your tenant application is being processed. Please wait for further instructions."
	auth := smtp.PlainAuth("", smptUser, smptPassword, smptHost)
	err := smtp.SendMail(smptHost+":"+smptPort, auth, smptUser, []string{recipient, ccEmail}, []byte("Subject: "+subject+"\n\n"+body))
	if err != nil {
		logs.Logs(logErr, fmt.Sprintf("Failed to send email: %s", err.Error()))
		return err
	}

	logs.Logs(logInfo, "Email sent successfully. Tenant notified that application is being processed.")
	return nil
}

func NotifyTenantNewAccount(tenantUsername, tenantPassword, roomType, moveInDate, rentDue, monthlyRent, currency string) error {
	if os.Getenv("LHP_EMAIL") == "" || os.Getenv("LHP_EMAIL_PASSWORD") == "" || os.Getenv("NOTIFY_LANDLORD_EMAIL") == "" {
		logs.Logs(logWarn, "Could not get email credentials from hosting platform. Loading from .env file...")
		err := env.LoadEnv("env/.env")
		if err != nil {
			logs.Logs(logErr, fmt.Sprintf("Unable to load environment variables: %s", err.Error()))
		}
	}

	// Update smptUser, smptPassword, recipient, and ccEmail variables
	smptUser = os.Getenv("LHP_EMAIL")
	smptPassword = os.Getenv("LHP_EMAIL_PASSWORD")
	ccEmail = os.Getenv("LHP_EMAIL")

	// set recipient to tenant email
	recipient := tenantUsername

	if smptUser == "" || smptPassword == "" || recipient == "" || ccEmail == "" {
		logs.Logs(logErr, "Email credentials are empty!")
		return fmt.Errorf("email credentials are empty")
	}

	if recipient == ccEmail {
		logs.Logs(logWarn, "Primary and secondary email addresses are the same, skipping CC")
		ccEmail = ""
	}

	// create email message
	subject := "Tenant Application Approved"
	body := fmt.Sprintf(`
Congratulations! Your application has been approved!

Here is your account information:
	Username: %s
	Password: %s

YOUR ROOM DETAILS:

	Room Type: %s
	Move-in Date: %s
	Rent Due: %s (and same date every month thereafter)
	Monthly Rent: %s %s per month

Please login to the tenant dashboard to view your account details.
	`, tenantUsername, tenantPassword, roomType, moveInDate, rentDue, monthlyRent, currency)

	auth := smtp.PlainAuth("", smptUser, smptPassword, smptHost)
	err := smtp.SendMail(smptHost+":"+smptPort, auth, smptUser, []string{recipient, ccEmail}, []byte("Subject: "+subject+"\n\n"+body))
	if err != nil {
		logs.Logs(logErr, fmt.Sprintf("Failed to send email: %s", err.Error()))
		return err
	}

	logs.Logs(logInfo, "Email sent successfully. Tenant notified that application is being processed.")
	return nil
}

func NotifyLandlordNewAccount(tenantUsername, tenantPassword, roomType, moveInDate, rentDue, monthlyRent, currency string) error {
	if os.Getenv("LHP_EMAIL") == "" || os.Getenv("LHP_EMAIL_PASSWORD") == "" || os.Getenv("NOTIFY_LANDLORD_EMAIL") == "" {
		logs.Logs(logWarn, "Could not get email credentials from hosting platform. Loading from .env file...")
		err := env.LoadEnv("env/.env")
		if err != nil {
			logs.Logs(logErr, fmt.Sprintf("Unable to load environment variables: %s", err.Error()))
		}
	}

	// Update smptUser, smptPassword, recipient, and ccEmail variables
	smptUser = os.Getenv("LHP_EMAIL")
	smptPassword = os.Getenv("LHP_EMAIL_PASSWORD")
	recipient = os.Getenv("NOTIFY_LANDLORD_EMAIL")
	ccEmail = os.Getenv("LHP_EMAIL")

	if smptUser == "" || smptPassword == "" || recipient == "" || ccEmail == "" {
		logs.Logs(logErr, "Email credentials are empty!")
		return fmt.Errorf("email credentials are empty")
	}

	if recipient == ccEmail {
		logs.Logs(logWarn, "Primary and secondary email addresses are the same, skipping CC")
		ccEmail = ""
	}

	// create email message
	subject := "New Tenant Application Approved"
	body := fmt.Sprintf(`
Your tenant's application has been approved!

	Here is your tenant's account information:
	Username: %s
	Password: %s

TENANT ROOM DETAILS:

	Room Type: %s
	Move-in Date: %s
	Rent Due: %s (and same date every month thereafter)
	Monthly Rent: %s %s per month

Please login to the landlord dashboard to view more details.
	`, tenantUsername, tenantPassword, roomType, moveInDate, rentDue, monthlyRent, currency)

	auth := smtp.PlainAuth("", smptUser, smptPassword, smptHost)
	err := smtp.SendMail(smptHost+":"+smptPort, auth, smptUser, []string{recipient, ccEmail}, []byte("Subject: "+subject+"\n\n"+body))
	if err != nil {
		logs.Logs(logErr, fmt.Sprintf("Failed to send email: %s", err.Error()))
		return err
	}

	logs.Logs(logInfo, "Email sent successfully. Landlord notified that new tenant account has been created.")
	return nil
}

func NotifyLandlordNewMessageFromTenant(tenantName, landlordEmail, messageFromTenant string) error {
	smptUser = os.Getenv("LHP_EMAIL")
	smptPassword = os.Getenv("LHP_EMAIL_PASSWORD")
	recipient = os.Getenv("NOTIFY_LANDLORD_EMAIL")
	ccEmail = os.Getenv("LHP_EMAIL")

	if smptUser == "" || smptPassword == "" || recipient == "" || ccEmail == "" {
		logs.Logs(logErr, "Email credentials are empty!")
		return fmt.Errorf("email credentials are empty")
	}

	if recipient == ccEmail {
		logs.Logs(logWarn, "Primary and secondary email addresses are the same, skipping CC")
		ccEmail = ""
	}

	subject := "New Message from " + tenantName
	body := fmt.Sprintf(`
You've received a new message from %s.

MESAGE CONTENT:
%s

Log in to your landlord dashboard to respond.
	`, tenantName, messageFromTenant)

	auth := smtp.PlainAuth("", smptUser, smptPassword, smptHost)
	err := smtp.SendMail(smptHost+":"+smptPort, auth, smptUser, []string{recipient, ccEmail}, []byte("Subject: "+subject+"\n\n"+body))
	if err != nil {
		logs.Logs(logErr, fmt.Sprintf("Failed to send email: %s", err.Error()))
		return err
	}

	logs.Logs(logInfo, "Email sent successfully. Landlord notified of new message from tenant.")
	return nil

}
