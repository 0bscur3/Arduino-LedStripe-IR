package web

import (
	"github.com/gorilla/websocket"
	"log"
	"net/http"
	"github.com/bitly/go-simplejson"
	"encoding/json"
)

var socketUpgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool { return true },
}

func socketHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.Error(w, "Not found", 404)
		return
	}

	if r.Method != "GET" {
		http.Error(w, "Method not allowed", 405)
		return
	}

	connection, err := socketUpgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Print("upgrade:", err)
		return
	}
	defer connection.Close()

	go readSocket(connection, out)
	writeSocket(connection, in)
}

func readSocket(connection *websocket.Conn, out chan string) {
	for {
		_, payload, err := connection.ReadMessage()

		if err != nil {
			log.Println("read:", err)
			break
		}

		message := parseIncoming(payload)
		out <- message
		log.Println(message)
	}
}

func writeSocket(connection *websocket.Conn, in chan string) {

	for {
		message := <-in
		err := connection.WriteMessage(websocket.TextMessage, []byte(message))

		if err != nil {
			log.Println("write:", err)
			break
		}

		log.Println("Sending Message: %s", message)
	}
}

func parseIncoming(payload []byte) string {
	js, err := simplejson.NewJson(payload)

	if err != nil {
		log.Println("json unmarshal:", err)
		return ""
	}

	event := js.Get("event").MustString()

	if event == "color" {
		values := js.Get("data").MustArray()
		return "color:("+string(values[0].(json.Number))+","+string(values[1].(json.Number))+","+string(values[2].(json.Number))+")"
	}

	if event == "cmd" {
		values := js.Get("data").MustString()
		return "cmd:"+values
	}

	return ""
}