package main

import (
	"flag"
	"net/http"
	"os"
	"time"

	"github.com/gorilla/mux"
	"github.com/willfantom/lu-covid-api/api"

	"github.com/willfantom/lu-covid-api/rates"

	log "github.com/sirupsen/logrus"
)

const (
	apiBase string = "/api/v1/"
)

func main() {
	log.Infoln("Lancaster Univerity covid cases API 🦠 (unofficial)")

	if _, exist := os.LookupEnv("DEBUG"); exist {
		log.SetLevel(log.DebugLevel)
	}

	justUpdate := flag.Bool("update-db", false, "just update the database and then exit")
	flag.Parse()

	log.Debugln("💬 initial scrape...")
	if err := rates.Scrape(); err != nil {
		log.Fatalln("🆘 initial scrape failed!")
	}
	if err := rates.WriteRates(true); err != nil {
		log.Fatalln("🆘 could not perm a database write for initial data")
	}
	log.Debugln("✅ scrape success")

	if *justUpdate {
		log.Infoln("📁 database updated")
		os.Exit(0)
	}

	go fetch()

	router := mux.NewRouter()
	router.HandleFunc("/", redirect)
	router.HandleFunc(apiBase+"today", api.CasesToday)
	router.HandleFunc(apiBase+"summary", api.Summary)
	router.HandleFunc(apiBase+"raw", api.Raw)
	log.Debugln("💬 running api...")
	log.Fatal(http.ListenAndServe(":8080", router))

}

func redirect(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, "https://portal.lancaster.ac.uk/intranet/cms/coronavirus/covid-19-statistics", 301)
}

func fetch() {
	for {
		rates.Scrape()
		rates.WriteRates(false)
		time.Sleep(time.Minute * 30)
	}
}
