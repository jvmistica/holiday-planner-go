package gcal

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestGetHolidays(t *testing.T) {
	t.Run("error parsing date", func(t *testing.T) {
		events := `{"summary": "Holidays in Austria",
		 "items": [
		     {
		         "start": {"date": "2023/05/18" }
		     },
		     {
		         "start": {"date": "2023/05/28" }
		     },
		     {
		         "start": {"date": "2023/05/29" }
		     }
		 ]}`

		var e *Events
		err := json.Unmarshal([]byte(events), &e)
		assert.Nil(t, err)

		holidays, err := getHolidays(e)
		assert.NotNil(t, err)
		assert.Nil(t, holidays)
	})

	t.Run("successful", func(t *testing.T) {
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
	})
}

func TestGetWeekends(t *testing.T) {
	tests := []struct {
		startDate     string
		endDate       string
		expectedCount int
		wantErr       bool
	}{
		{
			startDate:     "2023-05-01",
			endDate:       "2023-05-31",
			expectedCount: 8,
		},
		{
			startDate:     "2023-12-15",
			endDate:       "2024-01-15",
			expectedCount: 10,
		},
		{
			startDate: "2023/12/15",
			endDate:   "2024-01-15",
			wantErr:   true,
		},
		{
			startDate: "2023-12-15",
			endDate:   "2024/01/15",
			wantErr:   true,
		},
	}

	for _, tt := range tests {
		weekends, err := getWeekends(tt.startDate, tt.endDate)
		if tt.wantErr {
			assert.Nil(t, weekends)
			assert.NotNil(t, err)
		} else {
			assert.Nil(t, err)
			assert.Equal(t, tt.expectedCount, len(weekends))
		}
	}
}

func TestGetVacationsWithoutLeaves(t *testing.T) {
	dates := `["2023-05-15T00:00:00Z", "2023-05-27T00:00:00Z", "2023-05-28T00:00:00Z", "2023-05-29T00:00:00Z"]`

	var free []time.Time
	err := json.Unmarshal([]byte(dates), &free)
	assert.Nil(t, err)

	result := getVacationsWithoutLeaves(free)
	assert.Equal(t, 1, len(result))
	assert.Equal(t, 3, result[0].Count)
}

func TestGetSuggestions(t *testing.T) {
	t.Run("one pair", func(t *testing.T) {
		dates := `[{"start": "2023-05-24T00:00:00Z", "end": "2023-05-28T00:00:00Z"}]`

		var free []*Vacation
		err := json.Unmarshal([]byte(dates), &free)
		assert.Nil(t, err)

		result := getSuggestions(free)
		assert.Nil(t, result)
	})

	t.Run("two pairs", func(t *testing.T) {
		dates := `[{"start": "2023-12-23T00:00:00Z", "end": "2023-12-26T00:00:00Z"}, {"start": "2023-12-30T00:00:00Z", "end": "2024-01-01T00:00:00Z"}]`

		var free []*Vacation
		err := json.Unmarshal([]byte(dates), &free)
		assert.Nil(t, err)

		result := getSuggestions(free)
		assert.Equal(t, 1, len(result))
		assert.Equal(t, 10, result[0].Vacation)
		assert.Equal(t, 3, result[0].Leaves)
		assert.Equal(t, "2023-12-23", result[0].Start.Format(defaultTimeFormat))
		assert.Equal(t, "2024-01-01", result[0].End.Format(defaultTimeFormat))
	})

	t.Run("three pairs", func(t *testing.T) {
		dates := `[{"start": "2023-05-22T00:00:00Z", "end": "2023-05-23T00:00:00Z"}, {"start": "2023-05-24T00:00:00Z", "end": "2023-05-25T00:00:00Z"}, {"start": "2023-05-27T00:00:00Z", "end": "2023-05-28T00:00:00Z"}]`

		var free []*Vacation
		err := json.Unmarshal([]byte(dates), &free)
		assert.Nil(t, err)

		result := getSuggestions(free)
		assert.Equal(t, 2, len(result))
		assert.Equal(t, 4, result[0].Vacation)
		assert.Equal(t, 0, result[0].Leaves)
		assert.Equal(t, "2023-05-22", result[0].Start.Format(defaultTimeFormat))
		assert.Equal(t, "2023-05-25", result[0].End.Format(defaultTimeFormat))
		assert.Equal(t, 5, result[1].Vacation)
		assert.Equal(t, 1, result[1].Leaves)
		assert.Equal(t, "2023-05-24", result[1].Start.Format(defaultTimeFormat))
		assert.Equal(t, "2023-05-28", result[1].End.Format(defaultTimeFormat))
	})
}

