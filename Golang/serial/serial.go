package serial

import (
	"github.com/tarm/serial"
	"log"
	"os"
	"regexp"
	"strings"
)

var device string;
var commandRegex = `(cmd:([A-Z\_]{1,10});)|(color:\(([0-9]{1,3}),([0-9]{1,3}),([0-9]{1,3})\);)`

/* Commandlist:
- color:(r,g,b)
- cmd:FADE
- cmd:STROBE
- cmd:POWER_ON
- cmd:POWER_OFF
- cmd:FLASH
*/

func Serve(webWriter chan<- string, webReader <-chan string, done chan bool) {

	device = os.Args[1]
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
		log.Println(message)
		n, err := s.Write([]byte(message))
		if err != nil {
			log.Fatal(err)
		}
		log.Printf(string(n))
	}
}

func readSerial(s *serial.Port, output chan<- string) {
	buf := make([]byte, 128)
	re := regexp.MustCompile(commandRegex)
	var message string

	for {
		n, err := s.Read(buf)
		if err != nil {
			log.Fatal(err)
		}

		message += string(buf[:n])
		message = strings.Replace(message, "\r", "", -1)
		message = strings.Replace(message, "\n", "", -1)

		results := re.FindStringSubmatch(message)

		if (0 < len(results)) {
			log.Println(results[0])
			output <- results[0]
			message = ""
		}
	}
}
