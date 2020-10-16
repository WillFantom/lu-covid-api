package main

import (
	"flag"
	"net/http"
	"os"
	"time"

	"github.com/willfantom/lu-covid-api/db"
	"github.com/willfantom/lu-covid-api/graphs"

	"github.com/willfantom/lu-covid-api/telegram"

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

	if err := db.Init(); err != nil {
		log.Fatalln("🆘 no database connection could be made")
	}

	justUpdate := flag.Bool("update-db", false, "just update the database and then exit")
	flag.Parse()

	log.Debugln("💬 initial scrape...")
	if err := rates.Scrape(true, true); err != nil {
		log.Fatalln("🆘 initial scrape failed!")
	}
	log.Debugln("✅ scrape success")

	if *justUpdate {
		log.Infoln("📁 database updated")
		os.Exit(0)
	}

	go fetch()

	if token, exists := os.LookupEnv("TG_TOKEN"); exists {
		go telegram.Init(token)
	}

	router := mux.NewRouter()
	router.HandleFunc("/", redirect)

	api.API(router.PathPrefix("/api").Subrouter())
	graphs.API(router.PathPrefix("/graphs").Subrouter())

	log.Debugln("💬 running api...")
	log.Fatal(http.ListenAndServe(":8080", router))

}

func redirect(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, "https://portal.lancaster.ac.uk/intranet/cms/coronavirus/covid-19-statistics", 301)
}

func fetch() {
	for {
		rates.Scrape(true, false)
		time.Sleep(time.Minute * 30)
	}
}
