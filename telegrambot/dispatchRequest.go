package main

import (
	"bytes"
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
		b   bytes.Buffer
		err error
		rm  ResponseMessage
	)

	if u.Callback_query != nil {
		var (
			msg    = u.Callback_query.Message
			newmsg EditedMessage
		)

		newmsg.Chat_id = msg.Chat.ID
		newmsg.Message_id = msg.Message_id

		switch {
		case strings.HasPrefix(u.Callback_query.Data, "/twiki/more/"):
			const incr = 5
			var newindex int

			// Get the index
			n, err := strconv.Atoi(u.Callback_query.Data[12:])
			if err != nil {
				log.Println("strconv.Atoi():", err)
				return err
			}

			if n+incr >= nfeeds {
				newindex = nfeeds + 1
				newmsg.AddCallbackButton("Inizio", "/twiki/more/1")
			} else {
				newindex = incr + n
				newmsg.AddCallbackButton("Altro", "/twiki/more/"+strconv.Itoa(newindex))
			}

			text, err := listFeeds(n, newindex)
			if err != nil {
				log.Println("listFeeds():", err)
				return err
			}

			// Send edited message with next index
			newmsg.Send(text, msg.From.ID)
		default:
			log.Println("dispatchRequest(): invalid callback:", u.Callback_query.Data)
		}

		return nil
	}

	if u.Message.Text[:1] == "/" {
		u.Message.Text = strings.ToLower(u.Message.Text)
	}

	cmd := u.Message.Text[1:]
	chat := u.Message.Chat.ID
	user := u.Message.From.ID

	conn := pool.Get()
	defer conn.Close()
	_, err = conn.Do("SET", "tgbot:user:chat:"+strconv.Itoa(user), chat)
	if err != nil {
		log.Println(err)
		return err
	}

	// Handle commands for (un)subscribing to specific feeds
	switch {
	case strings.HasPrefix(u.Message.Text, "/c_all"):
		removeSubscriptions(user)
		return nil

	// Subscribe commands
	case strings.HasPrefix(u.Message.Text, "/s_a"):
		fallthrough
	case strings.HasPrefix(u.Message.Text, "/s_d"):
		fallthrough
	case strings.HasPrefix(u.Message.Text, "/s_t"):
		feed := u.Message.Text[3:]

		yes, err := isSubscribed(user, feed)
		if err != nil {
			return err
		}

		if yes {
			rm.Send("Sei già iscritto al feed.", user)
			return nil
		}

		text, err := subscribeFeed(feed, user)
		if err != nil {
			log.Println("subscribeFeed():", err)
			return err
		}

		rm.Send(text, user)
		return nil

	// Unsubscribe commands
	case strings.HasPrefix(u.Message.Text, "/u_a"):
		fallthrough
	case strings.HasPrefix(u.Message.Text, "/u_d"):
		fallthrough
	case strings.HasPrefix(u.Message.Text, "/u_t"):
		feed := u.Message.Text[3:]

		yes, err := isSubscribed(user, feed)
		if err != nil {
			return err
		}

		if !yes {
			rm.Send("Non sei iscritto a questo feed.", user)
			return nil
		}

		text, err := unsubscribeFeed(feed, user)
		if err != nil {
			log.Println("unsubscribeFeed():", err)
			return err
		}

		rm.Send(text, user)
		return nil
	default:
	}

	// Handle generic commands
	switch u.Message.Text {
	case "/annulla":
		lcmd, err := getLastCommand(user)
		if err != nil {
			return err
		}
		if lcmd == "" {
			rm.Send("Nessun comando da annullare.", user)
			return nil
		}

		// Cancel current operation
		rm.Send("Il comando è stato annullato.", user)
		err = setLastCommand("", user)
		if err != nil {
			return err
		}
	case "/avvisi":
		yes, err := isSubscribed(user, "avvisi")
		if err != nil {
			return err
		}

		if yes {
			rm.Send("Sei già iscritto al feed.", user)
			return nil
		}

		err = newSubscription(cmd, user)
		if err != nil {
			return err
		}
		rm.Send("Sei ora iscritto agli avvisi.", user)
	case "/cancella":
		err = removeSubscriptions(user)
		if err != nil {
			return err
		}
		rm.Send("Non segui più alcun feed.", user)
	case "/cerca":
		rm.Send("Scrivi il cognome del docente:", user)
		err = setLastCommand("/cerca", user)
		if err != nil {
			return err
		}
	case "/feed":
		text, err := getSubscriptions(user)
		if err != nil {
			log.Println("getSubscriptions():", err)
			return err
		}

		// rm.AddCallbackButton("Altro", "/feed/more/10")
		rm.Send(text, user)
	case "/help", "/start":
		if err = templates.ExecuteTemplate(&b, "start.tpl", nil); err != nil {
			return err
		}
		rm.Send(b.String(), user)
	case "/twiki":
		const incr = 5

		newindex := incr + 1
		text, err := listFeeds(1, newindex)
		if err != nil {
			log.Println("listFeeds():", err)
			return err
		}

		rm.AddCallbackButton("Altro", "/twiki/more/"+strconv.Itoa(newindex))
		rm.Send(text, user)
	// case "/segreteria":
	// 	err := newSubscription(conn, cmd, user)
	// 	if err != nil {
	// 		return err
	// 	}

	// 	if ok {
	// 		rm.Send("Sei ora iscritto al feed della segreteria.", user)
	// 	} else {
	// 		rm.Send("Sei già iscritto al feed della segreteria.", user)
	// 	}
	// If we're here the user has provided input for a previous command
	// or has entered a wrong commands
	default:
		lcmd, err := getLastCommand(user)
		if err != nil {
			return err
		}
		if lcmd == "" {
			rm.Send("Comando non riconosciuto.", user)
			return nil
		}

		switch lcmd {
		case "/cerca":
			res, err := getInfo(u.Message.Text)
			if err != nil {
				return err
			}
			if res == "" {
				rm.Send("Nessun docente trovato.", user)
				return nil
			}

			rm.Send(res, user)
		default:
			rm.Send("Comando non riconosciuto.", user)
			return nil
		}

		// Reset last command
		err = setLastCommand("", user)
		if err != nil {
			return err
		}
	}

	return nil
}
