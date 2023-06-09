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
	dates := `["2023-05-15T00:00:00Z", "2023-05-27T00:00:00Z", "2023-05-28T00:00:00Z", "2023-05-29T00:00:00Z"]`

	var free []time.Time
	err := json.Unmarshal([]byte(dates), &free)
	assert.Nil(t, err)

	result := getVacationsWithoutLeaves(free)
	assert.Equal(t, 1, len(result))
	assert.Equal(t, "3", result[0]["count"])
}

func TestGetSuggestions(t *testing.T) {
	t.Run("one pair", func(t *testing.T) {
		dates := `[{"start": "2023-05-24T00:00:00Z", "end": "2023-05-28T00:00:00Z"}]`

		var free []map[string]string
		err := json.Unmarshal([]byte(dates), &free)
		assert.Nil(t, err)

		result, err := getSuggestions(free)
		assert.Nil(t, err)
		assert.Nil(t, result)
	})

	t.Run("two pairs", func(t *testing.T) {
		dates := `[{"start": "2023-12-23T00:00:00Z", "end": "2023-12-26T00:00:00Z"}, {"start": "2023-12-30T00:00:00Z", "end": "2024-01-01T00:00:00Z"}]`

		var free []map[string]string
		err := json.Unmarshal([]byte(dates), &free)
		assert.Nil(t, err)

		result, err := getSuggestions(free)
		assert.Nil(t, err)
		assert.Equal(t, 1, len(result))
		assert.Equal(t, "10", result[0].Vacation)
		assert.Equal(t, "3", result[0].Leaves)
		assert.Equal(t, "2023-12-23T00:00:00Z", result[0].Start)
		assert.Equal(t, "2024-01-01T00:00:00Z", result[0].End)
	})

	t.Run("three pairs", func(t *testing.T) {
		dates := `[{"start": "2023-05-22T00:00:00Z", "end": "2023-05-23T00:00:00Z"}, {"start": "2023-05-24T00:00:00Z", "end": "2023-05-25T00:00:00Z"}, {"start": "2023-05-27T00:00:00Z", "end": "2023-05-28T00:00:00Z"}]`

		var free []map[string]string
		err := json.Unmarshal([]byte(dates), &free)
		assert.Nil(t, err)

		result, err := getSuggestions(free)
		assert.Nil(t, err)
		assert.Equal(t, 2, len(result))
		assert.Equal(t, "4", result[0].Vacation)
		assert.Equal(t, "0", result[0].Leaves)
		assert.Equal(t, "2023-05-22T00:00:00Z", result[0].Start)
		assert.Equal(t, "2023-05-25T00:00:00Z", result[0].End)
		assert.Equal(t, "5", result[1].Vacation)
		assert.Equal(t, "1", result[1].Leaves)
		assert.Equal(t, "2023-05-24T00:00:00Z", result[1].Start)
		assert.Equal(t, "2023-05-28T00:00:00Z", result[1].End)
	})
}
