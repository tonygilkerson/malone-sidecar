package main

import (
	"log"
	"net/http"
	"os"

	"github.com/tarm/serial"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)


//
// Main
//
func main() {

	// Log to the console with date, time and filename prepended
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	
	//
	// Get environment Variables
	//
	serialPort, exists := os.LookupEnv("MALONE_SERIAL_PORT")
	if !exists {
		log.Fatalln("MALONE_SERIAL_PORT environment variable not set")
	}

	log.Printf("Using MALONE_SERIAL_PORT=%s",serialPort)

	//
	// Start the serial server
	//
	go serialServer(serialPort)

	//
	// Server up API endpoints
	//
	http.Handle("/metrics", promhttp.Handler())
	http.ListenAndServe(":8080", nil)

}


func serialServer(serialPort string) {

	//
	// Define a counter to keep track of the number of times the post has been delivered
	// The device in the field will send a signal when it detect post delivery 
	// A device attached to the serial port of one of the nodes will receive the signal 
	// We will want to increment the counter when the signal is received
	var postDeliveryCounter = prometheus.NewCounter(
		prometheus.CounterOpts{
				Name: "post_delivered_count",
				Help: "No of times post has been delivered",
		},
	)
	prometheus.MustRegister(postDeliveryCounter)

	//
	// Open the serial device
	//
	log.Printf("[serialServer] Open serial port: %v", serialPort)
	// Device is something like "/dev/ttyUSB0"
	cfg := &serial.Config{Name: serialPort, Baud: 115200}
	port, err := serial.OpenPort(cfg)
	// log.Printf("[DEBUG] port: %v, err: %v", port, err)
	if err != nil {
		log.Fatalf("error trying to open serial port %q. %v", serialPort, err)
	}
	log.Printf("[serialServer] using serial port: %v", serialPort)
	
	//
	// Monitor the serial port forever
	//
	log.Printf("[serialServer] Start read loop for serial port: %q...", serialPort)
	for {
		var err error
		var n int

		buf := make([]byte, 128)
		n, err = port.Read(buf)
		if err != nil {
			log.Fatalf("error trying to read serial port %q. %v", serialPort, err)
		}
		log.Printf("%q", buf[:n])

		//
		// For now consider any input as post delivery 
		//
		postDeliveryCounter.Inc()

	}
}
