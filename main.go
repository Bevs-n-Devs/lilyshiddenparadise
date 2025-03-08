package main

import (
	"fmt"

	"github.com/Bevs-n-Devs/lilyshiddenparadise/db"
	"github.com/Bevs-n-Devs/lilyshiddenparadise/handlers"
	"github.com/Bevs-n-Devs/lilyshiddenparadise/logs"
)

const (
	logErr = 3
)

func main() {
	go logs.LogProcessor()

	err := db.ConnectDB()
	if err != nil {
		logs.Logs(logErr, fmt.Sprintf("Error connecting to database: %s", err.Error()))
		return
	}

	go handlers.StartHTTPServer()

	select {} // keeps the program running
}
