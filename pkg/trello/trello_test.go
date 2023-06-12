package trello

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCreateBoard(t *testing.T) {
	t.Run("invalid URL", func(t *testing.T) {
		origURL := createBoardURL
		createBoardURL = "testInvalidURL%s"
		defer func() {
			createBoardURL = origURL
		}()

		result, err := CreateBoard(defaultBoardName)
		assert.Equal(t, "", result)
		assert.NotNil(t, err.Error())
	})

	t.Run("unsupported protocol", func(t *testing.T) {
		origURL := createBoardURL
		createBoardURL = "testInvalidURL"
		defer func() {
			createBoardURL = origURL
		}()

		result, err := CreateBoard(defaultBoardName)
		assert.Equal(t, "", result)
		assert.NotNil(t, err.Error())
	})

	t.Run("empty response", func(t *testing.T) {
		ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
		}))
		defer ts.Close()

		origURL := createBoardURL
		createBoardURL = ts.URL
		defer func() {
			createBoardURL = origURL
		}()

		result, err := CreateBoard(defaultBoardName)
		assert.Equal(t, "", result)
		assert.NotNil(t, err)
	})

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
			_, err := w.Write([]byte(`{"id": "abc123a36eaf8d75e160000f"}`))
			assert.Nil(t, err)
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
	t.Run("invalid URL", func(t *testing.T) {
		origURL := createListURL
		createListURL = "testInvalidURL"
		defer func() {
			createListURL = origURL
		}()

		result, err := CreateList("abc123a36eaf8d75e160000f", "sample list unauthorized", "1")
		assert.Equal(t, "", result)
		assert.NotNil(t, err.Error())
	})

	t.Run("unsupported protocol", func(t *testing.T) {
		origURL := createListURL
		createListURL = "testInvalidURL%s"
		defer func() {
			createListURL = origURL
		}()

		result, err := CreateList("abc123a36eaf8d75e160000f", "sample list unauthorized", "1")
		assert.Equal(t, "", result)
		assert.NotNil(t, err.Error())
	})

	t.Run("empty response", func(t *testing.T) {
		ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
		}))
		defer ts.Close()

		ts.URL = ts.URL + "/%s"
		origURL := createListURL
		createListURL = ts.URL
		defer func() {
			createListURL = origURL
		}()

		result, err := CreateList("abc123a36eaf8d75e160000f", "sample list unauthorized", "1")
		assert.NotNil(t, err)
		assert.Equal(t, "", result)
	})

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

		result, err := CreateList("abc123a36eaf8d75e160000f", "sample list unauthorized", "1")
		assert.Equal(t, "", result)
		assert.Equal(t, "failed to create list - status code: 401", err.Error())
	})

	t.Run("successful", func(t *testing.T) {
		ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
			_, err := w.Write([]byte(`{"id": "abc123a36ech8d75e160000f"}`))
			assert.Nil(t, err)
		}))
		defer ts.Close()

		ts.URL = ts.URL + "/%s"
		origURL := createListURL
		createListURL = ts.URL
		defer func() {
			createListURL = origURL
		}()

		result, err := CreateList("abc123a36eaf8d75e160000f", "sample list", "1")
		assert.Equal(t, "abc123a36ech8d75e160000f", result)
		assert.Nil(t, err)
	})
}

func TestCreateCard(t *testing.T) {
	t.Run("invalid URL", func(t *testing.T) {
		origURL := createCardURL
		createCardURL = "testInvalidURL%s"
		defer func() {
			createCardURL = origURL
		}()

		result, err := CreateCard("abc123a36ech8d75e160000f", "sample card unauthorized")
		assert.Equal(t, "", result)
		assert.NotNil(t, err.Error())
	})

	t.Run("unsupported protocol", func(t *testing.T) {
		origURL := createCardURL
		createCardURL = "testInvalidURL"
		defer func() {
			createCardURL = origURL
		}()

		result, err := CreateCard("abc123a36ech8d75e160000f", "sample card unauthorized")
		assert.Equal(t, "", result)
		assert.NotNil(t, err.Error())
	})

	t.Run("empty response", func(t *testing.T) {
		ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
		}))
		defer ts.Close()

		origURL := createCardURL
		createCardURL = ts.URL
		defer func() {
			createCardURL = origURL
		}()

		result, err := CreateCard("abc123a36ech8d75e160000f", "sample card unauthorized")
		assert.NotNil(t, err)
		assert.Equal(t, "", result)
	})

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

		result, err := CreateCard("abc123a36ech8d75e160000f", "sample card unauthorized")
		assert.Equal(t, "", result)
		assert.Equal(t, "failed to create card - status code: 401", err.Error())
	})

	t.Run("successful", func(t *testing.T) {
		ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
			_, err := w.Write([]byte(`{"id": "abc123a36eaf8d78u160000f"}`))
			assert.Nil(t, err)
		}))
		defer ts.Close()

		origURL := createCardURL
		createCardURL = ts.URL
		defer func() {
			createCardURL = origURL
		}()

		result, err := CreateCard("abc123a36ech8d75e160000f", "sample card")
		assert.Equal(t, "abc123a36eaf8d78u160000f", result)
		assert.Nil(t, err)
	})
}
