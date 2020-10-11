package telegram

import (
	"fmt"
	"time"

	log "github.com/sirupsen/logrus"
	"github.com/willfantom/lu-covid-api/api"
	"github.com/willfantom/lu-covid-api/db"
	tb "gopkg.in/tucnak/telebot.v2"
)

const (
	genericErrorMessage string = `Sorry, something went wrong on my side. You can still check here though: 
	https://portal.lancaster.ac.uk/intranet/api/content/cms/coronavirus/covid-19-statistics`
)

var (
	tgBot *tb.Bot
)

func Init(token string) error {
	log.Infoln("🤖 creating telegram bot")
	var err error
	tgBot, err = tb.NewBot(tb.Settings{
		Token:  token,
		Poller: &tb.LongPoller{Timeout: 10 * time.Second},
	})
	if err != nil {
		log.Errorln("⚠️ could not create telegram bot")
		return err
	}

	tgBot.Handle("/start", start)
	tgBot.Handle("/today", today)
	tgBot.Handle("/recent", recent)
	tgBot.Handle("/totals", totals)

	tgBot.Start()

	return nil
}

func start(m *tb.Message) {
	log.Debugln("🤖 start message")
	message := ("Bot Functions:\n" +
		"/today | get today's case statistics\n" +
		"/recent | get the most recently published day's statistics")
	tgBot.Send(m.Sender, message)
	return
}

func today(m *tb.Message) {
	log.Debugln("🤖 getting rates today")
	rate, err := db.MostRecent()
	if err != nil {
		tgBot.Send(m.Sender, "Sorry, I can't get that data right now")
		return
	}
	if !api.IsTimeToday(rate.Date) {
		tgBot.Send(m.Sender, "No data has been added to show today's statistics")
		return
	}
	message := fmt.Sprintf("🦠 Todays Cases:\n"+
		"Total: %d\nStudents [campus]: %d\n Students [city]: %d\nStaff: %d",
		(rate.Staff + rate.Campus + rate.City), rate.Campus, rate.City, rate.Staff)
	tgBot.Send(m.Sender, message)
	return
}

func recent(m *tb.Message) {
	log.Debugln("🤖 getting most recent rate")
	rate, err := db.MostRecent()
	if err != nil {
		tgBot.Send(m.Sender, "Sorry, I can't get that data right now")
		return
	}
	message := fmt.Sprintf("🦠 Most Recent Daily Cases:\n Date: %s\n"+
		" Total: %d\n Students [campus]: %d\n Students [city]: %d\n Staff: %d",
		rate.Date.Format(time.RFC822), (rate.Staff + rate.Campus + rate.City), rate.Campus, rate.City, rate.Staff)
	tgBot.Send(m.Sender, message)
	return
}

func totals(m *tb.Message) {
}
