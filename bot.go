package main

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"io"
	"net/http"
)

const apiUrl = "https://api.telegram.org/" + "bot5453963529:AAFv-sJb6OZoKjgofFpxteqNEPYqGRGTla0"

func main() {
	router := mux.NewRouter()
	router.HandleFunc("/", IndexHandler)
	http.ListenAndServe("localhost:8080", router)
}

type UpdateResponse struct {
	Data []UpdateStruct `json:"data"`
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
	Id         int    `json:"id"`
	Is_bot     bool   `json:"is_Bot"`
	First_name string `json:"first_Name"`
	Username   string `json:"username"`
	Join       bool   `json:"can_join_groups"`
	Read       bool   `json:"can_read_all_group_messages"`
	Support    bool   `json:"supports_inline_queries"`
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

	respReady, err := json.Marshal(R.Result)
	if err != nil {
		panic(err)
	}

	w.Write([]byte(respReady))

	println("НАШИ ДАННЫЕ ПРОЧИТАНЫ! ПОЛНАЯ ГОТОВНОСТЬ У НАС ГОСТИ!")

	w.Write([]byte("Вывод успешно произведён!"))
}

func Update() {
	raw, err := http.Get(apiUrl + "/getUpdates")
	if err != nil {
		panic(err)
	}
	body, _ := io.ReadAll(raw.Body)

	var v []interface{}
	err = json.Unmarshal(body, &v)
	if err != nil {
		panic(err)
	}

	for _, ev := range v {
		t := ev.(UpdateStruct)
		txt := t.Message.Text
		if txt == "/privet" {

			http.Post(apiUrl+"/sendMessage", "application/json", nil)
		}
	}

}
