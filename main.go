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

type Config struct {
	DiscordToken   string `split_words:"true",required:"true"`
	DiscordChannel string `split_words:"true",required:"true"`
}

var (
	Conf      Config
	GitCommit string
)

func init() {
	err := envconfig.Process("SEMREH", &Conf)
	if err != nil {
		log.Fatal(err.Error())
	}

	version := flag.Bool("v", false, "Print the version of the application")
	notifAtStart := flag.Bool("almanax-at-start", false, "Force to send Almanax notification at start")
	flag.Parse()

	// Print version when -v is used
	if *version {
		fmt.Printf("Commit: %s\n", GitCommit)
		os.Exit(0)
	}

	if *notifAtStart {
		SendDailyAlmanaxMessage()
	}
}

func main() {
	log.Println("semreh started")

	c := cron.New()
	c.AddFunc("0 1 0 * * *", SendDailyAlmanaxMessage)
	c.Start()

	router := mux.NewRouter()
	router.HandleFunc("/_health_check", healthCheckHandler).Methods("GET")
	log.Fatal(http.ListenAndServe(":8000", router))
}
