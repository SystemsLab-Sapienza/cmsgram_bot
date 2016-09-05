package main

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/garyburd/redigo/redis"
)

func broadcastUpdate(w http.ResponseWriter, r *http.Request) error {
	const delay = 500 // Delay in ms
	var (
		payload = struct {
			Key   string
			Value string
		}{}
		rm ResponseMessage
	)

	// Read the request
	data, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return err
	}

	// Decode the JSON payload into the struct
	err = json.Unmarshal(data, &payload)
	if err != nil {
		return err
	}

	conn := pool.Get()
	defer conn.Close()

	switch payload.Key {
	case "crawler":
		var (
			b    bytes.Buffer
			data = struct {
				Title string
				URL   string
			}{}
		)

		values, err := redis.Values(conn.Do("HMGET", "crawler:news:"+payload.Value, "title", "url"))
		if err != err {
			return err
		}

		_, err = redis.Scan(values, &data.Title, &data.URL)
		if err != nil {
			return err
		}

		err = templates.ExecuteTemplate(&b, "news.tpl", data)
		if err != nil {
			return err
		}

		if config.TestRecipient != 0 {
			rm.Send(b.String(), config.TestRecipient)
			return nil
		}

		// Fetch list of recipients
		recipients, err := redis.Strings(conn.Do("SMEMBERS", "tgbot:feed:subscribers:avvisi"))
		if err != nil {
			return err
		}

		// Broadcast update to recipients
		for _, r := range recipients {
			chat, err := strconv.Atoi(r)
			if err != nil {
				return err
			}

			// Send update to user
			err = rm.Send(b.String(), chat)
			if err != nil {
				return err
			}
			time.Sleep(delay * time.Millisecond)
		}
	case "facebook":
	case "rss":
		var (
			b    bytes.Buffer
			data = struct {
				Name string
				URL  string
			}{}
		)

		values, err := redis.Values(conn.Do("HMGET", "rss:feed:"+payload.Value, "name", "url"))
		if err != nil {
			return err
		}

		_, err = redis.Scan(values, &data.Name, &data.URL)
		if err != nil {
			return err
		}

		data.URL = strings.Replace(data.URL, "WebRss?skin=rss", "WebHome", 1)

		err = templates.ExecuteTemplate(&b, "rss_update.tpl", data)
		if err != nil {
			return err
		}

		if config.TestRecipient != 0 {
			rm.Send(b.String(), config.TestRecipient)
			return nil
		}

		// Fetch list of recipients
		recipients, err := redis.Strings(conn.Do("SMEMBERS", "tgbot:feed:subscribers:t"+payload.Value))
		if err != nil {
			return err
		}

		// Broadcast update to recipients
		for _, r := range recipients {
			chat, err := strconv.Atoi(r)
			if err != nil {
				return err
			}

			// Send update to user
			err = rm.Send(b.String(), chat)
			if err != nil {
				return err
			}
			time.Sleep(delay * time.Millisecond)
		}
	default:
		log.Println("handlerBroadcastUpdate():", "Invalid payload:", payload)
	}

	return nil
}

func broadcastUpdateHandler(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	switch r.Method {
	case "POST":
		if err := broadcastUpdate(w, r); err != nil {
			log.Println("broadcastUpdate(): handling: ", r.RequestURI, err)
		}
	default:
		http.Error(w, "POST ONLY", http.StatusMethodNotAllowed)
	}
}
