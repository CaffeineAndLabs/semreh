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

const (
	feedURLDofusNews      = "http://www.dofus.com/fr/rss/news.xml"
	feedURLDofusDevBlog   = "http://www.dofus.com/fr/rss/devblog.xml"
	feedURLDofusChangelog = "http://www.dofus.com/fr/rss/changelog.xml"
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
	notifDofusNewsDevBlog := flag.Int("notif-dofus-devblog", 0, "Send to discord N last news from Dofus DevBlog")
	notifDofusNewsChangelog := flag.Int("notif-dofus-changelog", 0, "Send to discord N last news from Dofus Changelog")
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

	if *notifDofusNews > 0 || *notifDofusNewsDevBlog > 0 || *notifDofusNewsChangelog > 0 {
		sendLastNNews(feedURLDofusNews, *notifDofusNews)
		sendLastNNews(feedURLDofusDevBlog, *notifDofusNewsDevBlog)
		sendLastNNews(feedURLDofusChangelog, *notifDofusNewsChangelog)
		os.Exit(0)
	}
}

func main() {
	log.Println("semreh started")

	c := cron.New()
	almanaxCronTime := fmt.Sprintf("0 %s %s * * *", Conf.AlmanaxDailyReportMinute, Conf.AlmanaxDailyReportHour)
	c.AddFunc(almanaxCronTime, SendDailyAlmanaxMessage)
	c.AddFunc("0 * * * * *", func() {
		cronRSSNews(feedURLDofusNews)
	})
	c.AddFunc("0 * * * * *", func() {
		cronRSSNews(feedURLDofusDevBlog)
	})
	c.AddFunc("0 * * * * *", func() {
		cronRSSNews(feedURLDofusChangelog)
	})
	c.Start()

	router := mux.NewRouter()
	router.HandleFunc("/_health_check", healthCheckHandler).Methods("GET")
	log.Fatal(http.ListenAndServe(":8000", router))
}
