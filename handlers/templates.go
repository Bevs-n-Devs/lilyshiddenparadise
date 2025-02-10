package handlers

import (
	"fmt"
	"html/template"
	"os"

	"github.com/Bevs-n-Devs/lilyshiddenparadise/logs"
)

var Templates *template.Template

func InitTemplates() {
	var err error
	Templates, err = template.ParseGlob("./handlers/templates/*.html")
	if err != nil {
		logs.Logs(3, fmt.Sprintf("Error parsing templates: %s", err))
		os.Exit(1)
	}
}
