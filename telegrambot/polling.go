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
		response ResponseT
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
			return
		}

		// Read the response body
		message, err := ioutil.ReadAll(res.Body)
		res.Body.Close()
		if err != nil {
			log.Println("polling(): ioutil.ReadAll():", err)
			return
		}

		// Decode the JSON payload
		err = json.Unmarshal(message, &response)
		if err != nil {
			log.Println("polling(): json.Unmarshal():", err)
			return
		}

		if !response.Ok {
			log.Println("polling():", "Request not valid.")
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
