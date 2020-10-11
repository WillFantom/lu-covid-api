package rates

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
	log "github.com/sirupsen/logrus"
	"github.com/willfantom/lu-covid-api/db"
)

const shortForm = "2006-Jan-2"

func getContent(endpoint string) (*ResponseContent, error) {

	var content ResponseContent
	response, err := http.Get(endpoint)
	if err != nil {
		return nil, err
	}
	if response.StatusCode != 200 {
		log.Debugln("🆘 api request failed. perhaps the endpoint has changed?")
		return nil, fmt.Errorf("api request got a non-200 response")
	}
	defer response.Body.Close()
	var data WebResponse
	if err := json.NewDecoder(response.Body).Decode(&data); err != nil {
		log.Errorln("🆘 api request failed. perhaps the data layout has changed?")
		return nil, err
	}

	content = data.Content[0]
	return &content, nil
}

func getRates(htmlContent *goquery.Document) ([]db.Rate, error) {

	var rates []db.Rate

	var rows []*goquery.Selection

	htmlContent.Find("table").Each(func(tableVal int, tablehtml *goquery.Selection) {
		tablehtml.Find("tr").Each(func(rowVal int, rowhtml *goquery.Selection) {
			if rowVal > 0 {
				rows = append(rows, rowhtml)
			}
		})
	})

	for _, row := range rows {
		cells := row.Find("td")
		if cells.Length() != 4 {
			return rates, fmt.Errorf("4 cells per row expected")
		}
		var rate db.Rate
		cell := cells.First()
		for n := 0; n < cells.Length(); n++ {
			text := strings.TrimSpace(cell.Text())
			if n == 0 {
				parsedDate, err := parseDate(text)
				if err == nil {
					rate.Date = parsedDate
				}
			} else if n == 1 {
				rate.Campus, _ = strconv.ParseUint(text, 10, 32)
			} else if n == 2 {
				rate.City, _ = strconv.ParseUint(text, 10, 32)
			} else if n == 3 {
				rate.Staff, _ = strconv.ParseUint(text, 10, 32)
			}
			cell = cell.Next()
		}
		rates = append(rates, rate)
	}

	return rates, nil
}

func parseDate(dateStr string) (time.Time, error) {
	var date time.Time
	splitDateStr := strings.Split(dateStr, " ")
	formattedStr := "2020-" + splitDateStr[2][0:3] + "-" + splitDateStr[1]
	date, err := time.Parse(shortForm, formattedStr)
	if err != nil {
		return date, err
	}
	return date, nil
}
