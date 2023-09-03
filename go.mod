module github.com/tonygilkerson/serial-gateway

go 1.19

require (
	github.com/prometheus/client_golang v1.15.0
	github.com/tarm/serial v0.0.0-20180830185346-98f6abe2eb07
)

require (
	github.com/beorn7/perks v1.0.1 // indirect
	github.com/cespare/xxhash/v2 v2.2.0 // indirect
	github.com/golang/protobuf v1.5.3 // indirect
	github.com/matttproud/golang_protobuf_extensions v1.0.4 // indirect
	github.com/prometheus/client_model v0.3.0 // indirect
	github.com/prometheus/common v0.42.0 // indirect
	github.com/prometheus/procfs v0.9.0 // indirect
	github.com/tonygilkerson/mbx-iot v0.0.5
	golang.org/x/sys v0.6.0 // indirect
	google.golang.org/protobuf v1.30.0 // indirect
)

// DEVTODO - Remove me when done with development
replace github.com/tonygilkerson/mbx-iot v0.0.5 => ../mbx-iot
