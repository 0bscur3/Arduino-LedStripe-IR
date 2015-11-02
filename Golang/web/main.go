package web

import (
	"net/http"
	"log"
	"flag"
)

var out chan string
var in chan string

func Serve(serialInput chan string, webInput chan string, done chan bool) {
	addr := flag.String("addr", "localhost:8081", "*")
	in = webInput
	out = serialInput

	log.Println("# Listening on", *addr)
	staticHandler := http.FileServer(http.Dir("static"))
	http.Handle("/static/", http.StripPrefix("/static/", staticHandler))
	http.HandleFunc("/", socketHandler)
	err := http.ListenAndServe(*addr, nil)

	if err != nil {
		log.Fatal("# Error starting HTTP server: ", err)
	}

	done <- true
}

