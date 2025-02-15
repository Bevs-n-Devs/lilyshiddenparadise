package handlers

import (
	"fmt"
	"net/http"

	"github.com/Bevs-n-Devs/lilyshiddenparadise/logs"
)

const (
	httpServer = ":9001"
)

func StartHTTPServer() {
	logs.Logs(1, "Starting HTTP server...")

	InitTemplates()

	// Static file server for assets like CSS, JS, images
	var staticFiles = http.FileServer(http.Dir("./static"))
	http.Handle("/static/", http.StripPrefix("/static/", staticFiles))

	// define routes
	http.HandleFunc("/", Home)
	http.HandleFunc("/form", Tenancy)
	http.HandleFunc("/login", Login)
	http.HandleFunc("/submit", SubmitForm)
	http.HandleFunc("/submitLogin", SubmitLogin)

	logs.Logs(1, fmt.Sprintf("Server running on http://localhost%s", httpServer))
	err := http.ListenAndServe(httpServer, nil)
	if err != nil {
		logs.Logs(3, fmt.Sprintf("Error starting HTTP server: %s", err))
		return
	}
}
