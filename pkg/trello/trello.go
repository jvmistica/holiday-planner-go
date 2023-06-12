package trello

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
)

var (
	DefaultBoardName          = "Holidays"
	ListSuggestions           = "Leave suggestions"
	ListVacationWithoutLeaves = "Vacation without leaves"
	CreateBoardURL            = "https://api.trello.com/1/boards/"
	CreateCardURL             = "https://api.trello.com/1/cards"
	CreateListURL             = "https://api.trello.com/1/boards/%s/lists"

	trelloAPIKey   = os.Getenv("TRELLO_API_KEY")
	trelloAPIToken = os.Getenv("TRELLO_API_TOKEN")

	defaultBoardBackground = "sky"
)

// Response is the structure of the Calendar API's response
type Response struct {
	ID string `json:"id"`
}

// CreateBoard creates a board on Trello and returns the board ID
func CreateBoard(boardName string) (string, error) {
	client := &http.Client{}
	req, err := http.NewRequest(http.MethodPost, CreateBoardURL, nil)
	if err != nil {
		return "", err
	}

	q := req.URL.Query()
	q.Add("key", trelloAPIKey)
	q.Add("token", trelloAPIToken)
	q.Add("name", boardName)
	q.Add("prefs_background", defaultBoardBackground)
	req.URL.RawQuery = q.Encode()

	res, err := client.Do(req)
	if err != nil {
		return "", err
	}

	if res.StatusCode != http.StatusOK {
		return "", fmt.Errorf("failed to create board - status code: %d", res.StatusCode)
	}

	defer res.Body.Close()

	b, err := io.ReadAll(res.Body)
	if err != nil {
		return "", err
	}

	var response *Response
	if err := json.Unmarshal(b, &response); err != nil {
		return "", err
	}

	return response.ID, nil
}

// CreateList creates a list on Trello and returns the list ID
func CreateList(boardID, listName, position string) (string, error) {
	client := &http.Client{}
	url := fmt.Sprintf(CreateListURL, boardID)
	req, err := http.NewRequest(http.MethodPost, url, nil)
	if err != nil {
		return "", err
	}

	q := req.URL.Query()
	q.Add("key", trelloAPIKey)
	q.Add("token", trelloAPIToken)
	q.Add("name", listName)
	q.Add("pos", position) // order of the list
	req.URL.RawQuery = q.Encode()

	res, err := client.Do(req)
	if err != nil {
		return "", err
	}

	if res.StatusCode != http.StatusOK {
		return "", fmt.Errorf("failed to create list - status code: %d", res.StatusCode)
	}

	defer res.Body.Close()

	b, err := io.ReadAll(res.Body)
	if err != nil {
		return "", err
	}

	var response *Response
	if err := json.Unmarshal(b, &response); err != nil {
		return "", err
	}

	return response.ID, nil
}

// CreateCard creates a card on Trello and returns the card ID
func CreateCard(listID, cardName string) (string, error) {
	client := &http.Client{}
	req, err := http.NewRequest(http.MethodPost, CreateCardURL, nil)
	if err != nil {
		return "", err
	}

	q := req.URL.Query()
	q.Add("key", trelloAPIKey)
	q.Add("token", trelloAPIToken)
	q.Add("name", cardName)
	q.Add("idList", listID)
	req.URL.RawQuery = q.Encode()

	res, err := client.Do(req)
	if err != nil {
		return "", err
	}

	if res.StatusCode != http.StatusOK {
		return "", fmt.Errorf("failed to create card - status code: %d", res.StatusCode)
	}

	defer res.Body.Close()

	b, err := io.ReadAll(res.Body)
	if err != nil {
		return "", err
	}

	var response *Response
	if err := json.Unmarshal(b, &response); err != nil {
		return "", err
	}

	return response.ID, nil
}
