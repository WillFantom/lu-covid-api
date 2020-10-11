package db

import (
	"errors"
	"fmt"
	"time"

	log "github.com/sirupsen/logrus"
)

const (
	insertRateTemplate  string = "INSERT INTO rates(date, campus, city, staff) VALUES (?,?,?,?)"
	selectRatesTemplate string = "SELECT date, campus, city, staff FROM rates WHERE date BETWEEN date('%s') AND date('%s')"
	updateRateTemplate  string = "UPDATE rates SET staff = %d, city = %d, campus = %d WHERE date = date('%s')"
	deleteRatesTemplate string = "DELETE FROM rates WHERE date BETWEEN date('%s') AND date('%s')"
	selectMostRecent    string = "SELECT date, campus, city, staff FROM rates ORDER BY date DESC LIMIT 1"
	selectEarliest      string = "SELECT date, campus, city, staff FROM rates ORDER BY date ASC LIMIT 1"
)

// FetchInRange gets all rate records that are within a
// provided time range.
func FetchInRange(from time.Time, to time.Time) (*[]Rate, error) {
	to = dateAdjust(to)
	var rates []Rate
	statement := fmt.Sprintf(selectRatesTemplate, from.Format(time.RFC3339), to.Format(time.RFC3339))
	if err := connection.Select(&rates, statement); err != nil {
		log.Errorln("⚠️ could not query the rates table")
		return nil, err
	}
	return &rates, nil
}

// DeleteInRange deletes all rate records that are within a
// provided time range.
func DeleteInRange(from time.Time, to time.Time) error {
	to = dateAdjust(to)
	statement := fmt.Sprintf(deleteRatesTemplate, from.Format(time.RFC3339), to.Format(time.RFC3339))
	if _, err := connection.Exec(statement); err != nil {
		log.Errorln("⚠️ could not query the rates table")
		return err
	}
	return nil
}

// UpdateForDate updates a given record. Useful for not making
// crazy large ids.
func UpdateForDate(date time.Time, staff uint64, city uint64, campus uint64) error {
	statement := fmt.Sprintf(updateRateTemplate, staff, city, campus, date)
	if _, err := connection.Exec(statement); err != nil {
		log.Errorln("⚠️ could not query the rates table")
		return err
	}
	return nil
}

// AddRate creates a new rate in the database.
func AddRate(date time.Time, staff uint64, city uint64, campus uint64) error {
	arStatement, err := connection.Prepare(insertRateTemplate)
	if err != nil {
		log.Errorln("⚠️ could not create add rate statement")
		return err
	}
	_, err = arStatement.Exec(date, campus, city, staff)
	if err != nil {
		log.Errorln("⚠️ failed to add rate")
		return err
	}
	return nil
}

// MostRecent gets the most recent rate record.
func MostRecent() (*Rate, error) {
	var rates []Rate
	if err := connection.Select(&rates, selectMostRecent); err != nil {
		log.Errorln("⚠️ could not query the rates table")
		return nil, err
	}
	if len(rates) != 1 {
		return nil, errors.New("could not get most recent record")
	}
	return &(rates[0]), nil
}

// Earliest gets the first recent rate record found in the database.
func Earliest() (*Rate, error) {
	var rates []Rate
	if err := connection.Select(&rates, selectEarliest); err != nil {
		log.Errorln("⚠️ could not query the rates table")
		return nil, err
	}
	if len(rates) != 1 {
		return nil, errors.New("could not get earliest record")
	}
	return &(rates[0]), nil
}
