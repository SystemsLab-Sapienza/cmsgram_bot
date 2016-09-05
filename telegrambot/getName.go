package main

import (
	"errors"

	"github.com/garyburd/redigo/redis"
)

func getFullName(uid string) (string, error) {
	var name, lname string

	conn := pool.Get()
	defer conn.Close()

	// Get first and last name
	values, err := redis.Values(conn.Do("HMGET", "webapp:users:data:"+uid, "nome", "cognome"))
	if err != nil {
		return "", err
	}

	_, err = redis.Scan(values, &name, &lname)
	if err != nil {
		return "", err
	}

	return name + " " + lname, nil
}

func getFeedName(feed string) (string, error) {
	var (
		err             error
		errFeedNotValid = errors.New("Feed not valid")
		feedName        string
		ok              bool
	)

	switch feed[0] {
	case 'a':
		feedName = "avvisi"
	case 'd':
		feedName, err = getFullName(feed[1:])
		if err != nil {
			return "", err
		}
	case 't':
		// Check the feed exists in the global map
		feedName, ok = rssfeeds[feed]
		if !ok {
			return feedName, errFeedNotValid
		}
	default:
		return "", errFeedNotValid
	}

	return feedName, nil
}
