package main

import (
	"bytes"
	"encoding/json"
	"errors"
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
		b       bytes.Buffer
		message = struct {
			SenderID   string
			SenderName string
			Content    string
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
		return errors.New("Invalid payload")
	}

	conn := pool.Get()
	defer conn.Close()

	// Get message data
	values, err := redis.Values(conn.Do("HMGET", "webapp:messages:"+payload.Value, "user_id", "content"))
	if err != nil {
		return err
	}

	_, err = redis.Scan(values, &message.SenderID, &message.Content)
	if err != nil {
		return err
	}

	message.SenderName, err = getFullName(message.SenderID)
	if err != nil {
		return err
	}

	err = templates.ExecuteTemplate(&b, "message.tpl", message)
	if err != nil {
		return err
	}

	if config.TestRecipient != 0 {
		rm.Send(b.String(), config.TestRecipient)
		return nil
	}

	// Fetch list of recipients
	recipients, err := redis.Strings(conn.Do("SMEMBERS", "tgbot:feed:subscribers:d"+message.SenderID))
	if err != nil {
		return err
	}

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
