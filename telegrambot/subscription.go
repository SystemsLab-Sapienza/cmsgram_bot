package main

import (
	"bytes"
	"sort"
	"strconv"

	"github.com/garyburd/redigo/redis"
)

func newSubscription(feed string, uid int) (err error) {
	conn := pool.Get()
	defer conn.Close()

	conn.Send("MULTI")
	conn.Send("SADD", "tgbot:feed:subscribers:"+feed, uid)
	conn.Send("SADD", "tgbot:user:feeds:"+strconv.Itoa(uid), feed)
	conn.Do("EXEC")
	if err != nil {
		return
	}

	return
}

func removeSubscription(feed string, uid int) (err error) {
	conn := pool.Get()
	defer conn.Close()

	conn.Send("MULTI")
	conn.Send("SREM", "tgbot:feed:subscribers:"+feed, uid)
	conn.Send("SREM", "tgbot:user:feeds:"+strconv.Itoa(uid), feed)
	conn.Do("EXEC")
	if err != nil {
		return
	}

	return
}

func removeSubscriptions(uid int) error {
	var userFeeds = "tgbot:user:feeds:" + strconv.Itoa(uid)

	conn := pool.Get()
	defer conn.Close()

	// Get the user's feeds
	feeds, err := redis.Strings(conn.Do("SMEMBERS", userFeeds))
	if err != nil {
		return err
	}

	if len(feeds) != 0 {
		// Unsubscribe the user from all the feeds
		conn.Send("MULTI")
		for _, f := range feeds {
			conn.Send("SREM", "tgbot:feed:subscribers:"+f, uid)
		}
		_, err = conn.Do("EXEC")
		if err != nil {
			return err
		}

		// _, err = conn.Do("SREM", redis.Args{}.Add(user).AddFlat(feeds)...)
		_, err = conn.Do("DEL", userFeeds)
		if err != nil {
			return err
		}
	}

	return nil
}

func getSubscriptions(uid int) (string, error) {
	type feed struct {
		ID   string
		Name string
	}

	var (
		b    bytes.Buffer
		data = struct {
			Feeds []feed
		}{}
	)

	conn := pool.Get()
	defer conn.Close()

	// Get the user's feed
	feeds, err := redis.Strings(conn.Do("SMEMBERS", "tgbot:user:feeds:"+strconv.Itoa(uid)))
	if err != nil {
		return "", err
	}

	if len(feeds) == 0 {
		text := "Non sei iscritto ad alcun feed."
		return text, err
	}

	sort.Strings(feeds)
	for _, f := range feeds {
		feedName, err := getFeedName(f)
		if err != nil {
			return "", err
		}

		data.Feeds = append(data.Feeds, feed{f, feedName})
	}

	if err = templates.ExecuteTemplate(&b, "feeds.tpl", data); err != nil {
		return "", err
	}
	text := b.String()

	return text, err
}
