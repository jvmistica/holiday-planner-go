package main

import (
	"flag"
	"os"

	"github.com/jvmistica/gcal/pkg/query"
)

var (
	defaultCalendarID = "en.austrian#holiday@group.v.calendar.google.com"
	key               = os.Getenv("GCP_API_KEY")
)

// go run main.go -start=2023-05-01T00:00:00Z -end=2023-05-31T00:00:00Z
func main() {
	// Parse command-line arguments
	calendarID := flag.String("calendarId", defaultCalendarID, "the calendarID")
	start := flag.String("start", "", "the start date")
	end := flag.String("end", "", "the end date")
	flag.Parse()

	query.Query(key, start, end, calendarID)
}
