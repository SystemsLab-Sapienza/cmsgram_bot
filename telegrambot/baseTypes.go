package main

type From struct {
	ID         int
	First_name string
}

type Chat struct {
	ID         int64
	First_name string
	Type       string
}

type CallbackQuery struct {
	ID      string
	From    *From
	Message *Message
	Data    string
}

type Message struct {
	Message_id int
	From       *From
	Date       int
	Chat       *Chat
	Text       string
}

type Update struct {
	Update_id      int
	Message        *Message
	Callback_query *CallbackQuery
}

type Response struct {
	Ok     bool
	Result []Update
}

type InlineKeyboardButton struct {
	Text          string
	Callback_data string
}
