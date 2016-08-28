package main

import (
	"encoding/csv"
	"fmt"
	"io"
	"log"
	"os"
	"strconv"
	"time"
)

var (
	config = struct {
		BotAPIBaseURL string
		BotAPIToken   string

		RedisDomain      string
		RedisAddress     string
		RedisMaxIdle     int
		RedisIdleTimeout time.Duration

		TestRecipient    int
		WorkingDirectory string

		CrawlerEndpoint string
		WebappEndpoint  string
	}{}
)

// Set the default configuration
func init() {
	config.RedisDomain = "tcp"
	config.RedisAddress = "localhost:6379"
	config.RedisMaxIdle = 3
	config.RedisIdleTimeout = 240

	config.WorkingDirectory = "/usr/local/bin"

	config.CrawlerEndpoint = "/crawer"
	config.WebappEndpoint = "/webapp"
}

func loadConfig(filepath string) error {
	f, err := os.Open(filepath)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	r := csv.NewReader(f)
	r.Comma = ':'
	r.Comment = '#'
	r.FieldsPerRecord = 2
	r.LazyQuotes = true
	r.TrimLeadingSpace = true

	for {
		record, err := r.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatal("loadConfig():", err)
		}

		value := record[1]
		switch record[0] {
		case "botAPI_base_URL":
			config.BotAPIBaseURL = value
		case "botAPI_token":
			config.BotAPIToken = value
		case "redis_domain":
			config.RedisDomain = value
		case "redis_address":
			config.RedisAddress = value
		case "redis_max_idle":
			i, err := strconv.Atoi(value)
			if err != nil {
				fmt.Printf("redis_max_idle value '%s' not valid. Using default.\n", value)
			}

			config.RedisMaxIdle = i
		case "redis_idle_timeout":
			i, err := strconv.Atoi(value)
			if err != nil {
				fmt.Printf("redis_idle_timeout value '%s' not valid. Using default.\n", value)
			}

			config.RedisIdleTimeout = time.Duration(i)
		case "test_recipient":
			i, err := strconv.Atoi(value)
			if err != nil {
				fmt.Printf("test_recipient value '%s' not valid. Ignored.\n", value)
			}

			config.TestRecipient = i
		case "working_directory":
			config.WorkingDirectory = value
		case "webapp_endpoint":
			config.WebappEndpoint = value
		case "crawler_endpoint":
			config.CrawlerEndpoint = value
		default:
			fmt.Printf("Parameter '%s' in config file not valid. Ignored.\n", record[0])
		}
	}

	fmt.Printf("Server started with the following configuration:\n%-20s\t%s\n%-20s\t%s\n%-20s\t%s\n%-20s\t%s\n",
		"redis domain:", config.RedisDomain,
		"redis address:", config.RedisAddress,
		"Crawler endpoint:", config.CrawlerEndpoint,
		"Webapp endpont:", config.WebappEndpoint,
	)

	return err
}
