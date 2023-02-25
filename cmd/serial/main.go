package main

import (
	"fmt"
	"log"
	"os"

	"github.com/tarm/serial"
)

func main() {
	
	//
	// Get environment Variables
	//
	device, exists := os.LookupEnv("MALONE_DEVICE")

	if !exists {

		err := fmt.Errorf("MALONE_DEVICE is required")
		log.Fatal(err)
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)

	}


	for {
		var err error
		var n int

		log.Printf("[main] Open serial device: %v", device)
		// c := &serial.Config{Name: "/dev/ttyUSB0", Baud: 115200}
		// c := &serial.Config{Name: "/dev/ttys002", Baud: 115200}
		c := &serial.Config{Name: device, Baud: 115200}
		s, err := serial.OpenPort(c)
		if err != nil {
			log.Fatal(err)
		}

		// _, err = s.Write([]byte("test"))
		// if err != nil {
		// 	log.Fatal(err)
		// }

		log.Println("[main] Read...")
		buf := make([]byte, 128)
		n, err = s.Read(buf)
		if err != nil {
			log.Fatal(err)
		}
		log.Printf("%q", buf[:n])

	}

}
