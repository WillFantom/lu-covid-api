package db

import (
	"fmt"
	"os"
	"time"

	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"
	log "github.com/sirupsen/logrus"
)

// Get a connection to the db (new if required)
func Check(dbPath string) error {
	if !dbExists(dbPath) {
		log.Debugln("💬 db does not exist...")
		if err := dbCreate(dbPath); err != nil {
			return err
		}
	}

	return nil
}

func dbExists(path string) bool {
	if _, err := os.Stat(path); err == nil {
		return true
	}
	log.Debugln("💬 database does not exist")
	return false
}

func dbCreate(path string) error {
	log.Debugln("💬 creating database")
	dbFile, err := os.Create(path)
	if err != nil {
		log.Errorln("⚠️ could not create database file")
		return err
	}
	dbFile.Close()
	database, _ := sqlx.Open("sqlite3", path)
	defer database.Close()
	ratesTableSQL := `CREATE TABLE rates (
		"id" integer NOT NULL PRIMARY KEY AUTOINCREMENT, 
		"date" date NOT NULL,
		"campus" integer NOT NULL,
		"city" integer NOT NULL,
		"staff" integer NOT NULL
	);`
	ctStatement, err := database.Prepare(ratesTableSQL)
	if err != nil {
		log.Errorln("⚠️ could not create rates table statement")
		return err
	}
	_, err = ctStatement.Exec()
	if err != nil {
		log.Errorln("⚠️ failed to create table")
		return err
	}

	if err := fillMissedData(path); err != nil {
		log.Errorln("🆘 failed to add missed days to db")
		return err
	}

	return nil
}

func InsertNewRate(path string, date string, campus uint64, city uint64, staff uint64) error {
	if err := Check(path); err != nil {
		return err
	}
	database, err := sqlx.Open("sqlite3", path)
	if err != nil {
		log.Errorln("⚠️ failed to open database")
		return err
	}
	defer database.Close()
	addRateSQL := `INSERT INTO rates(date, campus, city, staff) VALUES (?,?,?,?)`
	arStatement, err := database.Prepare(addRateSQL)
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

func FetchRates(path string, from time.Time, to time.Time) (*[]Rate, error) {
	to = dateAdjust(to)
	var rates []Rate

	if err := Check(path); err != nil {
		return nil, err
	}
	database, err := sqlx.Open("sqlite3", path)
	if err != nil {
		log.Errorln("⚠️ failed to open database")
		return nil, err
	}
	defer database.Close()
	statement := fmt.Sprintf("SELECT * FROM rates WHERE date BETWEEN date('%s') AND date('%s')", from.Format(time.RFC3339), to.Format(time.RFC3339))
	if err := database.Select(&rates, statement); err != nil {
		log.Errorln("⚠️ could not query the rates table")
		return nil, err
	}

	return &rates, nil
}

func DeleteRates(path string, from time.Time, to time.Time) error {
	to = dateAdjust(to)

	if err := Check(path); err != nil {
		return err
	}
	database, err := sqlx.Open("sqlite3", path)
	if err != nil {
		log.Errorln("⚠️ failed to open database")
		return err
	}
	defer database.Close()
	statement := fmt.Sprintf("DELETE FROM rates WHERE date BETWEEN date('%s') AND date('%s')", from.Format(time.RFC3339), to.Format(time.RFC3339))
	if _, err := database.Exec(statement); err != nil {
		log.Errorln("⚠️ could not query the rates table")
		return err
	}

	return nil
}

func dateAdjust(date time.Time) time.Time {
	return date.Add(time.Hour * 24)
}

// FillMissedData adds in any missed days worth of rates
// The site only shows the past 7 days, this began at 9 days...
func fillMissedData(path string) error {
	log.Debugln("💬 adding in missed days (left the site prior to this app")
	first := time.Date(2020, time.Month(10), 1, 0, 0, 0, 0, time.UTC)
	if err := InsertNewRate(path, first.Format(time.RFC3339), 1, 2, 0); err != nil {
		return err
	}
	second := time.Date(2020, time.Month(10), 2, 0, 0, 0, 0, time.UTC)
	if err := InsertNewRate(path, second.Format(time.RFC3339), 4, 2, 0); err != nil {
		return err
	}
	return nil
}