func TestFormatFreeTime(t *testing.T) {
	holidays := `["2023-12-25T00:00:00Z", "2023-12-26T00:00:00Z", "2023-12-21T00:00:00Z", "2024-01-01T00:00:00Z"]`
	weekends := `["2023-12-23T00:00:00Z", "2023-12-24T00:00:00Z", "2023-12-30T00:00:00Z", "2023-12-31T00:00:00Z"]`

	var h []time.Time
	err := json.Unmarshal([]byte(holidays), &h)
	assert.Nil(t, err)

	var w []time.Time
	err = json.Unmarshal([]byte(weekends), &w)
	assert.Nil(t, err)

	result := formatFreeTime(h, w)
	assert.Equal(t, 8, len(result))
	assert.Equal(t, "2023-12-21", result[0].Format(defaultTimeFormat))
	assert.Equal(t, "2024-01-01", result[7].Format(defaultTimeFormat))
}

func TestQueryCalendarAPI(t *testing.T) {
	t.Run("failed to create JSON file", func(t *testing.T) {
		var events *Events
		events, err := queryCalendarAPI(events, "def", "test", "2023-08-01T00:00:00Z", "2023-09-30T00:00:00Z", "/not/exist")
		assert.NotNil(t, err)
		assert.Nil(t, events)
	})

	t.Run("error - unauthorized", func(t *testing.T) {
		ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusUnauthorized)
		}))
		defer ts.Close()

		ts.URL = ts.URL + "/%s?" + "key=abc&timeMin=2023-08-01T00:00:00Z&timeMax=2023-09-30T00:00:00Z"
		origURL := eventsListURL
		eventsListURL = ts.URL
		defer func() {
			eventsListURL = origURL
		}()

		var events *Events
		events, err := queryCalendarAPI(events, "def", "test", "2023-08-01T00:00:00Z", "2023-09-30T00:00:00Z", t.TempDir()+"test.json")
		assert.Equal(t, "unsuccessful - status code: 401", err.Error())
		assert.Nil(t, events)
	})

	t.Run("error parsing JSON response", func(t *testing.T) {
		ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(`invalid`))
		}))
		defer ts.Close()

		ts.URL = ts.URL + "/%s?" + "key=abc&timeMin=2023-08-01T00:00:00Z&timeMax=2023-09-30T00:00:00Z"
		origURL := eventsListURL
		eventsListURL = ts.URL
		defer func() {
			eventsListURL = origURL
		}()

		var events *Events
		events, err := queryCalendarAPI(events, "abc", "test", "2023-08-01T00:00:00Z", "2023-09-30T00:00:00Z", t.TempDir()+"test.json")
		assert.NotNil(t, err)
		assert.Nil(t, events)
	})

	t.Run("successful", func(t *testing.T) {
		ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(`{
				"summary": "Holidays in Austria",
				"nextSyncToken": "CMDu0emHs_8CEAAYASCn_tSAAg==",
				"items": [{
					"summary": "Assumption of Mary",
					"description": "Public holiday",
					"start": {
					    "date": "2023-08-15"
					}
				},
				{
					"summary": "Yom Kippur",
					"description": "Observance\nTo hide observances, go to Google Calendar Settings \u003e Holidays in Austria",
					"start": {
						"date": "2023-09-25"
					}
				}]}`))
		}))
		defer ts.Close()

		ts.URL = ts.URL + "/%s?" + "key=abc&timeMin=2023-08-01T00:00:00Z&timeMax=2023-09-30T00:00:00Z"
		origURL := eventsListURL
		eventsListURL = ts.URL
		defer func() {
			eventsListURL = origURL
		}()

		var events *Events
		events, err := queryCalendarAPI(events, "abc", "test", "2023-08-01T00:00:00Z", "2023-09-30T00:00:00Z", t.TempDir()+"test.json")
		assert.Nil(t, err)
		assert.Equal(t, "Holidays in Austria", events.Summary)
		assert.Equal(t, "Assumption of Mary", events.Items[0].Summary)
		assert.Equal(t, "Yom Kippur", events.Items[1].Summary)
	})
}

