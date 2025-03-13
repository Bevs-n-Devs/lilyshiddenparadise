package handlers

import (
	"fmt"
	"net/http"

	"github.com/Bevs-n-Devs/lilyshiddenparadise/logs"
)

func TenancyForm(w http.ResponseWriter, r *http.Request) {
	// collect any error messages
	ageError := r.URL.Query().Get("ageError")
	evictedError := r.URL.Query().Get("evictedError")
	convictedError := r.URL.Query().Get("convictedError")
	vehicleError := r.URL.Query().Get("vehicleError")
	childrenError := r.URL.Query().Get("childrenError")
	refusedRentError := r.URL.Query().Get("refusedRentError")
	unstableIncomeError := r.URL.Query().Get("unstableIncomeError")
	dbError := r.URL.Query().Get("dbError")
	emailError := r.URL.Query().Get("emailError")

	data := ErrorMessages{
		AgeError:            ageError,
		EvictedError:        evictedError,
		ConvictedError:      convictedError,
		VehicleError:        vehicleError,
		ChildrenError:       childrenError,
		RefusedRentError:    refusedRentError,
		UnstableIncomeError: unstableIncomeError,
		DatabaseError:       dbError,
		EmailError:          emailError,
	}

	// pass error messages to HTML template
	err := Templates.ExecuteTemplate(w, "tenancyForm.html", data)
	if err != nil {
		logs.Logs(logErr, fmt.Sprintf("Unable to load tenancy page: %s", err.Error()))
		http.Error(w, fmt.Sprintf("Unable to load tenancy page: %s", err.Error()), http.StatusInternalServerError)
		return
	}
}
