package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/willfantom/lu-covid-api/db"

	"github.com/willfantom/lu-covid-api/rates"

	log "github.com/sirupsen/logrus"
)

const expectedFormat = "2006-Jan-2"

func dateFromQuery(values url.Values) (*time.Time, error) {
	day, hasDay := values["day"]
	month, hasMonth := values["month"]
	year, hasYear := values["year"]
	if !(hasDay && hasMonth && hasYear) || !(len(year) == 1 && len(month) == 1 && len(day) == 1) {
		return nil, fmt.Errorf("%s", expectedFormat)
	}
	givenDate := strings.Join([]string{year[0], month[0], day[0]}, "-")
	date, err := time.Parse(expectedFormat, givenDate)
	if err != nil {
		return &date, fmt.Errorf(expectedFormat)
	}
	return &date, nil
}

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
			return
		}
		w.Write(jsonData)
	}
}

func CasesForDay(w http.ResponseWriter, r *http.Request) {
	log.Debugln("✉️ getting rates for given date")
	date, err := dateFromQuery(r.URL.Query())
	if err != nil {
		http.Error(w, "Expected Format: day=01&month=01&year=1970", 500)
		return
	}
	rates, err := rates.ForDateRange(*date, *date)
	if err != nil {
		http.Error(w, "information reqest failed", 500)
	} else if rates == nil {
		http.Error(w, "todays information has not yet been published", 204)
	} else {
		w.Header().Set("Content-Type", "application/json")
		jsonData, err := json.Marshal((*rates)[0])
		if err != nil {
			http.Error(w, "data returned can not be marshalled", 500)
			return
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
			return
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
			return
		}
		w.Write(jsonData)
	}
}
