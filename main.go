package main

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/kelseyhightower/envconfig"
	"github.com/robfig/cron"
)

type Config struct {
	DiscordToken   string `split_words:"true",required:"true"`
	DiscordChannel string `split_words:"true",required:"true"`
}

var (
	Conf Config
)

func init() {
	err := envconfig.Process("SEMREH", &Conf)
	if err != nil {
		log.Fatal(err.Error())
	}
}

func main() {
	c := cron.New()
	c.AddFunc("@every 15s", SendDailyAlmanaxMessage)
	c.Start()

	router := mux.NewRouter()
	router.HandleFunc("/_health_check", healthCheckHandler).Methods("GET")
	log.Fatal(http.ListenAndServe(":8000", router))
}
