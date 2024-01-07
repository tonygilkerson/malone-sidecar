# serial-gateway

A serial gateway that runs in the cluster. It will read a serial port mounted from the node and process the messages received

## Virtual Serial

Taken from this [stackoverflow post](https://stackoverflow.com/questions/22568878/emulate-serial-port)

To test, in one terminal run the following to create a virtual serial device.

```sh
socat PTY,link=./virtual-tty,raw,echo=0 -
```

Then you can run this in a different terminal to read from the device.  Now you can type input into the first terminal, hit enter to make it available in the serial gateway.

```sh
SERIAL_PORT=./virtual-tty go run cmd/serial/main.go
```

## Pub

The serial-gateway has a `/pub` endpoint that can be used for testing. A `POST` to `/pub` will result in the http body being written to the serial port on the host. As a result the LORA gateway will broadcast the contents for any LORA device to receive.

```sh
ssh -D 9995 weeble
kubectl ctx weeble
kubectl -n iot port-forward svc/serial-gateway 8080:8080 
curl -X POST "http://localhost:8080/pub" -d "a-message" 
```
