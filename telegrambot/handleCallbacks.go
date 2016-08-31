package main

import (
	"log"
	"strconv"
	"strings"
)

func handleCallbacks(cbq *CallbackQueryT) error {
	var (
		msg    = cbq.Message
		newmsg EditedMessage
	)

	newmsg.Chat_id = msg.Chat.ID
	newmsg.Message_id = msg.Message_id

	switch {
	case strings.HasPrefix(cbq.Data, "/twiki/more/"):
		const incr = 5
		var newindex int

		// Get the index
		n, err := strconv.Atoi(cbq.Data[12:])
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
		log.Println("dispatchRequest(): invalid callback:", cbq.Data)
	}

	return nil
}
