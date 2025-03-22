package main

import (
	"fmt"

	"github.com/Bevs-n-Devs/lilyshiddenparadise/db"
	"github.com/Bevs-n-Devs/lilyshiddenparadise/handlers"
	"github.com/Bevs-n-Devs/lilyshiddenparadise/logs"
)

const (
	logInfo  = 1
	logDbErr = 5
)

func main() {
	go logs.LogProcessor()
	logs.Logs(logInfo, "Welcome to Lily's Hidden Paradise, a web app to manage tenants and landlords.")

	err := db.ConnectDB()
	if err != nil {
		logs.Logs(logDbErr, fmt.Sprintf("Error connecting to database: %s", err.Error()))
		return
	}

	go handlers.StartHTTPServer()

	select {} // keeps the program running
}
