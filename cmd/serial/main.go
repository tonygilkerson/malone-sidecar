package main

import (
	"log"
	"net/http"
	"os"
	"strings"


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

///////////////////////////////////////////////////////////////////////////////
//
//			Functions
//
///////////////////////////////////////////////////////////////////////////////
func serialServer(serialPort string) {

	//
	// MailboxDoorOpened
	//
	var mbxMailboxDoorOpenedCount = prometheus.NewCounter(
		prometheus.CounterOpts{
			Name: "mbx_mailbox_door_opened_count",
			Help: "No of times the mailbox door has been opened",
		},
	)
	prometheus.MustRegister(mbxMailboxDoorOpenedCount)

	var mbxMailboxDoorOpenedHeartbeatCount = prometheus.NewCounter(
		prometheus.CounterOpts{
			Name: "mbx_mailbox_door_opened_heartbeat_count",
			Help: "Heartbeat counter for mbxMailboxDoorOpened",
		},
	)
	prometheus.MustRegister(mbxMailboxDoorOpenedHeartbeatCount)


	//
	// MuleAlarm
	//
	var mbxMuleAlarmCount = prometheus.NewCounter(
		prometheus.CounterOpts{
			Name: "mbx_mule_alarm_count",
			Help: "No of times the mule alarm has gone off",
		},
	)
	prometheus.MustRegister(mbxMuleAlarmCount)

	var mbxMuleAlarmHeartbeatCount = prometheus.NewCounter(
		prometheus.CounterOpts{
			Name: "mbx_mule_alarm_heartbeat_count",
			Help: "Heartbeat counter for mbxMuleAlarm",
		},
	)
	prometheus.MustRegister(mbxMuleAlarmHeartbeatCount)

	//
	// Mailbox Temperature
	//
	var mbxTemperatureFahrenheit = prometheus.NewGauge(
		prometheus.GaugeOpts{
			Name: "mbx_temperature_fahrenheit",
			Help: "The temperature reading in fahrenheit from the device on the mailbox",
		},
	)
	prometheus.MustRegister(mbxTemperatureFahrenheit)

	//
	// Define a counter to keep track of the number of mbx heartbeats  
	//
	var mbxRoadMainLoopHeartbeatCount = prometheus.NewCounter(
		prometheus.CounterOpts{
			Name: "mbx_road_main_loop_heartbeat_count",
			Help: "Heartbeat counter for the main loop for the device down on the road ",
		},
	)
	prometheus.MustRegister(mbxRoadMainLoopHeartbeatCount)

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
		
		switch {

		case strings.Contains(msg,"MailboxTemperature"):
			log.Printf("DEBUG: %v", msg) 
			
			parts := strings.Split(msg, ":")
			if len(parts[1]) > 0 {
				f := float64(len(parts[1]))
				mbxTemperatureFahrenheit.Set(f)
				log.Printf("set MailboxTemperature to: %v", f)
			} else {
				log.Printf("Temperature reading expected but not found for input message: %v", msg)
			}

		case msg == "MuleAlarm":
			mbxMuleAlarmCount.Inc()
			log.Println("increment mbxMuleAlarmCount")

		case msg == "MuleAlarmHeartbeat":
			mbxMuleAlarmHeartbeatCount.Inc()
			log.Println("increment mbxMuleAlarmHeartbeatCount")

		case msg == "MailboxDoorOpened":
			mbxMailboxDoorOpenedCount.Inc()
			log.Println("increment mbxMailboxDoorOpenedCount")

		case msg == "MailboxDoorOpenedHeartbeat":
			mbxMailboxDoorOpenedHeartbeatCount.Inc()
			log.Println("increment mbxMailboxDoorOpenedHeartbeatCount")

		case msg == "RoadMainLoopHeartbeat":
			mbxRoadMainLoopHeartbeatCount.Inc()
			log.Println("increment mbxRoadMainLoopHeartbeatCount")

		default:
			log.Printf("No-op: %s\n", msg )
		}
	}
}
