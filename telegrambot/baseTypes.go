package main

type FromT struct {
	ID         int
	First_name string
}

type ChatT struct {
	ID         int64
	First_name string
	Type       string
}

type CallbackQueryT struct {
	ID      string
	From    *FromT
	Message *MessageT
	Data    string
}

type MessageT struct {
	Message_id int
	From       *FromT
	Date       int
	Chat       *ChatT
	Text       string
}

type Update struct {
	Update_id      int
	Message        *MessageT
	Callback_query *CallbackQueryT
}

type ResponseT struct {
	Ok     bool
	Result []Update
}

type InlineKeyboardButtonT struct {
	Text          string
	Callback_data string
}
