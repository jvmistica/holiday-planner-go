package trello

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
)

var (
	defaultBoardName  = "Holidays"
	suggestion        = "Leave Suggestions"
	q1                = "Jan - Mar"
	q2                = "Apr - Jun"
	q3                = "Jul - Sep"
	q4                = "Oct - Dec"
	defaultBackground = "sky"
	trelloAPIKey      = os.Getenv("TRELLO_API_KEY")
	trelloAPIToken    = os.Getenv("TRELLO_API_TOKEN")
	createBoardURL    = "https://api.trello.com/1/boards/"
	createCardURL     = "https://api.trello.com/1/cards"
	createListURL     = "https://api.trello.com/1/boards/%s/lists"
)

type Response struct {
	ID string `json:"id"`
}

func CreateBoard(board string) (string, error) {
	client := &http.Client{}
	req, err := http.NewRequest(http.MethodPost, createBoardURL, nil)
	if err != nil {
		return "", err
	}

	q := req.URL.Query()
	q.Add("key", trelloAPIKey)
	q.Add("token", trelloAPIToken)
	q.Add("name", board)
	q.Add("prefs_background", defaultBackground)
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

func CreateList(board, list, position string) (string, error) {
	client := &http.Client{}
	url := fmt.Sprintf(createListURL, board)
	req, err := http.NewRequest(http.MethodPost, url, nil)
	if err != nil {
		return "", err
	}

	q := req.URL.Query()
	q.Add("key", trelloAPIKey)
	q.Add("token", trelloAPIToken)
	q.Add("name", list)
	q.Add("pos", position)
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

func CreateCard(list, card string) (string, error) {
	client := &http.Client{}
	req, err := http.NewRequest(http.MethodPost, createCardURL, nil)
	if err != nil {
		return "", err
	}

	q := req.URL.Query()
	q.Add("key", trelloAPIKey)
	q.Add("token", trelloAPIToken)
	q.Add("name", card)
	q.Add("idList", list)
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