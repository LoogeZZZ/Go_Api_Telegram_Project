package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"io"
	"net/http"
	"strconv"
	"time"
)

const apiUrl = "https://api.telegram.org/" + "bot5453963529:AAFv-sJb6OZoKjgofFpxteqNEPYqGRGTla0"

func main() {
	go UpdateLoop()
	router := mux.NewRouter()
	router.HandleFunc("/", IndexHandler)
	http.ListenAndServe("localhost:8080", router)
}

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

func IndexHandler(w http.ResponseWriter, _ *http.Request) {
	var R MainStru

	resp, err := http.Get(apiUrl + "/getMe")

	if err != nil {
		fmt.Println(err)
	}
	respBody, _ := io.ReadAll(resp.Body)
	fmt.Println(string(respBody))

	err = json.Unmarshal(respBody, &R) // заполнили перемнную р
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
	}

	R.Result.Abilites = append(R.Result.Abilites, "reacting to command /privet")

	respReady, err := json.Marshal(R.Result)
	if err != nil {
		panic(err)
	}

	w.Write([]byte(respReady))

	println("НАШИ ДАННЫЕ ПРОЧИТАНЫ! ПОЛНАЯ ГОТОВНОСТЬ У НАС ГОСТИ!")

	w.Write([]byte("Вывод успешно произведён!"))
}

func UpdateLoop() {
	lastId := 0
	for {
		lastId = Update(lastId)
		time.Sleep(5 * time.Second)
	}
}

func Update(lastId int) int {
	raw, err := http.Get(apiUrl + "/getUpdates?offset=" + strconv.Itoa(lastId))
	if err != nil {
		panic(err)
	}
	body, _ := io.ReadAll(raw.Body)

	var v UpdateResponse
	err = json.Unmarshal(body, &v)
	if err != nil {
		panic(err)
	}

	if len(v.Result) > 0 {
		ev := v.Result[len(v.Result)-1]
		txt := ev.Message.Text
		if txt == "/privet" {
			txtmsg := SendMessage{
				ChId:                ev.Message.Chat.Id,
				Text:                "ИДИ ОТ СЮДА, ЧИТАЙ ОПИСАНИЕ!",
				Reply_To_Message_Id: ev.Message.Id,
			}

			bytemsg, _ := json.Marshal(txtmsg)
			_, err = http.Post(apiUrl+"/sendMessage", "application/json", bytes.NewReader(bytemsg))
			if err != nil {
				fmt.Println(err)
				return lastId
			} else {
				return ev.Id + 1
			}

		}
	}
	return lastId
}
