package api

import (
	"errors"
	"net/url"
	"strings"
	"time"
)

const expectedFormat = "2006-Jan-2"

func isTimeToday(t time.Time) bool {
	if t.Year() == time.Now().Year() && t.YearDay() == time.Now().YearDay() {
		return true
	}
	return false
}

func dateFromQuery(values url.Values) (time.Time, error) {
	day, hasDay := values["day"]
	month, hasMonth := values["month"]
	year, hasYear := values["year"]
	if !(hasDay && hasMonth && hasYear) || !(len(year) == 1 && len(month) == 1 && len(day) == 1) {
		return time.Time{}, errors.New(expectedFormat)
	}
	givenDate := strings.Join([]string{year[0], month[0], day[0]}, "-")
	date, err := time.Parse(expectedFormat, givenDate)
	if err != nil {
		return time.Time{}, errors.New(expectedFormat)
	}
	return date, nil
}
