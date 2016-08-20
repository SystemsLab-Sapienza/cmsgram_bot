package main

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
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

func (rm *EditedMessage) AddCallbackButton(text, data string) {
	var (
		buttons [][]InlineKeyboardT
		row1    []InlineKeyboardT = []InlineKeyboardT{InlineKeyboardT{text, data}}
	)
	buttons = append(buttons, row1)

	rm.Reply_markup = &ReplyMarkupT{buttons}
}

// Sends message 'text' to the the specified chat (an ID)
func (rm *EditedMessage) Send(text string, to int) (err error) {
	var url = config.BotAPIBaseURL + config.BotAPIToken + "/editMessageText"
	conn := pool.Get()
	defer conn.Close()

	// chat, err := redis.Int64(conn.Do("GET", "tgbot:user:chat:"+strconv.Itoa(to)))
	// if err != nil {
	// 	return
	// }

	// Initialize message
	rm.Text = text
	rm.Disable_web_page_preview = true
	rm.Parse_mode = "HTML"

	// Encode data into JSON
	payload, err := json.Marshal(rm)
	if err != nil {
		return
	}

	// Send the payload to the BotAPI
	res, err := http.Post(url, "application/json", bytes.NewBuffer(payload))
	if err != nil {
		return
	}

	// TODO check response to be valid
	_, err = ioutil.ReadAll(res.Body)
	if err != nil {
		return
	}
	res.Body.Close()

	return
}
