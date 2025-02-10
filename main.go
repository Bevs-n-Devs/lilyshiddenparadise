package main

import (
	"github.com/Bevs-n-Devs/lilyshiddenparadise/handlers"
	"github.com/Bevs-n-Devs/lilyshiddenparadise/logs"
)

func main() {
	go logs.LogProcessor()
	go handlers.StartHTTPServer()

	select {}
}
