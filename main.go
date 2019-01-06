package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"github.com/kelseyhightower/envconfig"
	"github.com/robfig/cron"
)

var (
	Conf      Config
	GitCommit string
)

func init() {
	// Environment variable
	err := envconfig.Process("SEMREH", &Conf)
	if err != nil {
		log.Fatal(err.Error())
	}

	// CLI
	version := flag.Bool("v", false, "Print the version of the application")
	notifAlmanax := flag.Bool("notif-almanax", false, "Force to send Almanax notification at start")
	notifDofusNews := flag.Int("notif-dofus-news", 0, "Send to discord N last news from Dofus news RSS")
	flag.Parse()

	// Print version when -v is used
	if *version {
		fmt.Printf("Commit: %s\n", GitCommit)
		os.Exit(0)
	}

	if *notifAlmanax {
		SendDailyAlmanaxMessage()
		os.Exit(0)
	}

	if *notifDofusNews > 0 {
		sendLastNNews(*notifDofusNews)
		os.Exit(0)
	}
}

func main() {
	log.Println("semreh started")

	c := cron.New()
	almanaxCronTime := fmt.Sprintf("0 %s %s * * *", Conf.AlmanaxDailyReportMinute, Conf.AlmanaxDailyReportHour)
	c.AddFunc(almanaxCronTime, SendDailyAlmanaxMessage)
	c.AddFunc("0 * * * * *", cronRSSNews)
	c.Start()

	router := mux.NewRouter()
	router.HandleFunc("/_health_check", healthCheckHandler).Methods("GET")
	log.Fatal(http.ListenAndServe(":8000", router))
}
