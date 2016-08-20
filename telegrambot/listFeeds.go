package main

import (
	"bytes"
	"strconv"

	"github.com/garyburd/redigo/redis"
)

func listFeeds(i, max int) (string, error) {
	var (
		b    bytes.Buffer
		feed = struct {
			ID   int    `redis:"-"`
			URL  string `redis:"url"`
			Kind string `redis:"kind"`
			Name string `redis:"name"`
		}{}
	)

	conn := pool.Get()
	defer conn.Close()

	text := ""
	for ; i < max && i <= nfeeds; i++ {
		feed.ID = i

		data, err := redis.Values(conn.Do("HGETALL", "rss:feed:"+strconv.Itoa(i)))
		if err != nil {
			return "", err
		}

		err = redis.ScanStruct(data, &feed)
		if err != nil {
			return "", err
		}

		err = templates.ExecuteTemplate(&b, "twiki.tpl", feed)
		if err != nil {
			return "", err
		}

		text += b.String()
		b.Reset()
	}

	return text, nil
}
