package api

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/gorilla/mux"

	"github.com/willfantom/lu-covid-api/db"

	log "github.com/sirupsen/logrus"
)

const (
	apiBase string = "/api/v1/"
)

func API(router *mux.Router) {
	router.HandleFunc(apiBase+"today", today)
	router.HandleFunc(apiBase+"recent", recent)
	router.HandleFunc(apiBase+"day", forDay)
	router.HandleFunc(apiBase+"totals", totals)
}

func today(w http.ResponseWriter, r *http.Request) {
	log.Debugln("✉️ getting today's rates")
	rate, err := db.MostRecent()
	if err != nil {
		http.Error(w, "server encountered an issue", 500)
		return
	}
	if !isTimeToday(rate.Date) {
		http.Error(w, "the server does not yet have data for today", 204)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	jsonData, err := json.Marshal(rate)
	if err != nil {
		http.Error(w, "server data can not be marshalled", 500)
		return
	}
	w.Write(jsonData)
	return
}

func recent(w http.ResponseWriter, r *http.Request) {
	log.Debugln("✉️ getting most recent")
	rate, err := db.MostRecent()
	if err != nil {
		http.Error(w, "server encountered an issue", 500)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	jsonData, err := json.Marshal(rate)
	if err != nil {
		http.Error(w, "server data can not be marshalled", 500)
		return
	}
	w.Write(jsonData)
	return
}

func forDay(w http.ResponseWriter, r *http.Request) {
	log.Debugln("✉️ getting rates for given date")
	date, err := dateFromQuery(r.URL.Query())
	if err != nil {
		http.Error(w, "expected parameter format: day=01&month=01&year=1970", 500)
		return
	}
	rates, err := db.FetchInRange(date, date)
	if err != nil {
		http.Error(w, "server encountered an issue", 500)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	jsonData, err := json.Marshal((*rates)[0])
	if err != nil {
		http.Error(w, "server data can not be marshalled", 500)
		return
	}
	w.Write(jsonData)
	return
}

func totals(w http.ResponseWriter, r *http.Request) {
	log.Debugln("✉️ getting totla rates")
	earliest, err := db.Earliest()
	if err != nil {
		http.Error(w, "server encountered an issue", 500)
		return
	}
	recent, err := db.MostRecent()
	if err != nil {
		http.Error(w, "server encountered an issue", 500)
		return
	}
	rates, err := db.FetchInRange(earliest.Date, recent.Date)
	if err != nil {
		http.Error(w, "server encountered an issue", 500)
		return
	}
	var staff, students uint64
	for _, rate := range *rates {
		staff += (rate.Staff)
		students += (rate.Campus + rate.City)
	}
	data := map[string]interface{}{
		"starting date": earliest.Date.Format(time.RFC1123),
		"ending date":   recent.Date.Format(time.RFC1123),
		"staff total":   staff,
		"student total": students,
		"total cases":   students + staff,
	}
	jsonData, err := json.Marshal(data)
	if err != nil {
		http.Error(w, "server data can not be marshalled", 500)
		return
	}
	w.Write(jsonData)
	return
}

// func Raw(w http.ResponseWriter, r *http.Request) {
// 	log.Debugln("✉️ getting raw rates")
// 	startDate, err := db.GetStartDate(rates.DatabasePath)
// 	if err != nil {
// 		http.Error(w, "no data could be found in database", 500)
// 	}
// 	rates, err := rates.ForDateRange(*startDate, time.Now())
// 	if err != nil {
// 		http.Error(w, "information reqest failed", 500)
// 	} else if rates == nil {
// 		http.Error(w, "todays information has not yet been published", 204)
// 	} else {
// 		w.Header().Set("Content-Type", "application/json")
// 		jsonData, err := json.Marshal(rates)
// 		if err != nil {
// 			http.Error(w, "data returned can not be marshalled", 500)
// 			return
// 		}
// 		w.Write(jsonData)
// 	}
// }

// func Summary(w http.ResponseWriter, r *http.Request) {
// 	log.Debugln("✉️ getting summary")
// 	startDate, err := db.GetStartDate(rates.DatabasePath)
// 	if err != nil {
// 		http.Error(w, "no data could be found in database", 500)
// 	}
// 	log.Infoln(startDate)
// 	rates, err := rates.ForDateRange(*startDate, time.Now())
// 	if err != nil {
// 		http.Error(w, "information reqest failed", 500)
// 	} else if rates == nil {
// 		http.Error(w, "todays information has not yet been published", 204)
// 	} else {
// 		w.Header().Set("Content-Type", "application/json")
// 		var total, staff, students uint64
// 		for _, rate := range *rates {
// 			staff += (rate.Staff)
// 			students += (rate.Campus + rate.City)
// 		}
// 		total = students + staff
// 		data := map[string]uint64{
// 			"Total Cases":   (total),
// 			"Student Cases": (students),
// 			"Staff Cases":   (staff),
// 		}
// 		jsonData, err := json.Marshal(data)
// 		if err != nil {
// 			http.Error(w, "data returned can not be marshalled", 500)
// 			return
// 		}
// 		w.Write(jsonData)
// 	}
// }
