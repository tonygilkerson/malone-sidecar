package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/tarm/serial"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)


var postDeliveryCounter = prometheus.NewCounter(
	prometheus.CounterOpts{
			Name: "post_delivered_count",
			Help: "No of times post has been delivered",
	},
)

func monitorSerial(device string) {
	for {
		var err error
		var n int

		log.Printf("[monitorSerial] Open serial device: %v", device)
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

		log.Println("[monitorSerial] Read...")
		buf := make([]byte, 128)
		n, err = s.Read(buf)
		if err != nil {
			log.Fatal(err)
		}
		log.Printf("%q", buf[:n])

		//
		// For now consider any input as post delivery 
		//
		postDeliveryCounter.Inc()

	}
}

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
	go monitorSerial(device)

	//
	// Register metric
	//
	prometheus.MustRegister(postDeliveryCounter)

	//
	// Server up API endpoints
	//
	http.Handle("/metrics", promhttp.Handler())
	http.ListenAndServe(":8080", nil)

}
