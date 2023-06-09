package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/jvmistica/gcal/pkg/gcal"
	// "github.com/jvmistica/gcal/pkg/trello"
)

var (
	defaultCalendarID = "en.austrian#holiday@group.v.calendar.google.com"
	key               = os.Getenv("GCP_API_KEY")
	suggestion        = "Leave Suggestions"
	q1                = "Jan - Mar"
	q2                = "Apr - Jun"
	q3                = "Jul - Sep"
	q4                = "Oct - Dec"
)

// go run main.go -start=2023-05-01T00:00:00Z -end=2023-05-31T00:00:00Z
func main() {
	// Parse command-line arguments
	calendarID := flag.String("calendarId", defaultCalendarID, "the calendarID")
	start := flag.String("start", "", "the start date")
	end := flag.String("end", "", "the end date")
	flag.Parse()

	events, suggestions, err := gcal.GetCalendarEvents(key, *start, *end, *calendarID)
	if err != nil {
		log.Fatal(err)
	}

	e, err := json.MarshalIndent(events, "", "    ")
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("No leaves")
	fmt.Println(string(e))

	s, err := json.MarshalIndent(suggestions, "", "    ")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Suggestions")
	fmt.Println(string(s))

	// boardID, err := trello.CreateBoard("Holidays")
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// fmt.Println("Board ID:", boardID)

	// suggestListID, err := trello.CreateList(boardID, suggestion, "1")
	// if err != nil {
	// 	log.Fatal(err)
	// }

	// q1ListID, err := trello.CreateList(boardID, q1, "2")
	// if err != nil {
	// 	log.Fatal(err)
	// }

	// // q2ListID, err := trello.CreateList(boardID, q2, "3")
	// // if err != nil {
	// // 	log.Fatal(err)
	// // }

	// // q3ListID, err := trello.CreateList(boardID, q3, "4")
	// // if err != nil {
	// // 	log.Fatal(err)
	// // }

	// // q4ListID, err := trello.CreateList(boardID, q4, "5")
	// // if err != nil {
	// // 	log.Fatal(err)
	// // }

	// _, _ = trello.CreateCard(q1ListID, string(e))
	// _, _ = trello.CreateCard(suggestListID, string(s))
}
