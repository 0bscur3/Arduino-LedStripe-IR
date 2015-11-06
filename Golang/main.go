package main

import (
	"LedServer/serial"
	"LedServer/web"
	"log"
)

func main() {
	log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)
	log.Println("=== LED Server v0.1 ===")
	log.Println("# Starting Server ...")
	webWriter := make(chan string)
	webReader := make(chan string)
	done1 := make(chan bool)
	done2 := make(chan bool)
	go web.Serve(webWriter, webReader, done1)
	go serial.Serve(webWriter, webReader, done2)

	for i := 0; i < 2; i++ {
		select {
		case <-done1:
			log.Println("# HTTP Routine ended!")
		case <-done2:
			log.Println("# Serial Routine ended!")
		}
	}

	log.Println("# Finishing, good bye!")
}
