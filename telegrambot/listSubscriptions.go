package main

import (
	"bytes"
	"sort"
	"strconv"

	"github.com/garyburd/redigo/redis"
)

func setLastIndex(user, index int) error {
	conn := pool.Get()
	defer conn.Close()

	_, err := conn.Do("SET", "tgbot:user:feeds:index:"+strconv.Itoa(user), index)
	if err != nil {
		return err
	}

	return nil
}

func listSubscriptions(uid int) (string, int, error) {
	type feed struct {
		ID   string
		Name string
	}

	var (
		b    bytes.Buffer
		data = struct {
			Feeds []feed
		}{}
		i     int
		last  int
		feeds []string
		user  = strconv.Itoa(uid)
	)

	conn := pool.Get()
	defer conn.Close()

	last, err := redis.Int(conn.Do("GET", "tgbot:user:feeds:index:"+user))
	if err != nil && err != redis.ErrNil {
		return "", 0, err
	}

	if last == 0 {
		feeds, err = redis.Strings(conn.Do("SMEMBERS", "tgbot:user:feeds:"+user))
		if err != nil {
			return "", 0, err
		}

		conn.Send("MULTI")
		conn.Send("DEL", "tgbot:user:feeds:cached:"+user)
		conn.Send("SADD", redis.Args{}.Add("tgbot:user:feeds:cached:"+user).AddFlat(feeds)...)
		conn.Do("EXEC")
		if err != nil {
			return "", 0, err
		}
	} else {
		feeds, err = redis.Strings(conn.Do("SMEMBERS", "tgbot:user:feeds:cached:"+user))
		if err != nil {
			return "", 0, err
		}
	}

	if len(feeds) == 0 {
		text := "Non sei iscritto ad alcun feed."
		return text, 0, err
	}

	sort.Strings(feeds)
	for i = last; i < len(feeds) && i < last+5; i++ {
		feedName, err := getFeedName(feeds[i])
		if err != nil {
			return "", 0, err
		}

		data.Feeds = append(data.Feeds, feed{feeds[i], feedName})
	}

	if err = templates.ExecuteTemplate(&b, "feeds.tpl", data); err != nil {
		return "", 0, err
	}

	if i == len(feeds) {
		i = 0
	}

	if err := setLastIndex(uid, i); err != nil {
		return "", 0, err
	}

	return b.String(), i, nil
}
