package serial

import (
	"github.com/tarm/serial"
	"log"
)

var device = "/dev/pts/10"

/* Commandlist:
	- color:(r,g,b)
	- cmd:FADE
	- cmd:STROBE
	- cmd:POWER_ON
	- cmd:POWER_OFF
	- cmd:FLASH
 */

func Serve(serialInput chan string, webInput chan string, done chan bool) {
	log.Println("# Starting Serial Listener on", device)

	c := &serial.Config{Name: device, Baud: 9600}
	s, err := serial.OpenPort(c)
	if err != nil {
		log.Fatal(err)
	}

	go readSerial(s, webInput)
	writeSerial(s, serialInput)
	done <- true
}

func writeSerial (s *serial.Port, serialInput chan string) {

	for {
		message := <-serialInput
		message += "\n"
		n, err := s.Write([]byte(message))
		if err != nil {
			log.Fatal(err)
		}
		log.Printf(string(n))
	}
}

func readSerial(s *serial.Port, webInput chan string) {
	buf := make([]byte, 128)
	for {
		n, err := s.Read(buf)
		if err != nil {
			log.Fatal(err)
		}
		webInput <- string(buf[:n])
		log.Printf("# Received Info:")
		log.Printf("%q", buf[:n])
	}
}
