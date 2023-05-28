package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
)

var (
	defaultCalendarID = "en.austrian#holiday@group.v.calendar.google.com"
	key = os.Getenv("GCP_API_KEY")
)

type Events struct {
	Summary       string  `json:"summary,omitempty"`
	NextSyncToken string  `json:"nextSyncToken,omitempty"`
	Items         []*Item `json:"items,omitempty"`
}

type Item struct {
	Summary     string `json:"summary,omitempty"`
	Description string `json:"description,omitempty"`
	Start struct {
		Date string `json:"date,omitempty"`
	} `json:"start,omitempty"`
}

func main() {
	// Parse command-line arguments
	calendarID := flag.String("calendarId", defaultCalendarID, "the calendarID")
	start := flag.String("start", "", "the start date")
	end := flag.String("end", "", "the end date")
	flag.Parse()

	id := url.QueryEscape(*calendarID)
	query := fmt.Sprintf("key=%s&timeMin=%s&timeMax=%s", key, *start, *end)
	url := fmt.Sprintf("https://www.googleapis.com/calendar/v3/calendars/%s/events?" + query, id)

	resp, err := http.Get(url)
	if err != nil {
		log.Fatal(err)
	}

	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}

	var events *Events
	if err := json.Unmarshal(body, &events); err != nil {
		log.Fatal(err)
	}

	holidays := getHolidays(events)
	fmt.Println(holidays)
}

func getHolidays(events *Events) map[string]string {
	holidays := make(map[string]string)
	for _, item := range events.Items {
		holidays[item.Summary] = item.Start.Date
	}

	return holidays
}
