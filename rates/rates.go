package rates

import (
	"fmt"
	"strings"
	"time"

	"github.com/willfantom/lu-covid-api/db"

	"github.com/PuerkitoBio/goquery"
	log "github.com/sirupsen/logrus"
)

const (
	endpoint     string = "https://portal.lancaster.ac.uk/intranet/api/content/cms/coronavirus/covid-19-statistics"
	databasePath string = "cases.db"
)

var (
	//StartDate is the first date that any data was made avaliable
	StartDate   time.Time = time.Date(2020, time.Month(10), 1, 0, 0, 0, 0, time.UTC)
	currentData []Rate
)

// Scrape gets all the data from the table given out by the university.
func Scrape() error {
	log.Debugln("💬 attempting to scrape for new data")

	cmsContent, err := getContent(endpoint)
	if err != nil {
		log.Errorln("⚠️ could not get a response from the api")
		return err
	}

	htmlContent, err := goquery.NewDocumentFromReader(strings.NewReader(cmsContent.Content))
	if err != nil {
		log.Errorln("⚠️ api content does not seem to have valid html")
		return err
	}

	rates, err := getRates(htmlContent)
	if err != nil {
		log.Errorln("⚠️ table of rates could not be parsed")
		return err
	} else if len(rates) != 7 {
		log.Errorln("⚠️ data from rates table seems to be wrong")
		return fmt.Errorf("0 rates were parsed from table (should be 7)")
	}

	log.Debugln("💬 found rates via scrape")
	currentData = rates

	return nil
}

// WriteRates add the most recently scraped data to the database.
// if update is set to True, it will update pre existing dates with
// the most recently scraped data.
func WriteRates(updateExisting bool) error {

	for _, rate := range currentData {
		currData, err := db.FetchRates(databasePath, rate.Date, rate.Date)
		if err != nil {
			return err
		}
		if currData != nil && !updateExisting {
			continue
		}
		if err := db.DeleteRates(databasePath, rate.Date, rate.Date); err != nil {
			return err
		}
		if err := db.InsertNewRate(databasePath, rate.Date.Format(time.RFC3339), rate.CampusStudents, rate.CityStudents, rate.Staff); err != nil {
			return err
		}
	}

	return nil
}

// Today returns the given rates for the current day.
// If these have not yet been provided, no error will be flagged,
// yet rate will be nil.
func Today() (*db.Rate, error) {
	rates, err := ForDateRange(time.Now(), time.Now())
	if err != nil {
		return nil, err
	}
	if rates != nil {
		if len((*rates)) != 1 {
			return nil, nil
		}
		return &(*rates)[0], nil
	}
	return nil, nil
}

// ForDateRange returns the given rates a range of days.
func ForDateRange(from time.Time, to time.Time) (*[]db.Rate, error) {
	rates, err := db.FetchRates(databasePath, from, to)
	if err != nil {
		return nil, err
	}
	if len((*rates)) < 1 {
		return nil, nil
	}
	return rates, nil
}
