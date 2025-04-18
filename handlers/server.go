package handlers

import (
	"fmt"
	"net/http"
	"os"

	"github.com/Bevs-n-Devs/lilyshiddenparadise/logs"
	"github.com/Bevs-n-Devs/lilyshiddenparadise/utils"
)

func StartHTTPServer() {
	logs.Logs(logInfo, "Starting HTTP server...")

	InitTemplates()

	// initialise encryption functions
	err := utils.InitEncryption()
	if err != nil {
		logs.Logs(logErr, fmt.Sprintf("Error initialising encryption functions: %s", err.Error()))
	}

	// Static file server for assets like CSS, JS, images
	var staticFiles = http.FileServer(http.Dir("./static"))
	http.Handle("/static/", http.StripPrefix("/static/", staticFiles))

	// define routes
	http.HandleFunc("/", Home)
	http.HandleFunc("/tenancy-form", TenancyForm)
	http.HandleFunc("/tenancy-form/submit", SubmitTenantForm)
	http.HandleFunc("/new/landlord", NewLandlord)
	http.HandleFunc("/new/landlord/submit", SubmitNewLandlord)
	http.HandleFunc("/login/landlord", LoginLandlord)
	http.HandleFunc("/login/landlord/submit", SubmitLoginLandlord)
	http.HandleFunc("/login/tenant", LoginTenant)
	http.HandleFunc("/login/tenant/submit", SubmitLoginTenant)

	// protected landlord routes
	http.HandleFunc("/logout-landlord", LogoutLandlord)
	http.HandleFunc("/landlord/dashboard", LandlordDashboard)
	http.HandleFunc("/landlord/dashboard/tenants", LandlordDashboardTenants)
	http.HandleFunc("/landlord/dashboard/tenant-applications", LandlordTenantApplications)
	http.HandleFunc("/landlord/dashboard/manage-applications", LandlordManageApplications)
	http.HandleFunc("/landlord/dashboard/new-tenant", LandlordNewTenant)
	http.HandleFunc("/landlord/dashboard/new-tenant/submit", LandlordSubmitNewTenant)
	http.HandleFunc("/landlord/dashboard/messages", LandlordMessages)
	http.HandleFunc("/landlord/dashboard/messages/tenant/", LandlordTenantMessages)
	http.HandleFunc("/landlord/send-message/", SendMessageToTenant)

	// protected tenant routes
	http.HandleFunc("/tenant/dashboard", TenantDashboard)
	http.HandleFunc("/logout-tenant", LogoutTenant)
	http.HandleFunc("/tenant/dashboard/account", TenantAccount)
	http.HandleFunc("/tenant/update-password", UpdateTenantPassword)
	http.HandleFunc("/tenant/dashboard/messages", TenantMessages)
	http.HandleFunc("/tenant/send-message", SendMessageToLandlord)

	// initialise port for application
	httpPort := os.Getenv("PORT") // attempt to get port from hosting platform

	// start server on local machine if hosting platform port is not set
	if httpPort == "" {
		logs.Logs(logWarn, fmt.Sprintf("Could not get PORT from hosting platform. Defaulting to http://localhost:%s...", localPort))
		httpPort = localPort
		err := http.ListenAndServe("localhost:"+localPort, nil)
		if err != nil {
			logs.Logs(logErr, fmt.Sprintf("Failed to start HTTP server: %s", err.Error()))
		}
	}

	// start server on hosting platform port
	logs.Logs(logInfo, fmt.Sprintf("HTTP server running on http://localhost:%s", httpPort))
	err = http.ListenAndServe(":"+httpPort, nil)
	if err != nil {
		logs.Logs(logErr, fmt.Sprintf("Error starting HTTP server: %s", err.Error()))
	}
}
