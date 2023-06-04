package main

import (
	"flag"
	"os"

	"github.com/jvmistica/gcal/pkg/query"
)

var (
	defaultCalendarID = "en.austrian#holiday@group.v.calendar.google.com"
	defaultBoardName  = "Holidays"
	suggestion        = "Leave Suggestions"
	q1                = "Jan - Mar"
	q2                = "Apr - Jun"
	q3                = "Jul - Sep"
	q4                = "Oct - Dec"
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

	// boardID, err := query.CreateBoard("Holidays")
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// fmt.Println("Board ID:", boardID)

	// suggestListID, err := query.CreateList(boardID, suggestion, "1")
	// if err != nil {
	// 	log.Fatal(err)
	// }

	// q1ListID, err := query.CreateList(boardID, q1, "2")
	// if err != nil {
	// 	log.Fatal(err)
	// }

	// q2ListID, err := query.CreateList(boardID, q2, "3")
	// if err != nil {
	// 	log.Fatal(err)
	// }

	// q3ListID, err := query.CreateList(boardID, q3, "4")
	// if err != nil {
	// 	log.Fatal(err)
	// }

	// q4ListID, err := query.CreateList(boardID, q4, "5")
	// if err != nil {
	// 	log.Fatal(err)
	// }
}
