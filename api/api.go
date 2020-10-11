package api

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/willfantom/lu-covid-api/db"

	"github.com/willfantom/lu-covid-api/rates"

	log "github.com/sirupsen/logrus"
)

func CasesToday(w http.ResponseWriter, r *http.Request) {
	log.Debugln("✉️ getting rates today")
	rate, err := rates.Today()
	if err != nil {
		http.Error(w, "information reqest failed", 500)
	} else if rate == nil {
		http.Error(w, "todays information has not yet been published", 204)
	} else {
		w.Header().Set("Content-Type", "application/json")
		jsonData, err := json.Marshal(rate)
		if err != nil {
			http.Error(w, "data returned can not be marshalled", 500)
		}
		w.Write(jsonData)
	}
}

func Raw(w http.ResponseWriter, r *http.Request) {
	log.Debugln("✉️ getting raw rates")
	startDate, err := db.GetStartDate(rates.DatabasePath)
	if err != nil {
		http.Error(w, "no data could be found in database", 500)
	}
	rates, err := rates.ForDateRange(*startDate, time.Now())
	if err != nil {
		http.Error(w, "information reqest failed", 500)
	} else if rates == nil {
		http.Error(w, "todays information has not yet been published", 204)
	} else {
		w.Header().Set("Content-Type", "application/json")
		jsonData, err := json.Marshal(rates)
		if err != nil {
			http.Error(w, "data returned can not be marshalled", 500)
		}
		w.Write(jsonData)
	}
}

func Summary(w http.ResponseWriter, r *http.Request) {
	log.Debugln("✉️ getting summary")
	startDate, err := db.GetStartDate(rates.DatabasePath)
	if err != nil {
		http.Error(w, "no data could be found in database", 500)
	}
	rates, err := rates.ForDateRange(*startDate, time.Now())
	if err != nil {
		http.Error(w, "information reqest failed", 500)
	} else if rates == nil {
		http.Error(w, "todays information has not yet been published", 204)
	} else {
		w.Header().Set("Content-Type", "application/json")
		var total, staff, students uint64
		for _, rate := range *rates {
			staff += (rate.Staff)
			students += (rate.Campus + rate.City)
		}
		total = students + staff
		data := map[string]uint64{
			"Total Cases":   (total),
			"Student Cases": (students),
			"Staff Cases":   (staff),
		}
		jsonData, err := json.Marshal(data)
		if err != nil {
			http.Error(w, "data returned can not be marshalled", 500)
		}
		w.Write(jsonData)
	}
}
