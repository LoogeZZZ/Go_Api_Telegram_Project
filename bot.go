package main

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"io"
	"net/http"
)

func main() {
	router := mux.NewRouter()
	router.HandleFunc("/", IndexHandler)
	http.ListenAndServe("localhost:8080", router)
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

	tgtoken := "bot5453963529:AAFv-sJb6OZoKjgofFpxteqNEPYqGRGTla0"

	apiUrl := "https://api.telegram.org/" + tgtoken
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
}
