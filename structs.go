package main

type UpdateResponse struct {
	Ok bool `json:"ok"`

	Result []UpdateStruct `json:"result"`
}

type User struct {
	Id       int    `json:"id"`
	Is_bot   bool   `json:"is_bot"`
	Username string `json:"username"`
	IsPrem   bool   `json:"is_prem"`
}

type Chat struct {
	Id   int    `json:"id"`
	Type string `json:"type"`
}

type Message struct {
	Id   int    `json:"message_id"`
	User User   `json:"from"`
	Date int    `json:"date"`
	Chat Chat   `json:"chat"`
	Text string `json:"text"`
}

type SendMessage struct {
	ChId                int    `json:"chat_id"`
	Text                string `json:"text"`
	Reply_To_Message_Id int    `json:"reply_to_message_id"`
}

type UpdateStruct struct {
	Id               int     `json:"update_id"`
	Message          Message `json:"message"`
	EditedMessage    Message `json:"edited_message"`
	ChannalPost      Message `json:"channa_lPost"`
	EditedChanelPost Message `json:"edited_chanel_post"`
}

type MainStru struct {
	Ok     bool   `json:"ok"`
	Result Result `json:"result"`
}
type Result struct {
	Id         int      `json:"id"`
	Is_bot     bool     `json:"is_Bot"`
	First_name string   `json:"first_Name"`
	Username   string   `json:"username"`
	Join       bool     `json:"can_join_groups"`
	Read       bool     `json:"can_read_all_group_messages"`
	Support    bool     `json:"supports_inline_queries"`
	Abilites   []string `json:"abilites"`
}
