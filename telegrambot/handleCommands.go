package main

import (
	"bytes"
	"log"
	"strconv"

	"github.com/garyburd/redigo/redis"
)

func handleCommands(cmd string, user int) error {
	var (
		b   bytes.Buffer
		err error
		rm  ResponseMessage
	)

	if len(cmd) > 3 {
		prefix := cmd[:3]

		if prefix == "/s_" || prefix == "/u_" {
			return handleSubscriptions(cmd, user)
		}
	}

	switch cmd {
	case "/annulla":
		lcmd, err := getLastCommand(user)
		if err != nil && err != redis.ErrNil {
			return err
		}
		if len(lcmd) == 0 {
			rm.Send("Nessun comando da annullare.", user)
			return nil
		}

		// Cancel current operation
		err = setLastCommand("", user)
		if err != nil {
			return err
		}

		rm.Send("Il comando è stato annullato.", user)
	case "/avvisi":
		yes, err := isSubscribed(user, "avvisi")
		if err != nil {
			return err
		}

		if yes {
			rm.Send("Sei già iscritto al feed.", user)
			return nil
		}

		err = newSubscription(cmd[1:], user)
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

		if err := setLastCommand(cmd, user); err != nil {
			return err
		}
	case "/feed":
		if err := setLastIndex(user, 0); err != nil {
			return err
		}

		text, i, err := listSubscriptions(user)
		if err != nil {
			log.Println("listSubscriptions():", err)
			return err
		}

		if i != 0 {
			rm.AddCallbackButton("Altro", "/feed/more/")
		}

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
	default:
		rm.Send("Comando non riconosciuto.", user)
		return nil
	}

	return nil
}
