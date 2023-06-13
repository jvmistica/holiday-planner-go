package main

import (
	"flag"
	"log"
	"os"

	"github.com/jvmistica/holiday-planner-go/pkg/suggestion"
)

var (
	defaultCalendarID = "en.austrian#holiday@group.v.calendar.google.com"
	gcpAPIKey         = os.Getenv("GCP_API_KEY")
)

func init() {
	gcpAPIKey := os.Getenv("GCP_API_KEY")
	if gcpAPIKey == "" {
		log.Fatal("missing environment variable GCP_API_KEY")
	}

	trelloAPIKey := os.Getenv("TRELLO_API_KEY")
	if trelloAPIKey == "" {
		log.Fatal("missing environment variable TRELLO_API_KEY")
	}

	trelloAPIToken := os.Getenv("TRELLO_API_TOKEN")
	if trelloAPIToken == "" {
		log.Fatal("missing environment variable TRELLO_API_TOKEN")
	}
}

func main() {
	calendarID := flag.String("calendarId", defaultCalendarID, "the calendarID")
	start := flag.String("start", "", "the start date")
	end := flag.String("end", "", "the end date")
	flag.Parse()

	if err := suggestion.GenerateSuggestions(gcpAPIKey, *start, *end, *calendarID); err != nil {
		log.Fatalf("failed to generate suggestions - %s", err.Error())
	}
}
