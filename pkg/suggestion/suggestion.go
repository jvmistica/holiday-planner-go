package suggestion

import (
	"fmt"

	"github.com/jvmistica/gcal/pkg/gcal"
	"github.com/jvmistica/gcal/pkg/trello"
)

var (
	defaultBoardName          = "Holidays"
	defaultTimeFormat         = "2006-01-02"
	listSuggestions           = "Leave suggestions"
	listVacationWithoutLeaves = "Vacation without leaves"
)

// GenerateSuggestions queries Google Calendar for holidays and generates a list of long weekends and suggested leaves on Trello
func GenerateSuggestions(gcpAPIKey, start, end, calendarID string) error {
	vacationWithoutLeaves, suggestions, err := gcal.GetCalendarEvents(gcpAPIKey, start, end, calendarID)
	if err != nil {
		return err
	}

	boardID, err := trello.CreateBoard(defaultBoardName)
	if err != nil {
		return err
	}

	vacationListID, err := trello.CreateList(boardID, listVacationWithoutLeaves, "1")
	if err != nil {
		return err
	}

	suggestionListID, err := trello.CreateList(boardID, listSuggestions, "2")
	if err != nil {
		return err
	}

	for _, i := range vacationWithoutLeaves {
		name := fmt.Sprintf("%s - %s -> %d days", i.Start.Format(defaultTimeFormat), i.End.Format(defaultTimeFormat), i.Count)
		if _, err := trello.CreateCard(vacationListID, name); err != nil {
			return err
		}
	}

	for _, i := range suggestions {
		name := fmt.Sprintf("%s - %s -> %d leaves / %d days", i.Start.Format(defaultTimeFormat), i.End.Format(defaultTimeFormat), i.Leaves, i.Vacation)
		if _, err := trello.CreateCard(suggestionListID, name); err != nil {
			return err
		}
	}

	return nil
}
