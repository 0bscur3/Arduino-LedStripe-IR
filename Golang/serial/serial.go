package serial

import (
	"github.com/tarm/serial"
	"log"
	"os"
)

var device string;

/* Commandlist:
- color:(r,g,b)
- cmd:FADE
- cmd:STROBE
- cmd:POWER_ON
- cmd:POWER_OFF
- cmd:FLASH
*/

func Serve(webWriter chan<- string, webReader <-chan string, done chan bool) {

	device := os.Args[1]
	log.Println("# Starting Serial Listener on", device)

	c := &serial.Config{Name: device, Baud: 9600}
	s, err := serial.OpenPort(c)
	if err != nil {
		log.Fatal(err)
	}

	go readSerial(s, webWriter)
	writeSerial(s, webReader)
	done <- true
}

func writeSerial(s *serial.Port, input <-chan string) {

	for {
		message := <-input
		message += "\n"
		n, err := s.Write([]byte(message))
		if err != nil {
			log.Fatal(err)
		}
		log.Printf(string(n))
	}
}

func readSerial(s *serial.Port, output chan<- string) {
	buf := make([]byte, 128)
	for {
		n, err := s.Read(buf)
		if err != nil {
			log.Fatal(err)
		}
		output <- string(buf[:n])
		log.Printf("# Received Info:")
		log.Printf("%q", buf[:n])
	}
}
