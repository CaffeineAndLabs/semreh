package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
	"time"
)

// sanitizeJSON is a workaround due to a bad format of the "API" used in the project
// TODO: parse website and generate a proper API
func sanitizeJSON(s []byte) []byte {
	json := string(s)
	json = strings.TrimPrefix(json, "<pre>")
	json = strings.TrimSuffix(json, "</pre>")

	return []byte(json)
}

func getAlmanaxCalendar() AlmanaxCalendar {
	almanaxCalJSON := "https://almanax.ordre2vlad.fr/source.php?lang=fr"
	resp, err := http.Get(almanaxCalJSON)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}

	body = sanitizeJSON(body)

	var cal AlmanaxCalendar
	err = json.Unmarshal(body, &cal)
	if err != nil {
		log.Fatal(err)
	}

	return cal
}

func getDailyAlmanax() AlmanaxEvent {
	cal := getAlmanaxCalendar()
	now := time.Now()
	_, month, day := now.Date()
	today := fmt.Sprintf("%02d/%02d", day, month)
	fmt.Println(today)
	event := cal[today]

	return event
}
