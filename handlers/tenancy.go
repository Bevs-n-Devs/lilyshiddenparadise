package handlers

import (
	"fmt"
	"net/http"

	"github.com/Bevs-n-Devs/lilyshiddenparadise/logs"
)

func Tenancy(w http.ResponseWriter, r *http.Request) {
	// collect any error messages
	evictedError := r.URL.Query().Get("evictedError")
	convictedError := r.URL.Query().Get("convictedError")
	vehicleError := r.URL.Query().Get("vehicleError")
	childrenError := r.URL.Query().Get("childrenError")
	refusedRentError := r.URL.Query().Get("refusedRentError")
	unstableIncomeError := r.URL.Query().Get("unstableIncomeError")

	data := FormError{
		EvictedError:        evictedError,
		ConvictedError:      convictedError,
		VehicleError:        vehicleError,
		ChildrenError:       childrenError,
		RefusedRentError:    refusedRentError,
		UnstableIncomeError: unstableIncomeError,
	}

	// pass error messages to HTML template
	err := Templates.ExecuteTemplate(w, "tenancy.html", data)
	if err != nil {
		logs.Logs(logErr, fmt.Sprintf("Unable to load tenancy page: %s", err.Error()))
		http.Error(w, fmt.Sprintf("Unable to load tenancy page: %s", err.Error()), http.StatusInternalServerError)
		return
	}
}
