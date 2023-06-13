package suggestion

import (
	"fmt"

	"github.com/jvmistica/gcal/pkg/gcal"
	"github.com/jvmistica/gcal/pkg/trello"
)

// GenerateSuggestions queries Google Calendar for holidays and generates a trello.List of long weekends and suggested leaves on Trello
func GenerateSuggestions(gcpAPIKey, start, end, calendarID string) error {
	vacationWithoutLeaves, suggestions, err := gcal.GetCalendarEvents(gcpAPIKey, start, end, calendarID)
	if err != nil {
		return err
	}

	boardID, err := trello.CreateBoard(trello.DefaultBoardName)
	if err != nil {
		return err
	}

	// create vacation list on the first column
	vacationListID, err := trello.CreateList(boardID, trello.ListVacationWithoutLeaves, "1")
	if err != nil {
		return err
	}

	// create suggestion list on the second column
	suggestionListID, err := trello.CreateList(boardID, trello.ListSuggestions, "2")
	if err != nil {
		return err
	}

	for _, i := range vacationWithoutLeaves {
		name := fmt.Sprintf("%s - %s -> %d days", i.Start.Format(gcal.DefaultTimeFormat), i.End.Format(gcal.DefaultTimeFormat), i.Count)
		if _, err := trello.CreateCard(vacationListID, name); err != nil {
			return err
		}
	}

	for _, i := range suggestions {
		name := fmt.Sprintf("%s - %s -> %d leaves / %d days", i.Start.Format(gcal.DefaultTimeFormat), i.End.Format(gcal.DefaultTimeFormat), i.Leaves, i.Vacation)
		if _, err := trello.CreateCard(suggestionListID, name); err != nil {
			return err
		}
	}

	return nil
}
