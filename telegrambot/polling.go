package main

import (
	"encoding/json"
	"errors"
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

	client := &http.Client{
		Timeout: 0,
	}

	for {
		if offset == 0 {
			req = URL + "/getUpdates"
		} else {
			req = URL + "/getUpdates?offset=" + strconv.Itoa(offset+1)
		}

		err := func() error {
			// Get updates from the BotAPI
			res, err := client.Get(req)
			if err != nil {
				return err
			}
			defer res.Body.Close()

			// Read the response body
			body, err := ioutil.ReadAll(res.Body)
			if err != nil {
				return err
			}

			// Decode the JSON payload
			err = json.Unmarshal(body, &response)
			if err != nil {
				return errors.New(string(body))
			}

			if !response.Ok {
				return errors.New("Request not valid.")
			}

			return nil
		}()

		if err != nil {
			log.Println("polling():", err)
			time.Sleep(time.Second)
			continue
		}

		// Process each update in its own goroutine
		if len(response.Result) != 0 {
			for _, r := range response.Result {
				if r.Update_id > offset {
					offset = r.Update_id
				}

				go func() {
					if err := dispatchRequest(r); err != nil {
						log.Println(err)
					}
				}()
			}
		}

		time.Sleep(delay * time.Millisecond)
	}
}
