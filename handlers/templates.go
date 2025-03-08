package handlers

import (
	"fmt"
	"html/template"
	"os"

	"github.com/Bevs-n-Devs/lilyshiddenparadise/logs"
)

func InitTemplates() {
	var err error
	Templates, err = template.ParseGlob("./templates/*.html")
	if err != nil {
		logs.Logs(logErr, fmt.Sprintf("Error parsing templates: %s", err.Error()))
		os.Exit(1)
	}
}
