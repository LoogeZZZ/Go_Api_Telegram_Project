package main

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/mattn/go-sqlite3"
	"io"
	"math/rand"
	"net/http"
	"strconv"
	"strings"
	"time"
)

const apiUrl = "https://api.telegram.org/" + "bot5453963529:AAFv-sJb6OZoKjgofFpxteqNEPYqGRGTla0"

func main() {

	sql.Register("sqlite3_with_extensions",
		&sqlite3.SQLiteDriver{
			Extensions: []string{
				"sqlite3_mod_regexp",
			},
		})

	_, err := sql.Open("sqlite3", "APIBOTSTATUS.sql")
	if err != nil {
		panic(err)
	}

	go UpdateLoop()
	router := mux.NewRouter()
	router.HandleFunc("/api", IndexHandler)
	router.HandleFunc("/botName", NameHandler)
	router.HandleFunc("/lastId", LastIdHandler)
	router.PathPrefix("/").Handler(http.FileServer(http.Dir("./static/")))
	router.HandleFunc("/login", IndexLogin)
	router.HandleFunc("/register", IndexRegister)
	_ = http.ListenAndServe("localhost:8000", router)
}

func IndexLogin(w http.ResponseWriter, _ *http.Request) {

}

func IndexRegister(w http.ResponseWriter, _ *http.Request) {

}

func IndexHandler(w http.ResponseWriter, _ *http.Request) {
	var R MainStru

	Ping() /// - Страница посещена

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

	R.Result.Abilites = append(R.Result.Abilites, "reacting to commands")

	respReady, err := json.Marshal(R.Result)
	if err != nil {
		panic(err)
	}

	_, _ = w.Write(respReady)

	println("НАШИ ДАННЫЕ ПРОЧИТАНЫ! ПОЛНАЯ ГОТОВНОСТЬ У НАС ГОСТИ!")

	w.Write([]byte("Вывод успешно произведён!"))
}

func NameHandler(w http.ResponseWriter, _ *http.Request) {
	db, err := sql.Open("sqlite3", "APIBOTSTATUS.sql")
	if err != nil {
		panic(err)
	}
	defer db.Close()
	var gotname string
	var resp sql.NullString // для результата
	err = db.QueryRow("SELECT name FROM bot_status").Scan(&resp)
	if err != nil {
		fmt.Println(err)
	}
	if resp.Valid { // если результат валид
		gotname = resp.String // берём оттуда обычный string
	}
	w.Write([]byte(gotname))
}

// Обращение//////////////////////////////////
//var appeal = "олежа"

func LastIdHandler(w http.ResponseWriter, _ *http.Request) {
	db, err := sql.Open("sqlite3", "APIBOTSTATUS.sql")
	if err != nil {
		panic(err)
	}
	defer db.Close()
	var gotlastid string
	var resp sql.NullString // для результата
	err = db.QueryRow("SELECT lastid FROM bot_status").Scan(&resp)
	if err != nil {
		fmt.Println(err)
	}
	if resp.Valid { // если результат валид
		gotlastid = resp.String // берём оттуда обычный string
	}
	w.Write([]byte(gotlastid))
}

func UpdateLoop() {
	db, err := sql.Open("sqlite3", "APIBOTSTATUS.sql")
	if err != nil {
		panic(err)
	}
	defer db.Close() //закрывает коннект при закрытии программы
	lastId := 0
	var nickname1 string
	err = db.QueryRow(`select name from bot_status`).Scan(&nickname1)
	if err != nil {
		fmt.Println(err)
	}

	for {
		newId := Update(lastId, &nickname1)
		if lastId != newId {
			lastId = newId
			db.Exec(`update bot_status set lastid = $1`, lastId)
		}
		time.Sleep(5 * time.Millisecond)
	}
}

func Update(lastId int, nickname *string) int {
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
		txt := strings.ToLower(ev.Message.Text)
		if txt == "/privet" {
			txtmsg := SendMessage{
				ChId:                ev.Message.Chat.Id,
				Text:                "Hello!",
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
		/////////////////////////// 22.10.22
		if strings.Split(txt, ", ")[0] == *nickname {

			switch strings.Split(strings.Split(txt, ", ")[1], ": ")[0] {
			case "расскажи анекдот":
				{
					return Anek(lastId, ev)
				}
			case "сгенерируй число ":
				{
					return RandGen(lastId, ev, txt)
				}
			case "измени обращение на":
				{
					if strings.Contains(txt, ": ") {
						return ChangeName(lastId, ev, txt, nickname)
					} else {
						fmt.Println("error")
					}
				}
			}

		}
	}
	return lastId
}

func Anek(lastId int, ev UpdateStruct) int {
	txtmsg := SendMessage{
		ChId: ev.Message.Chat.Id,
		Text: "Пьяный пьяный ежик влез на провода, током пиз**нуло пьяного ежа.",
	}

	bytemsg, _ := json.Marshal(txtmsg)
	_, err := http.Post(apiUrl+"/sendMessage", "application/json", bytes.NewReader(bytemsg))
	if err != nil {
		fmt.Println(err)
		return lastId
	} else {
		return ev.Id + 1
	}
}

func RandGen(lastId int, ev UpdateStruct, txt string) int {
	fmt.Println("Randgen")
	retotal := strings.Split(txt, "до ")[1]
	s, err := strconv.Atoi(retotal)
	if err != nil {
		panic(err)
	}
	fmt.Println(s)
	num := strconv.Itoa(rand.Intn(s))
	txtmsg := SendMessage{
		ChId: ev.Message.Chat.Id,
		Text: "Сгенерированное число: " + num,
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

func ChangeName(lastId int, ev UpdateStruct, txt string, nickname *string) int {
	new := strings.Split(txt, "измени обращение на: ")
	*nickname = new[1]
	fmt.Println(nickname)
	db, err := sql.Open("sqlite3", "APIBOTSTATUS.sql")
	if err != nil {
		panic(err)
	}
	defer db.Close()                                     //закрывает коннект при закрытии программы или выходе из зоны видимости
	db.Exec(`UPDATE bot_status set name =$1`, *nickname) // новое имя в таблицу bot_status
	txtmsg := SendMessage{
		ChId: ev.Message.Chat.Id,
		Text: "Обращение изменено на: " + *nickname,
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

func Ping() {
	txtmsg := SendMessage{
		ChId: 690215801,
		Text: "Страница посещена",
	}

	bytemsg, _ := json.Marshal(txtmsg)
	_, err := http.Post(apiUrl+"/sendMessage", "application/json", bytes.NewReader(bytemsg))
	if err != nil {
		fmt.Println(err)
	}
}

func AuthCheak() {

}
