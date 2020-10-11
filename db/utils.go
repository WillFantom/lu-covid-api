package db

import (
	"errors"
	"io"
	"net/http"
	"os"
	"time"

	"github.com/jmoiron/sqlx"
	log "github.com/sirupsen/logrus"
)

const (
	downloadURL      string = "https://github.com/WillFantom/lu-covid-api/blob/main/database/cases.db?raw=true"
	ratesTableCreate string = `CREATE TABLE rates (
		"id" integer NOT NULL PRIMARY KEY AUTOINCREMENT, 
		"date" date NOT NULL,
		"campus" integer NOT NULL,
		"city" integer NOT NULL,
		"staff" integer NOT NULL
	);`
)

func exists() bool {
	if _, err := os.Stat(dbPath); err == nil {
		return true
	}
	return false
}

func downloadFromGitHub() error {
	log.Debugln("💬 attempting to fetch database from github")
	response, err := http.Get(downloadURL)
	if err != nil {
		return err
	}
	defer response.Body.Close()
	if response.StatusCode != 200 {
		return errors.New("non-200 status code, download failed")
	}
	os.MkdirAll(dbDir, os.ModePerm)
	dbFile, err := os.Create(dbPath)
	if err != nil {
		return err
	}
	defer dbFile.Close()
	_, err = io.Copy(dbFile, response.Body)
	if err != nil {
		return err
	}
	log.Debugln("✅ downloaded database from github")
	return nil
}

func createEmpty() error {
	log.Debugln("💬 creating empty database")
	os.MkdirAll(dbDir, os.ModePerm)
	dbFile, err := os.Create(dbPath)
	if err != nil {
		log.Errorln("⚠️ could not create database file")
		return err
	}
	dbFile.Close()
	connection, err := getConnection()
	if err != nil {
		log.Errorln("⚠️ could not create connection to new database file")
		return err
	}
	defer connection.Close()
	ctStatement, err := connection.Prepare(ratesTableCreate)
	if err != nil {
		log.Errorln("⚠️ could not create rates table statement")
		return err
	}
	_, err = ctStatement.Exec()
	if err != nil {
		log.Errorln("⚠️ failed to create table")
		return err
	}
	return nil
}

func getConnection() (*sqlx.DB, error) {
	connection, err := sqlx.Open("sqlite3", dbPath)
	if err != nil {
		log.Errorln("⚠️ failed to open a connection to the database")
		return nil, err
	}
	return connection, nil
}

// dateAdjust adds 1 day onto a date. This helps when using
// sqlite's BETWEEN operator
func dateAdjust(date time.Time) time.Time {
	return date.Add(time.Hour * 24)
}
