package main

import (
	"log"
	"LedServer/web"
	"LedServer/serial"
)

func main() {
	log.SetFlags(log.Ldate|log.Ltime|log.Lshortfile)
	log.Println("=== LED Server v0.1 ===")
	log.Println("# Starting Server ...")
	serialInput := make(chan string)
	webInput := make(chan string)
	done1 := make(chan bool)
	done2 := make(chan bool)
	go web.Serve(serialInput, webInput, done1)
	go serial.Serve(serialInput, webInput, done2)

	for i := 0; i < 2; i++ {
		select {
		case <- done1:
			log.Println("# HTTP Routine ended!")
		case <- done2:
			log.Println("# Serial Routine ended!")
		}
	}

	log.Println("# Finishing, good bye!")
}
