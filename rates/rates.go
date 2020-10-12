package rates

import (
	"fmt"
	"strings"

	"github.com/willfantom/lu-covid-api/db"

	"github.com/PuerkitoBio/goquery"
	log "github.com/sirupsen/logrus"
)

const (
	// Endpoint is the url where theuniversity provide the table
	Endpoint string = "https://portal.lancaster.ac.uk/intranet/api/content/cms/coronavirus/covid-19-statistics"
	//DatabasePath is file path for sqlite db
	DatabasePath string = "./database/cases.db"
)

// Scrape gets all the data from the table given out by the university.
func Scrape(write bool, updateExisting bool) error {
	log.Debugln("💬 attempting to scrape for new data")

	cmsContent, err := getContent(Endpoint)
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

	log.Debugln("💬 found rates via scrape, count:", len(rates))

	if write {
		if err := writeRates(updateExisting, rates); err != nil {
			log.Debugln("💬⚠️ ", err.Error())
			return err
		}
	}

	return nil
}

// WriteRates add the most recently scraped data to the database.
// if update is set to True, it will update pre existing dates with
// the most recently scraped data.
func writeRates(updateExisting bool, data []db.Rate) error {

	for _, rate := range data {
		currData, err := db.FetchInRange(rate.Date, rate.Date)
		if err != nil {
			return err
		}
		if len(*currData) > 0 && !updateExisting {
			continue
		} else if len(*currData) > 0 && updateExisting {
			if err := db.UpdateForDate(rate.Date, rate.Staff, rate.City, rate.Campus); err != nil {
				return err
			}
		} else if len(*currData) <= 0 {
			if err := db.AddRate(rate.Date, rate.Staff, rate.City, rate.Campus); err != nil {
				return err
			}
		}
	}

	return nil
}
