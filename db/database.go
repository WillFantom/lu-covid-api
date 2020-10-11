package db

import (
	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3" //sqlite 3 driver
	log "github.com/sirupsen/logrus"
)

var (
	dbPath     string = "./database/cases.db"
	connection *sqlx.DB
)

// Init setups up a connection to the database if it exists. If not, it attempts
// to download a recent version from GitHub. If all else fails, it creates a new
// database. If that fails, the program exists.
func Init() error {
	if !exists() {
		if err := downloadFromGitHub(); err != nil {
			log.Errorln("⚠️ failed to get a copy of the full database from github!")
			log.Warnln("⚠️ you will be starting with a blank dataset")
			if err := createEmpty(); err != nil {
				log.Fatalln("🆘 could not even create a blank database... i give up...")
			}
		}
	}
	var err error
	connection, err = getConnection()
	if err != nil {
		log.Fatalln("🆘 could not create an initial connection to the database")
	}
	return nil
}
