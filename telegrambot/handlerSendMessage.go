package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/garyburd/redigo/redis"
)

func sendMessage(w http.ResponseWriter, r *http.Request) error {
	const delay = 500 // Delay in ms
	var (
		news = struct {
			user_id string
			content string
		}{}
		payload = struct {
			Key   string
			Value string
		}{}
		rm ResponseMessage
	)

	// Read request
	data, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return err
	}

	// Decode the JSON payload
	err = json.Unmarshal(data, &payload)
	if err != nil {
		return err
	}

	// Check payload is valid
	if payload.Key != "message" {
		return nil // TODO check
	}

	conn := pool.Get()
	defer conn.Close()

	// Get news data
	values, err := redis.Values(conn.Do("HMGET", "webapp:messages:"+payload.Value, "user_id", "content"))
	if err != nil {
		return err
	}

	_, err = redis.Scan(values, &news.user_id, &news.content)
	if err != nil {
		return err
	}

	if config.TestRecipient != 0 {
		rm.Send(news.content, config.TestRecipient)
		return nil
	}

	// Fetch list of recipients
	recipients, err := redis.Strings(conn.Do("SMEMBERS", "tgbot:feed:subscribers:d"+news.user_id))
	if err != nil {
		return err
	}

	for _, r := range recipients {
		chat, err := strconv.Atoi(r)
		if err != nil {
			return err
		}

		// Send update to user
		err = rm.Send(news.content, chat)
		if err != nil {
			return err
		}
		time.Sleep(delay * time.Millisecond)
	}

	return nil
}

func sendMessageHandler(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	switch r.Method {
	case "POST":
		if err := sendMessage(w, r); err != nil {
			log.Println("sendMessage(): handling:", r.RequestURI, err)
		}
	default:
		http.Error(w, "POST ONLY", http.StatusMethodNotAllowed)
	}
}
