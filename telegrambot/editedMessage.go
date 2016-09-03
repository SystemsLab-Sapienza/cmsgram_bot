package main

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
)

type EditedMessage struct {
	Chat_id                  int64         `json:"chat_id"`
	Message_id               int           `json:"message_id"`
	Text                     string        `json:"text"`
	Parse_mode               string        `json:"parse_mode"`
	Disable_web_page_preview bool          `json:"disable_web_page_preview"`
	Reply_markup             *ReplyMarkupT `json:"reply_markup,omitempty"`
}

func (em *EditedMessage) AddCallbackButton(text, data string) {
	var (
		buttons [][]InlineKeyboardT
		row1    []InlineKeyboardT = []InlineKeyboardT{{text, data}}
	)
	buttons = append(buttons, row1)

	em.Reply_markup = &ReplyMarkupT{buttons}
}

// Sends message 'text' to the the specified chat (an ID)
func (em *EditedMessage) Send(text string, to int) (err error) {
	var (
		response = struct {
			Ok     bool
			Result MessageT
		}{}
		url = config.BotAPIBaseURL + config.BotAPIToken + "/editMessageText"
	)

	conn := pool.Get()
	defer conn.Close()

	// Initialize message
	em.Text = text
	em.Disable_web_page_preview = true
	em.Parse_mode = "HTML"

	// Encode data into JSON
	payload, err := json.Marshal(em)
	if err != nil {
		return
	}

	// Send the payload to the BotAPI
	res, err := http.Post(url, "application/json", bytes.NewBuffer(payload))
	if err != nil {
		return
	}

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return
	}
	res.Body.Close()

	// Decode the JSON payload
	err = json.Unmarshal(body, &response)
	if err != nil {
		log.Println("Send(): json.Unmarshal():", err)
		return
	}

	if !response.Ok {
		log.Println("Send(): Invalid request", response)
		return
	}

	return
}
