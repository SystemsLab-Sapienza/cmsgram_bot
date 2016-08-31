package main

import "log"

func handleSubscriptions(cmd string, user int) error {
	var rm ResponseMessage

	prefix := cmd[:4]

	switch prefix {
	case "/s_a":
		fallthrough
	case "/s_d":
		fallthrough
	case "/s_t":
		feed := cmd[3:]

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
	case "/u_a":
		fallthrough
	case "/u_d":
		fallthrough
	case "/u_t":
		feed := cmd[3:]

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
		rm.Send("Comando non riconosciuto", user)
	}

	return nil
}
