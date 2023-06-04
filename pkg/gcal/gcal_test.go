package gcal

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestGetHolidays(t *testing.T) {
	events := `{"summary": "Holidays in Austria",
	 "items": [
	     {
	         "start": {"date": "2023-05-18" }
	     },
	     {
	         "start": {"date": "2023-05-28" }
	     },
	     {
	         "start": {"date": "2023-05-29" }
	     }
	 ]}`

	var e *Events
	err := json.Unmarshal([]byte(events), &e)
	assert.Nil(t, err)

	holidays, err := getHolidays(e)
	assert.Nil(t, err)
	assert.Equal(t, 3, len(holidays))
}

func TestGetWeekends(t *testing.T) {
	tests := []struct {
		startDate     string
		endDate       string
		expectedCount int
	}{
		{
			startDate:     "2023-05-01T00:00:00Z",
			endDate:       "2023-05-31T00:00:00Z",
			expectedCount: 8,
		},
		{
			startDate:     "2023-12-15T00:00:00Z",
			endDate:       "2024-01-15T00:00:00Z",
			expectedCount: 10,
		},
	}

	for _, tt := range tests {
		weekends, err := getWeekends(tt.startDate, tt.endDate)
		assert.Nil(t, err)
		assert.Equal(t, tt.expectedCount, len(weekends))

	}
}

func TestGetVacationsWithoutLeaves(t *testing.T) {
	dates := `["2023-05-28T00:00:00Z", "2023-05-29T00:00:00Z", "2023-05-27T00:00:00Z", "2023-05-28T00:00:00Z"]`

	var free []time.Time
	err := json.Unmarshal([]byte(dates), &free)
	assert.Nil(t, err)

	result := getVacationsWithoutLeaves(free)
	assert.Equal(t, 1, len(result))
	assert.Equal(t, "3", result[0]["count"])
}
