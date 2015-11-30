package web

import (
	"encoding/json"
	"github.com/bitly/go-simplejson"
	"github.com/gorilla/websocket"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"
	"flag"
)

var socketUpgrader = websocket.Upgrader{
	CheckOrigin:     func(r *http.Request) bool { return true },
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}
const (
	writeWait = 10 * time.Second
	pongWait = 60 * time.Second
	pingPeriod = (pongWait * 9) / 10
	maxMessageSize = 512
)

type Socket struct {
	connection *websocket.Conn
	write chan string
	read chan string
}

type SocketRegistry struct {
	items []*Socket
	count int
}

var Registry *SocketRegistry

var readChannel chan string

func Serve(webWriter chan string, webReader chan string, done chan bool) {
	addr := flag.String("addr", "localhost:8081", "*")
	readChannel = webReader

	Registry = &SocketRegistry{}
	go Registry.notifyListener(webWriter)

	log.Println("# Listening on", *addr)
	staticHandler := http.FileServer(http.Dir("static"))
	http.Handle("/static/", http.StripPrefix("/static/", staticHandler))
	http.HandleFunc("/", socketHandler)
	err := http.ListenAndServe(*addr, nil)

	if err != nil {
		log.Fatal("! Error starting HTTP server: ", err)
	}

	done <- true
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
	connection.SetReadLimit(maxMessageSize)
	connection.SetReadDeadline(time.Now().Add(pongWait))
	connection.SetPongHandler(func(string) error { connection.SetReadDeadline(time.Now().Add(pongWait)); return nil })

	if err != nil {
		log.Print("! Upgrade Error:", err)
		return
	}

	defer connection.Close()

	write := make(chan string)

	s := &Socket{
		connection,
		write,
		readChannel,
	}

	Registry.addSocket(s)
	go s.writeSocket()
	s.readSocket()
}

func (r *SocketRegistry) addSocket (s *Socket) {
	r.items = append(r.items, s)
	r.count++
	log.Println("# Added Socket No.", r.count)
}

func (r *SocketRegistry) notifyListener(in chan string) {
	for {
		message := <- in
		log.Printf("# Broadcasting Message to %d sockets...", r.count)
		for index, item := range r.items {
			log.Printf("# Sent Message to Socket %d: %s", index, message)
			item.write <- message
		}
	}
}

func (s *Socket) readSocket() {
	for {
		_, payload, err := s.connection.ReadMessage()

		if err != nil {
			log.Println("! Error during read:", err)
			break
		}
		message := parseIncoming(payload)
		log.Println(message)
		s.read <- message
	}
}

func (s *Socket) writeSocket() {

	ticker := time.NewTicker(pingPeriod)

	defer func() {
		ticker.Stop()
		s.connection.Close()
	}()

	for {
		select {
		case message, ok := <-s.write:
			if !ok {
				s.writeMessage(websocket.CloseMessage, []byte{})
				return
			}

			json := parseOutgoing(message)
			payload, err := json.Encode()

			if err != nil {
				log.Fatal("! Json Encode Error:", err)
				break
			}

			if err := s.writeMessage(websocket.TextMessage, payload); err != nil {
				log.Fatal(err)
				return
			}

			log.Println("# Sending Message: ", message)

		case <-ticker.C:
			if err := s.writeMessage(websocket.PingMessage, []byte{}); err != nil {
				return
			}
		}
	}

}

func (s *Socket) writeMessage(mt int, payload []byte) error {
	s.connection.SetWriteDeadline(time.Now().Add(writeWait))
	return s.connection.WriteMessage(mt, payload)
}

func parseIncoming(payload []byte) string {
	js, err := simplejson.NewJson(payload)

	if err != nil {
		log.Println("!JSON Unmarshal Error:", err)
		return ""
	}

	event := js.Get("event").MustString()

	if event == "color" {
		values := js.Get("data").MustArray()
		return "color:(" + string(values[0].(json.Number)) + "," + string(values[1].(json.Number)) + "," + string(values[2].(json.Number)) + ")"
	}

	if event == "cmd" {
		values := js.Get("data").MustString()
		return "cmd:" + values
	}

	return ""
}

func parseOutgoing(payload string) *simplejson.Json {
	command := string(payload)

	index := strings.Index(command, ":")
	event := command[:index]

	document := simplejson.New()
	document.Set("event", event)

	if event == "cmd" {
		document.Set("data", command[index+1:len(command)-1])
	}

	if event == "color" {
		values := strings.Split(command[index+2:len(command)-2], ",")
		colors := [3]int{}

		for i := 0; i < 3; i++ {
			colors[i], _ = strconv.Atoi(values[i])
		}

		document.Set("data", colors)
	}

	return document
}
