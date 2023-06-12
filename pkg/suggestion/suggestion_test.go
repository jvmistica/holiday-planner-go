package suggestion

import (
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/jvmistica/gcal/pkg/gcal"
	"github.com/jvmistica/gcal/pkg/trello"
	"github.com/stretchr/testify/assert"
)

func TestGenerateSuggestions(t *testing.T) {
	t.Run("path error, file does not exist", func(t *testing.T) {
		err := GenerateSuggestions("testKey", "2023-05-01", "2023-06-31", t.TempDir())
		assert.NotNil(t, err)
	})

	t.Run("failed to create board", func(t *testing.T) {
		tmpDir := t.TempDir()
		origDir := gcal.DefaultFilePath
		gcal.DefaultFilePath = tmpDir + "%s"
		defer func() {
			gcal.DefaultFilePath = origDir
		}()

		f, err := os.Create(tmpDir + "test")
		assert.Nil(t, err)
		defer f.Close()

		_, err = f.Write([]byte(`{
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
		assert.Nil(t, err)

		// mock board creation response
		ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusUnauthorized)
		}))
		defer ts.Close()

		origURL := trello.CreateBoardURL
		trello.CreateBoardURL = ts.URL
		defer func() {
			trello.CreateBoardURL = origURL
		}()

		err = GenerateSuggestions("testKey", "2023-08-01", "2023-10-31", "test")
		assert.Equal(t, "failed to create board - status code: 401", err.Error())
	})

	t.Run("failed to create list", func(t *testing.T) {
		tmpDir := t.TempDir()
		origDir := gcal.DefaultFilePath
		gcal.DefaultFilePath = tmpDir + "%s"
		defer func() {
			gcal.DefaultFilePath = origDir
		}()

		f, err := os.Create(tmpDir + "test")
		assert.Nil(t, err)
		defer f.Close()

		_, err = f.Write([]byte(`{
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
		assert.Nil(t, err)

		// mock board creation response
		ts1 := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
			_, err := w.Write([]byte(`{"id": "abc123a36eaf8d75e160000f"}`))
			assert.Nil(t, err)
		}))
		defer ts1.Close()

		origURL1 := trello.CreateBoardURL
		trello.CreateBoardURL = ts1.URL
		defer func() {
			trello.CreateBoardURL = origURL1
		}()

		// mock list creation response
		ts2 := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusUnauthorized)
		}))
		defer ts2.Close()

		ts2.URL = ts2.URL + "/%s"
		origURL2 := trello.CreateListURL
		trello.CreateListURL = ts2.URL
		defer func() {
			trello.CreateListURL = origURL2
		}()

		err = GenerateSuggestions("testKey", "2023-08-01", "2023-10-31", "test")
		assert.Equal(t, "failed to create list - status code: 401", err.Error())
	})

	t.Run("failed to create card", func(t *testing.T) {
		tmpDir := t.TempDir()
		origDir := gcal.DefaultFilePath
		gcal.DefaultFilePath = tmpDir + "%s"
		defer func() {
			gcal.DefaultFilePath = origDir
		}()

		f, err := os.Create(tmpDir + "test")
		assert.Nil(t, err)
		defer f.Close()

		_, err = f.Write([]byte(`{
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
		assert.Nil(t, err)

		// mock board creation response
		ts1 := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
			_, err := w.Write([]byte(`{"id": "abc123a36eaf8d75e160000f"}`))
			assert.Nil(t, err)
		}))
		defer ts1.Close()

		origURL1 := trello.CreateBoardURL
		trello.CreateBoardURL = ts1.URL
		defer func() {
			trello.CreateBoardURL = origURL1
		}()

		// mock list creation response
		ts2 := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
			_, err := w.Write([]byte(`{"id": "abc123a36ech8d75e160000f"}`))
			assert.Nil(t, err)
		}))
		defer ts2.Close()

		ts2.URL = ts2.URL + "/%s"
		origURL2 := trello.CreateListURL
		trello.CreateListURL = ts2.URL
		defer func() {
			trello.CreateListURL = origURL2
		}()

		// mock card creation response
		ts3 := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusUnauthorized)
		}))
		defer ts3.Close()

		origURL3 := trello.CreateCardURL
		trello.CreateCardURL = ts3.URL
		defer func() {
			trello.CreateCardURL = origURL3
		}()

		err = GenerateSuggestions("testKey", "2023-08-01", "2023-10-31", "test")
		assert.Equal(t, "failed to create card - status code: 401", err.Error())
	})

	t.Run("successful", func(t *testing.T) {
		tmpDir := t.TempDir()
		origDir := gcal.DefaultFilePath
		gcal.DefaultFilePath = tmpDir + "%s"
		defer func() {
			gcal.DefaultFilePath = origDir
		}()

		f, err := os.Create(tmpDir + "test")
		assert.Nil(t, err)
		defer f.Close()

		data, err := os.ReadFile("fixtures/test_gcal_response.json")
		assert.Nil(t, err)

		_, err = f.Write(data)
		assert.Nil(t, err)

		// mock board creation response
		ts1 := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
			_, err := w.Write([]byte(`{"id": "abc123a36eaf8d75e160000f"}`))
			assert.Nil(t, err)
		}))
		defer ts1.Close()

		origURL1 := trello.CreateBoardURL
		trello.CreateBoardURL = ts1.URL
		defer func() {
			trello.CreateBoardURL = origURL1
		}()

		// mock list creation response
		ts2 := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
			_, err := w.Write([]byte(`{"id": "abc123a36ech8d75e160000f"}`))
			assert.Nil(t, err)
		}))
		defer ts2.Close()

		ts2.URL = ts2.URL + "/%s"
		origURL2 := trello.CreateListURL
		trello.CreateListURL = ts2.URL
		defer func() {
			trello.CreateListURL = origURL2
		}()

		// mock card creation response
		ts3 := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
			_, err := w.Write([]byte(`{"id": "abc123a36eaf8d78u160000f"}`))
			assert.Nil(t, err)
		}))
		defer ts3.Close()

		origURL3 := trello.CreateCardURL
		trello.CreateCardURL = ts3.URL
		defer func() {
			trello.CreateCardURL = origURL3
		}()

		err = GenerateSuggestions("testKey", "2023-08-01", "2023-10-31", "test")
		assert.Nil(t, err)
	})
}
