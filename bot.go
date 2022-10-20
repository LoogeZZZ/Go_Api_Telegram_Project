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
	lastid := 0
	nickname := "олежа"
	for {
		lastid = Update(lastid, &nickname)
		time.Sleep(5 * time.Second)
	}
}

func Update(lastid int, nickname *string) int {
	raw, err := http.Get(apiUrl + "/getUpdates?offset=" + strconv.Itoa(lastid))
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
		txtmsg := SendMessage{
			ChId:                ev.Message.Chat.Id,
			Text:                "Непонял обратись нормально! Меня зовут " + *nickname,
			Reply_To_Message_Id: ev.Message.Id,
		}

		if strings.Contains(strings.ToLower(txt), *nickname) {
			if strings.Contains(strings.ToLower(txt), "по пивку") {
				txtmsg = SendMessage{
					ChId:                ev.Message.Chat.Id,
					Text:                "да давай",
					Reply_To_Message_Id: ev.Message.Id,
				}
			} else if strings.Contains(strings.ToLower(txt), "расскажи анекдот") {
				txtmsg = SendMessage{
					ChId:                ev.Message.Chat.Id,
					Text:                "Пьяный пьяный ежик влез на провода, током пиз**нуло пьного ежа.",
					Reply_To_Message_Id: ev.Message.Id,
				}
			} else {
				txtmsg = SendMessage{
					ChId:                ev.Message.Chat.Id,
					Text:                "Чё надо?",
					Reply_To_Message_Id: ev.Message.Id,
				}
			}
		}

		if strings.Contains(strings.ToLower(txt), "я буду называть тебя ") {
			if len(strings.SplitAfter(txt, "я буду называть тебя ")) > 1 {
				*nickname = strings.SplitAfter(txt, "я буду называть тебя ")[1]
				txtmsg = SendMessage{
					ChId:                ev.Message.Chat.Id,
					Text:                "Теперь я " + *nickname,
					Reply_To_Message_Id: ev.Message.Id,
				}
			} else {
				txtmsg = SendMessage{
					ChId:                ev.Message.Chat.Id,
					Text:                "Нормально назови",
					Reply_To_Message_Id: ev.Message.Id,
				}
			}

		}

		bytemsg, _ := json.Marshal(txtmsg)
		_, err := http.Post(apiUrl+"/sendMessage", "application/json", bytes.NewReader(bytemsg))
		if err != nil {
			fmt.Println(err)
			return lastid
		} else {
			return ev.Id + 1
		}

	}

	return lastid
}
