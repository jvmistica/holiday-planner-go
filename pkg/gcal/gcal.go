package gcal

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"sort"
	"time"
)

var (
	defaultTimeFormat          = "2006-01-02"
	defaultFilePath            = "./pkg/gcal/data/%s.json"
	defaultMinDaysWithoutLeave = 3
	eventsListURL              = "https://www.googleapis.com/calendar/v3/calendars/%s/events?"
)

// Events is the structure of the response from the Google Calendar API
type Events struct {
	Summary       string  `json:"summary,omitempty"`
	NextSyncToken string  `json:"nextSyncToken,omitempty"`
	Items         []*Item `json:"items,omitempty"`
}

// Item is the structure of each event
type Item struct {
	Summary     string `json:"summary,omitempty"`
	Description string `json:"description,omitempty"`
	Start       struct {
		Date string `json:"date,omitempty"`
	} `json:"start,omitempty"`
}

// Suggestion contains the details of suggested vacation dates
type Suggestion struct {
	Vacation int
	Leaves   int
	Start    time.Time
	End      time.Time
}

// Vacation contains the details of vacation dates (long weekends, etc.)
type Vacation struct {
	Start time.Time
	End   time.Time
	Count int
}

// GetCalendarEvents returns all holidays, weekends, and suggested vacation leaves
func GetCalendarEvents(key, start, end, calendarID string) ([]*Vacation, []*Suggestion, error) {
	var events *Events
	filePath := fmt.Sprintf(defaultFilePath, calendarID)

	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		fmt.Println("Initiating GET request..")

		var err error
		events, err = queryCalendarAPI(events, key, calendarID, start, end, filePath)
		if err != nil {
			return nil, nil, err
		}
	} else {
		fmt.Println("Skipping GET request..")

		data, err := os.ReadFile(filePath)
		if err != nil {
			return nil, nil, err
		}

		if err := json.Unmarshal(data, &events); err != nil {
			return nil, nil, err
		}
	}

	holidays, err := getHolidays(events)
	if err != nil {
		return nil, nil, err
	}

	weekends, err := getWeekends(start, end)
	if err != nil {
		return nil, nil, err
	}

	freeTime := formatFreeTime(holidays, weekends)
	vacationWithoutLeaves := getVacationsWithoutLeaves(freeTime)
	suggestions := getSuggestions(vacationWithoutLeaves)

	return vacationWithoutLeaves, suggestions, nil
}

// getHolidays returns a map of holidays and their date
func getHolidays(events *Events) ([]time.Time, error) {
	var holidays []time.Time
	for _, item := range events.Items {
		start, err := time.Parse(defaultTimeFormat, item.Start.Date)
		if err != nil {
			return holidays, err
		}
		holidays = append(holidays, start)
	}

	return holidays, nil
}

// getWeekends returns a list of dates that fall on Saturdays and Sundays
func getWeekends(startDate, endDate string) ([]time.Time, error) {
	var weekends []time.Time
	start, err := time.Parse(defaultTimeFormat, startDate)
	if err != nil {
		return weekends, err
	}

	end, err := time.Parse(defaultTimeFormat, endDate)
	if err != nil {
		return weekends, err
	}

	for d := start; !d.After(end); d = d.AddDate(0, 0, 1) {
		if d.Weekday().String() == "Saturday" || d.Weekday().String() == "Sunday" {
			weekends = append(weekends, d)
		}
	}

	return weekends, nil
}

// getVacationsWithoutLeaves returns free time of 3 (default) or more days where filing a vacation leave
// is not needed (i.e. long weekends)
func getVacationsWithoutLeaves(freeTime []time.Time) []*Vacation {
	var toDate time.Time
	days := 0
	fromDate := freeTime[0]
	var dates []*Vacation
	i := 0

	for i < len(freeTime) {
		if days == 0 {
			fromDate = freeTime[i]
		}

		days += 1
		if i == len(freeTime)-1 || freeTime[i].AddDate(0, 0, 1) != freeTime[i+1] {
			toDate = freeTime[i]
			if days >= defaultMinDaysWithoutLeave {
				date := &Vacation{
					Start: fromDate,
					End:   toDate,
					Count: days,
				}
				dates = append(dates, date)
			}
			days = 0
		}
		i += 1
	}

	return dates
}

// getSuggestions returns a list of suggested vacation dates
func getSuggestions(pairs []*Vacation) []*Suggestion {
	var suggestions []*Suggestion
	for i, d := range pairs {
		if i >= len(pairs)-1 {
			continue
		}

		start := d.Start
		end := d.End
		nextStart := pairs[i+1].Start
		nextEnd := pairs[i+1].End
		leaves := int((nextStart.Sub(end).Hours() / 24) - 1)
		if leaves <= 5 {
			vacation := int(nextEnd.Sub(start).Hours() / 24)
			if vacation-leaves > 1 {
				suggestions = append(suggestions,
					&Suggestion{
						Vacation: vacation + 1,
						Leaves:   leaves,
						Start:    d.Start,
						End:      pairs[i+1].End,
					})
			}
		}
	}

	return suggestions
}

// formatFreeTime returns a sorted list of holidays and weekends combined
func formatFreeTime(holidays, weekends []time.Time) []time.Time {
	var freeTime []time.Time
	freeTime = append(freeTime, holidays...)
	freeTime = append(freeTime, weekends...)

	sort.Slice(freeTime, func(i, j int) bool {
		return freeTime[i].Before(freeTime[j])
	})

	var newList []time.Time
	for k, v := range freeTime {
		if k == len(freeTime)-1 || v != freeTime[k+1] {
			newList = append(newList, v)
		}
	}

	return newList
}

// queryCalendarAPI gets the list of holidays from the Calendar API and writes it into a JSON file
func queryCalendarAPI(events *Events, key, calendarID, start, end, filePath string) (*Events, error) {
	id := url.QueryEscape(calendarID)
	query := fmt.Sprintf("key=%s&timeMin=%sT00:00:00Z&timeMax=%sT00:00:00Z", key, start, end)
	url := fmt.Sprintf(eventsListURL+query, id)

	f, err := os.Create(filePath)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unsuccessful - status code: %d", resp.StatusCode)
	}

	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if err := json.Unmarshal(body, &events); err != nil {
		return nil, err
	}

	s, err := json.MarshalIndent(events, "", "    ")
	if err != nil {
		return nil, err
	}

	if _, err := f.Write(s); err != nil {
		return nil, err
	}

	return events, nil
}
