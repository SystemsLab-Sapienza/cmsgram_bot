package main

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strconv"

	"github.com/garyburd/redigo/redis"
)

type InlineKeyboardT struct {
	Text          string `json:"text"`
	Callback_data string `json:"callback_data"`
}

type ReplyMarkupT struct {
	Inline_keyboard [][]InlineKeyboardT `json:"inline_keyboard"`
}

// This file defines the type ResponseMessage and associated methods.
type ResponseMessage struct {
	Chat_id                  int64         `json:"chat_id"`
	Text                     string        `json:"text"`
	Disable_web_page_preview bool          `json:"disable_web_page_preview"`
	Parse_mode               string        `json:"parse_mode"`
	Reply_markup             *ReplyMarkupT `json:"reply_markup,omitempty"`
}

func (rm *ResponseMessage) AddCallbackButton(text, data string) {
	var (
		buttons [][]InlineKeyboardT
		row1    []InlineKeyboardT = []InlineKeyboardT{InlineKeyboardT{text, data}}
	)
	buttons = append(buttons, row1)

	rm.Reply_markup = &ReplyMarkupT{buttons}
}

// Sends message 'text' to the the specified chat (an ID)
func (rm *ResponseMessage) Send(text string, to int) (err error) {
	var url = config.BotAPIBaseURL + config.BotAPIToken + "/sendMessage"
	conn := pool.Get()
	defer conn.Close()

	chat, err := redis.Int64(conn.Do("GET", "tgbot:user:chat:"+strconv.Itoa(to)))
	if err != nil {
		return
	}

	// Initialize message
	// *rm = ResponseMessage{chat, text, true, "HTML", &ReplyMarkupT{buttons}}
	// *rm = ResponseMessage{Chat_id: chat, Text: text, Disable_web_page_preview: true, Parse_mode: "HTML"}
	rm.Chat_id = chat
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
