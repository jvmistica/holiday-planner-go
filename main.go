package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/jvmistica/gcal/pkg/gcal"
	"github.com/jvmistica/gcal/pkg/trello"
)

var (
	defaultTimeFormat         = "2006-01-02"
	defaultCalendarID         = "en.austrian#holiday@group.v.calendar.google.com"
	key                       = os.Getenv("GCP_API_KEY")
	listSuggestions           = "Leave Suggestions"
	listVacationWithoutLeaves = "Vacation without leaves"
)

// go run main.go -start=2023-05-01 -end=2023-05-31
func main() {
	// Parse command-line arguments
	calendarID := flag.String("calendarId", defaultCalendarID, "the calendarID")
	start := flag.String("start", "", "the start date")
	end := flag.String("end", "", "the end date")
	flag.Parse()

	vacationWithoutLeaves, suggestions, err := gcal.GetCalendarEvents(key, *start, *end, *calendarID)
	if err != nil {
		log.Fatal(err)
	}

	boardID, err := trello.CreateBoard("Holidays")
	if err != nil {
		log.Fatal(err)
	}

	vacationListID, err := trello.CreateList(boardID, listVacationWithoutLeaves, "1")
	if err != nil {
		log.Fatal(err)
	}

	suggestionListID, err := trello.CreateList(boardID, listSuggestions, "2")
	if err != nil {
		log.Fatal(err)
	}

	for _, i := range vacationWithoutLeaves {
		name := fmt.Sprintf("%s - %s -> %d days", i.Start.Format(defaultTimeFormat), i.End.Format(defaultTimeFormat), i.Count)
		if _, err := trello.CreateCard(vacationListID, name); err != nil {
			log.Fatal(err)
		}
	}

	for _, i := range suggestions {
		name := fmt.Sprintf("%s - %s -> %d leaves / %d days", i.Start.Format(defaultTimeFormat), i.End.Format(defaultTimeFormat), i.Leaves, i.Vacation)
		if _, err := trello.CreateCard(suggestionListID, name); err != nil {
			log.Fatal(err)
		}
	}
}
