package main

import (
	"strconv"

	"github.com/garyburd/redigo/redis"
)

func isSubscribed(uid int, feed string) (yes bool, err error) {
	conn := pool.Get()
	defer conn.Close()

	yes, err = redis.Bool(conn.Do("SISMEMBER", "tgbot:feed:subscribers:"+feed, uid))
	if err != nil {
		return
	}

	return
}

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