func TestGetCalendarEvents(t *testing.T) {
	t.Run("file does not exist", func(t *testing.T) {
		tmpDir := t.TempDir()
		origDir := defaultResultDir
		defaultResultDir = tmpDir + "/%s"
		defer func() {
			defaultResultDir = origDir
		}()

		ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(`{
				"summary": "Holidays in Austria",
				"nextSyncToken": "CMDu0emHs_8CEAAYASCn_tSAAg==",
				"items": [{
					"summary": "Assumption of Mary",
					"description": "Public holiday",
					"start": {
					    "date": "2023-08-15"
					}
				},
				{
					"summary": "Yom Kippur",
					"description": "Observance\nTo hide observances, go to Google Calendar Settings \u003e Holidays in Austria",
					"start": {
						"date": "2023-09-25"
					}
				}]}`))
		}))
		defer ts.Close()

		ts.URL = ts.URL + "/%s?" + "key=abc&timeMin=2023-08-01T00:00:00Z&timeMax=2023-09-30T00:00:00Z"
		origURL := eventsListURL
		eventsListURL = ts.URL
		defer func() {
			eventsListURL = origURL
		}()

		v, s, err := GetCalendarEvents("abc", "2023-08-01", "2023-09-30", "test")
		assert.Nil(t, err)
		assert.Equal(t, 1, len(v))
		assert.Equal(t, 3, v[0].Count)
		assert.Equal(t, "2023-09-23", v[0].Start.Format(defaultTimeFormat))
		assert.Equal(t, "2023-09-25", v[0].End.Format(defaultTimeFormat))
		assert.Nil(t, s)
	})

	t.Run("file exists", func(t *testing.T) {
		tmpDir := t.TempDir()
		origDir := defaultResultDir
		defaultResultDir = tmpDir + "/%s"
		defer func() {
			defaultResultDir = origDir
		}()

		f, err := os.Create(tmpDir + "/test.json")
		assert.Nil(t, err)
		defer f.Close()

		f.Write([]byte(`{
			"summary": "Holidays in Austria",
			"nextSyncToken": "CMDu0emHs_8CEAAYASCn_tSAAg==",
			"items": [{
				"summary": "Assumption of Mary",
				"description": "Public holiday",
				"start": {
				    "date": "2023-08-15"
				}
			},
			{
				"summary": "Yom Kippur",
				"description": "Observance\nTo hide observances, go to Google Calendar Settings \u003e Holidays in Austria",
				"start": {
					"date": "2023-09-25"
				}
			}]}`))

		v, s, err := GetCalendarEvents("abc", "2023-08-01", "2023-09-30", "test")
		assert.Nil(t, err)
		assert.Equal(t, 1, len(v))
		assert.Equal(t, 3, v[0].Count)
		assert.Equal(t, "2023-09-23", v[0].Start.Format(defaultTimeFormat))
		assert.Equal(t, "2023-09-25", v[0].End.Format(defaultTimeFormat))
		assert.Nil(t, s)
	})

	t.Run("error parsing JSON file", func(t *testing.T) {
		tmpDir := t.TempDir()
		origDir := defaultResultDir
		defaultResultDir = tmpDir + "/%s"
		defer func() {
			defaultResultDir = origDir
		}()

		f, err := os.Create(tmpDir + "/test.json")
		assert.Nil(t, err)
		defer f.Close()

		f.Write([]byte(`invalid`))

		v, s, err := GetCalendarEvents("abc", "2023-08-01T00:00:00Z", "2023-09-30T00:00:00Z", "test")
		assert.NotNil(t, err)
		assert.Nil(t, v)
		assert.Nil(t, s)
	})

	t.Run("error querying calendar API", func(t *testing.T) {
		tmpDir := "/not/exist"
		origDir := defaultResultDir
		defaultResultDir = tmpDir + "/%s"
		defer func() {
			defaultResultDir = origDir
		}()

		v, s, err := GetCalendarEvents("abc", "2023-08-01T00:00:00Z", "2023-09-30T00:00:00Z", "test")
		assert.NotNil(t, err)
		assert.Nil(t, v)
		assert.Nil(t, s)
	})

	t.Run("error parsing date while getting holidays", func(t *testing.T) {
		tmpDir := t.TempDir()
		origDir := defaultResultDir
		defaultResultDir = tmpDir + "/%s"
		defer func() {
			defaultResultDir = origDir
		}()

		f, err := os.Create(tmpDir + "/test.json")
		assert.Nil(t, err)
		defer f.Close()

		f.Write([]byte(`{
			"summary": "Holidays in Austria",
			"nextSyncToken": "CMDu0emHs_8CEAAYASCn_tSAAg==",
			"items": [{
				"summary": "Assumption of Mary",
				"description": "Public holiday",
				"start": {
				    "date": "2023/08/15"
				}
			},
			{
				"summary": "Yom Kippur",
				"description": "Observance\nTo hide observances, go to Google Calendar Settings \u003e Holidays in Austria",
				"start": {
					"date": "2023/09/25"
				}
			}]}`))

		v, s, err := GetCalendarEvents("abc", "2023-08-01T00:00:00Z", "2023-09-30T00:00:00Z", "test")
		assert.NotNil(t, err)
		assert.Nil(t, v)
		assert.Nil(t, s)
	})

	t.Run("error parsing date while getting weekends", func(t *testing.T) {
		tmpDir := t.TempDir()
		origDir := defaultResultDir
		defaultResultDir = tmpDir + "/%s"
		defer func() {
			defaultResultDir = origDir
		}()

		f, err := os.Create(tmpDir + "/test.json")
		assert.Nil(t, err)
		defer f.Close()

		f.Write([]byte(`{
			"summary": "Holidays in Austria",
			"nextSyncToken": "CMDu0emHs_8CEAAYASCn_tSAAg==",
			"items": [{
				"summary": "Assumption of Mary",
				"description": "Public holiday",
				"start": {
				    "date": "2023-08-15"
				}
			},
			{
				"summary": "Yom Kippur",
				"description": "Observance\nTo hide observances, go to Google Calendar Settings \u003e Holidays in Austria",
				"start": {
					"date": "2023-09-25"
				}
			}]}`))

		v, s, err := GetCalendarEvents("abc", "2023/08/01T00:00:00Z", "2023/09/30T00:00:00Z", "test")
		assert.NotNil(t, err)
		assert.Nil(t, v)
		assert.Nil(t, s)
	})
}
