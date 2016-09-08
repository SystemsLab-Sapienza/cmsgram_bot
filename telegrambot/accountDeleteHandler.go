package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/garyburd/redigo/redis"
)

func accountDelete(w http.ResponseWriter, r *http.Request) error {
	var (
		payload = struct {
			Key   string
			Value string
		}{}
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

	userID := payload.Value
	switch payload.Key {
	case "user":
		conn := pool.Get()
		defer conn.Close()

		// Fetch list of subscribers
		susbcribers, err := redis.Strings(conn.Do("SMEMBERS", "tgbot:feed:subscribers:d"+userID))
		if err != nil {
			return err
		}

		// Unsubscribe all users from deleted user
		conn.Send("MULTI")
		for _, s := range susbcribers {
			conn.Send("SREM", "tgbot:user:feeds:"+s, "d"+userID)
		}
		conn.Send("DEL", "tgbot:feed:subscribers:d"+userID)
		_, err = conn.Do("EXEC")
		if err != nil {
			return err
		}
	default:
		log.Println("accountDelete():", "Invalid payload:", payload)
		return nil
	}

	conn := pool.Get()
	defer conn.Close()
	return nil
}

func accountDeleteHandler(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	switch r.Method {
	case "POST":
		if err := accountDelete(w, r); err != nil {
			log.Println("accountDelete(): handling: ", r.RequestURI, err)
		}
	default:
		http.Error(w, "POST ONLY", http.StatusMethodNotAllowed)
	}
}
