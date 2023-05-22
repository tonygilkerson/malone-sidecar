package main

import (
	"log"
	"net/http"
	"os"

	"github.com/tonygilkerson/marty/pkg/marty"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/tarm/serial"
)

// Main
func main() {

	// Log to the console with date, time and filename prepended
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	//
	// Get environment Variables
	//
	serialPort, exists := os.LookupEnv("SERIAL_PORT")
	if !exists {
		log.Fatalln("SERIAL_PORT environment variable not set")
	}

	log.Printf("Using SERIAL_PORT=%s", serialPort)

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

	// Define a counter to keep track of the number of times the post has been delivered
	var mbxPostDeliveryCounter = prometheus.NewCounter(
		prometheus.CounterOpts{
			Name: "mvx_post_delivered_count",
			Help: "No of times post has been delivered",
		},
	)
	prometheus.MustRegister(mbxPostDeliveryCounter)

	// Define a counter to keep track of the number of times a car has arrived
	var mbxArrivedCounter = prometheus.NewCounter(
		prometheus.CounterOpts{
			Name: "mbx_arrived_count",
			Help: "No of times a car has arrived",
		},
	)
	prometheus.MustRegister(mbxArrivedCounter)

	//
	// Define a counter to keep track of the number of times a car has departed
	//
	var mbxDepartedCounter = prometheus.NewCounter(
		prometheus.CounterOpts{
			Name: "mbx_departed_count",
			Help: "No of times a car has departed",
		},
	)
	prometheus.MustRegister(mbxDepartedCounter)

	//
	// Define a counter to keep track of the number of times a car has departed
	//
	var mbxErrorCounter = prometheus.NewCounter(
		prometheus.CounterOpts{
			Name: "mbx_error_count",
			Help: "No of times an error has occurred counting cars",
		},
	)
	prometheus.MustRegister(mbxErrorCounter)

	//
	// Define a counter to keep track of the number of false alarms while counting cars  
	//
	var mbxFalseAlarmCounter = prometheus.NewCounter(
		prometheus.CounterOpts{
			Name: "mbx_false_count",
			Help: "No of times an error has occurred counting cars",
		},
	)
	prometheus.MustRegister(mbxFalseAlarmCounter)

	//
	// Define a counter to keep track of the number of times a sound is heard
	//
	var mbxHeardSoundCounter = prometheus.NewCounter(
		prometheus.CounterOpts{
			Name: "mbx_heard_sound_count",
			Help: "No of times a sound is heard",
		},
	)
	prometheus.MustRegister(mbxHeardSoundCounter)

	//
	// Define a counter to keep track of the number of mbx heartbeats  
	//
	var mbxHeartbeatCounter = prometheus.NewCounter(
		prometheus.CounterOpts{
			Name: "mbx_heartbeat_count",
			Help: "No of times a heartbeat has occurred",
		},
	)
	prometheus.MustRegister(mbxHeartbeatCounter)

	//
	// Open the serial device
	//
	log.Printf("Open serial port: %v", serialPort)

	// Device is something like "/dev/ttyUSB0"
	cfg := &serial.Config{Name: serialPort, Baud: 115200}
	port, err := serial.OpenPort(cfg)
	
	if err != nil {
		log.Fatalf("error trying to open serial port %q. %v", serialPort, err)
	}
	log.Printf("Using serial port: %v", serialPort)

	//
	// Monitor the serial port forever
	//
	buf := make([]byte, 128)
	log.Printf("Start read loop for serial port: %q...", serialPort)

	for {

		var err error
		var n int
		var msg string

		n, err = port.Read(buf)
		if err != nil {
			log.Fatalf("error trying to read serial port %q. %v\n", serialPort, err)
		}
		msg = string(buf[:n])
		
		switch msg {
		case string(marty.Arrived):
			mbxArrivedCounter.Inc()
			log.Println("increment mbxArrivedCounter")

		case string(marty.Departed):
			mbxDepartedCounter.Inc()
			log.Println("increment mbxDepartedCounter")

		case string(marty.Error):
			mbxErrorCounter.Inc()
			log.Println("increment mbxErrorCounter")

		case string(marty.FalseAlarm):
			mbxFalseAlarmCounter.Inc()
			log.Println("increment mbxFalseAlarmCounter")

		case "HeardSound":
			mbxHeardSoundCounter.Inc()
			log.Println("increment mbxHeardSoundCounter")

		case "MBX-HEARTBEAT":
			mbxHeartbeatCounter.Inc()
			log.Println("increment mbxHeartbeatCounter")

		default:
			log.Printf("No-op serial input: %s\n", msg )
		}
	}
}
