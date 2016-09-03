package main

import (
	"log"
	"strconv"
	"strings"

	"github.com/garyburd/redigo/redis"
)

func setLastCommand(cmd string, uid int) (err error) {
	conn := pool.Get()
	defer conn.Close()

	_, err = conn.Do("SET", "tgbot:last_command:"+strconv.Itoa(uid), cmd, "EX", 60*10)

	return
}

func getLastCommand(uid int) (cmd string, err error) {
	conn := pool.Get()
	defer conn.Close()

	cmd, err = redis.String(conn.Do("GET", "tgbot:last_command:"+strconv.Itoa(uid)))

	return
}

func dispatchRequest(u Update) error {
	var (
		cmd string
	)

	if u.Callback_query != nil {
		return handleCallbacks(u.Callback_query)
	}

	chat := u.Message.Chat.ID
	user := u.Message.From.ID

	conn := pool.Get()
	defer conn.Close()

	_, err := conn.Do("SET", "tgbot:user:chat:"+strconv.Itoa(user), chat)
	if err != nil {
		log.Println(err)
		return err
	}

	if u.Message.Text[:1] == "/" {
		cmd = strings.ToLower(u.Message.Text)
	}

	if len(cmd) != 0 {
		if err := handleCommands(cmd, user); err != nil {
			log.Println(err)
		}
	} else {
		if err := handleInput(u.Message.Text, user); err != nil {
			log.Println(err)
		}
	}

	return nil
}
