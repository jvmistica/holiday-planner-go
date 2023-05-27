package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
)

type Events struct {
	Kind          string  `json:"kind,omitempty"`
	Etag          string  `json:"etag,omitempty"`
	Summary       string  `json:"summary,omitempty"`
	Updated       string  `json:"updated,omitempty"`
	TimeZone      string  `json:"timeZone,omitempty"`
	NextSyncToken string  `json:"nextSyncToken,omitempty"`
	Items         []*Item `json:"items,omitempty"`
}
type Item struct {
	Kind        string `json:"kind,omitempty"`
	Etag        string `json:"etag,omitempty"`
	ID          string `json:"id,omitempty"`
	Status      string `json:"status,omitempty"`
	HTMLLink    string `json:"htmlLink,omitempty"`
	Created     string `json:"created,omitempty"`
	Updated     string `json:"updated,omitempty"`
	Summary     string `json:"summary,omitempty"`
	Description string `json:"description,omitempty"`
	Creator     struct {
		Email       string `json:"email,omitempty"`
		DisplayName string `json:"displayName,omitempty"`
		Self        bool   `json:"self,omitempty"`
	} `json:"creator,omitempty"`
	Organizer struct {
		Email       string `json:"email,omitempty"`
		DisplayName string `json:"displayName,omitempty"`
		Self        bool   `json:"self,omitempty"`
	} `json:"organizer,omitempty"`
	Start struct {
		Date string `json:"date,omitempty"`
	} `json:"start,omitempty"`
	End struct {
		Date string `json:"date,omitempty"`
	} `json:"end,omitempty"`
	Transparency string `json:"transparency,omitempty"`
	Visibility   string `json:"visibility,omitempty"`
	ICalUID      string `json:"iCalUID,omitempty"`
	Sequence     int    `json:"sequence,omitempty"`
	EventType    string `json:"eventType,omitempty"`
}

func main() {
	key := os.Getenv("GCP_API_KEY")
	timeMin := "2023-01-01T00:00:00Z"
	timeMax := "2024-01-01T00:00:00Z"
	query := fmt.Sprintf("key=%s&timeMin=%s&timeMax=%s", key, timeMin, timeMax)
	url := "https://www.googleapis.com/calendar/v3/calendars/en.austrian%23holiday%40group.v.calendar.google.com/events?" + query

	resp, _ := http.Get(url)

	defer resp.Body.Close()
	body, _ := io.ReadAll(resp.Body)

	var events *Events
	json.Unmarshal(body, &events)

	for _, item := range events.Items {
		fmt.Println(fmt.Sprintf("%s (%s) - %s", item.Summary, item.Description, item.Start.Date))
	}
}
