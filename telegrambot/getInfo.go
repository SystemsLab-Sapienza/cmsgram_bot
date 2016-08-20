package main

import (
	"bytes"
	"strings"

	"github.com/garyburd/redigo/redis"
)

func getInfoNew(name string) (string, error) {
	var (
		b    bytes.Buffer
		user = struct {
			ID string

			Email      string   `redis:"email"`
			EmailAltre []string `redis:"-"`
			Nome       string   `redis:"nome"`
			Cognome    string   `redis:"cognome"`
			Indirizzo  string   `redis:"indirizzo"`
			Telefono   string   `redis:"telefono"`
			Sito       string   `redis:"sito"`
			SitoAltri  []string `redis:"-"`
		}{}
	)

	conn := pool.Get()
	defer conn.Close()

	lname := strings.ToLower(name)

	uids, err := redis.Strings(conn.Do("SMEMBERS", "webapp:users:info:"+lname))
	if err != nil && err != redis.ErrNil {
		return "", err
	}

	text := ""
	for _, u := range uids {
		data, err := redis.Values(conn.Do("HGETALL", "webapp:users:data:"+u))
		if err != nil {
			return "", err
		}

		err = redis.ScanStruct(data, &user)
		if err != nil {
			return "", err
		}

		user.ID = u

		user.EmailAltre, err = redis.Strings(conn.Do("LRANGE", "webapp:users:data:email:"+u, 0, -1))
		if err != nil && err != redis.ErrNil {
			return "", err
		}

		user.SitoAltri, err = redis.Strings(conn.Do("LRANGE", "webapp:users:data:url:"+u, 0, -1))
		if err != nil && err != redis.ErrNil {
			return "", err
		}

		if err = templates.ExecuteTemplate(&b, "get_info_new.tpl", user); err != nil {
			return "", err
		}

		text += b.String()
		b.Reset()
	}

	return text, nil
}

func getInfoLegacy(name string) (string, error) {
	type userT struct {
		Tipo     string `redis:"tipo"`
		Nome     string `redis:"nome"`
		Cognome  string `redis:"cognome"`
		Email    string `redis:"email"`
		Telefono string `redis:"telefono"`
		URL      string `redis:"url"`
	}

	var (
		lname = strings.ToLower(name)
		user  userT
		users = make([]userT, 0)
	)

	conn := pool.Get()
	defer conn.Close()

	people, err := redis.Strings(conn.Do("SMEMBERS", "webapp:docenti:"+lname))
	if err != nil {
		return "", err
	}

	for _, p := range people {
		data, err := redis.Values(conn.Do("HGETALL", "webapp:docenti:"+p))
		if err != nil {
			return "", err
		}

		err = redis.ScanStruct(data, &user)
		if err != nil {
			return "", err
		}

		users = append(users, user)
	}

	var b bytes.Buffer
	err = templates.ExecuteTemplate(&b, "get_info_legacy.tpl", users)
	if err != nil {
		return "", err
	}

	return b.String(), nil
}

func getInfo(name string) (string, error) {
	text, err := getInfoNew(name)
	if err != nil {
		return "", err
	}

	if text == "" {
		return getInfoLegacy(name)
	}

	return text, nil
}
