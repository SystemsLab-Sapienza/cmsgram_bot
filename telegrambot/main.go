package main

import (
	"flag"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"time"

	"bitbucket.org/ansijax/rfidlab_telegramdi_backend/auth"

	"github.com/garyburd/redigo/redis"
)

var (
	pool      *redis.Pool
	templates *template.Template

	flagNewToken   bool
	flagUsePolling bool
)

func init() {
	flag.BoolVar(&flagNewToken, "g", false, "Generates a new pseudorandom token.")
	flag.BoolVar(&flagUsePolling, "p", false, "Tells the bot to interface with the Telegram Bot API through polling.")

	// Create a thread-safe connection pool for redis
	pool = &redis.Pool{
		MaxIdle:     config.RedisMaxIdle,
		IdleTimeout: config.RedisIdleTimeout * time.Second,
		Dial: func() (redis.Conn, error) {
			c, err := redis.Dial(config.RedisDomain, config.RedisAddress)
			if err != nil {
				return nil, err
			}

			return c, err
		},
	}
}

func main() {
	flag.Parse()

	if flagNewToken {
		fmt.Println("Token->%s", auth.NewBase36(32))
		return
	}

	if len(os.Args) == 2 {
		fmt.Println("No config file provided, exiting.")
		return
	}

	loadConfig(os.Args[2])

	// Change root directory
	if err := os.Chdir(config.WorkingDirectory); err != nil {
		log.Fatal(err)
	}

	// Parse templates
	templates = template.Must(template.ParseFiles(
		"templates/feeds.tpl",
		"templates/get_info_legacy.tpl",
		"templates/get_info_new.tpl",
		"templates/news.tpl",
		"templates/rss_update.tpl",
		"templates/start.tpl",
		"templates/subscribe.tpl",
		"templates/twiki.tpl",
		"templates/unsubscribe.tpl"))

	// Set up endpoints
	http.HandleFunc(config.CrawlerEndpoint, broadcastUpdateHandler)
	http.HandleFunc(config.WebappEndpoint, sendMessageHandler)

	err := initRSS()
	if err != nil {
		return
	}

	if flagUsePolling {
		go polling()
	} else {
		return
	}

	err = http.ListenAndServe(":8443", nil)
	if err != nil {
		log.Fatal("ListenAndServe:", err)
	}
}
