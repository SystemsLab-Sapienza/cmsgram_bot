package main

import (
	"log"
	"strconv"

	"github.com/garyburd/redigo/redis"
)

// rssfeeds maps a feed ID to its name
var (
	nfeeds   int
	rssfeeds map[string]string
)

// initRSS initializes the global map rssfeeds
func initRSS() error {
	conn := pool.Get()
	defer conn.Close()

	n, err := redis.Int(conn.Do("GET", "rss:feed:counter"))
	if err != nil {
		log.Println("initRSS():", err)
		return err
	}

	// Initialize global index
	nfeeds = n

	// Initialize map
	rssfeeds = make(map[string]string)

	// Map feeds' IDs to their respective name
	for i := 1; i <= nfeeds; i++ {
		id := strconv.Itoa(i)
		name, err := redis.String(conn.Do("HGET", "rss:feed:"+id, "name"))
		if err != nil {
			log.Println("initRSS():", err)
			return err
		}

		rssfeeds["t"+id] = name
	}

	return nil
}
