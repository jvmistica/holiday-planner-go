package query

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
)

var (
	defaultBackground = "sky"
	trelloAPIKey      = os.Getenv("TRELLO_API_KEY")
	trelloAPIToken    = os.Getenv("TRELLO_API_TOKEN")
)

func createBoard(board string) string {
	client := &http.Client{}
	url := "https://api.trello.com/1/boards/"
	req, err := http.NewRequest(http.MethodPost, url, nil)
	if err != nil {
		log.Fatal(err)
	}

	q := req.URL.Query()
	q.Add("key", trelloAPIKey)
	q.Add("token", trelloAPIToken)
	q.Add("name", board)
	q.Add("prefs_background", defaultBackground)
	req.URL.RawQuery = q.Encode()

	res, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	}

	defer res.Body.Close()
	body, err := io.ReadAll(res.Body)
	if err != nil {
		log.Fatal(err)
	}

	return string(body)
}

func createList(board, list, position string) {
	client := &http.Client{}
	url := fmt.Sprintf("https://api.trello.com/1/boards/%s/lists", board)
	req, err := http.NewRequest(http.MethodPost, url, nil)
	if err != nil {
		log.Fatal(err)
	}

	q := req.URL.Query()
	q.Add("key", trelloAPIKey)
	q.Add("token", trelloAPIToken)
	q.Add("name", list)
	q.Add("pos", position)
	req.URL.RawQuery = q.Encode()

	res, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(res, err)
}

func createCard(list, card string) {
	client := &http.Client{}
	url := "https://api.trello.com/1/cards"
	req, err := http.NewRequest(http.MethodPost, url, nil)
	if err != nil {
		log.Fatal(err)
	}

	q := req.URL.Query()
	q.Add("key", trelloAPIKey)
	q.Add("token", trelloAPIToken)
	q.Add("name", card)
	q.Add("idList", list)
	req.URL.RawQuery = q.Encode()

	res, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(res, err)
}
