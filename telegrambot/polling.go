package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"time"
)

func polling() {
	const delay time.Duration = 250
	var (
		response Response
		offset   int
		req      string

		URL = config.BotAPIBaseURL + config.BotAPIToken
	)

	for {
		if offset == 0 {
			req = URL + "/getUpdates"
		} else {
			req = URL + "/getUpdates?offset=" + strconv.Itoa(offset+1)
		}

		// Get updates from the BotAPI
		res, err := http.Get(req)
		if err != nil {
			log.Println("polling(): http.Get():", err)
			continue
		}

		// Read the response body
		body, err := ioutil.ReadAll(res.Body)
		if err != nil {
			log.Println("polling(): ioutil.ReadAll():", err)
			continue
		}
		res.Body.Close()

		// Decode the JSON payload
		err = json.Unmarshal(body, &response)
		if err != nil {
			log.Println("polling(): json.Unmarshal():", err)
			continue
		}

		if !response.Ok {
			log.Println("polling():", "Request not valid.")
			continue
		}

		// Process each update in its own goroutine
		if len(response.Result) != 0 {
			for _, r := range response.Result {
				if r.Update_id > offset {
					offset = r.Update_id
				}
				go dispatchRequest(r)
			}
		}

		time.Sleep(delay * time.Millisecond)
	}
}
