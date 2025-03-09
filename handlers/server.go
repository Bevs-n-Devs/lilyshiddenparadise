package handlers

import (
	"fmt"
	"net/http"
	"os"

	"github.com/Bevs-n-Devs/lilyshiddenparadise/logs"
)

func StartHTTPServer() {
	logs.Logs(logInfo, "Starting HTTP server...")

	InitTemplates()

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
	// http.HandleFunc("/login/tenant/submit", SubmitLoginTenant)

	// protected routes
	http.HandleFunc("/landlord/dashboard", LandlordDashboard)
	// http.HandleFunc("/tenant/dashboard", TenantDashboard)

	// initialise port for application
	httpPort := os.Getenv("PORT") // attempt to get port from hosting platform

	// start server on local machine if hosting platform port is not set
	if httpPort == "" {
		logs.Logs(logWarn, fmt.Sprintf("Could not get PORT from hosting platform. Defaulting to http://localhost%s...", localPort))
		httpPort = localPort
		err := http.ListenAndServe(localPort, nil)
		if err != nil {
			logs.Logs(logErr, fmt.Sprintf("Failed to start HTTP server: %s", err.Error()))
		}
	}

	// start server on hosting platform port
	logs.Logs(logInfo, fmt.Sprintf("HTTP server running on http://localhost%s", httpPort))
	err := http.ListenAndServe(httpPort, nil)
	if err != nil {
		logs.Logs(logErr, fmt.Sprintf("Error starting HTTP server: %s", err.Error()))
	}
}
