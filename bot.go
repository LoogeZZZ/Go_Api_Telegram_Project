package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"io"
	"net/http"
	"strconv"
	"strings"
	"time"
)

const apiUrl = "https://api.telegram.org/" + "bot5453963529:AAFv-sJb6OZoKjgofFpxteqNEPYqGRGTla0"

func main() {
	go UpdateLoop()
	router := mux.NewRouter()
	router.HandleFunc("/", IndexHandler)
	http.ListenAndServe("localhost:8080", router)
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
		if txt == "/easter_egg" {
			txtmsg := SendMessage{
				ChId:                ev.Message.Chat.Id,
				Text:                "https://www.youtube.com/watch?v=lIxM2rGKEV4",
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
		var Bot_Name = "Олежа"
		txt1 := ev.Message.Text

		if strings.Contains(txt1, Bot_Name) {
			if strings.Contains(txt1, "Расскажи анекдот") {
			}
			txtmsg := SendMessage{
				ChId:                ev.Message.Chat.Id,
				Text:                "Пьяный пьяный ежик влез на провода, током пиз**нуло пьного ежа.",
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
		txt2 := ev.Message.Text
		if strings.Contains(txt2, Bot_Name) {
			if strings.Contains(txt2, "кто ты?") {
			}
			txtmsg := SendMessage{
				ChId:                ev.Message.Chat.Id,
				Text:                "Зовут Олежа, немного о себе. Парень сипотяга, по жизни бродяга, походка городская, жизнь воровскааая",
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
