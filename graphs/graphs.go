package graphs

import (
	"net/http"

	"github.com/gorilla/mux"

	"github.com/guptarohit/asciigraph"
	"github.com/willfantom/lu-covid-api/api"
	"github.com/willfantom/lu-covid-api/db"

	log "github.com/sirupsen/logrus"
)

// THIS WHOLE FILE NEEDS REFACTORING... I WILL DO IT AT SOMEPOINT
// Much copy-pasta

func API(router *mux.Router) {
	router.HandleFunc("/total", totalCases)
	router.HandleFunc("/students", studentCases)
	router.HandleFunc("/staff", staffCases)
}

func totalCases(w http.ResponseWriter, r *http.Request) {
	log.Debug("📈 getting ascii graph for total cases")
	earliestDate, errA := api.GetEarliestDate()
	recentDate, errB := api.GetRecentDate()
	if errA != nil || errB != nil {
		http.Error(w, "server encountered an issue", 500)
		return
	}
	rates, err := db.FetchInRange(earliestDate, recentDate)
	if err != nil {
		http.Error(w, "server encountered an issue", 500)
		return
	}
	w.Header().Set("Content-Type", "application/text")
	data := make([]float64, len(*rates), len(*rates))
	for idx, rate := range *rates {
		data[idx] = float64(rate.Campus + rate.City + rate.Staff)
		if idx > 0 {
			data[idx] += data[idx-1]
		}
	}
	asciigraph.Width(100)
	w.Write([]byte(asciigraph.Plot(data, asciigraph.Width(100), asciigraph.Height(15), asciigraph.Caption("Total Cases"))))
	return
}

func studentCases(w http.ResponseWriter, r *http.Request) {
	log.Debug("📈 getting ascii graph for student cases")
	earliestDate, errA := api.GetEarliestDate()
	recentDate, errB := api.GetRecentDate()
	if errA != nil || errB != nil {
		http.Error(w, "server encountered an issue", 500)
		return
	}
	rates, err := db.FetchInRange(earliestDate, recentDate)
	if err != nil {
		http.Error(w, "server encountered an issue", 500)
		return
	}
	w.Header().Set("Content-Type", "application/text")
	data := make([]float64, len(*rates), len(*rates))
	for idx, rate := range *rates {
		data[idx] = float64(rate.Campus + rate.City)
		if idx > 0 {
			data[idx] += data[idx-1]
		}
	}
	asciigraph.Width(100)
	w.Write([]byte(asciigraph.Plot(data, asciigraph.Width(100), asciigraph.Height(15), asciigraph.Caption("Student Cases"))))
	return
}

func staffCases(w http.ResponseWriter, r *http.Request) {
	log.Debug("📈 getting ascii graph for student cases")
	earliestDate, errA := api.GetEarliestDate()
	recentDate, errB := api.GetRecentDate()
	if errA != nil || errB != nil {
		http.Error(w, "server encountered an issue", 500)
		return
	}
	rates, err := db.FetchInRange(earliestDate, recentDate)
	if err != nil {
		http.Error(w, "server encountered an issue", 500)
		return
	}
	w.Header().Set("Content-Type", "application/text")
	data := make([]float64, len(*rates), len(*rates))
	for idx, rate := range *rates {
		data[idx] = float64(rate.Staff)
		if idx > 0 {
			data[idx] += data[idx-1]
		}
	}
	asciigraph.Width(100)
	w.Write([]byte(asciigraph.Plot(data, asciigraph.Width(100), asciigraph.Height(15), asciigraph.Caption("Staff Cases"))))
	return
}
