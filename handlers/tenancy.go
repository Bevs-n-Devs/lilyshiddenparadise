package handlers

import (
	"fmt"
	"net/http"

	"github.com/Bevs-n-Devs/lilyshiddenparadise/logs"
)

func Tenancy(w http.ResponseWriter, r *http.Request) {
	evictedError := r.URL.Query().Get("evictedError")
	convictedError := r.URL.Query().Get("convictedError")
	vehicleError := r.URL.Query().Get("vehicleError")
	childrenError := r.URL.Query().Get("childrenError")
	refusedRentError := r.URL.Query().Get("refusedRentError")
	unstableIncomeError := r.URL.Query().Get("unstableIncomeError")

	data := struct {
		EvictedError        string
		ConvictedError      string
		VehicleError        string
		ChildrenError       string
		RefusedRentError    string
		UnstableIncomeError string
	}{
		EvictedError:        evictedError,
		ConvictedError:      convictedError,
		VehicleError:        vehicleError,
		ChildrenError:       childrenError,
		RefusedRentError:    refusedRentError,
		UnstableIncomeError: unstableIncomeError,
	}

	err := Templates.ExecuteTemplate(w, "tenancy.html", data)
	if err != nil {
		logs.Logs(3, fmt.Sprintf("Unable to load tenancy page: %s", err))
		http.Error(w, "Unable to load tenancy page: "+err.Error(), http.StatusInternalServerError)
		return
	}
}
