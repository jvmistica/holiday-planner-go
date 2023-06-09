package trello

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCreateBoard(t *testing.T) {
	t.Run("error - unauthorized board creation", func(t *testing.T) {
		ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusUnauthorized)
		}))
		defer ts.Close()

		origURL := createBoardURL
		createBoardURL = ts.URL
		defer func() {
			createBoardURL = origURL
		}()

		result, err := CreateBoard(defaultBoardName)
		assert.Equal(t, "", result)
		assert.Equal(t, "failed to create board - status code: 401", err.Error())
	})

	t.Run("successful", func(t *testing.T) {
		ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(`{"id": "abc123a36eaf8d75e160000f"}`))
		}))
		defer ts.Close()

		origURL := createBoardURL
		createBoardURL = ts.URL
		defer func() {
			createBoardURL = origURL
		}()

		result, err := CreateBoard(defaultBoardName)
		assert.Equal(t, "abc123a36eaf8d75e160000f", result)
		assert.Nil(t, err)
	})
}

func TestCreateList(t *testing.T) {
	t.Run("error - unauthorized list creation", func(t *testing.T) {
		ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusUnauthorized)
		}))
		defer ts.Close()

		ts.URL = ts.URL + "/%s"
		origURL := createListURL
		createListURL = ts.URL
		defer func() {
			createListURL = origURL
		}()

		result, err := CreateList(defaultBoardName, "sample list unauthorized", "1")
		assert.Equal(t, "", result)
		assert.Equal(t, "failed to create list - status code: 401", err.Error())
	})

	t.Run("successful", func(t *testing.T) {
		ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(`{"id": "abc123a36ech8d75e160000f"}`))
		}))
		defer ts.Close()

		ts.URL = ts.URL + "/%s"
		origURL := createListURL
		createListURL = ts.URL
		defer func() {
			createListURL = origURL
		}()

		result, err := CreateList(defaultBoardName, "sample list", "1")
		assert.Equal(t, "abc123a36ech8d75e160000f", result)
		assert.Nil(t, err)
	})
}

func TestCreateCard(t *testing.T) {
	t.Run("error - unauthorized card creation", func(t *testing.T) {
		ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusUnauthorized)
		}))
		defer ts.Close()

		origURL := createCardURL
		createCardURL = ts.URL
		defer func() {
			createCardURL = origURL
		}()

		result, err := CreateCard("sample list", "sample card unauthorized")
		assert.Equal(t, "", result)
		assert.Equal(t, "failed to create card - status code: 401", err.Error())
	})

	t.Run("successful", func(t *testing.T) {
		ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(`{"id": "abc123a36eaf8d78u160000f"}`))
		}))
		defer ts.Close()

		origURL := createCardURL
		createCardURL = ts.URL
		defer func() {
			createCardURL = origURL
		}()

		result, err := CreateCard("sample list", "sample card")
		assert.Equal(t, "abc123a36eaf8d78u160000f", result)
		assert.Nil(t, err)
	})
}
